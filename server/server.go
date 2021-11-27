package main

import (
	"context"
	"errors"
	"io"
	"log"
	"net"
	"os"

	pb "github.com/beesaferoot/grpc-lnd/lnd"
	"github.com/jackc/pgx/v4"
	"google.golang.org/grpc"
)

const (
	port = ":52000"
)

var globContext = context.Background()

func NewLNDServer() *LNDServer {
	return &LNDServer{}
}

type LNDServer struct {
	pb.UnimplementedLNDServer
	conn              *pgx.Conn
	shouldCreatetable bool
}

func (s *LNDServer) GetNodesListByStatus(ctx context.Context, status *pb.Status) (*pb.NodeList, error) {
	var nodeList *pb.NodeList = &pb.NodeList{}
	rows, err := s.conn.Query(context.Background(), "select * from nodes where node.status = "+status.Value.String())
	if err != nil {
		return nil, err
	}
	// get node details
	for rows.Next() {
		nodeDetail := &pb.NodeDetail{}
		var status int32
		err = rows.Scan(&nodeDetail.Id, &nodeDetail.Nodename, &nodeDetail.IP, &nodeDetail.UserId.Value, &status)
		if err != nil {
			return nil, err
		}
		nodeList.Nodes = append(nodeList.Nodes, nodeDetail)
	}
	return nodeList, nil
}

func (s *LNDServer) DestroyNode(ctx context.Context, id *pb.NodeId) (*pb.NodeDetail, error) {
	//TODO
	return &pb.NodeDetail{Nodename: "node1", IP: "ip", UserId: &pb.UUID{Value: "user-id"}, CreateAt: &pb.Date{Year: 2021, Month: 11, Day: 27}}, nil
}

// SpawnNodes accepts a stream, store retrieved node details into db
// return Node list if no occured
func (s *LNDServer) SpawnNodes(stream pb.LND_SpawnNodesServer) error {
	var nodeList *pb.NodeList = &pb.NodeList{}
	for {
		nodeDetail, err := stream.Recv()
		if err == io.EOF {
			rows, err := s.conn.Query(context.Background(), "select * from nodes")
			if err != nil {
				return err
			}
			defer rows.Close()
			for rows.Next() {
				nodeDetail := &pb.NodeDetail{}
				var status int32
				err = rows.Scan(&nodeDetail.Id, &nodeDetail.Nodename, &nodeDetail.IP, &nodeDetail.UserId.Value, &status)
				if err != nil {
					return err
				}
				nodeList.Nodes = append(nodeList.Nodes, nodeDetail)
			}
			return stream.SendAndClose(nodeList)
		}

		if err != nil {
			return err
		}
		_, err = s.conn.Exec(globContext, "insert into nodes(name, ip, status)", nodeDetail.Nodename, nodeDetail.IP, nodeDetail.Status, nil)
		if err != nil {
			log.Printf("Error during table insert %v ", err)
			return errors.New("Error persiting nodes")
		}
	}
}

func (s *LNDServer) Run() error {
	lis, err := net.Listen("tcp", port)
	if err != nil {
		log.Fatalf("failed to listen: %v", err)
	}

	grpcServer := grpc.NewServer()
	pb.RegisterLNDServer(grpcServer, s)
	log.Printf("server listening at %v", lis.Addr())
	return grpcServer.Serve(lis)
}

func createNodeTable(ctx context.Context, conn *pgx.Conn) {
	createtableSQL := `
	create table if not exists node (
			id SERIAL PRIMARY KEY,
			name VAR(255) NOT NULL,
			ip VAR(255) NOT NULL,
			user_id VAR(255) NOT NULL,
			status int, 
			created_at date default CURRENT_TIMESTAMP, 
		);
	`
	_, err := conn.Exec(ctx, createtableSQL)
	if err != nil {
		log.Fatalf("Table creation failed: %v\n", err)
	}

}

func main() {
	db_url := os.Getenv("DATAB_URL")
	var lnd_node_server *LNDServer = NewLNDServer()
	conn, err := pgx.Connect(globContext, db_url)
	if err != nil {
		log.Fatalf("Unable to establish connection: %v", err)
	}
	// create Nodes table
	createNodeTable(globContext, conn)
	defer conn.Close(globContext)
	if err := lnd_node_server.Run(); err != nil {
		log.Fatalf("failed to serve: %v", err)
	}
}

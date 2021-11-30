package main

import (
	"context"
	"errors"
	"io"
	"log"
	"net"
	"os"
	"time"

	pb "github.com/beesaferoot/grpc-lnd/lnd"
	"github.com/jackc/pgx/v4"
	"google.golang.org/grpc"
	"google.golang.org/grpc/reflection"
)

const (
	port = ":52000"
)

var globContext context.Context = context.Background()

func NewLNDServer(conn *pgx.Conn) *LNDServer {
	return &LNDServer{conn: conn}
}

type LNDServer struct {
	pb.UnimplementedLNDServer
	conn              *pgx.Conn
	shouldCreatetable bool
}

func (s *LNDServer) GetNodesListByStatus(ctx context.Context, status *pb.Status) (*pb.NodeList, error) {
	var nodeList *pb.NodeList = &pb.NodeList{}
	rows, err := s.conn.Query(context.Background(), "select * from node where node.status=$1 ", status.Value)
	if err != nil {
		return nil, err
	}
	// get node details
	for rows.Next() {
		nodeDetail := &pb.NodeDetail{}
		var currentTime time.Time
		err = rows.Scan(&nodeDetail.Id, &nodeDetail.Nodename, &nodeDetail.IP, &nodeDetail.UserId, &nodeDetail.Status, &currentTime)
		if err != nil {
			return nil, err
		}
		nodeDetail.CreateAt = currentTime.Format("2000-01-01")
		nodeList.Nodes = append(nodeList.Nodes, nodeDetail)
	}
	return nodeList, nil
}

func (s *LNDServer) DestroyNode(ctx context.Context, id *pb.NodeId) (*pb.NodeDetail, error) {
	tag, err := s.conn.Exec(ctx, "delete from node where node.id=$1", id.Value)
	if err != nil {
		return nil, err
	}
	if tag.RowsAffected() != 1 {
		return nil, errors.New("node with id " + id.String() + " does not exist.")
	}

	return &pb.NodeDetail{Id: id.Value}, nil
}

// SpawnNodes accepts a stream, store retrieved node details into db
// return Node list if no error occured
func (s *LNDServer) SpawnNodes(stream pb.LND_SpawnNodesServer) error {
	var nodeList *pb.NodeList = &pb.NodeList{}
	for {
		nodeDetail, err := stream.Recv()
		if err == io.EOF {
			rows, err := s.conn.Query(context.Background(), "select * from node")
			if err != nil {
				return err
			}
			defer rows.Close()
			for rows.Next() {
				nodeDetail := &pb.NodeDetail{}
				var currentTime time.Time
				err = rows.Scan(&nodeDetail.Id, &nodeDetail.Nodename, &nodeDetail.IP, &nodeDetail.UserId, &nodeDetail.Status, &currentTime)
				nodeDetail.CreateAt = currentTime.Format("2000-01-01")
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
		_, err = s.conn.Exec(globContext, "insert into node(name, ip, user_id, status, created_at) values ($1,$2,$3,$4,$5)", nodeDetail.GetNodename(), nodeDetail.GetIP(), nodeDetail.GetUserId(), nodeDetail.GetStatus(), time.Now())
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
	reflection.Register(grpcServer)
	pb.RegisterLNDServer(grpcServer, s)
	log.Printf("server listening at %v", lis.Addr())
	return grpcServer.Serve(lis)
}

func createNodeTable(ctx context.Context, conn *pgx.Conn) {
	createtableSQL := `
	create table if not exists node (
			id SERIAL PRIMARY KEY,
			name text NOT NULL,
			ip text NOT NULL,
			user_id text NOT NULL,
			status int NOT NULL, 
			created_at date not NULL
		);
	`
	_, err := conn.Exec(ctx, createtableSQL)
	if err != nil {
		log.Fatalf("Table creation failed: %v\n", err)
	}

}

func main() {
	db_url := os.Getenv("DATAB_URL")
	conn, err := pgx.Connect(globContext, db_url)
	if err != nil {
		log.Fatalf("Unable to establish connection: %v", err)
	}
	var lnd_node_server *LNDServer = NewLNDServer(conn)
	// create Nodes table
	createNodeTable(globContext, conn)
	defer conn.Close(globContext)
	if err := lnd_node_server.Run(); err != nil {
		log.Fatalf("failed to serve: %v", err)
	}
}

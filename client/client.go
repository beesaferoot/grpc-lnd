package main

import (
	"context"
	"log"
	"time"

	pb "github.com/beesaferoot/grpc-lnd/lnd"
	"google.golang.org/grpc"
)

const (
	address = "localhost:52000"
)

var globContext = context.Background()

func main() {
	conn, err := grpc.Dial(address, grpc.WithInsecure(), grpc.WithBlock())
	if err != nil {
		log.Fatalf("failed to connect: %v", err)
	}
	defer conn.Close()
	client := pb.NewLNDClient(conn)
	ctx, cancel := context.WithTimeout(globContext, time.Second)
	defer cancel()
	nodeList := &pb.NodeList{}
	nodeDetail := &pb.NodeDetail{IP: "client-node-ip", UserId: &pb.UUID{Value: "client-id"}}
	nodeList.Nodes = append(nodeList.Nodes, nodeDetail)
	SpawnNodes(ctx, client, nodeList)
	GetNodesListByStatus(ctx, client, &pb.Status{Value: 1}) // get running nodes
	GetNodesListByStatus(ctx, client, &pb.Status{Value: 1}) // get failed nodes
	DestroyNode(ctx, client, &pb.NodeId{Value: 4})
}

func SpawnNodes(cxt context.Context, client pb.LNDClient, nodeList *pb.NodeList) {

	stream, err := client.SpawnNodes(cxt)
	if err != nil {
		log.Fatalf("%v.SpawnNodes(_) = %v", client, err)
	}
	for _, nodeDetail := range nodeList.Nodes {
		if err = stream.Send(nodeDetail); err != nil {
			log.Fatalf("%v.Send(%v) = %v", stream, nodeDetail, err)
		}
	}
	availNodes, err := stream.CloseAndRecv()
	if err != nil {
		log.Fatalf("%v.CloseAndRecv() got error %v, want %v", stream, err, nil)
	}
	log.Printf("List of Nodes: %v", availNodes)
}

func GetNodesListByStatus(ctx context.Context, client pb.LNDClient, status *pb.Status) {

	var nodeList *pb.NodeList
	nodeList, err := client.GetNodesListByStatus(ctx, status)
	if err != nil {
		log.Fatalf("Error retrieving Nodes: %v", err)
	}
	log.Printf("Nodes: %v", nodeList.Nodes)
}

func DestroyNode(ctx context.Context, client pb.LNDClient, nodeId *pb.NodeId) {
	var nodeDetail *pb.NodeDetail
	nodeDetail, err := client.DestroyNode(ctx, nodeId)
	if err != nil {
		log.Fatalf("Error destorying Node: %v", err)
	}
	log.Printf("Node detail: %v", nodeDetail)
}

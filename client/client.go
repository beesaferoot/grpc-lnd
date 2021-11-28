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
	nodeList.Nodes = append(nodeList.Nodes, &pb.NodeDetail{Nodename: "node1", IP: "client-node-ip", UserId: "client-id", Status: pb.Status_RUNNING})
	nodeList.Nodes = append(nodeList.Nodes, &pb.NodeDetail{Nodename: "node2", IP: "client-node-ip1", UserId: "client-id1", Status: pb.Status_RUNNING})
	nodeList.Nodes = append(nodeList.Nodes, &pb.NodeDetail{Nodename: "node3", IP: "client-node-ip3", UserId: "client-id3", Status: pb.Status_FAILED})
	SpawnNodes(ctx, client, nodeList)
	GetNodesListByStatus(ctx, client, &pb.Status{Value: pb.Status_RUNNING}) // get running nodes
	GetNodesListByStatus(ctx, client, &pb.Status{Value: pb.Status_FAILED})  // get failed nodes
	DestroyNode(ctx, client, &pb.NodeId{Value: 1})
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
	_, err = stream.CloseAndRecv()
	if err != nil {
		log.Fatalf("%v.CloseAndRecv() got error %v, want %v", stream, err, nil)
	}
}

func GetNodesListByStatus(ctx context.Context, client pb.LNDClient, status *pb.Status) {

	var nodeList *pb.NodeList
	nodeList, err := client.GetNodesListByStatus(ctx, status)
	if err != nil {
		log.Fatalf("Error retrieving Nodes: %v", err)
	}
	log.Printf("Nodes: %v\n\n\n", nodeList.Nodes)
}

func DestroyNode(ctx context.Context, client pb.LNDClient, nodeId *pb.NodeId) {
	var nodeDetail *pb.NodeDetail
	nodeDetail, err := client.DestroyNode(ctx, nodeId)
	if err != nil {
		log.Fatalf("Error destorying Node: %v", err)
	}
	log.Printf("Node detail: %v", nodeDetail)
}

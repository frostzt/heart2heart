package main

import (
	"context"
	"log"
	"net"

	pb "keeper/proto"

	"google.golang.org/grpc"
)

type KeeperAPIGRPCServer struct {
	pb.KeeperAPIServiceServer
}

func (s *KeeperAPIGRPCServer) HealthCheck(ctx context.Context, req *pb.HealthCheckRequest) (*pb.HealthCheckResponse, error) {
	return &pb.HealthCheckResponse{
		Status:  true,
		Message: "Service is healthy",
	}, nil
}

func InitGRPCServer() {
	listen, err := net.Listen("tcp", ":50051")
	if err != nil {
		log.Fatalf("Failed to listen: %v", err)
	}

	server := grpc.NewServer()
	pb.RegisterKeeperAPIServiceServer(server, &KeeperAPIGRPCServer{})

	log.Println("🏗️ gRPC server listening on :50051")
	if err := server.Serve(listen); err != nil {
		log.Fatalf("Failed to serve: %v", err)
	}
}

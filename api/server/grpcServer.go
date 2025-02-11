package server

import (
	"context"
	"fmt"
	"google.golang.org/grpc"
	"log"
	"net"
	pb "userService/generated/proto"
	"userService/internal/user"
)

type userServiceServer struct {
	pb.UnimplementedUserServiceServer
}

func StartRpc() {
	lis, err := net.Listen("tcp", ":50051")
	if err != nil {
		log.Fatalf("failed to listen: %v", err)
	}

	server := grpc.NewServer()
	pb.RegisterUserServiceServer(server, &userServiceServer{})

	fmt.Print("grpc server connect success\n")
	if err := server.Serve(lis); err != nil {
		log.Fatalf("failed to serve: %v", err)
	}
}

func (s *userServiceServer) GetUser(ctx context.Context, req *pb.GetUserRequest) (*pb.GetUserResponse, error) {
	return &pb.GetUserResponse{
		UserId: req.UserId,
		Name:   user.GetUserByID(req.UserId).Name,
	}, nil
}

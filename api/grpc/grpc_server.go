package grpc

import (
	"context"
	"fmt"
	"google.golang.org/grpc"
	"log"
	"net"
	pb "userService/generated/proto"
	"userService/internal/repository"
)

type UserServiceServer struct {
	pb.UnimplementedUserServiceServer
	userService repository.IUserService
}

func StartRpc(userService repository.IUserService) error {
	lis, err := net.Listen("tcp", ":50051")
	if err != nil {
		log.Fatalf("failed to listen: %v", err)
		return err
	}

	server := grpc.NewServer()
	pb.RegisterUserServiceServer(server, &UserServiceServer{
		userService: userService,
	})

	fmt.Print("grpc server connect success\n")
	if err := server.Serve(lis); err != nil {
		log.Fatalf("failed to serve: %v", err)
		return err
	}

	return nil
}

func (s *UserServiceServer) GetUser(ctx context.Context, req *pb.GetUserRequest) (*pb.GetUserResponse, error) {
	user, err := s.userService.GetUserByID(req.UserId)
	if err != nil {
		return nil, err
	}
	return &pb.GetUserResponse{
		UserId: req.UserId,
		Name:   user.Name,
	}, nil
}

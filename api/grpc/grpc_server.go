package grpc

import (
	"context"
	"fmt"
	"log"
	"net"

	"google.golang.org/grpc"
	pb "userService/generated/proto"
	"userService/internal/repository"
)

type Service struct {
	pb.UnimplementedUserServiceServer
	userService repository.UserService
}

func Listen(userService repository.UserService) error {
	listen, err := net.Listen("tcp", ":50051")
	if err != nil {
		return fmt.Errorf("failed to listen: %v", err)
	}

	server := grpc.NewServer()
	pb.RegisterUserServiceServer(server, &Service{
		userService: userService,
	})

	log.Print("grpc server connect success\n")

	if errServ := server.Serve(listen); errServ != nil {
		return fmt.Errorf("failed to serve: %v", errServ)
	}

	return nil
}

func (s *Service) GetUser(ctx context.Context, req *pb.GetUserRequest) (*pb.GetUserResponse, error) {
	user, err := s.userService.GetUserByID(ctx, req.UserId)
	if err != nil {
		return nil, fmt.Errorf("failed to get user: %v", err)
	}

	return &pb.GetUserResponse{
		UserId: req.UserId,
		Name:   user.Name,
	}, nil
}

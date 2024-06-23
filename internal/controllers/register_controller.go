package controllers

import (
	"context"

	"github.com/adi-kmt/ai-streak-backend-go/proto"
)

func (s *Server) RegisterHandler(ctx context.Context, request *proto.RegisterRequest) (*proto.AuthResponse, error) {
	token, err := s.Service.AddUser(request.Username, request.Password, request.Isstandingforelection)
	if err != nil {
		return &proto.AuthResponse{
			Token: "",
			Error: err.Message,
		}, nil
	}
	return &proto.AuthResponse{
		Token: token,
		Error: "",
	}, nil
}

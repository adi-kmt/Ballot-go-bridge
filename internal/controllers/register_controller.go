package controllers

import (
	"context"

	"github.com/adi-kmt/ai-streak-backend-go/proto"
)

func (s *Server) RegisterHandler(ctx context.Context, request *proto.RegisterRequest) (*proto.AuthResponse, error) {
	token, err := s.Service.AddUser(request.Username, request.Password, request.Isstandingforelection)
	return &proto.AuthResponse{
		Token: token,
		Error: err.Message,
	}, nil
}

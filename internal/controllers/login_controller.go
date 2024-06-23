package controllers

import (
	"context"

	"github.com/adi-kmt/ai-streak-backend-go/proto"
)

func (s *Server) LoginHandler(ctx context.Context, request *proto.LoginRequest) (*proto.AuthResponse, error) {
	token, err := s.Service.ValidateUser(request.Username, request.Password)

	return &proto.AuthResponse{
		Token: token,
		Error: err.Message,
	}, nil
}

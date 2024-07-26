package controllers

import (
	"github.com/adi-kmt/voting-backend-go-grpc/internal/services"
	"github.com/adi-kmt/voting-backend-go-grpc/proto"
)

type Server struct {
	proto.UnimplementedAuthServiceServer
	Service *services.UserService
}

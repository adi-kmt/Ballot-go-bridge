package controllers

import (
	"github.com/adi-kmt/ai-streak-backend-go/internal/services"
	"github.com/adi-kmt/ai-streak-backend-go/proto"
)

type Server struct {
	proto.UnimplementedAuthServiceServer
	Service services.UserService
}

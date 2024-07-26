package injection

import (
	"context"

	"github.com/adi-kmt/voting-backend-go-grpc/internal/config"
	"github.com/adi-kmt/voting-backend-go-grpc/internal/repositories"
	"github.com/adi-kmt/voting-backend-go-grpc/internal/services"
)

// Simple DI component
func Injector(ctx context.Context) (*services.UserService, *services.VotingService) {
	redisClient := config.GetRedisClient(":6379", "", 0)
	repository := repositories.NewRedisRepository(redisClient, ctx)
	userService := services.NewUserService(repository)
	votingService := services.NewVotingService(repository)

	return userService, votingService
}

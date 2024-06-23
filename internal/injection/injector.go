package injection

import (
	"context"

	"github.com/adi-kmt/ai-streak-backend-go/internal/config"
	"github.com/adi-kmt/ai-streak-backend-go/internal/repositories"
	"github.com/adi-kmt/ai-streak-backend-go/internal/services"
)

// Simple DI component
func Injector(ctx context.Context) (*services.UserService, *services.VotingService) {
	redisClient := config.GetRedisClient(":6379", "", 0)
	repository := repositories.NewRedisRepository(redisClient, ctx)
	userService := services.NewUserService(repository)
	votingService := services.NewVotingService(repository)

	return userService, votingService
}

package injection

import (
	"github.com/adi-kmt/ai-streak-backend-go/internal/config"
	"github.com/adi-kmt/ai-streak-backend-go/internal/repositories"
	"github.com/adi-kmt/ai-streak-backend-go/internal/services"
)

func Injector() (*services.UserService, *services.VotingService) {
	redisClient := config.GetRedisClient(":5432", "", 0)
	repository := repositories.NewRedisRepository(redisClient)
	userService := services.NewUserService(repository)
	votingService := services.NewVotingService(repository)

	return userService, votingService
}

package repository

import (
	"context"

	"github.com/redis/go-redis/v9"
)

// Using an interface for easily testing the service
type IRepo interface {
	GetCurrentVoteSapshot() ([]leaderBoardItem, error)
	AddUser(userName, password string, isStandingForElection bool) error
	AddVote(userName, VotedForUserName string) error
	CheckUserExists(userName string) (bool, error)
	CheckUserIsStandingForElection(userName string) (bool, error)
	CheckUserHasVoted(userName string) (bool, error)
	GetUserPassword(userName string) (string, error)
}

// Redis Repository which implements all the IRepo functions
type RedisRepository struct {
	ctx         context.Context
	redisClient *redis.Client
}

package repositories

import (
	"context"

	"github.com/adi-kmt/ai-streak-backend-go/internal/entities"
	"github.com/adi-kmt/ai-streak-backend-go/internal/messages"
	"github.com/redis/go-redis/v9"
)

const (
	listOfUserNames                = "userNames"
	listOfUserStandingForElection  = "usersStandingForElection"
	sortedSetContainingLeaderBoard = "leaderBoard"

	passwordKey              = "password"
	isStandingForElectionKey = "isStandingForElection"
	hasVotedKey              = "hasVoted"
)

// Redis Repository which implements all the IRepo functions
type RedisRepository struct {
	ctx         context.Context
	redisClient *redis.Client
}

func NewRedisRepository(redisClient *redis.Client) *RedisRepository {
	return &RedisRepository{
		redisClient: redisClient,
	}
}

// Creating a hashset with the user key, and if standing in the election add in that
func (r *RedisRepository) AddUser(userName, password string, isStandingForElection bool) *messages.AppError {
	pipe := r.redisClient.TxPipeline()
	// First creating a user hashset
	err := pipe.HSet(r.ctx, userName, passwordKey, password, isStandingForElectionKey, isStandingForElection, hasVotedKey, false).Err()
	if err != nil {
		pipe.Discard()
		return messages.InternalServerError("Failed to add user")
	}
	// Then adding to the usernames hashset to keep usernames unique
	err1 := pipe.SAdd(r.ctx, listOfUserNames, userName).Err()
	if err1 != nil {
		pipe.Discard()
		return messages.InternalServerError("Failed to validate username is unique")
	}
	// if user is standing for election, then we are adding to th users standing for election hashset
	if isStandingForElection {
		err0 := pipe.SAdd(r.ctx, listOfUserStandingForElection, userName).Err()
		if err0 != nil {
			pipe.Discard()
			return messages.InternalServerError("Failed to make user standing for election")
		}
	}
	_, err = pipe.Exec(r.ctx)
	if err != nil {
		pipe.Discard()
		return messages.InternalServerError("Failed to add user")
	}
	return nil
}

// this is to check if the username is unique
func (r *RedisRepository) CheckUserExists(userName string) (bool, *messages.AppError) {
	bool, err := r.redisClient.SIsMember(r.ctx, listOfUserNames, userName).Result()
	if err != nil {
		return false, messages.InternalServerError("Failed to check if user exists")
	}
	return bool, nil
}

// this is to check if the user is standing for election
func (r *RedisRepository) CheckUserIsStandingForElection(userName string) (bool, *messages.AppError) {
	bool, err := r.redisClient.SIsMember(r.ctx, listOfUserStandingForElection, userName).Result()
	if err != nil {
		return false, messages.InternalServerError("Failed to check if user is standing for election")
	}
	return bool, nil
}

// this is to check if the user has already voted
func (r *RedisRepository) CheckUserHasVoted(userName string) (bool, *messages.AppError) {
	hasVotedString, err := r.redisClient.HGet(r.ctx, userName, hasVotedKey).Result()
	if err != nil {
		return false, messages.InternalServerError("Failed to check if user has voted")
	}
	return hasVotedString == "true", nil
}

// this is to get the password of the user with particular username
func (r *RedisRepository) GetUserPassword(userName string) (string, *messages.AppError) {
	passeord, err := r.redisClient.HGet(r.ctx, userName, passwordKey).Result()
	if err != nil {
		return "", messages.InternalServerError("Failed to get user password")
	}
	return passeord, nil
}

// this is to add vote to the user
func (r *RedisRepository) AddVote(userName, VotedForUserName string) *messages.AppError {
	pipe := r.redisClient.TxPipeline()

	// We set the voter's has voted to true
	_, err0 := pipe.HSet(r.ctx, userName, "hasVoted", true).Result()
	if err0 != nil {
		pipe.Discard()
		return messages.InternalServerError("Failed to Check if user has voted")
	}
	// We increment the vote count list to keep track of the votes
	err1 := pipe.ZIncrBy(r.ctx, sortedSetContainingLeaderBoard, 1, VotedForUserName).Err()
	if err1 != nil {
		pipe.Discard()
		return messages.InternalServerError("Failed to increase vote count")
	}
	_, err := pipe.Exec(r.ctx)
	if err != nil {
		pipe.Discard()
		return messages.InternalServerError("Failed to add vote")
	}
	return nil
}

// this is to get the current vote snapshot
func (r *RedisRepository) GetCurrentVoteSapshot() ([]entities.LeaderBoardItem, *messages.AppError) {
	leadersWithScores, err := r.redisClient.ZRevRangeWithScores(r.ctx, sortedSetContainingLeaderBoard, 0, -1).Result()
	if err != nil {
		return nil, messages.InternalServerError("Failed to get current vote snapshot")
	}

	var leaderBoard []entities.LeaderBoardItem
	for _, leader := range leadersWithScores {
		leaderBoard = append(leaderBoard, entities.LeaderBoardItem{
			LeaderName: leader.Member.(string),
			Score:      int(leader.Score),
		})
	}
	return leaderBoard, nil
}

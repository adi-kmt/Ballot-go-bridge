package repository

import (
	"context"

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

type leaderBoardItem struct {
	leaderName string
	score      int
}

type Repository struct {
	ctx         context.Context
	redisClient *redis.Client
}

func NewRepository(redisClient *redis.Client) *Repository {
	return &Repository{
		redisClient: redisClient,
	}
}

// Creating a hashset with the user key, and if standing in the election add in that
func (r *Repository) AddUser(userName, password string, isStandingForElection bool) error {
	pipe := r.redisClient.TxPipeline()
	// First creating a user hashset
	err := pipe.HSet(r.ctx, userName, passwordKey, password, isStandingForElectionKey, isStandingForElection, hasVotedKey, false).Err()
	// Then adding to the usernames hashset to keep usernames unique
	err1 := pipe.SAdd(r.ctx, listOfUserNames, userName).Err()
	if err1 != nil {
		pipe.Discard()
		return err1
	}
	// if user is standing for election, then we are adding to th users standing for election hashset
	if isStandingForElection {
		err0 := pipe.SAdd(r.ctx, listOfUserStandingForElection, userName).Err()
		if err0 != nil {
			pipe.Discard()
			return err0
		}
	}
	if err != nil {
		pipe.Discard()
		return err
	}
	_, err = pipe.Exec(r.ctx)
	return err
}

// this is to check if the username is unique
func (r *Repository) CheckUserExists(userName string) (bool, error) {
	return r.redisClient.SIsMember(r.ctx, listOfUserNames, userName).Result()
}

// this is to check if the user is standing for election
func (r *Repository) CheckUserIsStandingForElection(userName string) (bool, error) {
	return r.redisClient.SIsMember(r.ctx, listOfUserStandingForElection, userName).Result()
}

// this is to check if the user has already voted
func (r *Repository) CheckUserHasVoted(userName string) (bool, error) {
	hasVotedString, err := r.redisClient.HGet(r.ctx, userName, hasVotedKey).Result()
	if err != nil {
		return false, err
	}
	return hasVotedString == "true", nil
}

// this is to get the password of the user with particular username
func (r *Repository) GetUserPassword(userName string) (string, error) {
	return r.redisClient.HGet(r.ctx, userName, passwordKey).Result()
}

// this is to add vote to the user
func (r *Repository) AddVote(userName, VotedForUserName string) error {
	pipe := r.redisClient.TxPipeline()

	// We set the voter's has voted to true
	_, err0 := pipe.HSet(r.ctx, userName, "hasVoted", true).Result()
	if err0 != nil {
		pipe.Discard()
		return err0
	}
	// We increment the vote count list to keep track of the votes
	err1 := pipe.ZIncrBy(r.ctx, sortedSetContainingLeaderBoard, 1, VotedForUserName).Err()
	if err1 != nil {
		pipe.Discard()
		return err1
	}
	_, err := pipe.Exec(r.ctx)
	return err
}

// this is to get the current vote snapshot
func (r *Repository) GetCurrentVoteSapshot() ([]leaderBoardItem, error) {
	leadersWithScores, err := r.redisClient.ZRevRangeWithScores(r.ctx, sortedSetContainingLeaderBoard, 0, -1).Result()
	if err != nil {
		return nil, err
	}

	var leaderBoard []leaderBoardItem
	for _, leader := range leadersWithScores {
		leaderBoard = append(leaderBoard, leaderBoardItem{
			leaderName: leader.Member.(string),
			score:      int(leader.Score),
		})
	}
	return leaderBoard, nil
}

package services

import (
	"github.com/adi-kmt/voting-backend-go-grpc/internal/entities"
	"github.com/adi-kmt/voting-backend-go-grpc/internal/messages"
)

var sampleUserName = []string{"user1", "user2", "user3"}
var samplePasswords = []string{"password1", "password2", "password3"}

type FakeRepo struct{}

func (f *FakeRepo) AddUser(userName, password string, isStandingForElection bool) *messages.AppError {
	return nil
}

func (f *FakeRepo) AddVote(userName, VotedForUserName string) *messages.AppError {
	return nil
}

func (f *FakeRepo) CheckUserExists(userName string) (bool, *messages.AppError) {
	for _, v := range sampleUserName {
		if v == userName {
			return true, nil
		}
	}
	return false, nil
}

func (f *FakeRepo) CheckUserIsStandingForElection(userName string) (bool, *messages.AppError) {
	if userName == "user2" {
		return true, nil
	}
	return false, nil
}

func (f *FakeRepo) CheckUserHasVoted(userName string) (bool, *messages.AppError) {
	return true, nil
}

func (f *FakeRepo) GetUserPassword(userName string) (string, *messages.AppError) {
	for i, v := range sampleUserName {
		if v == userName {
			return samplePasswords[i], nil
		}
	}
	return "", nil
}

func (f *FakeRepo) SubscribeToLeaderBoardUpdates(sub *entities.LeaderBoardSubscription) {
	return
}

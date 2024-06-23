package services

import (
	"fmt"
	"testing"
)

func getVotingServiceWithFakeRepo() *VotingService {
	return NewVotingService(&FakeRepo{})
}

// Only testing the add vote, since the get snapshot is just returning the repo functoin

func TestAddVoteUserDoesNotExist(t *testing.T) {
	service := getVotingServiceWithFakeRepo()
	err := service.AddVote("userName", "userName2")
	if err == nil {
		t.Errorf("User does not exist")
	} else if err.Message == "User userName does not exist" {
		fmt.Printf("User userName does not exist\n")
	} else {
		t.Errorf("Testcase failed")
	}
}

func TestAddVoteUserIsNotStandingForElection(t *testing.T) {
	service := getVotingServiceWithFakeRepo()
	err := service.AddVote("user1", "userName2")
	if err == nil {
		t.Errorf("User is not standing for election")
	} else if err.Message == "User userName2 is not standing for election" {
		fmt.Printf("User userName2 is not standing for election\n")
	} else {
		t.Errorf("Testcase failed")
	}
}

func TestAddVote(t *testing.T) {
	service := getVotingServiceWithFakeRepo()
	err := service.AddVote("user1", "user2")
	if err != nil {
		t.Errorf("Testcase failed")
	}
	fmt.Println("User voted successfully")
}

package repositories

import "github.com/adi-kmt/ai-streak-backend-go/internal/entities"

// Using an interface for easily testing the service
type IRepo interface {
	GetCurrentVoteSapshot() ([]entities.LeaderBoardItem, error)
	AddUser(userName, password string, isStandingForElection bool) error
	AddVote(userName, VotedForUserName string) error
	CheckUserExists(userName string) (bool, error)
	CheckUserIsStandingForElection(userName string) (bool, error)
	CheckUserHasVoted(userName string) (bool, error)
	GetUserPassword(userName string) (string, error)
}

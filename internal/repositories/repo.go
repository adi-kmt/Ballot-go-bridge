package repositories

import (
	"github.com/adi-kmt/ai-streak-backend-go/internal/entities"
	"github.com/adi-kmt/ai-streak-backend-go/internal/messages"
)

// Using an interface for easily testing the service
type IRepo interface {
	SubscribeToLeaderBoardUpdates(sub *entities.LeaderBoardSubscription)
	AddUser(userName, password string, isStandingForElection bool) *messages.AppError
	AddVote(userName, VotedForUserName string) *messages.AppError
	CheckUserExists(userName string) (bool, *messages.AppError)
	CheckUserIsStandingForElection(userName string) (bool, *messages.AppError)
	CheckUserHasVoted(userName string) (bool, *messages.AppError)
	GetUserPassword(userName string) (string, *messages.AppError)
}

package entities

type LeaderBoardSubscription struct {
	UpdateChan chan []LeaderBoardItem
	StopChan   chan bool
}

func NewLeaderBoardSubscription() *LeaderBoardSubscription {
	return &LeaderBoardSubscription{
		UpdateChan: make(chan []LeaderBoardItem),
		StopChan:   make(chan bool),
	}
}

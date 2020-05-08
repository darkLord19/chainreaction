package constants

const (
	// MoveRcvMsg is msg type of received move msgs from players
	MoveRcvMsg int = iota
	// MoveBcastMsg is msg type of move to be broadcasted to players
	MoveBcastMsg
	// StateUpBcastMsg is msg type of upadated game state to be broadcasted to players
	StateUpBcastMsg
	// UserWonMsg is msg type of winner to be broadcasted to players
	UserWonMsg
	// InvalidMoveMsg is msg type of invalid move to be unicasted to player who originated move
	InvalidMoveMsg
)

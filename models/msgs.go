package models

// MoveMsg struct is used to get Move messages from websocket client
type MoveMsg struct {
	MsgType        int    `json:"msg_type"`
	XPos           int    `json:"xpos"`
	YPos           int    `json:"ypos"`
	Color          string `json:"color"`
	PlayerUserName string `json:"player_username"`
}

// NewStateMsg struct is used to represent board update for websocket broadcast
type NewStateMsg struct {
	MsgType     int       `json:"msg_type"`
	NewCurrTurn string    `json:"new_currturn"`
	NewBoard    [][]Pixel `json:"new_board"`
}

// WinnerMsg struct is used to send winner notification to users
type WinnerMsg struct {
	MsgType int    `json:"msg_type"`
	Winner  Player `json:"winner"`
}

// ErrMsg struct is used to notify user with err msgs
type ErrMsg struct {
	MsgType int    `json:"msg_type"`
	ErrStr  string `json:"errstr"`
}

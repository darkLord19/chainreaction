package models

// Move struct is used to get Move messages from websocket client
type Move struct {
	MsgType        int    `json:"msg_type"`
	XPos           int    `json:"xpos"`
	YPos           int    `json:"ypos"`
	PlayerUserName string `json:"player_username"`
}

// NewState struct is used to represent board update for websocket broadcast
type NewState struct {
	MsgType     int       `json:"msg_type"`
	NewCurrTurn string    `json:"new_currturn"`
	NewBoard    [][]Pixel `json:"new_board"`
}

// Winner struct is used to send winner notification to users
type Winner struct {
	MsgType int    `json:"msg_type"`
	Winner  Player `json:"winner"`
}

// Err struct is used to notify user with err msgs
type Err struct {
	MsgType int    `json:"msg_type"`
	ErrStr  string `json:"errstr"`
}

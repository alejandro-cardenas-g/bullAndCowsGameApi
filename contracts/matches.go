package contracts

import "github.com/alejandro-cardenas-g/bullAndCowsApp/internal/domain"

type CreateRoomCommand struct {
	Username string `json:"username" validate:"required"`
}
type CreateRoomResponse struct {
	RoomId string         `json:"room_id"`
	Player PlayerResponse `json:"player"`
}

type JoinRoomCommand struct {
	Username string `json:"username" validate:"required"`
	RoomId   string
}

type JoinRoomResponse struct {
	RoomId string         `json:"room_id"`
	Player PlayerResponse `json:"player"`
}

type SetCombinationCommand struct {
	PlayerId    string `json:"player_id" validate:"required"`
	Combination int    `json:"combination" validate:"required"`
	RoomId      string
}

type SuccessResponse struct {
	Success bool `json:"success"`
}

type StartMatchResponse struct {
	IsTurnOf string `json:"is_turn_of"`
}

type MakeGuessCommand struct {
	Guess    int    `json:"guess"`
	PlayerId string `json:"player_id"`
	RoomId   string
}

type MakeGuessResponse struct {
	IsWinner bool                `json:"is_winner"`
	Guesses  domain.MatchGuesses `json:"guesses"`
}

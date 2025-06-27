package contracts

import "context"

type IMatchesService interface {
	CreateRoom(ctx context.Context, createRoomCommand CreateRoomCommand) (*CreateRoomResponse, error)
	JoinRoom(ctx context.Context, joinRoomCommand JoinRoomCommand) (*JoinRoomResponse, error)
	SetCombination(ctx context.Context, setCombinationCommand SetCombinationCommand) (*SuccessResponse, error)
	StartGame(ctx context.Context, roomId string) (*StartMatchResponse, error)
	MakeGuess(ctx context.Context, command MakeGuessCommand) (*MakeGuessResponse, error)
	RestartGame(ctx context.Context, roomId string) (*SuccessResponse, error)
}

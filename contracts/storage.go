package contracts

import (
	"context"

	"github.com/alejandro-cardenas-g/bullAndCowsApp/internal/domain"
)

type SetPlayersCommand struct {
	RoomId  string
	Players domain.MatchPlayers
}

type SetOpponentCombinationsCommand struct {
	RoomId       string
	Combinations domain.MatchOpponentCombinations
}

type IMatchesRepository interface {
	CreateMatch(ctx context.Context, player domain.Player) (*domain.Match, error)
	GetRoomPlayers(ctx context.Context, roomId string) (domain.MatchPlayers, error)
	SetPlayersAndFillRoom(ctx context.Context, command SetPlayersCommand) error
	GetMatchStatusById(ctx context.Context, roomId string) (domain.MatchStatus, error)
	SetPlayerCombination(ctx context.Context, command SetOpponentCombinationsCommand) error
	GetPlayersAndCombinations(ctx context.Context, roomId string) (*domain.Match, error)
	GetAllButGuesses(ctx context.Context, roomId string) (*domain.Match, error)
	ChangeStatusAndTurn(ctx context.Context, roomId string, status domain.MatchStatus, isTurnOf string) error
}

type Storage struct {
	MatchesRepository IMatchesRepository
}

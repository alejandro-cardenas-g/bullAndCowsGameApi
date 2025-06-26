package services

import (
	"context"
	"errors"
	"fmt"

	"github.com/alejandro-cardenas-g/bullAndCowsApp/contracts"
	"github.com/alejandro-cardenas-g/bullAndCowsApp/internal/domain"
)

var (
	ErrCanNotAddAnotherPlayer = fmt.Errorf("can not add another player to this room")
	ErrMatchNotFullRoom       = fmt.Errorf("match is being played already or room is not completed")
	ErrInvalidCombination     = fmt.Errorf("invalid combination")
	ErrMatchNotFound          = fmt.Errorf("match not found")
	ErrExpectingCombinations  = fmt.Errorf("can not start game until players set combinations")
)

type MatchesService struct {
	storage contracts.Storage
}

func NewMatchesService(storage contracts.Storage) contracts.IMatchesService {
	return &MatchesService{
		storage: storage,
	}
}

func (s *MatchesService) CreateRoom(ctx context.Context, command contracts.CreateRoomCommand) (*contracts.CreateRoomResponse, error) {

	playerId := domain.GeneratePlayerId()
	match, err := s.storage.MatchesRepository.CreateMatch(ctx, domain.Player{Id: playerId, Username: command.Username})

	if err != nil {
		return nil, err
	}

	resp := &contracts.CreateRoomResponse{
		RoomId: match.RoomId,
		Player: contracts.PlayerResponse{
			Username: command.Username,
			Id:       playerId,
		},
	}

	return resp, nil
}

func (s *MatchesService) JoinRoom(ctx context.Context, joinRoomCommand contracts.JoinRoomCommand) (*contracts.JoinRoomResponse, error) {
	players, err := s.storage.MatchesRepository.GetRoomPlayers(ctx, joinRoomCommand.RoomId)
	if err != nil {
		if err == domain.ErrEmptyResult {
			return nil, ErrMatchNotFound
		}
		return nil, err
	}

	if len(players) != 1 {
		return nil, ErrCanNotAddAnotherPlayer
	}

	newPlayer := domain.Player{
		Id:       domain.GeneratePlayerId(),
		Username: joinRoomCommand.Username,
	}

	players[newPlayer.Id] = newPlayer

	if err = s.storage.MatchesRepository.SetPlayersAndFillRoom(ctx, contracts.SetPlayersCommand{
		RoomId:  joinRoomCommand.RoomId,
		Players: players,
	}); err != nil {
		return nil, err
	}

	return &contracts.JoinRoomResponse{
		RoomId: joinRoomCommand.RoomId,
		Player: contracts.PlayerResponse{
			Id:       newPlayer.Id,
			Username: newPlayer.Username,
		},
	}, nil
}

func (s *MatchesService) SetCombination(ctx context.Context, command contracts.SetCombinationCommand) (*contracts.SuccessResponse, error) {
	strCombination := fmt.Sprint(command.Combination)
	if err := domain.ValidateCombination(strCombination); err != nil {
		switch err {
		case domain.ErrInvalidCombination, domain.ErrInvalidUniqueCombination:
			return nil, fmt.Errorf("%w: "+err.Error(), ErrInvalidCombination)
		}
		return nil, err
	}

	status, err := s.storage.MatchesRepository.GetMatchStatusById(ctx, command.RoomId)
	if err != nil {
		if err == domain.ErrEmptyResult {
			return nil, ErrMatchNotFound
		}
		return nil, err
	}

	if status != domain.MatchStateFullRoom {
		return nil, ErrMatchNotFullRoom
	}

	match, err := s.storage.MatchesRepository.GetPlayersAndCombinations(ctx, command.RoomId)

	if err != nil {
		return nil, err
	}

	for key := range match.Players {
		if key == command.PlayerId {
			continue
		}
		match.OpponentsCombinations[key] = strCombination
		break
	}

	if err := s.storage.MatchesRepository.SetPlayerCombination(ctx, contracts.SetOpponentCombinationsCommand{
		RoomId:       command.RoomId,
		Combinations: match.OpponentsCombinations,
	}); err != nil {
		return nil, err
	}

	return &contracts.SuccessResponse{
		Success: true,
	}, nil
}

func (s *MatchesService) StartGame(ctx context.Context, roomId string) (*contracts.SuccessResponse, error) {

	match, err := s.storage.MatchesRepository.GetAllButGuesses(ctx, roomId)

	if err != nil {
		if errors.Is(err, domain.ErrEmptyResult) {
			return nil, ErrMatchNotFound
		}
		return nil, err
	}

	if match.Status != domain.MatchStateFullRoom {
		return nil, ErrMatchNotFullRoom
	}

	if len(match.OpponentsCombinations) != 2 {
		return nil, ErrExpectingCombinations
	}

	isTurnOf, err := match.GetRandomUser()

	if err != nil {
		return nil, ErrMatchNotFullRoom
	}

	if err := s.storage.MatchesRepository.ChangeStatusAndTurn(ctx, roomId, domain.MatchStatePlaying, isTurnOf); err != nil {
		return nil, err
	}

	return &contracts.SuccessResponse{
		Success: true,
	}, nil
}

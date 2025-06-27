package store

import (
	"context"
	"encoding/json"
	"fmt"
	"time"

	"github.com/alejandro-cardenas-g/bullAndCowsApp/contracts"
	"github.com/alejandro-cardenas-g/bullAndCowsApp/internal/domain"
	"github.com/alejandro-cardenas-g/bullAndCowsApp/internal/utils"
	"github.com/redis/go-redis/v9"
)

const (
	CREATE_OR_UPDATE_MATCH_EXP = time.Hour * 1
)

type MatchesRepository struct {
	rdb *redis.Client
}

func newMatchesRepository(rdb *redis.Client) *MatchesRepository {
	return &MatchesRepository{
		rdb: rdb,
	}
}

func (r *MatchesRepository) CreateMatch(ctx context.Context, player domain.Player) (*domain.Match, error) {
	roomId, err := domain.GenerateMatchId()
	if err != nil {
		return nil, err
	}

	match := &domain.Match{
		RoomId:                roomId,
		Players:               make(domain.MatchPlayers),
		OpponentsCombinations: make(domain.MatchOpponentCombinations),
		Guesses:               make(domain.MatchGuesses),
		Status:                domain.MatchStateWaiting,
		IsTurnOf:              player.Id,
	}

	match.Players[player.Id] = player

	playersJSON, _ := json.Marshal(match.Players)
	opponentsJSON, _ := json.Marshal(match.OpponentsCombinations)
	guessesJSON, _ := json.Marshal(match.Guesses)

	payload := map[string]interface{}{
		"Players":               string(playersJSON),
		"OpponentsCombinations": string(opponentsJSON),
		"Guesses":               string(guessesJSON),
		"Status":                string(match.Status),
		"IsTurnOf":              match.IsTurnOf,
	}

	key := getKeyById(roomId)

	if err := r.rdb.HSet(ctx, key, payload).Err(); err != nil {
		return nil, err
	}

	if err := r.rdb.Expire(ctx, key, CREATE_OR_UPDATE_MATCH_EXP).Err(); err != nil {
		return nil, err
	}

	return match, nil
}

func (r *MatchesRepository) SetPlayersAndFillRoom(ctx context.Context, command contracts.SetPlayersCommand) error {
	key := getKeyById(command.RoomId)

	plainPlayers, _ := json.Marshal(command.Players)

	payload := map[string]interface{}{
		"Players": string(plainPlayers),
		"Status":  string(domain.MatchStateFullRoom),
	}

	return r.rdb.HSet(ctx, key, payload).Err()
}

func (r *MatchesRepository) GetRoomPlayers(ctx context.Context, roomId string) (domain.MatchPlayers, error) {
	key := getKeyById(roomId)
	result, err := r.rdb.HGet(ctx, key, "Players").Result()
	if err != nil {
		if err == redis.Nil {
			return nil, domain.ErrEmptyResult
		}
		return nil, err
	}

	players := &domain.MatchPlayers{}
	if result != "" {
		err := json.Unmarshal([]byte(result), players)
		if err != nil {
			return nil, err
		}
	}

	return *players, nil
}

func (r *MatchesRepository) GetMatchStatusById(ctx context.Context, roomId string) (domain.MatchStatus, error) {
	key := getKeyById(roomId)
	result, err := r.rdb.HGet(ctx, key, "Status").Result()
	if err != nil {
		if err == redis.Nil {
			return "", domain.ErrEmptyResult
		}
		return "", err
	}

	matchStatus := domain.MatchStatus(result)

	return matchStatus, nil
}

func (r *MatchesRepository) GetPlayersAndCombinations(ctx context.Context, roomId string) (*domain.Match, error) {
	key := getKeyById(roomId)
	results, err := r.rdb.HMGet(ctx, key, "Players", "OpponentsCombinations").Result()

	if err != nil {
		return nil, err
	}

	if results[0] == nil || results[1] == nil {
		return nil, domain.ErrEmptyResult
	}

	var players domain.MatchPlayers

	err = json.Unmarshal([]byte(results[0].(string)), &players)
	if err != nil {
		return nil, err
	}

	var combinations domain.MatchOpponentCombinations

	err = json.Unmarshal([]byte(results[1].(string)), &combinations)
	if err != nil {
		return nil, err
	}

	match := &domain.Match{
		Players:               players,
		OpponentsCombinations: combinations,
	}

	return match, nil
}

func (r *MatchesRepository) SetPlayerCombination(ctx context.Context, command contracts.SetOpponentCombinationsCommand) error {
	key := getKeyById(command.RoomId)

	opponentsJSON, _ := json.Marshal(command.Combinations)

	payload := map[string]interface{}{
		"OpponentsCombinations": string(opponentsJSON),
	}

	return r.rdb.HSet(ctx, key, payload).Err()
}

func (r *MatchesRepository) GetAllButGuesses(ctx context.Context, roomId string) (*domain.Match, error) {
	key := getKeyById(roomId)
	results, err := r.rdb.HMGet(ctx, key, "Players", "OpponentsCombinations", "Status", "IsTurnOf").Result()

	if err != nil {
		return nil, err
	}

	if utils.IsSliceWithNilValues(results) {
		return nil, domain.ErrEmptyResult
	}

	var players domain.MatchPlayers

	if err := json.Unmarshal([]byte(results[0].(string)), &players); err != nil {
		return nil, err
	}

	var combinations domain.MatchOpponentCombinations
	if err := json.Unmarshal([]byte(results[1].(string)), &combinations); err != nil {
		return nil, err
	}

	status := domain.MatchStatus(results[2].(string))

	isTurnOf := results[3].(string)

	match := &domain.Match{
		Players:               players,
		OpponentsCombinations: combinations,
		Status:                status,
		IsTurnOf:              isTurnOf,
	}

	return match, nil
}

func (r *MatchesRepository) ChangeStatusAndTurn(ctx context.Context, roomId string, status domain.MatchStatus, isTurnOf string) error {
	key := getKeyById(roomId)

	payload := map[string]interface{}{
		"Status":   string(status),
		"IsTurnOf": isTurnOf,
	}

	return r.rdb.HSet(ctx, key, payload).Err()
}

func (r *MatchesRepository) GetAll(ctx context.Context, roomId string) (*domain.Match, error) {
	key := getKeyById(roomId)
	results, err := r.rdb.HMGet(ctx, key, "Players", "OpponentsCombinations", "Guesses", "Status", "IsTurnOf").Result()

	if err != nil {
		return nil, err
	}

	if utils.IsSliceWithNilValues(results) {
		return nil, domain.ErrEmptyResult
	}

	var players domain.MatchPlayers

	if err := json.Unmarshal([]byte(results[0].(string)), &players); err != nil {
		return nil, err
	}

	var combinations domain.MatchOpponentCombinations
	if err := json.Unmarshal([]byte(results[1].(string)), &combinations); err != nil {
		return nil, err
	}

	var guesses domain.MatchGuesses
	if err := json.Unmarshal([]byte(results[2].(string)), &guesses); err != nil {
		return nil, err
	}

	status := domain.MatchStatus(results[3].(string))

	isTurnOf := results[4].(string)

	match := &domain.Match{
		Players:               players,
		OpponentsCombinations: combinations,
		Status:                status,
		IsTurnOf:              isTurnOf,
		Guesses:               guesses,
	}

	return match, nil
}

func (r *MatchesRepository) SetNewGuess(ctx context.Context, command contracts.SetNewGuessCommand) error {
	key := getKeyById(command.RoomId)

	guessesJSON, _ := json.Marshal(command.Guesses)

	payload := map[string]interface{}{
		"Guesses":  string(guessesJSON),
		"IsTurnOf": command.IsTurnOf,
	}

	if command.IsWinner {
		payload["Status"] = string(domain.MatchStateFinished)
	}

	return r.rdb.HSet(ctx, key, payload).Err()
}

func (r *MatchesRepository) Exists(ctx context.Context, roomId string) error {
	key := getKeyById(roomId)
	return r.rdb.Exists(ctx, key).Err()
}

func (r *MatchesRepository) Restart(ctx context.Context, roomId string) error {
	key := getKeyById(roomId)

	match := &domain.Match{
		OpponentsCombinations: make(domain.MatchOpponentCombinations),
		Guesses:               make(domain.MatchGuesses),
		Status:                domain.MatchStateFullRoom,
	}

	opponentsJSON, _ := json.Marshal(match.OpponentsCombinations)
	guessesJSON, _ := json.Marshal(match.Guesses)

	payload := map[string]interface{}{
		"OpponentsCombinations": string(opponentsJSON),
		"Guesses":               string(guessesJSON),
		"Status":                string(match.Status),
	}

	resErr := r.rdb.HSet(ctx, key, payload).Err()

	if err := r.rdb.Expire(ctx, key, CREATE_OR_UPDATE_MATCH_EXP).Err(); err != nil {
		return err
	}

	return resErr
}

func getKeyById(roomId string) string {
	return fmt.Sprintf("room:%v", roomId)
}

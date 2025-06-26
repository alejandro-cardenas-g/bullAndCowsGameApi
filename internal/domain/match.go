package domain

import (
	"crypto/rand"
	"fmt"
	"math/big"
	"time"
)

var (
	ErrInvalidCombination       = fmt.Errorf("combination must have only 4 digits")
	ErrInvalidUniqueCombination = fmt.Errorf("number can not be repeated")
)

type MatchPlayers map[string]Player
type MatchOpponentCombinations map[string]string
type MatchGuesses map[string][]GuessesHistoryItem

type MatchStatus string

const (
	MatchStateWaiting  = MatchStatus("Waiting")
	MatchStateFullRoom = MatchStatus("FullRoom")
	MatchStatePlaying  = MatchStatus("Playing")
	MatchStateFinished = MatchStatus("Finished")
)

type Match struct {
	RoomId                string
	Players               MatchPlayers
	OpponentsCombinations MatchOpponentCombinations
	Guesses               MatchGuesses
	Status                MatchStatus
	IsTurnOf              string
}

func GenerateMatchId() (string, error) {
	const lenght = 7
	const charset = "abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ0123456789"

	result := make([]byte, lenght)
	for i := 0; i < lenght; i++ {
		num, err := rand.Int(rand.Reader, big.NewInt(int64(len(charset))))
		if err != nil {
			return "", err
		}
		result[i] = charset[num.Int64()]
	}
	return string(result), nil
}

func ValidateCombination(combination string) error {
	if len(combination) != 4 {
		return ErrInvalidCombination
	}

	previousValues := make(map[rune]rune)

	for _, digit := range combination {
		if _, exists := previousValues[digit]; exists {
			return ErrInvalidUniqueCombination
		}
		previousValues[digit] = digit
	}
	return nil
}

func (m *Match) GetRandomUser() (string, error) {
	if len(m.Players) != 2 {
		return "", fmt.Errorf("")
	}

	now := time.Now().UnixNano()
	index := now % 2

	values := make([]string, 0, len(m.Players))
	for _, player := range m.Players {
		values = append(values, player.Id)
	}

	selected := values[index]
	return selected, nil
}

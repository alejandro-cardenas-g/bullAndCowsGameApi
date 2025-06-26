package domain

type GuessesHistoryItem struct {
	Guess               []BullAndCowGuess
	IsWinnerCombination bool
}

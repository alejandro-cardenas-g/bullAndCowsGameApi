package domain

type BullAndCowType string

const (
	Bull BullAndCowType = "bull"
	Cow  BullAndCowType = "cow"
	None BullAndCowType = "none"
)

type BullAndCowGuess struct {
	Value string
	Type  BullAndCowType
}

func newBullAndCowGuess(value string, guessType BullAndCowType) BullAndCowGuess {
	return BullAndCowGuess{
		Value: value,
		Type:  guessType,
	}
}

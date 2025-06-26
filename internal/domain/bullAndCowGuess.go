package domain

type BullAndCowType string

const (
	Bull BullAndCowType = "bull"
	Cow  BullAndCowType = "cow"
	None BullAndCowType = "none"
)

type BullAndCowGuess struct {
	Value rune
	Type  BullAndCowType
}

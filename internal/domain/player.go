package domain

import "github.com/google/uuid"

type Player struct {
	Id       string
	Username string
}

func GeneratePlayerId() string {
	return uuid.NewString()
}

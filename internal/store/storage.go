package store

import (
	"github.com/alejandro-cardenas-g/bullAndCowsApp/contracts"
	"github.com/redis/go-redis/v9"
)

func NewRedisStorage(rdb *redis.Client) contracts.Storage {

	matchesRepository := newMatchesRepository(rdb)

	return contracts.Storage{
		MatchesRepository: matchesRepository,
	}
}

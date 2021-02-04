package memorystore

import (
	"github.com/go-redis/redis/v7"
)

type Store struct {
	client *redis.Client

	tokenRepository *TokenRepository
}

func New(redisUrl string) (*Store, error) {
	client, clientErr := newClient(redisUrl)
	if clientErr != nil {
		return nil, clientErr
	}

	return &Store{
		client: client,
	}, nil
}

func (s *Store) Token() *TokenRepository {
	if s.tokenRepository != nil {
		return s.tokenRepository
	}

	s.tokenRepository = &TokenRepository{
		store: s,
	}

	return s.tokenRepository
}

func newClient(memoryStoreUrl string) (*redis.Client, error) {
	client := redis.NewClient(&redis.Options{
		Addr: memoryStoreUrl,
	})

	if _, err := client.Ping().Result(); err != nil {
		return nil, err
	}

	return client, nil
}

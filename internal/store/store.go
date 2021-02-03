package store

import (
	"github.com/jmoiron/sqlx"
	"godmin/internal/store/repository"
)

type Store struct {
	db             *sqlx.DB
	userRepository *repository.User
}

func New(db *sqlx.DB) *Store {
	return &Store{
		db: db,
	}
}

func (s *Store) User() *repository.User {
	if s.userRepository != nil {
		return s.userRepository
	}

	s.userRepository = repository.NewUser(s.db)

	return s.userRepository
}

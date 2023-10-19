package storage

import (
	"context"
	"time"

	"github.com/jackc/pgx/v5/pgxpool"
)

type UserStorage struct {
	Pool *pgxpool.Pool
}

type UserData struct {
	ID   int    `json:"id"`
	Name string `json:"name"`
}

func NewUserStorage(pool *pgxpool.Pool) *UserStorage {
	return &UserStorage{Pool: pool}
}

// Creates a new user and returns the user once done!
func (s *UserStorage) CreateNewUser(data UserData) (int, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	_, err := s.Pool.Query(ctx, `INSERT INTO "users" (id, name) VALUES (1, 'Sourav');`)
	if err != nil {
		return 1, err
	}

	return 0, nil
}

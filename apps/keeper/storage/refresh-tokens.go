package storage

import (
	"context"
	"fmt"
	"time"

	"github.com/google/uuid"
	"github.com/jackc/pgx/v5/pgxpool"
)

type RefreshTokenStorage struct {
	Pool *pgxpool.Pool
}

// Used to construct a refresh token object
type RefreshTokenObject struct {
	UserID  int        `json:"user_id"`
	Token   string     `json:"refresh_token"`
	Expires *time.Time `json:"expires"`
}

func NewRefreshTokenStorage(pool *pgxpool.Pool) *RefreshTokenStorage {
	return &RefreshTokenStorage{Pool: pool}
}

// Generates and inserts the refresh token for the user
func (s *RefreshTokenStorage) GenerateAndInsertRefreshToken(data *SensitiveUserData) (*RefreshTokenObject, error) {
	refreshToken := uuid.New()
	expirationTime := time.Now().Add(168 * time.Hour) // 7 Days

	query := fmt.Sprintf(`INSERT INTO "refresh_tokens" ("user_id", "refresh_token", "expires") VALUES (%d, '%s', '%s');`, data.UserID, refreshToken.String(), expirationTime.UTC().Format("2006/01/02 15:04:05"))

	_, err := s.Pool.Query(context.Background(), query)
	if err != nil {
		return nil, err
	}

	return &RefreshTokenObject{
		UserID:  data.UserID,
		Token:   refreshToken.String(),
		Expires: &expirationTime,
	}, nil
}

// Finds if a given refresh token exists
func (s *RefreshTokenStorage) FindExistingRefreshToken(rt string) (*RefreshTokenObject, error) {
	query := fmt.Sprintf(`SELECT "refresh_tokens"."user_id", "refresh_tokens"."refresh_token", "refresh_tokens"."expires" FROM "refresh_tokens" WHERE "refresh_tokens"."refresh_token" = '%s';`, rt)

	rows, err := s.Pool.Query(context.Background(), query)
	if err != nil {
		return nil, err
	}

	defer rows.Close()

	var rto RefreshTokenObject
	for rows.Next() {
		err = rows.Scan(&rto.UserID, &rto.Token, &rto.Expires)
		if err != nil {
			return nil, err
		}
	}

	return &rto, nil
}

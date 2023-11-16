package storage

import (
	"context"
	"errors"
	"fmt"
	"time"

	"github.com/jackc/pgx/v5/pgxpool"
	"golang.org/x/crypto/bcrypt"
)

type UserStorage struct {
	Pool *pgxpool.Pool
}

// `UserData` struct is used when creating and registering the user with the system
type UserData struct {
	Name      string `json:"name" validate:"required"`
	BirthDate string `json:"birth_date" validate:"required"`
	Username  string `json:"username" validate:"required"`
	Email     string `json:"email" validate:"required"`
	Password  string `json:"password" validate:"required"`
}

type LoginUser struct {
	Username string `json:"username" validate:"required"`
	Password string `json:"password" validate:"required"`
}

// `SensitiveUserData` is used when the user is fetched from the database with password
type SensitiveUserData struct {
	UserID    int
	Name      string
	BirthDate *time.Time
	Username  string
	Email     string
	Password  string
}

// `NonSensitiveUserData` is used when the user is fetched from the database with password
type NonSensitiveUserData struct {
	UserID    int
	Name      string
	BirthDate *time.Time
	Username  string
	Email     string
}

func NewUserStorage(pool *pgxpool.Pool) *UserStorage {
	return &UserStorage{Pool: pool}
}

// Find a user with a given `uid` Primary key
func (s *UserStorage) FindUserWithUID(uid int) (NonSensitiveUserData, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 6*time.Second)
	defer cancel()

	// Construct the query
	query := fmt.Sprintf(`SELECT "users"."uid", "users"."name", "users"."birth_date", "users"."username", "users"."email" FROM "users" WHERE "users"."uid" = %d;`, uid)

	// Query and get results from db
	rows, err := s.Pool.Query(ctx, query)
	if err != nil {
		return NonSensitiveUserData{}, err
	}

	defer rows.Close()

	var user NonSensitiveUserData
	for rows.Next() {
		err = rows.Scan(&user.UserID, &user.Name, &user.BirthDate, &user.Username, &user.Email)
		if err != nil {
			return NonSensitiveUserData{}, err
		}
	}

	return user, nil
}

// Find a user with a given `username`
func (s *UserStorage) FindUserWithUsername(username string, withPassword bool) (*SensitiveUserData, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 6*time.Second)
	defer cancel()

	// Construct the query
	var query string

	if withPassword {
		query = fmt.Sprintf(`SELECT "users"."uid", "users"."name", "users"."birth_date", "users"."username", "users"."email", "users"."password" FROM "users" WHERE "users"."username" = '%s';`, username)
	} else {
		query = fmt.Sprintf(`SELECT "users"."uid", "users"."name", "users"."birth_date", "users"."username", "users"."email" FROM "users" WHERE "users"."username" = '%s';`, username)
	}

	// Query and get results from db
	rows, err := s.Pool.Query(ctx, query)
	if err != nil {
		return nil, err
	}

	defer rows.Close()

	var user SensitiveUserData
	for rows.Next() {
		userDataSlice, err := rows.Values()

		userid, ok := userDataSlice[0].(int32)
		if ok {
			user.UserID = int(userid)
		} else {
			newError := errors.New("error converting 'id' to integer")
			return nil, newError
		}

		name, ok := userDataSlice[1].(string)
		if ok {
			user.Name = name
		} else {
			newError := errors.New("error converting 'name' to string")
			return nil, newError
		}

		birthdate, ok := userDataSlice[2].(time.Time)
		if ok {
			user.BirthDate = &birthdate
		} else {
			newError := errors.New("error converting 'birthdate' to *time.Time")
			return nil, newError
		}

		username, ok := userDataSlice[3].(string)
		if ok {
			user.Username = username
		} else {
			newError := errors.New("error converting 'username' to string")
			return nil, newError
		}

		email, ok := userDataSlice[4].(string)
		if ok {
			user.Email = email
		} else {
			newError := errors.New("error converting 'email' to string")
			return nil, newError
		}

		if withPassword {
			password, ok := userDataSlice[5].(string)
			if ok {
				user.Password = password
			} else {
				newError := errors.New("error converting 'password' to string")
				return nil, newError
			}
		} else {
			user.Password = ""
		}

		if err != nil {
			return nil, err
		}

	}

	return &user, nil
}

// Find a user with provided `email`
func (s *UserStorage) FindUserWithEmail(email string) (NonSensitiveUserData, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 6*time.Second)
	defer cancel()

	// Construct the query
	query := fmt.Sprintf(`SELECT "users"."uid", "users"."name", "users"."birth_date", "users"."username", "users"."email" FROM "users" WHERE "users"."email" = '%s';`, email)

	// Query and get results from db
	rows, err := s.Pool.Query(ctx, query)
	if err != nil {
		return NonSensitiveUserData{}, err
	}

	defer rows.Close()

	var user NonSensitiveUserData
	for rows.Next() {
		err = rows.Scan(&user.UserID, &user.Name, &user.BirthDate, &user.Username, &user.Email)
		if err != nil {
			return NonSensitiveUserData{}, err
		}
	}

	return user, nil
}

// Creates a new user and returns the user once done!
func (s *UserStorage) CreateNewUser(data UserData) (int, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 6*time.Second)
	defer cancel()

	password := []byte(data.Password)
	hashedPassword, err := bcrypt.GenerateFromPassword(password, bcrypt.DefaultCost)
	if err != nil {
		return 2, err
	}

	query := fmt.Sprintf(
		`INSERT INTO "users" (name, birth_date, username, email, password) VALUES ('%s', '%s', '%s', '%s', '%s');`,
		data.Name,
		data.BirthDate,
		data.Username,
		data.Email,
		hashedPassword,
	)

	_, err = s.Pool.Query(ctx, query)
	if err != nil {
		return 1, err
	}

	return 0, nil
}

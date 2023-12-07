package storage

import (
	"context"
	"errors"
	"time"

	"github.com/jackc/pgx/v4"
	"github.com/jackc/pgx/v5/pgxpool"
	"golang.org/x/crypto/bcrypt"
)

type UserStorage struct {
	Pool *pgxpool.Pool
}

func NewUserStore(pool *pgxpool.Pool) *UserStorage {
	return &UserStorage{
		Pool: pool,
	}
}

// UserData is used while registering a new user
type UserData struct {
	Name      string `json:"name" validate:"required"`
	BirthDate string `json:"birth_date" validate:"required"`
	Username  string `json:"username" validate:"required"`
	Email     string `json:"email" validate:"required"`
	Password  string `json:"password" validate:"required"`
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

// LoginUser is struct that we will receive when the user tries to login
type LoginUser struct {
	Username string `json:"username" validate:"required"`
	Password string `json:"password" validate:"required"`
}

// Finds a user with a provided UID and returns the user with Password field removed
func (s *UserStorage) FindUserWithUID(uid int) (*NonSensitiveUserData, error) {
	query := `SELECT "users"."uid", "users"."name", "users"."birth_date", "users"."username", "users"."email" FROM "users" WHERE "users"."uid" = $1;`

	// Query and get results from db
	rows, err := s.Pool.Query(context.Background(), query, uid)
	if err != nil {
		return nil, err
	}

	defer rows.Close()

	var user NonSensitiveUserData
	found := false

	for rows.Next() {
		found = true
		err = rows.Scan(&user.UserID, &user.Name, &user.BirthDate, &user.Username, &user.Email)
		if err != nil {
			return nil, err
		}
	}

	// If no rows were found return ErrNoRows
	if !found {
		return nil, pgx.ErrNoRows
	}

	return &user, nil
}

// Finds a user with a provided username and returns the user with either password field removed or attached
func (s *UserStorage) FindUserWithUsername(username string, withPassword bool) (*SensitiveUserData, error) {
	// Construct the query
	var query string

	if withPassword {
		query = `SELECT "users"."uid", "users"."name", "users"."birth_date", "users"."username", "users"."email", "users"."password" FROM "users" WHERE "users"."username" = $1;`
	} else {
		query = `SELECT "users"."uid", "users"."name", "users"."birth_date", "users"."username", "users"."email" FROM "users" WHERE "users"."username" = $1;`
	}

	// Query and get results from db
	rows, err := s.Pool.Query(context.Background(), query, username)
	if err != nil {
		return nil, err
	}

	defer rows.Close()

	var user SensitiveUserData
	found := false

	for rows.Next() {
		found = true
		userDataSlice, err := rows.Values()

		userid, ok := userDataSlice[0].(int32)
		if ok {
			user.UserID = int(userid)
		} else {
			return nil, errors.New("error converting 'id' to integer")
		}

		name, ok := userDataSlice[1].(string)
		if ok {
			user.Name = name
		} else {
			return nil, errors.New("error converting 'name' to string")
		}

		birthdate, ok := userDataSlice[2].(time.Time)
		if ok {
			user.BirthDate = &birthdate
		} else {
			return nil, errors.New("error converting 'birthdate' to *time.Time")
		}

		username, ok := userDataSlice[3].(string)
		if ok {
			user.Username = username
		} else {
			return nil, errors.New("error converting 'username' to string")
		}

		email, ok := userDataSlice[4].(string)
		if ok {
			user.Email = email
		} else {
			return nil, errors.New("error converting 'email' to string")
		}

		if withPassword {
			password, ok := userDataSlice[5].(string)
			if ok {
				user.Password = password
			} else {
				return nil, errors.New("error converting 'password' to string")
			}
		} else {
			user.Password = ""
		}

		if err != nil {
			return nil, err
		}

	}

	// If no rows were found return ErrNoRows
	if !found {
		return nil, pgx.ErrNoRows
	}

	return &user, nil
}

// Creates a new user, if successful returns 0
func (s *UserStorage) CreateNewUser(data UserData) (int, error) {
	if err := s.CheckDuplicateUser(data.Username, data.Email); err != nil {
		return 1, err
	}

	password := []byte(data.Password)
	hashedPassword, err := bcrypt.GenerateFromPassword(password, bcrypt.DefaultCost)
	if err != nil {
		return 2, err
	}

	query := `INSERT INTO "users" (name, birth_date, username, email, password) VALUES ($1, $2, $3, $4, $5);`
	_, err = s.Pool.Query(context.Background(), query, data.Name, data.BirthDate, data.Username, data.Email, hashedPassword)
	if err != nil {
		return 3, err
	}

	return 0, nil
}

// Checks for duplicate user with same email or username
func (s *UserStorage) CheckDuplicateUser(username, email string) error {
	query := `SELECT COUNT(uid) FROM "users" WHERE username = $1 OR email = $2;`

	var count int
	err := s.Pool.QueryRow(context.Background(), query, username, email).Scan(&count)
	if err != nil {
		return err
	}

	if count > 0 {
		return errors.New("duplicate username or email")
	}

	return nil
}

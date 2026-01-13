package models

import (
	"database/sql"
	"errors"
	"strings"

	"github.com/google/uuid"
	"github.com/jackc/pgx/v5/pgconn"
	"golang.org/x/crypto/bcrypt"
)

type User struct {
	ID             uuid.UUID
	Username       string
	Email          string
	HashedPassword []byte
}

type UserModel struct {
	DB *sql.DB
}

func (m *UserModel) Insert(username, email, password string) error {
	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
	if err != nil {
		return err
	}

	stmt := `INSERT INTO users (username, email, hashed_password) VALUES ($1, $2, $3)`
	_, err = m.DB.Exec(stmt, username, email, hashedPassword)
	if err != nil {
		var pgErr *pgconn.PgError

		if errors.As(err, &pgErr) {
			if pgErr.Code == "23505" {
				if strings.Contains(pgErr.ConstraintName, "username") {
					return ErrUsernameAlreadyExists
				} else if strings.Contains(pgErr.ConstraintName, "email") {
					return ErrEmailAlreadyExists
				}
			}
		}
		return err
	}

	return nil
}

func (m *UserModel) Authenticate(email, password string) (uuid.UUID, error) {
	var id uuid.UUID
	var hashedPassword []byte
	stmt := `SELECT id, hashed_password FROM users WHERE email = $1`

	err := m.DB.QueryRow(stmt, email).Scan(&id, &hashedPassword)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return uuid.UUID{}, ErrInvalidCredentials
		}

		return uuid.UUID{}, err
	}

	err = bcrypt.CompareHashAndPassword(hashedPassword, []byte(password))
	if err != nil {
		if errors.Is(err, bcrypt.ErrMismatchedHashAndPassword) {
			return uuid.UUID{}, ErrInvalidCredentials
		}

		return uuid.UUID{}, err
	}

	return id, nil
}

func (m *UserModel) Exists(id uuid.UUID) (bool, error) {
	var exists bool
	stmt := "SELECT EXISTS(SELECT true FROM users WHERE id = $1)"

	err := m.DB.QueryRow(stmt, id).Scan(&exists)

	return exists, err
}

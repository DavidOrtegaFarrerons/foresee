package models

import (
	"database/sql"
	"errors"
	"strings"
	"time"

	"github.com/google/uuid"
	"github.com/jackc/pgx/v5/pgconn"
	"golang.org/x/crypto/bcrypt"
)

type User struct {
	ID             uuid.UUID
	Username       string
	Email          string
	HashedPassword []byte
	LastClaimedAt  sql.NullTime
	Balance        int
}

type UserModel struct {
	DB *sql.DB
}

const initialBalance int = 1000

func (m *UserModel) Insert(username, email, password string) error {
	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
	if err != nil {
		return err
	}

	stmt := `INSERT INTO users (username, email, hashed_password, balance) VALUES ($1, $2, $3, $4)`
	_, err = m.DB.Exec(stmt, username, email, hashedPassword, initialBalance)
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

func (m *UserModel) GetTemplateInfo(id uuid.UUID) (int, sql.NullTime, error) {
	var balance int
	var lastClaimedAt sql.NullTime
	stmt := "SELECT balance, last_daily_claim FROM users WHERE id = $1"
	err := m.DB.QueryRow(stmt, id).Scan(&balance, &lastClaimedAt)
	if err != nil {
		return 0, sql.NullTime{}, err
	}

	return balance, lastClaimedAt, nil
}

func (m *UserModel) SelectForUpdate(tx *sql.Tx, id uuid.UUID) (User, error) {
	var user User
	stmt := `SELECT id, balance, last_daily_claim
			FROM users
			WHERE id = $1
			FOR UPDATE
			`

	err := tx.QueryRow(stmt, id).Scan(&user.ID, &user.Balance, &user.LastClaimedAt)
	if err != nil {
		return User{}, err
	}

	return user, nil
}

func (m *UserModel) ApplyDailyClaim(tx *sql.Tx, id uuid.UUID, balance int, lastClaimedAt time.Time) error {
	stmt := `UPDATE users SET balance = $1, last_daily_claim = $2 WHERE id = $3`
	_, err := tx.Exec(stmt, balance, lastClaimedAt, id)
	if err != nil {
		return err
	}

	return nil
}

func (m *UserModel) DecreaseBalanceBy(tx *sql.Tx, id uuid.UUID, amount int) error {
	stmt := `UPDATE users SET balance = balance - $1 WHERE id = $2`
	_, err := tx.Exec(stmt, amount, id)
	return err
}

func (m *UserModel) IncreaseBalanceBy(tx *sql.Tx, userID uuid.UUID, amount int) error {
	stmt := `UPDATE users SET balance = balance + $1 WHERE id = $2`
	_, err := tx.Exec(stmt, amount, userID)
	return err
}

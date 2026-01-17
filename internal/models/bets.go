package models

import (
	"database/sql"
	"errors"

	"github.com/google/uuid"
	"github.com/jackc/pgx/v5/pgconn"
)

type Bet struct {
	ID        uuid.UUID
	UserID    uuid.UUID
	MarketID  uuid.UUID
	OutcomeID uuid.UUID
	Amount    int
	CreatedAt sql.NullTime
}

type BetModel struct {
	DB *sql.DB
}

const MinimumBetAmount int = 100

func (m *BetModel) Place(tx *sql.Tx, userID uuid.UUID, marketID uuid.UUID, outcomeID uuid.UUID, amount int) error {
	stmt := `INSERT INTO bets (user_id, market_id, outcome_id, amount) VALUES ($1, $2, $3, $4)`

	_, err := tx.Exec(stmt, userID, marketID, outcomeID, amount)
	if err != nil {
		var pgErr *pgconn.PgError
		ok := errors.As(err, &pgErr)
		if ok && pgErr.Code == "23505" && pgErr.ConstraintName == "bets_market_id_user_id_key" {
			return ErrUserAlreadyBetOnMarket
		}
	}

	return err
}

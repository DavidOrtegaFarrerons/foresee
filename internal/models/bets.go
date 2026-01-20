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

func (m *BetModel) ForMarketForUpdate(tx *sql.Tx, marketID uuid.UUID) ([]Bet, error) {
	stmt := `SELECT id, user_id, amount, outcome_id FROM bets WHERE market_id = $1 FOR UPDATE`

	rows, err := tx.Query(stmt, marketID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var bets []Bet

	for rows.Next() {
		var b Bet
		err = rows.Scan(&b.ID, &b.UserID, &b.Amount, &b.OutcomeID)
		if err != nil {
			return nil, err
		}
		bets = append(bets, b)
	}

	if err = rows.Err(); err != nil {
		return nil, err
	}

	return bets, nil
}

func (m *BetModel) SetPayout(tx *sql.Tx, betID uuid.UUID, payout int) error {
	stmt := `UPDATE bets SET payout_amount = $1, settled_at = NOW() WHERE id = $2`
	_, err := tx.Exec(stmt, payout, betID)
	return err
}

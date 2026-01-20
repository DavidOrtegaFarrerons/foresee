package models

import (
	"database/sql"
	"errors"

	"github.com/google/uuid"
	"github.com/lib/pq"
)

type Outcome struct {
	ID         uuid.UUID
	MarketID   uuid.UUID
	Label      string
	createdAt  sql.NullTime
	PoolAmount int
}

type OutcomeModel struct {
	DB *sql.DB
}

const yesLabel = "yes"
const noLabel = "no"

func (m *OutcomeModel) CreateWithYesNo(tx *sql.Tx, marketId uuid.UUID) error {
	stmt := `INSERT INTO outcomes (
            	market_id, label                
			) VALUES ($1, $2)`

	_, err := tx.Exec(stmt, marketId, yesLabel)
	if err != nil {
		return err
	}

	_, err = tx.Exec(stmt, marketId, noLabel)
	if err != nil {
		return err
	}

	return nil
}

func (m *OutcomeModel) ForMarkets(ids []uuid.UUID) (map[uuid.UUID][]Outcome, error) {
	stmt := `SELECT id, market_id, label, pool_amount FROM outcomes WHERE market_id = ANY($1)`

	rows, err := m.DB.Query(stmt, pq.Array(ids))
	if err != nil {
		return nil, err
	}

	defer rows.Close()

	outcomes := make(map[uuid.UUID][]Outcome)

	for rows.Next() {
		var o Outcome
		err = rows.Scan(&o.ID, &o.MarketID, &o.Label, &o.PoolAmount)
		if err != nil {
			return nil, err
		}
		outcomes[o.MarketID] = append(outcomes[o.MarketID], o)
	}

	if err = rows.Err(); err != nil {
		return nil, err
	}

	return outcomes, nil
}

func (m *OutcomeModel) ForMarket(id uuid.UUID) ([]Outcome, error) {
	stmt := `SELECT id, market_id, label, pool_amount FROM outcomes WHERE market_id = $1`

	rows, err := m.DB.Query(stmt, id)
	if err != nil {
		return nil, err
	}

	defer rows.Close()

	var outcomes []Outcome

	for rows.Next() {
		var o Outcome
		err = rows.Scan(&o.ID, &o.MarketID, &o.Label, &o.PoolAmount)
		if err != nil {
			return nil, err
		}
		outcomes = append(outcomes, o)
	}

	if err = rows.Err(); err != nil {
		return nil, err
	}

	return outcomes, nil
}

func (m *OutcomeModel) SelectForUpdate(tx *sql.Tx, id uuid.UUID) (Outcome, error) {
	stmt := `SELECT id, market_id, pool_amount FROM outcomes WHERE id = $1 FOR UPDATE`

	var o Outcome
	err := tx.QueryRow(stmt, id).Scan(&o.ID, &o.MarketID, &o.PoolAmount)
	if err != nil {
		return Outcome{}, err
	}

	return o, nil
}

func (m *OutcomeModel) AddPoolAmount(tx *sql.Tx, id uuid.UUID, amount int) error {
	stmt := `UPDATE outcomes SET pool_amount = pool_amount + $1 WHERE id = $2`
	_, err := tx.Exec(stmt, amount, id)
	return err
}

func (m *OutcomeModel) ExistsForMarket(tx *sql.Tx, outcomeID uuid.UUID, marketID uuid.UUID) (bool, error) {
	stmt := `SELECT 1 FROM outcomes WHERE id = $1 AND market_id = $2`

	var exists int
	err := tx.QueryRow(stmt, outcomeID, marketID).Scan(&exists)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return false, nil
		}
		return false, err
	}

	return true, nil
}

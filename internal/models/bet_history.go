package models

import (
	"time"

	"github.com/google/uuid"
)

type BetHistoryRow struct {
	BetID           uuid.UUID
	MarketID        uuid.UUID
	MarketTitle     string
	MarketStatus    string
	OutcomeID       uuid.UUID
	OutcomeLabel    string
	Amount          int
	Payout          *int
	Result          string
	BetCreatedAt    time.Time
	MarketExpiresAt time.Time
}

func (m *BetModel) GetUserBetHistory(userID uuid.UUID) ([]BetHistoryRow, error) {
	stmt := `
	SELECT
		b.id,
		m.id,
		m.title,
		m.status,
		o.id,
		o.label,
		b.amount,
		b.payout_amount,
		b.created_at,
		m.expires_at
	FROM bets b
	JOIN markets m ON m.id = b.market_id
	JOIN outcomes o ON o.id = b.outcome_id
	WHERE b.user_id = $1
	ORDER BY b.created_at DESC`

	rows, err := m.DB.Query(stmt, userID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var betHistoryRows []BetHistoryRow

	for rows.Next() {
		var row BetHistoryRow

		err = rows.Scan(
			&row.BetID,
			&row.MarketID,
			&row.MarketTitle,
			&row.MarketStatus,
			&row.OutcomeID,
			&row.OutcomeLabel,
			&row.Amount,
			&row.Payout,
			&row.BetCreatedAt,
			&row.MarketExpiresAt,
		)
		if err != nil {
			return nil, err
		}

		if row.MarketStatus != "resolved" {
			row.Result = "pending"
		} else if row.Payout != nil && *row.Payout > 0 {
			row.Result = "win"
		} else {
			row.Result = "lose"
		}

		betHistoryRows = append(betHistoryRows, row)
	}

	if err = rows.Err(); err != nil {
		return nil, err
	}

	return betHistoryRows, nil
}

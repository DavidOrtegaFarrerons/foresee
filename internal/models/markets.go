package models

import (
	"database/sql"
	"time"

	"github.com/google/uuid"
)

type Category string

const (
	CategoryFriends Category = "friends"
	CategoryCrypto  Category = "crypto"
)

type ResolverType string

const (
	ResolverCreator ResolverType = "creator"
	ResolverAdmin   ResolverType = "admin"
)

func AllCategories() []Category {
	return []Category{
		CategoryFriends,
		CategoryCrypto,
	}
}

func AllResolverTypes() []ResolverType {
	return []ResolverType{
		ResolverCreator,
		ResolverAdmin,
	}
}

type Market struct {
	ID           uuid.UUID
	Title        string
	Description  string
	Category     Category
	ResolverType ResolverType
	ResolverRef  *uuid.UUID
	ExpiresAt    time.Time
	Status       string
	CreatedBy    uuid.UUID
	Outcomes     []Outcome
}

type MarketModel struct {
	DB *sql.DB
}

func (m *MarketModel) Insert(
	tx *sql.Tx,
	title string,
	description string,
	category Category,
	resolverType ResolverType,
	resolverRef *uuid.UUID,
	expiresAt time.Time,
	userID uuid.UUID,
) (uuid.UUID, error) {
	stmt := `INSERT INTO markets
		(title, description, category, resolver_type, resolver_ref, expires_at, status, created_by)
		VALUES ($1, $2, $3, $4, $5, $6, $7, $8)
		RETURNING id`

	var id uuid.UUID
	err := tx.QueryRow(stmt,
		title,
		description,
		category,
		resolverType,
		resolverRef,
		expiresAt,
		"open",
		userID,
	).Scan(&id)

	if err != nil {
		return uuid.UUID{}, err
	}

	return id, nil
}

func (m *MarketModel) Latest() ([]*Market, error) {
	stmt := `SELECT id, title, description, category, resolver_type, resolver_ref, expires_at, status, created_by
			 FROM markets WHERE expires_at > NOW() ORDER BY expires_at DESC LIMIT 10`

	rows, err := m.DB.Query(stmt)
	if err != nil {
		return nil, err
	}

	//We have to do this, otherwise the db connection stays open. We also have to do it after checking for an error
	//Otherwise, the application will panic
	defer rows.Close()

	markets := make([]*Market, 0, 10)

	for rows.Next() {
		market := &Market{}
		err = rows.Scan(&market.ID,
			&market.Title,
			&market.Description,
			&market.Category,
			&market.ResolverType,
			&market.ResolverRef,
			&market.ExpiresAt,
			&market.Status,
			&market.CreatedBy,
		)
		if err != nil {
			return nil, err
		}

		markets = append(markets, market)
	}

	if err = rows.Err(); err != nil {
		return nil, err
	}

	return markets, nil
}

func (m *MarketModel) Get(id uuid.UUID) (Market, error) {
	stmt := `SELECT id, title, description, category, resolver_type, resolver_ref, expires_at, status, created_by
			 FROM markets WHERE id = $1`

	var market Market
	err := m.DB.QueryRow(stmt, id).Scan(&market.ID,
		&market.Title,
		&market.Description,
		&market.Category,
		&market.ResolverType,
		&market.ResolverRef,
		&market.ExpiresAt,
		&market.Status,
		&market.CreatedBy,
	)
	if err != nil {
		return Market{}, err
	}

	return market, nil
}

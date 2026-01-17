package services

import (
	"foresee/internal/models"
	"time"

	"github.com/google/uuid"
)

type MarketService struct {
	Markets        *models.MarketModel
	OutcomeService OutcomeService
}

func (s *MarketService) Create(
	title string,
	description string,
	categoryStr string,
	resolverTypeStr string,
	expiresAtStr string,
	userID uuid.UUID,
) error {
	category := models.Category(categoryStr)
	resolverType := models.ResolverType(resolverTypeStr)

	loc, err := time.LoadLocation("Europe/Madrid")
	if err != nil {
		return err
	}

	expiresAt, err := time.ParseInLocation("2006-01-02T15:04", expiresAtStr, loc)

	resolverRef := uuid.UUID{}
	if resolverType == models.ResolverCreator {
		resolverRef = userID
	}

	tx, err := s.Markets.DB.Begin()
	if err != nil {
		return err
	}

	defer tx.Rollback()

	id, err := s.Markets.Insert(
		tx,
		title,
		description,
		category,
		resolverType,
		&resolverRef,
		expiresAt,
		userID,
	)

	if err != nil {
		return err
	}

	err = s.OutcomeService.CreateDefaultForMarket(tx, id)
	if err != nil {
		return err
	}

	return tx.Commit()
}

func (s *MarketService) Get(id uuid.UUID) (models.Market, error) {
	return s.Markets.Get(id)
}

func (s *MarketService) Latest() ([]*models.Market, error) {
	markets, err := s.Markets.Latest()
	if err != nil {
		return nil, err
	}

	if len(markets) == 0 {
		return markets, nil
	}

	marketIDs := make([]uuid.UUID, len(markets))
	for i, m := range markets {
		marketIDs[i] = m.ID
	}

	outcomesByMarket, err := s.OutcomeService.ForMarkets(marketIDs)
	if err != nil {
		return nil, err
	}

	for _, m := range markets {
		m.Outcomes = outcomesByMarket[m.ID]
	}

	return markets, nil
}

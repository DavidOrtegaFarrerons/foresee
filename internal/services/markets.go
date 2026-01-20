package services

import (
	"foresee/internal/models"
	"time"

	"github.com/google/uuid"
)

type MarketService struct {
	Markets        *models.MarketModel
	BetService     BetService
	OutcomeService OutcomeService
	UserService    UserService
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
	if err != nil {
		return err
	}

	var resolverRef *uuid.UUID
	if resolverType == models.ResolverCreator {
		resolverRef = &userID
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
		resolverRef,
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
	m, err := s.Markets.Get(id)
	if err != nil {
		return models.Market{}, err
	}

	o, err := s.OutcomeService.ForMarket(m.ID)
	if err != nil {
		return models.Market{}, err
	}

	m.Outcomes = o
	return m, nil
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

func (s *MarketService) PendingResolution(userID uuid.UUID) ([]models.Market, error) {
	return s.Markets.PendingResolution(userID)
}

func (s *MarketService) ResolveMarket(marketID uuid.UUID, userID uuid.UUID, outcomeID uuid.UUID) error {
	tx, err := s.Markets.DB.Begin()
	if err != nil {
		return err
	}

	defer tx.Rollback()

	m, err := s.Markets.SelectForUpdate(tx, marketID)
	if err != nil {
		return err
	}

	if m.ResolverType != models.ResolverCreator {
		return models.ErrUserNotAuthorized
	}

	if m.ResolverRef == nil || *m.ResolverRef != userID {
		return models.ErrUserNotAuthorized
	}

	if m.ResolvedOutcomeID != nil {
		return models.ErrMarketAlreadyResolved
	}

	if time.Now().Before(m.ExpiresAt) {
		return models.ErrMarketNotExpired
	}

	outcomeInMarket, err := s.OutcomeService.ExistsForMarketTx(tx, outcomeID, marketID)
	if err != nil {
		return err
	}

	if !outcomeInMarket {
		return models.ErrOutcomeDoesNotBelongToMarket
	}

	bets, err := s.BetService.ForMarketForUpdate(tx, marketID)
	if err != nil {
		return err
	}

	totalPool := 0
	winningPool := 0

	for _, b := range bets {
		totalPool += b.Amount
		if b.OutcomeID == outcomeID {
			winningPool += b.Amount
		}
	}

	distributed := 0
	var firstWinner *models.Bet

	for _, b := range bets {
		if b.OutcomeID != outcomeID {
			err = s.BetService.SetPayout(tx, b.ID, 0)
			if err != nil {
				return err
			}
			continue
		}

		if winningPool == 0 {
			err = s.BetService.SetPayout(tx, b.ID, 0)
			if err != nil {
				return err
			}
			continue
		}

		payout := b.Amount * totalPool / winningPool
		distributed += payout

		if firstWinner == nil {
			firstWinner = &b
		}

		err = s.BetService.SetPayout(tx, b.ID, payout)
		if err != nil {
			return err
		}

		err = s.UserService.IncreaseBalanceBy(tx, b.UserID, payout)
		if err != nil {
			return err
		}
	}

	leftover := totalPool - distributed
	if leftover > 0 && firstWinner != nil {
		err = s.UserService.IncreaseBalanceBy(tx, firstWinner.UserID, leftover)
		if err != nil {
			return err
		}
	}

	err = s.Markets.ResolveMarket(tx, marketID, userID, outcomeID)
	if err != nil {
		return err
	}

	return tx.Commit()
}

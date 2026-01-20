package services

import (
	"database/sql"
	"errors"
	"foresee/internal/models"
	"time"

	"github.com/google/uuid"
)

type BetService struct {
	Bets          *models.BetModel
	UserService   *UserService
	MarketService *MarketService
	Outcome       *models.OutcomeModel
}

var ErrInsufficientBalance = errors.New("you cannot place a bet that is higher than your current balance")
var ErrMarketExpired = errors.New("you cannot place a bet in an expired market")
var ErrMarketNotOpen = errors.New("you cannot place a bet in a market that is not open")
var ErrOutcomeNotFound = errors.New("the selected outcome does not exist in the selected market")

func (s BetService) Place(userID uuid.UUID, marketID uuid.UUID, outcomeID uuid.UUID, amount int) error {
	tx, err := s.Bets.DB.Begin()
	if err != nil {
		return err
	}

	defer tx.Rollback()

	user, err := s.UserService.Users.SelectForUpdate(tx, userID)
	if err != nil {
		return err
	}

	if user.Balance < amount {
		return ErrInsufficientBalance
	}

	market, err := s.MarketService.Get(marketID)
	if err != nil {
		return err
	}

	if time.Now().After(market.ExpiresAt) {
		return ErrMarketExpired
	}

	if market.Status != "open" {
		return ErrMarketNotOpen
	}

	outcome, err := s.Outcome.SelectForUpdate(tx, outcomeID)
	if err != nil {
		return err
	}

	if marketID != outcome.MarketID {
		return ErrOutcomeNotFound
	}

	err = s.Bets.Place(tx, userID, marketID, outcomeID, amount)
	if err != nil {
		return err
	}

	user.Balance = user.Balance - amount

	err = s.UserService.DecreaseBalanceBy(tx, userID, amount)
	if err != nil {
		return err
	}

	err = s.Outcome.AddPoolAmount(tx, outcomeID, amount)
	if err != nil {
		return err
	}

	return tx.Commit()
}

func (s BetService) GetUserBetHistory(userID uuid.UUID) ([]models.BetHistoryRow, error) {
	return s.Bets.GetUserBetHistory(userID)
}

func (s BetService) ForMarketForUpdate(tx *sql.Tx, marketID uuid.UUID) ([]models.Bet, error) {
	return s.Bets.ForMarketForUpdate(tx, marketID)
}

func (s BetService) SetPayout(tx *sql.Tx, betID uuid.UUID, payout int) error {
	return s.Bets.SetPayout(tx, betID, payout)
}

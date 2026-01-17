package services

import (
	"database/sql"
	"foresee/internal/models"

	"github.com/google/uuid"
)

type OutcomeService struct {
	Outcomes *models.OutcomeModel
}

func (s *OutcomeService) CreateDefaultForMarket(tx *sql.Tx, marketID uuid.UUID) error {
	return s.Outcomes.CreateWithYesNo(tx, marketID)
}

func (s *OutcomeService) ForMarkets(ids []uuid.UUID) (map[uuid.UUID][]models.Outcome, error) {
	return s.Outcomes.ForMarkets(ids)
}

func (s *OutcomeService) AddPoolAmount(tx *sql.Tx, id uuid.UUID, amount int) error {
	return s.Outcomes.AddPoolAmount(tx, id, amount)
}

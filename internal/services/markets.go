package services

import (
	"foresee/internal/models"
	"time"

	"github.com/google/uuid"
)

type MarketService struct {
	Markets *models.MarketModel
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

	err = s.Markets.Insert(
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

	return nil
}

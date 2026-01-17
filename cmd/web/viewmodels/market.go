package viewmodels

import (
	"foresee/internal/models"
	"time"
)

type MarketView struct {
	ID          string
	Title       string
	Description string
	Category    string
	Resolver    string
	ExpiresAt   string
	Status      string
	CreatedBy   string
	Outcomes    []OutcomeView
	TotalPool   int
}

func NewMarketView(m models.Market, loc *time.Location) MarketView {
	outcomes := make([]OutcomeView, len(m.Outcomes))
	totalPool := 0
	for i, o := range m.Outcomes {
		outcomes[i] = NewOutcomeView(o)
		totalPool += o.PoolAmount
	}

	return MarketView{
		ID:          m.ID.String(),
		Title:       m.Title,
		Description: m.Description,
		Category:    string(m.Category),
		Resolver:    string(m.ResolverType),
		ExpiresAt:   m.ExpiresAt.In(loc).Format("2006-01-02 15:04"),
		Status:      m.Status,
		Outcomes:    outcomes,
		TotalPool:   totalPool,
	}
}

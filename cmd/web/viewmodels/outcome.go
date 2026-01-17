package viewmodels

import "foresee/internal/models"

type OutcomeView struct {
	ID         string
	Label      string
	PoolAmount int
}

func NewOutcomeView(outcome models.Outcome) OutcomeView {
	return OutcomeView{
		ID:         outcome.ID.String(),
		Label:      outcome.Label,
		PoolAmount: outcome.PoolAmount,
	}
}

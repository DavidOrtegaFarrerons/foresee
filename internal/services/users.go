package services

import (
	"database/sql"
	"errors"
	"foresee/internal/models"
	"time"

	"github.com/google/uuid"
)

type UserService struct {
	Users *models.UserModel
}

var ErrDailyRewardNotAvailable = errors.New("daily reward already claimed")

const DailyRewardAmmount = 1000

func (s *UserService) ClaimDailyReward(id uuid.UUID) error {
	tx, err := s.Users.DB.Begin()
	t := time.Now().UTC()
	if err != nil {
		return err
	}
	defer tx.Rollback()

	user, err := s.Users.SelectForUpdate(tx, id)
	if err != nil {
		return err
	}

	if user.LastClaimedAt.Valid {
		if !CanClaimReward(t, user.LastClaimedAt.Time) {
			return ErrDailyRewardNotAvailable
		}
	}

	err = s.Users.ApplyDailyClaim(tx, id, user.Balance+DailyRewardAmmount, t)
	if err != nil {
		return err
	}

	return tx.Commit()
}

func (s *UserService) DecreaseBalanceBy(tx *sql.Tx, id uuid.UUID, amount int) error {
	return s.Users.DecreaseBalanceBy(tx, id, amount)
}

func CanClaimReward(now time.Time, lastClaimedAt time.Time) bool {
	return now.After(lastClaimedAt.Add(24 * time.Hour))
}

package populate

import (
	"context"
	"errors"
	"fmt"

	"github.com/brianvoe/gofakeit/v7"
	"github.com/google/uuid"

	"bum-service/internal/domain"
	"bum-service/internal/service/headmaster"
)

func (s *Service) generateHeadmaster(
	ctx context.Context, schoolID uuid.UUID, userID uuid.UUID,
) (*domain.Headmaster, error) {
	var (
		createdHeadmaster domain.Headmaster
		err               error
	)

	for try := 1; try <= 10; try++ {
		var (
			phone = gofakeit.Phone()
			email = gofakeit.Email()
		)

		createdHeadmaster, err = s.headmasterService().AddHeadmaster(ctx, headmaster.AddHeadmasterArgs{
			UserID:   userID,
			SchoolID: schoolID,
			Phone:    &phone,
			Email:    &email,
		})
		if err != nil {
			switch {
			case // continue unless we create a unique headmaster.
				errors.Is(err, domain.ErrUserEmailAlreadyExists),
				errors.Is(err, domain.ErrUserPhoneAlreadyExists):
				continue
			}

			return nil, fmt.Errorf("failed to create headmaster: %w", err)
		}

		break
	}

	return &createdHeadmaster, nil
}

package populate

import (
	"context"
	"errors"
	"fmt"

	"github.com/brianvoe/gofakeit/v7"
	"github.com/google/uuid"

	"bum-service/internal/domain"
	"bum-service/internal/service/director"
)

func (s *Service) generateDirector(
	ctx context.Context, schoolID uuid.UUID, userID uuid.UUID,
) (domain.Director, error) {
	var (
		createdDirector domain.Director
		err             error
	)

	for try := 1; try < 10; try++ {
		var (
			phone = gofakeit.Phone()
			email = gofakeit.Email()
		)

		createdDirector, err = s.directorService().AddDirector(ctx, director.AddDirectorArgs{
			UserID:   userID,
			SchoolID: schoolID,
			Phone:    &phone,
			Email:    &email,
		})
		if err != nil {
			switch {
			case // continue unless we create a unique director
				errors.Is(err, domain.ErrUserEmailAlreadyExists),
				errors.Is(err, domain.ErrUserPhoneAlreadyExists):
				continue
			}

			return domain.Director{}, fmt.Errorf("failed to create director: %w", err)
		}

		break
	}

	return createdDirector, nil
}

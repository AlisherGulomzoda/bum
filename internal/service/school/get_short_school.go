package school

import (
	"context"
	"fmt"

	"github.com/google/uuid"

	"bum-service/internal/domain"
)

// SchoolShortByIDs get school short list by ids.
func (s Service) SchoolShortByIDs(ctx context.Context, ids []uuid.UUID) (domain.SchoolShortInfos, error) {
	schoolInfos, err := s.schoolRepo.SchoolShortByIDsTx(ctx, ids)
	if err != nil {
		return nil, fmt.Errorf("failed to get school short info by ids from database: %w", err)
	}

	return schoolInfos, nil
}

// SchoolShortByID get school short by id.
func (s Service) SchoolShortByID(ctx context.Context, id uuid.UUID) (domain.SchoolShortInfo, error) {
	schoolInfo, err := s.schoolRepo.SchoolShortByIDTx(ctx, id)
	if err != nil {
		return domain.SchoolShortInfo{}, fmt.Errorf("failed to get school short info by id from database: %w", err)
	}

	return schoolInfo, nil
}

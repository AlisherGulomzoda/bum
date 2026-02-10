package school

import (
	"context"
	"fmt"

	"github.com/google/uuid"

	"bum-service/internal/domain"
)

// CreateSchoolArgs is a list of arguments to create a new school.
type CreateSchoolArgs struct {
	Name           string
	OrganizationID uuid.UUID
	Location       string
	Phone          *string
	Email          *string
}

// AddSchool adds a new school.
func (s Service) AddSchool(ctx context.Context, arg CreateSchoolArgs) (domain.School, error) {
	schoolEntity := domain.NewSchool(
		arg.Name,
		arg.OrganizationID,
		arg.Location,
		arg.Phone,
		arg.Email,
		s.now,
	)

	err := s.checkSchoolPhoneAndEmail(ctx, schoolEntity.OrganizationID, schoolEntity.Email, schoolEntity.Phone)
	if err != nil {
		return domain.School{}, err
	}

	err = s.schoolRepo.CreateSchoolTx(ctx, schoolEntity)
	if err != nil {
		return domain.School{}, fmt.Errorf("failed to create a new school to database: %w", err)
	}

	schoolEntity, err = s.SchoolByID(ctx, schoolEntity.ID)
	if err != nil {
		return domain.School{}, fmt.Errorf("failed to get school by id: %w", err)
	}

	return schoolEntity, nil
}

func (s Service) checkSchoolPhoneAndEmail(ctx context.Context, organizationID uuid.UUID, email, phone *string) error {
	if email != nil && *email != "" {
		schools, err := s.schoolRepo.SchoolListTx(
			ctx,
			domain.NewSchoolFilters(
				domain.NewListFilter(
					domain.SortOrderASC,
					domain.NewPagination(domain.PaginationDefaultPage, domain.PaginationDefaultLimit),
				),
				[]string{*email}, nil, nil,
			),
		)
		if err != nil {
			return fmt.Errorf("failed to get school list: %w", err)
		}

		if len(schools) > 0 && schools[0].OrganizationID != organizationID {
			return domain.ErrSchoolEmailAlreadyExists
		}
	}

	if phone != nil && *phone != "" {
		schools, err := s.schoolRepo.SchoolListTx(
			ctx,
			domain.NewSchoolFilters(
				domain.NewListFilter(
					domain.SortOrderASC,
					domain.NewPagination(domain.PaginationDefaultPage, domain.PaginationDefaultLimit),
				),
				nil, []string{*phone}, nil,
			),
		)
		if err != nil {
			return fmt.Errorf("failed to get school list: %w", err)
		}

		if len(schools) > 0 && schools[0].OrganizationID != organizationID {
			return domain.ErrSchoolPhoneAlreadyExists
		}
	}

	return nil
}

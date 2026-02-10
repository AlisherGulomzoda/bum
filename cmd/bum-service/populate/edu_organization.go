//nolint:gocognit,gosec // it's ok here
package populate

import (
	"context"
	"errors"
	"fmt"
	"github.com/brianvoe/gofakeit/v7"

	"bum-service/internal/domain"
	eduorganization "bum-service/internal/service/edu-organization"
)

const (
	eduOrganizationCount = 5

	schoolPerOrganizationMinCount = 1
	schoolPerOrganizationMaxCount = 10

	studentsPerSchoolMinCount = 200
	studentsPerSchoolMaxCount = 500

	studentsPerGroupMinCount = 15
	studentsPerGroupMaxCount = 40
)

func (s *Service) generateOrganizations(ctx context.Context) error {
	for count := range eduOrganizationCount {
		_ = count

		_, err := s.generateOrganization(ctx)
		if err != nil {
			return err
		}
	}

	return nil
}

func (s *Service) generateOrganization(ctx context.Context) (domain.EduOrganization, error) {
	var (
		organization domain.EduOrganization
		err          error
	)

	for try := 1; try <= 10; try++ {
		var (
			name = gofakeit.Company()
			logo = gofakeit.URL()
		)

		// fist we create an organization.
		organization, err = s.eduOrganizationService().CreateEduOrganization(ctx, eduorganization.CreateEduOrganizationArgs{
			Name: name,
			Logo: &logo,
		})
		if err != nil {
			// continue unless we create a unique organization.
			if errors.Is(err, domain.ErrEduOrganizationAlreadyExists) {
				continue
			}

			return domain.EduOrganization{}, fmt.Errorf("failed to create organization: %w", err)
		}

		break
	}

	// then we create school for the organization.
	err = s.generateSchools(ctx, organization)
	if err != nil {
		return domain.EduOrganization{}, err
	}

	return organization, nil
}

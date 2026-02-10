package subject

import (
	"context"
	"fmt"

	"bum-service/internal/domain"
)

// CreateSubjectArgs is a request for creating a new subject.
type CreateSubjectArgs struct {
	Name        string  `json:"name"`
	Description *string `json:"description"`
}

// CreateSubject creates a new subject.
func (s Service) CreateSubject(ctx context.Context, req CreateSubjectArgs) (createdSubject domain.Subject, err error) {
	createdSubject = domain.NewSubject(req.Name, req.Description, s.now)

	err = s.subjectRepo.CreateSubjectTx(ctx, createdSubject)
	if err != nil {
		err = fmt.Errorf("failed to create a new subject to database: %w", err)
		return domain.Subject{}, err
	}

	return createdSubject, nil
}

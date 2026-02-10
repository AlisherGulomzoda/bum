package school

import (
	"context"
	"fmt"

	"github.com/google/uuid"

	"bum-service/internal/domain"
	"bum-service/pkg/transaction"
)

// GroupSubjectByID returns group subject by ID.
func (s Service) GroupSubjectByID(ctx context.Context, id uuid.UUID) (domain.GroupSubject, error) {
	groupSubject, err := s.groupSubjectsRepo.GroupSubjectByIDTx(ctx, id)
	if err != nil {
		return domain.GroupSubject{}, fmt.Errorf("failed to get group subject by id from database: %w", err)
	}

	if groupSubject.HasTeacher() {
		teacher, err := s.teacherService.TeacherByID(ctx, *groupSubject.TeacherID)
		if err != nil {
			return domain.GroupSubject{}, fmt.Errorf("failed to get teacher by id: %w", err)
		}

		groupSubject.SetTeacher(teacher)
	}

	return groupSubject, nil
}

// GroupSubjectList returns a list of subjects assigned to a group.
func (s Service) GroupSubjectList(ctx context.Context, groupID uuid.UUID) (domain.GroupSubjects, error) {
	list, err := s.groupSubjectsRepo.GroupSubjectListTx(ctx, groupID)
	if err != nil {
		return domain.GroupSubjects{}, fmt.Errorf("failed to get school subject list from database: %w", err)
	}

	teacherIDs := list.TeacherIDs()
	schoolSubjectIDs := list.SchoolSubjectIDs()

	schoolSubjects, err := s.SchoolSubjectByIDs(ctx, schoolSubjectIDs)
	if err != nil {
		return nil, fmt.Errorf("failed to get school subjects by ids: %w", err)
	}

	teachers, err := s.teacherService.TeachersByIDs(ctx, teacherIDs)
	if err != nil {
		return domain.GroupSubjects{}, fmt.Errorf("failed to get school teachers by ids: %w", err)
	}

	list.SetSchoolSubjects(schoolSubjects)
	list.SetTeachers(teachers)

	return list, nil
}

// AddGroupSubjectArgs is a struct for adding a subject to a group.
type AddGroupSubjectArgs struct {
	SchoolSubjectID uuid.UUID
	TeacherID       *uuid.UUID
	Count           *int16
}

// AddGroupSubject adds subjects to a group.
func (s Service) AddGroupSubject(
	ctx context.Context, _, groupID uuid.UUID, args AddGroupSubjectArgs,
) (domain.GroupSubject, error) {
	group, err := s.groupRepo.GroupByIDTx(ctx, groupID)
	if err != nil {
		return domain.GroupSubject{}, fmt.Errorf("failed to get group from database: %w", err)
	}

	newGroupSubject := domain.NewGroupSubject(
		args.SchoolSubjectID,
		group.ID,
		args.TeacherID,
		args.Count,

		s.now,
	)

	// TODO: Implement the school subjects, teachers validate logic.
	txCtx, tx, err := s.sessionAdapter.Begin(ctx)
	if err != nil {
		return domain.GroupSubject{}, fmt.Errorf("failed to begin transaction : %w", err)
	}

	defer func(tx transaction.SessionSolver) {
		errEnd := s.sessionAdapter.End(tx, err)
		if errEnd != nil {
			err = fmt.Errorf(
				"failed to end transaction on adding group subjects : %w: %w", domain.ErrInternalServerError, errEnd,
			)
		}
	}(tx)

	err = s.groupSubjectsRepo.AddGroupSubjectsTx(txCtx, newGroupSubject)
	if err != nil {
		return domain.GroupSubject{}, fmt.Errorf("failed add group subjects : %w", err)
	}

	newGroupSubject, err = s.GroupSubjectByID(txCtx, newGroupSubject.ID)
	if err != nil {
		return domain.GroupSubject{}, fmt.Errorf(
			"failed to get group subject by id=%s: %w", newGroupSubject.ID, err,
		)
	}

	return newGroupSubject, nil
}

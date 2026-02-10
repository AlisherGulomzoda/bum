package student

import (
	"context"
	"fmt"

	"github.com/google/uuid"

	"bum-service/internal/domain"
	"bum-service/pkg/transaction"
)

// AssignStudentGuardianArgs is arguments for AssignStudentGuardian method.
type AssignStudentGuardianArgs struct {
	StudentID uuid.UUID
	UserID    uuid.UUID
	Relation  string
	SchoolID  uuid.UUID
}

// AssignStudentGuardian assigns a guardian to a student.
func (s Service) AssignStudentGuardian(
	ctx context.Context,
	args AssignStudentGuardianArgs,
) (studentGuardian domain.StudentGuardian, err error) {
	newStudentGuardian, err := domain.NewStudentGuardian(
		args.StudentID,
		args.UserID,
		domain.StudentGuardianRelation(args.Relation),
		args.SchoolID,
		s.now,
	)
	if err != nil {
		return domain.StudentGuardian{}, fmt.Errorf("failed to create new student guardian domain: %w", err)
	}

	txCtx, tx, err := s.sessionAdapter.Begin(ctx)
	if err != nil {
		return domain.StudentGuardian{}, fmt.Errorf("failed to begin transaction : %w", err)
	}

	defer func(tx transaction.SessionSolver) {
		errEnd := s.sessionAdapter.End(tx, err)
		if errEnd != nil {
			err = fmt.Errorf(
				"failed to end transaction on assign student guardian: %w: %w",
				domain.ErrInternalServerError, errEnd,
			)
		}
	}(tx)

	_, err = s.userService.AddRoleToUser(txCtx, args.UserID, domain.RoleGuardian, nil, nil)
	if err != nil {
		return domain.StudentGuardian{}, fmt.Errorf("failed to add role guardian to user: %w", err)
	}

	err = s.studentRepo.AddStudentGuardianTx(txCtx, newStudentGuardian)
	if err != nil {
		return domain.StudentGuardian{}, fmt.Errorf("failed to create a new student guardian to database: %w", err)
	}

	studentGuardian, err = s.studentRepo.StudentGuardianByIDTx(txCtx, newStudentGuardian.ID)
	if err != nil {
		return domain.StudentGuardian{}, fmt.Errorf("failed to get student guardian by id: %w", err)
	}

	user, err := s.userInfoService.UserByID(ctx, studentGuardian.UserID)
	if err != nil {
		return domain.StudentGuardian{}, fmt.Errorf("failed to get user info by id: %w", err)
	}

	student, err := s.studentRepo.StudentByIDTx(ctx, studentGuardian.StudentID)
	if err != nil {
		return domain.StudentGuardian{}, fmt.Errorf("failed to get student info by id: %w", err)
	}

	studentGuardian.SetUser(user)
	studentGuardian.SetStudent(student)

	return studentGuardian, nil
}

// StudentGuardians returns a list of guardians for a student.
func (s Service) StudentGuardians(ctx context.Context, studentID uuid.UUID) (domain.StudentGuardians, error) {
	list, err := s.studentRepo.StudentGuardiansByStudentIDTx(ctx, studentID)
	if err != nil {
		return nil, fmt.Errorf("failed get student guardians by student id database: %w", err)
	}

	userIDs := list.UserIDs()
	studentIDs := list.StudentIDs()

	users, err := s.userInfoService.UsersByIDs(ctx, userIDs)
	if err != nil {
		return nil, fmt.Errorf("failed get users info by ids: %w", err)
	}

	students, err := s.studentRepo.StudentsByIDsTx(ctx, studentIDs)
	if err != nil {
		return nil, fmt.Errorf("failed get students info by ids: %w", err)
	}

	list.SetUsers(users)
	list.SetStudents(students)

	return list, nil
}

// StudentGuardianByUserID returns a student guardian by guardian user_id.
func (s Service) StudentGuardianByUserID(ctx context.Context, userID uuid.UUID) (domain.Guardian, error) {
	user, err := s.userInfoService.UserByID(ctx, userID)
	if err != nil {
		return domain.Guardian{}, fmt.Errorf("failed get user info by id: %w", err)
	}

	guardianStudentList, err := s.studentRepo.StudentGuardianByUserIDTx(ctx, userID)
	if err != nil {
		return domain.Guardian{}, fmt.Errorf("failed get guardian's students s by guardian user_id from database: %w", err)
	}

	studentIDs := guardianStudentList.StudentIDs()

	students, err := s.studentRepo.StudentsByIDsTx(ctx, studentIDs)
	if err != nil {
		return domain.Guardian{}, fmt.Errorf("failed get students info by ids: %w", err)
	}

	guardianStudentList.SetStudents(students)

	var guardian domain.Guardian
	guardian.SetUser(user)
	guardian.SetStudents(guardianStudentList)

	return guardian, nil
}

// StudentGuardianList student guardian list.
func (s Service) StudentGuardianList(
	ctx context.Context,
	filters domain.StudentGuardianListFilter,
) (domain.StudentGuardians, int, error) {
	list, err := s.studentRepo.StudentGuardianListTx(ctx, filters)
	if err != nil {
		return nil, 0, fmt.Errorf("failed get student guardian list from database: %w", err)
	}

	var (
		userIDs    = list.UserIDs()
		studentIDs = list.StudentIDs()
	)

	users, err := s.userInfoService.UsersByIDs(ctx, userIDs)
	if err != nil {
		return nil, 0, fmt.Errorf("failed get user info from userinfo service: %w", err)
	}

	students, err := s.studentRepo.StudentsByIDsTx(ctx, studentIDs)
	if err != nil {
		return nil, 0, fmt.Errorf("failed get student info by ids from database: %w", err)
	}

	list.SetUsers(users)
	list.SetStudents(students)

	total, err := s.studentRepo.StudentGuardianListCountTx(ctx, filters)
	if err != nil {
		return nil, 0, fmt.Errorf("failed get student guardian list count from database: %w", err)
	}

	return list, total, nil
}

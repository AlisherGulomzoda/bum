package handlers

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"

	"bum-service/internal/controller/http/handlers/request"
	"bum-service/internal/controller/http/handlers/response"
	"bum-service/internal/domain"
	"bum-service/internal/service/student"
	"bum-service/pkg/liblog"
)

// Student is a handler for Students.
type Student struct {
	studentService IStudentService
}

// NewStudent creates a new Subject handler.
func NewStudent(
	studentService IStudentService,
) *Student {
	return &Student{
		studentService: studentService,
	}
}

// AddStudent adds a new Student.
func (s Student) AddStudent(c *gin.Context) {
	var (
		ctx    = c.Request.Context()
		logger = liblog.Must(ctx)
		req    request.AddStudent
	)

	if err := c.ShouldBindJSON(&req); err != nil {
		logger.Errorf("failed to bind: %v", c.Error(domain.NewBadRequest(err.Error())))
		return
	}

	logger = logger.WithFields(liblog.Fields{"request": req})

	addedStudent, err := s.studentService.AddStudent(
		ctx,
		student.AddStudentArgs{
			UserID:   req.UserID,
			GroupID:  req.GroupID,
			SchoolID: req.SchoolID,
		})
	if err != nil {
		logger.Errorf("failed to add student: %v", c.Error(err))
		return
	}

	c.JSON(http.StatusCreated, response.NewStudent(addedStudent))
}

// StudentByID get student by id.
func (s Student) StudentByID(c *gin.Context) {
	var (
		ctx       = c.Request.Context()
		logger    = liblog.Must(ctx)
		studentID = request.GetStudentIDPathVar(c)
	)

	logger = logger.WithFields(liblog.Fields{"student_id": studentID})

	studentUUID, err := uuid.Parse(studentID)
	if err != nil {
		logger.Errorf("failed to parse student id to uuid: %v %v", err, c.Error(domain.ErrBadRequest))
		return
	}

	addedStudent, err := s.studentService.StudentByID(ctx, studentUUID)
	if err != nil {
		logger.Errorf("failed to create subject: %v", c.Error(err))
		return
	}

	c.JSON(http.StatusOK, response.NewStudent(addedStudent))
}

// StudentList get student list.
func (s Student) StudentList(c *gin.Context) {
	var (
		ctx    = c.Request.Context()
		logger = liblog.Must(ctx)
		req    request.StudentList
	)

	if err := c.ShouldBindQuery(&req); err != nil {
		logger.Errorf("failed to bind: %v", c.Error(domain.NewBadRequest(err.Error())))
		return
	}

	logger = logger.WithFields(liblog.Fields{
		"request": req,
	})

	list, total, err := s.studentService.StudentList(
		ctx,
		domain.NewStudentListFilter(
			domain.NewDateFilter(req.CreatedDate.DateFrom(), req.CreatedDate.DateTill()),
			domain.NewListFilter(req.SortOrder, domain.NewPagination(req.Page, req.PerPage)),

			req.GroupUUIDs(),
			req.SchoolUUIDs(),
			req.OrganizationUUIDs(),
		),
	)
	if err != nil {
		logger.Errorf("failed to get students list: %v", c.Error(err))
		return
	}

	c.JSON(http.StatusOK, response.NewStudentList(list, response.Pagination{
		Page:    req.Page,
		PerPage: req.PerPage,
		Total:   total,
	}))
}

// AssignStudentGuardian assign guardian to student.
func (s Student) AssignStudentGuardian(c *gin.Context) {
	var (
		ctx       = c.Request.Context()
		logger    = liblog.Must(ctx)
		studentID = request.GetStudentIDPathVar(c)
		req       request.AssignStudentGuardian
	)

	if err := c.ShouldBindJSON(&req); err != nil {
		logger.Errorf("failed to bind: %v", c.Error(domain.NewBadRequest(err.Error())))
		return
	}

	studentUUID, err := uuid.Parse(studentID)
	if err != nil {
		logger.Errorf("failed to parse student id to uuid: %v %v", err, c.Error(domain.ErrBadRequest))
		return
	}

	logFields := req.LogFields()
	logFields["student_id"] = studentUUID

	logger = logger.WithFields(logFields)

	studentGuardian, err := s.studentService.AssignStudentGuardian(ctx,
		student.AssignStudentGuardianArgs{
			StudentID: studentUUID,
			UserID:    req.UserID,
			Relation:  req.Relation,
			SchoolID:  req.SchoolID,
		})
	if err != nil {
		logger.Errorf("failed to assign student guardian: %v", c.Error(err))
		return
	}

	c.JSON(http.StatusOK, response.NewStudentGuardian(studentGuardian))
}

// StudentGuardians get student guardians by student id.
func (s Student) StudentGuardians(c *gin.Context) {
	var (
		ctx       = c.Request.Context()
		logger    = liblog.Must(ctx)
		studentID = request.GetStudentIDPathVar(c)
	)

	logger = logger.WithFields(liblog.Fields{"student_id": studentID})

	studentUUID, err := uuid.Parse(studentID)
	if err != nil {
		logger.Errorf("failed to parse student id to uuid: %v %v", err, c.Error(domain.ErrBadRequest))
		return
	}

	studentGuardians, err := s.studentService.StudentGuardians(ctx, studentUUID)
	if err != nil {
		logger.Errorf("failed to get student guardians: %v", c.Error(err))
		return
	}

	c.JSON(http.StatusOK, response.NewStudentGuardiansForStudent(studentGuardians))
}

// StudentGuardianByUserID get student guardian by guardian user_id.
func (s Student) StudentGuardianByUserID(c *gin.Context) {
	var (
		ctx      = c.Request.Context()
		logger   = liblog.Must(ctx)
		reqParam = request.GetUserIDPathVar(c)
		userID   uuid.UUID
		err      error
	)

	logger = logger.WithFields(liblog.Fields{"request": liblog.Fields{"teacher_id": reqParam}})
	ctx = liblog.With(ctx, logger)

	if userID, err = uuid.Parse(reqParam); err != nil {
		logger.Errorf("failed to parse to uuid: %v", c.Error(domain.NewBadRequest(err.Error())))
		return
	}

	logger = logger.WithFields(liblog.Fields{
		"request": userID,
	})

	guardian, err := s.studentService.StudentGuardianByUserID(ctx, userID)
	if err != nil {
		logger.Errorf("failed to get student guardian list: %v", c.Error(err))
		return
	}

	c.JSON(http.StatusOK, response.NewGuardian(guardian))
}

// StudentGuardianList get student guardians by student id.
func (s Student) StudentGuardianList(c *gin.Context) {
	var (
		ctx    = c.Request.Context()
		logger = liblog.Must(ctx)
		req    request.StudentGuardianList
	)

	if err := c.ShouldBindQuery(&req); err != nil {
		logger.Errorf("failed to bind: %v", c.Error(domain.NewBadRequest(err.Error())))
		return
	}

	logger = logger.WithFields(liblog.Fields{
		"request": req,
	})

	list, total, err := s.studentService.StudentGuardianList(
		ctx,
		domain.NewStudentGuardianListFilter(
			domain.NewListFilter(req.SortOrder, domain.NewPagination(req.Page, req.PerPage)),
			domain.NewDateFilter(req.CreatedDate.DateFrom(), req.CreatedDate.DateTill()),
			req.GroupUUIDs(),
			req.SchoolUUIDs(),
			req.OrganizationUUIDs(),
		),
	)
	if err != nil {
		logger.Errorf("failed to get student guardian list: %v", c.Error(err))
		return
	}

	c.JSON(http.StatusOK, response.NewStudentGuardianList(list, response.Pagination{
		Page:    req.Page,
		PerPage: req.PerPage,
		Total:   total,
	}))
}

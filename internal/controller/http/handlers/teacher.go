package handlers

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"

	"bum-service/internal/controller/http/handlers/request"
	"bum-service/internal/controller/http/handlers/response"
	"bum-service/internal/domain"
	"bum-service/internal/service/teacher"
	"bum-service/pkg/liblog"
)

// Teacher is teacher handler.
type Teacher struct {
	teacherSvc ITeacherService
}

// NewTeacher creates a new teacher handler.
func NewTeacher(
	teacherUseCase ITeacherService,
) *Teacher {
	return &Teacher{
		teacherSvc: teacherUseCase,
	}
}

// AddTeacher creates a new teacher.
func (t Teacher) AddTeacher(c *gin.Context) {
	var (
		ctx    = c.Request.Context()
		logger = liblog.Must(ctx)
		req    request.AddTeacher
	)

	if err := c.ShouldBindJSON(&req); err != nil {
		logger.Errorf("failed to bind: %v", c.Error(domain.NewBadRequest(err.Error())))
		return
	}

	logger = logger.WithFields(liblog.Fields{"request": req.LogFields()})
	ctx = liblog.With(ctx, logger)

	entity, err := t.teacherSvc.AddTeacher(
		ctx,
		teacher.AddTeacherArgs{
			UserID:   req.UserID,
			SchoolID: req.SchoolID,
			Phone:    req.Phone,
			Email:    req.Email,
		})
	if err != nil {
		logger.Errorf("failed to add a new teacher: %v", c.Error(err))
		return
	}

	c.JSON(http.StatusCreated, response.NewTeacher(entity))
}

// TeacherByID get teacher by id.
func (t Teacher) TeacherByID(c *gin.Context) { //nolint:revive // It's handler
	var (
		ctx       = c.Request.Context()
		logger    = liblog.Must(ctx)
		reqParam  = request.GetTeacherIDPathVar(c)
		teacherID uuid.UUID
		err       error
	)

	logger = logger.WithFields(liblog.Fields{"request": liblog.Fields{"teacher_id": reqParam}})
	ctx = liblog.With(ctx, logger)

	if teacherID, err = uuid.Parse(reqParam); err != nil {
		logger.Errorf("failed to parse to uuid: %v", c.Error(domain.NewBadRequest(err.Error())))
		return
	}

	teacherEntity, err := t.teacherSvc.TeacherByID(ctx, teacherID)
	if err != nil {
		logger.Errorf("failed to get teacher by id: %v", c.Error(err))
		return
	}

	c.JSON(http.StatusOK, response.NewTeacher(teacherEntity))
}

// ListTeacher endpoint for getting a list of teachers.
func (t Teacher) ListTeacher(c *gin.Context) {
	var (
		ctx    = c.Request.Context()
		logger = liblog.Must(ctx)
		req    request.TeacherList
	)

	if err := c.ShouldBindQuery(&req); err != nil {
		logger.Errorf("failed to bind: %v", c.Error(domain.NewBadRequest(err.Error())))
		return
	}

	logger = logger.WithFields(liblog.Fields{
		"request": req,
	})

	list, total, err := t.teacherSvc.TeacherList(
		ctx,
		domain.NewTeacherListFilter(
			domain.NewListFilter(req.SortOrder, domain.NewPagination(req.Page, req.PerPage)),
			domain.NewDateFilter(req.CreatedDate.DateFrom(), req.CreatedDate.DateTill()),
			req.GroupUUIDs(),
			req.SchoolUUIDs(),
			req.OrganizationUUIDs(),
		),
	)
	if err != nil {
		logger.Errorf("failed to get teacher list: %v", c.Error(err))
		return
	}

	c.JSON(http.StatusOK, response.NewTeacherList(list, response.Pagination{
		Page:    req.Page,
		PerPage: req.PerPage,
		Total:   total,
	}))
}

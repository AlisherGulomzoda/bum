package handlers

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"

	"bum-service/internal/controller/http/handlers/request"
	"bum-service/internal/controller/http/handlers/response"
	"bum-service/internal/domain"
	"bum-service/internal/service/subject"
	"bum-service/pkg/liblog"
)

// Subject is a handler for Subject.
type Subject struct {
	subjectService ISubjectService
}

// NewSubject creates a new Subject handler.
func NewSubject(subjectService ISubjectService) *Subject {
	return &Subject{
		subjectService: subjectService,
	}
}

// CreateSubject creates a new Subject.
func (h Subject) CreateSubject(c *gin.Context) {
	var (
		ctx    = c.Request.Context()
		logger = liblog.Must(ctx)
		req    request.CreateSubject
	)

	if err := c.ShouldBindJSON(&req); err != nil {
		logger.Errorf("failed to bind: %v", c.Error(domain.NewBadRequest(err.Error())))
		return
	}

	logger = logger.WithFields(liblog.Fields{"request": req.LogFields()})

	createdSubject, err := h.subjectService.CreateSubject(
		ctx,
		subject.CreateSubjectArgs{
			Name:        req.Name,
			Description: req.Description,
		})
	if err != nil {
		logger.Errorf("failed to create subject: %v", c.Error(err))
		return
	}

	c.JSON(http.StatusCreated, response.NewSubject(createdSubject))
}

// SubjectByID returns a Subject by id.
func (h Subject) SubjectByID(c *gin.Context) {
	var (
		ctx       = c.Request.Context()
		logger    = liblog.Must(ctx)
		reqParam  = request.GetSubjectIDPathVar(c)
		subjectID uuid.UUID
		err       error
	)

	logger = logger.WithFields(liblog.Fields{"request": liblog.Fields{"subject_id": reqParam}})

	if subjectID, err = uuid.Parse(reqParam); err != nil {
		logger.Errorf("failed to bind: %v", c.Error(domain.NewBadRequest(err.Error())))
		return
	}

	subjectData, err := h.subjectService.SubjectByID(ctx, subjectID)
	if err != nil {
		logger.Errorf("failed to get subject: %v", c.Error(err))
		return
	}

	c.JSON(http.StatusOK, response.NewSubject(subjectData))
}

// SubjectList returns a list of Subjects.
func (h Subject) SubjectList(c *gin.Context) {
	var (
		ctx    = c.Request.Context()
		logger = liblog.Must(ctx)
		req    request.SubjectList
	)

	if err := c.ShouldBindQuery(&req); err != nil {
		logger.Errorf("failed to bind: %v", c.Error(domain.NewBadRequest(err.Error())))
		return
	}

	logger = logger.WithFields(liblog.Fields{"request": req})

	subjectList, count, err := h.subjectService.SubjectList(
		ctx,
		domain.NewSubjectListFilter(
			domain.NewListFilter(
				req.SortOrder,
				domain.NewPagination(
					req.Page,
					req.PerPage,
				),
			),
		),
	)
	if err != nil {
		logger.Errorf("failed to get subject list: %v", c.Error(err))
		return
	}

	c.JSON(http.StatusOK, response.NewSubjectList(subjectList, response.Pagination{
		Page:    req.Page,
		PerPage: req.PerPage,
		Total:   count,
	}))
}

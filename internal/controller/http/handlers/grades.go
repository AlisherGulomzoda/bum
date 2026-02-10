package handlers

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"

	"bum-service/internal/controller/http/handlers/request"
	"bum-service/internal/controller/http/handlers/response"
	"bum-service/internal/domain"
	grades "bum-service/internal/service/grade-standard"
	"bum-service/pkg/liblog"
)

// Grades is grade-standard handler.
type Grades struct {
	gradesService IGradesService
}

// NewGrades creates a new grade-standard handler.
func NewGrades(
	gradesService IGradesService,
) *Grades {
	return &Grades{
		gradesService: gradesService,
	}
}

// CreateGradeStandard creates a new grade standard.
func (h Grades) CreateGradeStandard(c *gin.Context) {
	var (
		ctx    = c.Request.Context()
		logger = liblog.Must(ctx)
		req    request.CreateGradeStandard
	)

	if err := c.ShouldBindJSON(&req); err != nil {
		logger.Errorf("failed to bind: %v", c.Error(domain.NewBadRequest(err.Error())))
		return
	}

	logger = logger.WithFields(liblog.Fields{"request": req.LogFields()})
	ctx = liblog.With(ctx, logger)

	gradesArg := gradesInRequest(req)

	entity, err := h.gradesService.CreateGradeStandard(
		ctx,
		grades.CreateGradeStandardArgs{
			OrganizationID: req.OrganizationID,
			Name:           req.Name,
			EducationYears: req.EducationYears,
			Description:    req.Description,
			Grades:         gradesArg,
		})
	if err != nil {
		logger.Errorf("failed to create a new grade standard: %v", c.Error(err))
		return
	}

	c.JSON(http.StatusCreated, response.NewGradeStandard(entity))
}

// gradesInRequest get grades from request.
func gradesInRequest(req request.CreateGradeStandard) []grades.CreateGradeArgs {
	gradesArg := make([]grades.CreateGradeArgs, len(req.Grades))

	for i, grade := range req.Grades {
		gradesArg[i] = grades.CreateGradeArgs{
			Name:          grade.Name,
			EducationYear: grade.EducationYear,
		}
	}

	return gradesArg
}

// GradeStandardByID get grade standard by id.
func (h Grades) GradeStandardByID(c *gin.Context) {
	var (
		ctx             = c.Request.Context()
		logger          = liblog.Must(ctx)
		reqParam        = request.GetGradeStandardIDPathVar(c)
		gradeStandardID uuid.UUID
		err             error
	)

	logger = logger.WithFields(liblog.Fields{"request": liblog.Fields{"grade_standard_id": reqParam}})
	ctx = liblog.With(ctx, logger)

	if gradeStandardID, err = uuid.Parse(reqParam); err != nil {
		logger.Errorf("failed to bind: %v", c.Error(domain.NewBadRequest(err.Error())))
		return
	}

	gradeStandard, err := h.gradesService.GradeStandardByID(ctx, gradeStandardID)
	if err != nil {
		logger.Errorf("failed to get grade standard by id: %v", c.Error(err))
		return
	}

	c.JSON(http.StatusOK, response.NewGradeStandard(gradeStandard))
}

// GradeStandardList endpoint for getting a list of grade standard.
func (h Grades) GradeStandardList(c *gin.Context) {
	var (
		ctx    = c.Request.Context()
		logger = liblog.Must(ctx)
		req    request.GradeStandardList
	)

	if err := c.ShouldBindQuery(&req); err != nil {
		logger.Errorf("failed to bind: %v", c.Error(domain.NewBadRequest(err.Error())))
		return
	}

	logger = logger.WithFields(liblog.Fields{"request": req})
	ctx = liblog.With(ctx, logger)

	list, total, err := h.gradesService.GradeStandardList(
		ctx,
		domain.NewGradeStandardListFilter(
			domain.NewListFilter(
				req.SortOrder,
				domain.NewPagination(req.Page, req.PerPage),
			),
		),
	)
	if err != nil {
		logger.Errorf("failed to get grade standard list: %v", c.Error(err))
		return
	}

	c.JSON(
		http.StatusOK,
		response.NewGradeStandardList(
			list,
			response.Pagination{
				Page:    req.Page,
				PerPage: req.PerPage,
				Total:   total,
			},
		),
	)
}

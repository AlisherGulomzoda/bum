package handlers

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"

	"bum-service/internal/controller/http/handlers/request"
	"bum-service/internal/controller/http/handlers/response"
	"bum-service/internal/domain"
	"bum-service/internal/service/school"
	"bum-service/pkg/liblog"
)

// CreateAuditorium creates a new auditorium.
func (s School) CreateAuditorium(c *gin.Context) {
	var (
		ctx             = c.Request.Context()
		logger          = liblog.Must(ctx)
		req             request.CreateAuditorium
		schoolIDPathVar = request.GetSchoolIDPathVar(c)
		schoolID        uuid.UUID
		err             error
	)

	if schoolID, err = uuid.Parse(schoolIDPathVar); err != nil {
		logger.Errorf("failed to parse uuid: %v", c.Error(domain.NewBadRequest(err.Error())))
		return
	}

	if err = c.ShouldBindJSON(&req); err != nil {
		logger.Errorf("failed to bind: %v", c.Error(domain.NewBadRequest(err.Error())))
		return
	}

	logger = logger.WithFields(liblog.Fields{
		"request":   req.LogFields(),
		"school_id": schoolID,
	})

	createdAuditorium, err := s.schoolService.CreateAuditorium(
		ctx,
		schoolID,
		school.CreateAuditoriumArgs{
			Name:            req.Name,
			SchoolSubjectID: req.SchoolSubjectsID,
			Description:     req.Description,
		})
	if err != nil {
		logger.Errorf("failed to create school auditorium: %v", c.Error(err))
		return
	}

	c.JSON(http.StatusCreated, response.NewSchoolAuditorium(createdAuditorium))
}

// AuditoriumList returns a list of auditoriums.
func (s School) AuditoriumList(c *gin.Context) {
	var (
		ctx             = c.Request.Context()
		logger          = liblog.Must(ctx)
		req             request.AuditoriumList
		schoolIDPathVar = request.GetSchoolIDPathVar(c)
		schoolID        uuid.UUID
		err             error
	)

	if schoolID, err = uuid.Parse(schoolIDPathVar); err != nil {
		logger.Errorf("failed to parse uuid: %v", c.Error(domain.NewBadRequest(err.Error())))
		return
	}

	if err = c.ShouldBindQuery(&req); err != nil {
		logger.Errorf("failed to bind: %v", c.Error(domain.NewBadRequest(err.Error())))
		return
	}

	logger = logger.WithFields(liblog.Fields{
		"request":   req,
		"school_id": schoolIDPathVar,
	})
	ctx = liblog.With(ctx, logger)

	list, total, err := s.schoolService.AuditoriumList(
		ctx,
		domain.NewAuditoriumListFilters(
			domain.NewListFilter(req.SortOrder, domain.NewPagination(req.Page, req.PerPage)),
			schoolID,
		),
	)
	if err != nil {
		logger.Errorf("failed to get auditoriums: %v", c.Error(err))
		return
	}

	c.JSON(http.StatusOK, response.NewSchoolAuditoriumList(list, response.Pagination{
		Page:    req.Page,
		PerPage: req.PerPage,
		Total:   total,
	}))
}

// AuditoriumByIDAndSchoolID returns an auditorium by ID and school ID.
func (s School) AuditoriumByIDAndSchoolID(c *gin.Context) {
	var (
		ctx                 = c.Request.Context()
		logger              = liblog.Must(ctx)
		schoolIDPathVar     = request.GetSchoolIDPathVar(c)
		auditoriumIDPathVar = request.GetSchoolAuditoriumIDPathVar(c)
		schoolID            uuid.UUID
		auditoriumID        uuid.UUID
		err                 error
	)

	logger = logger.WithFields(liblog.Fields{
		"school_id":     schoolIDPathVar,
		"auditorium_id": auditoriumIDPathVar,
	})
	ctx = liblog.With(ctx, logger)

	logger = logger.WithFields(liblog.Fields{"school_id": schoolIDPathVar})
	ctx = liblog.With(ctx, logger)

	if schoolID, err = uuid.Parse(schoolIDPathVar); err != nil {
		logger.Errorf("failed to parse uuid: %v", c.Error(domain.NewBadRequest(err.Error())))
		return
	}

	if auditoriumID, err = uuid.Parse(auditoriumIDPathVar); err != nil {
		logger.Errorf("failed to parse uuid: %v", c.Error(domain.NewBadRequest(err.Error())))
		return
	}

	auditoriumDomain, err := s.schoolService.AuditoriumByIDAndSchoolID(ctx, auditoriumID, schoolID)
	if err != nil {
		logger.Errorf("failed to create a new school subject: %v", c.Error(err))
		return
	}

	c.JSON(http.StatusOK, response.NewSchoolAuditorium(auditoriumDomain))
}

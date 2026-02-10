package handlers

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"

	"bum-service/internal/controller/http/handlers/request"
	"bum-service/internal/controller/http/handlers/response"
	"bum-service/internal/domain"
	eduorganization "bum-service/internal/service/edu-organization"
	"bum-service/pkg/liblog"
)

// EduOrganization is educational organization handler.
type EduOrganization struct {
	eduOrganizationService IEduOrganizationService
}

// NewEduOrganization creates a new educational organization handler.
func NewEduOrganization(eduOrganizationService IEduOrganizationService) EduOrganization {
	return EduOrganization{
		eduOrganizationService: eduOrganizationService,
	}
}

// CreateEduOrganization creates a new educational organization.
func (s EduOrganization) CreateEduOrganization(c *gin.Context) {
	var (
		ctx    = c.Request.Context()
		logger = liblog.Must(ctx)
		req    request.CreateEduOrganizational
	)

	if err := c.ShouldBindJSON(&req); err != nil {
		logger.Errorf("failed to bind: %v", c.Error(domain.NewBadRequest(err.Error())))
		return
	}

	logger = logger.WithFields(liblog.Fields{"request": req.LogFields()})
	ctx = liblog.With(ctx, logger)

	entity, err := s.eduOrganizationService.CreateEduOrganization(
		ctx,
		eduorganization.CreateEduOrganizationArgs{
			Name:        req.Name,
			Logo:        req.Logo,
			Description: req.Description,
		},
	)
	if err != nil {
		logger.Errorf("failed to create a new educational organization: %v", c.Error(err))
		return
	}

	c.JSON(http.StatusCreated, response.NewEduOrganization(entity))
}

// EduOrganizationByID get educational organization by id.
func (s EduOrganization) EduOrganizationByID(c *gin.Context) {
	var (
		ctx            = c.Request.Context()
		logger         = liblog.Must(ctx)
		id             = request.GetEduOrganizationPathVar(c)
		err            error
		organizationID uuid.UUID
	)

	logger = logger.WithFields(liblog.Fields{"request": liblog.Fields{"edu_organization_id": id}})
	ctx = liblog.With(ctx, logger)

	if organizationID, err = uuid.Parse(id); err != nil {
		logger.Errorf("failed to bind: %v", c.Error(domain.NewBadRequest(err.Error())))
		return
	}

	entity, err := s.eduOrganizationService.EduOrganizationByID(ctx, organizationID)
	if err != nil {
		logger.Errorf("failed to get educational organization by id: %v", c.Error(err))
		return
	}

	c.JSON(http.StatusOK, response.NewEduOrganization(entity))
}

// UpdateEduOrganizationByID update educational organization by id.
func (s EduOrganization) UpdateEduOrganizationByID(c *gin.Context) {
	var (
		ctx            = c.Request.Context()
		logger         = liblog.Must(ctx)
		id             = request.GetEduOrganizationPathVar(c)
		req            request.UpdateEduOrganizational
		organizationID uuid.UUID
		err            error
	)

	if organizationID, err = uuid.Parse(id); err != nil {
		logger.Errorf("failed to bind: %v", c.Error(domain.NewBadRequest(err.Error())))
		return
	}

	if err := c.ShouldBindJSON(&req); err != nil {
		logger.Errorf("failed to bind: %v", c.Error(domain.NewBadRequest(err.Error())))
		return
	}

	logger = logger.WithFields(liblog.Fields{"request": req, "edu_organization_id": organizationID})
	ctx = liblog.With(ctx, logger)

	entity, err := s.eduOrganizationService.UpdateEduOrganizationByID(
		ctx,
		eduorganization.UpdateEduOrganizationArgs{
			ID:   organizationID,
			Name: req.Name,
			Logo: req.Logo,
		},
	)
	if err != nil {
		logger.Errorf("failed to update educational organization by id: %v", c.Error(err))
		return
	}

	c.JSON(http.StatusOK, response.NewEduOrganization(entity))
}

// EduOrganizationList returns a list of educational organizations.
func (s EduOrganization) EduOrganizationList(c *gin.Context) {
	var (
		ctx    = c.Request.Context()
		logger = liblog.Must(ctx)
		req    request.GetEduOrganizationalList
		err    error
	)

	if err = c.ShouldBindQuery(&req); err != nil {
		logger.Errorf("failed to bind: %v", c.Error(domain.NewBadRequest(err.Error())))
		return
	}

	logger = logger.WithFields(liblog.Fields{"request": req})
	ctx = liblog.With(ctx, logger)

	eduOrganizationEntities, count, err := s.eduOrganizationService.EduOrganizationList(
		ctx,
		domain.NewEduOrganizationFilters(
			domain.NewListFilter(
				req.SortOrder,
				domain.NewPagination(
					req.Pagination.Page,
					req.Pagination.PerPage,
				),
			),
		),
	)
	if err != nil {
		logger.Errorf("failed to get a list of educational organizations: %v", c.Error(err))
		return
	}

	c.JSON(http.StatusOK, response.NewEduOrganizationsList(
		eduOrganizationEntities,
		response.Pagination{
			Page:    req.Page,
			PerPage: req.Pagination.PerPage,
			Total:   count,
		},
	))
}

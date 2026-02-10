package handlers

import (
	"net/http"

	"bum-service/internal/controller/http/handlers/request"
	"bum-service/internal/controller/http/handlers/response"
	"bum-service/internal/domain"
	"bum-service/internal/service/director"
	"bum-service/pkg/liblog"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
)

// Director is director handler.
type Director struct {
	directorService IDirectorService
}

// NewDirector creates a new director handler.
func NewDirector(directorService IDirectorService) *Director {
	return &Director{
		directorService: directorService,
	}
}

// AddDirector adds a new director.
func (h Director) AddDirector(c *gin.Context) {
	var (
		ctx    = c.Request.Context()
		logger = liblog.Must(ctx)
		req    request.AddDirector
	)

	if err := c.ShouldBindJSON(&req); err != nil {
		logger.Errorf("failed to bind: %v", c.Error(domain.NewBadRequest(err.Error())))
		return
	}

	logger = logger.WithFields(liblog.Fields{"request": req.LogFields()})
	ctx = liblog.With(ctx, logger)

	directorEntity, err := h.directorService.AddDirector(
		ctx,
		director.AddDirectorArgs{
			UserID:   req.UserID,
			SchoolID: req.SchoolID,
			Phone:    req.Phone,
			Email:    req.Email,
		})
	if err != nil {
		logger.Errorf("failed to create director: %v", c.Error(err))
		return
	}

	c.JSON(http.StatusCreated, response.NewDirector(directorEntity))
}

// DirectorByID get director by id.
func (h Director) DirectorByID(c *gin.Context) { //nolint:revive // It's handler
	var (
		ctx        = c.Request.Context()
		logger     = liblog.Must(ctx)
		reqParam   = request.GetDirectorIDPathVar(c)
		directorID uuid.UUID
		err        error
	)

	logger = logger.WithFields(liblog.Fields{"request": liblog.Fields{"director_id": reqParam}})
	ctx = liblog.With(ctx, logger)

	if directorID, err = uuid.Parse(reqParam); err != nil {
		logger.Errorf("failed to bind: %v", c.Error(domain.NewBadRequest(err.Error())))
		return
	}

	directorEntity, err := h.directorService.DirectorByID(ctx, directorID)
	if err != nil {
		logger.Errorf("failed to get director by id: %v", c.Error(err))
		return
	}

	c.JSON(http.StatusOK, response.NewDirector(directorEntity))
}

// DirectorList endpoint for getting a list of directors.
func (h Director) DirectorList(c *gin.Context) {
	var (
		ctx    = c.Request.Context()
		logger = liblog.Must(ctx)
		req    request.DirectorList
	)

	if err := c.ShouldBindQuery(&req); err != nil {
		logger.Errorf("failed to bind: %v", c.Error(domain.NewBadRequest(err.Error())))
		return
	}

	logger = logger.WithFields(liblog.Fields{"request": req})
	ctx = liblog.With(ctx, logger)

	list, total, err := h.directorService.DirectorList(
		ctx,
		domain.NewDirectorListFilter(
			domain.NewDateFilter(req.CreatedDate.DateFrom(), req.CreatedDate.DateTill()),
			domain.NewListFilter(req.SortOrder, domain.NewPagination(req.Page, req.PerPage)),
			req.SchoolUUIDs(),
		),
	)
	if err != nil {
		logger.Errorf("failed to get director list: %v", c.Error(err))
		return
	}

	c.JSON(http.StatusOK, response.NewDirectorList(list, response.Pagination{
		Page:    req.Page,
		PerPage: req.PerPage,
		Total:   total,
	}))
}

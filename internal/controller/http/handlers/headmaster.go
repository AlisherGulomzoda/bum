//nolint:dupl // it's ok
package handlers

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"

	"bum-service/internal/controller/http/handlers/request"
	"bum-service/internal/controller/http/handlers/response"
	"bum-service/internal/domain"
	"bum-service/internal/service/headmaster"
	"bum-service/pkg/liblog"
)

// Headmaster is headmaster handler.
type Headmaster struct {
	headmasterService IHeadmasterService
}

// NewHeadmaster creates a new headmaster handler.
func NewHeadmaster(headmasterService IHeadmasterService) *Headmaster {
	return &Headmaster{
		headmasterService: headmasterService,
	}
}

// AddHeadmaster adds a new headmaster.
func (h Headmaster) AddHeadmaster(c *gin.Context) {
	var (
		ctx    = c.Request.Context()
		logger = liblog.Must(ctx)
		req    request.AddHeadmaster
	)

	if err := c.ShouldBindJSON(&req); err != nil {
		logger.Errorf("failed to bind: %v", c.Error(domain.NewBadRequest(err.Error())))
		return
	}

	logger = logger.WithFields(liblog.Fields{"request": req.LogFields()})
	ctx = liblog.With(ctx, logger)

	entity, err := h.headmasterService.AddHeadmaster(
		ctx,
		headmaster.AddHeadmasterArgs{
			UserID:   req.UserID,
			SchoolID: req.SchoolID,
			Phone:    req.Phone,
			Email:    req.Email,
		})
	if err != nil {
		logger.Errorf("failed to create a new headmaster: %v", c.Error(err))
		return
	}

	c.JSON(http.StatusCreated, response.NewHeadmaster(entity))
}

// HeadmasterByID get headmaster by id.
func (h Headmaster) HeadmasterByID(c *gin.Context) { //nolint:revive // It's handler
	var (
		ctx          = c.Request.Context()
		logger       = liblog.Must(ctx)
		reqParam     = request.GetHeadmasterIDPathVar(c)
		headmasterID uuid.UUID
		err          error
	)

	logger = logger.WithFields(liblog.Fields{"request": liblog.Fields{"headmaster_id": reqParam}})
	ctx = liblog.With(ctx, logger)

	if headmasterID, err = uuid.Parse(reqParam); err != nil {
		logger.Errorf("failed to bind: %v", c.Error(domain.NewBadRequest(err.Error())))
		return
	}

	headmasterEntity, err := h.headmasterService.HeadmasterByID(ctx, headmasterID)
	if err != nil {
		logger.Errorf("failed to get headmaster by id: %v", c.Error(err))
		return
	}

	c.JSON(http.StatusOK, response.NewHeadmaster(headmasterEntity))
}

// HeadmasterList endpoint for getting a list of headmasters.
func (h Headmaster) HeadmasterList(c *gin.Context) {
	var (
		ctx    = c.Request.Context()
		logger = liblog.Must(ctx)
		req    request.HeadmasterList
	)

	if err := c.ShouldBindQuery(&req); err != nil {
		logger.Errorf("failed to bind: %v", c.Error(domain.NewBadRequest(err.Error())))
		return
	}

	logger = logger.WithFields(liblog.Fields{"request": req})
	ctx = liblog.With(ctx, logger)

	list, total, err := h.headmasterService.HeadmasterList(
		ctx,
		domain.NewHeadmasterListFilter(
			domain.NewDateFilter(req.CreatedDate.DateFrom(), req.CreatedDate.DateTill()),
			domain.NewListFilter(req.SortOrder, domain.NewPagination(req.Page, req.PerPage)),
			req.SchoolUUIDs(),
		),
	)
	if err != nil {
		logger.Errorf("failed to get headmaster list: %v", c.Error(err))
		return
	}

	c.JSON(http.StatusOK, response.NewHeadmasterList(list, response.Pagination{
		Page:    req.Page,
		PerPage: req.PerPage,
		Total:   total,
	}))
}

package handlers

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"

	"bum-service/internal/controller/http/handlers/request"
	"bum-service/internal/controller/http/handlers/response"
	"bum-service/internal/domain"
	"bum-service/internal/service/owner"
	"bum-service/pkg/liblog"
)

// Owner is owner handler.
type Owner struct {
	ownerService IOwnerService
}

// NewOwner creates a new owner handler.
func NewOwner(ownerService IOwnerService) *Owner {
	return &Owner{ownerService: ownerService}
}

// AddOwner endpoint for adding a new owner.
func (o *Owner) AddOwner(c *gin.Context) {
	var (
		ctx    = c.Request.Context()
		logger = liblog.Must(ctx)
		req    request.AddOwner
	)

	if err := c.ShouldBindJSON(&req); err != nil {
		logger.Errorf("failed to bind: %v", c.Error(domain.NewBadRequest(err.Error())))
		return
	}

	logger = logger.WithFields(liblog.Fields{"request": req.LogFields()})
	ctx = liblog.With(ctx, logger)

	ownerEntity, err := o.ownerService.AddOwner(
		ctx,
		owner.AddOwnerArgs{
			UserID:         req.UserID,
			OrganizationID: req.OrganizationID,
			Phone:          req.Phone,
			Email:          req.Email,
		})
	if err != nil {
		logger.Errorf("failed to add owner: %v", c.Error(err))
		return
	}

	c.JSON(http.StatusCreated, response.NewOwner(ownerEntity))
}

// OwnerByID get owner by id.
func (o *Owner) OwnerByID(c *gin.Context) { //nolint:revive // It's handler
	var (
		ctx      = c.Request.Context()
		logger   = liblog.Must(ctx)
		reqParam = request.GetOwnerIDPathVar(c)
		ownerID  uuid.UUID
		err      error
	)

	logger = logger.WithFields(liblog.Fields{"request": liblog.Fields{"owner_id": reqParam}})
	ctx = liblog.With(ctx, logger)

	if ownerID, err = uuid.Parse(reqParam); err != nil {
		logger.Errorf("failed to bind: %v", c.Error(domain.NewBadRequest(err.Error())))
		return
	}

	ownerEntity, err := o.ownerService.OwnerByID(ctx, ownerID)
	if err != nil {
		logger.Errorf("failed to get owner by id: %v", c.Error(err))
		return
	}

	c.JSON(http.StatusOK, response.NewOwner(ownerEntity))
}

// OwnerByUserIDAndSchoolID get owner by user_id and school id.
func (o *Owner) OwnerByUserIDAndSchoolID(c *gin.Context) { //nolint:revive // It's handler
	var (
		ctx         = c.Request.Context()
		logger      = liblog.Must(ctx)
		userID      = MustGetUserID(c)
		schoolIDStr = request.GetSchoolIDHeader(c)
		schoolID    uuid.UUID
		err         error
	)

	logger = logger.WithFields(liblog.Fields{"request": liblog.Fields{"user_id": userID, "school_id": schoolID}})
	ctx = liblog.With(ctx, logger)

	if schoolID, err = uuid.Parse(schoolIDStr); err != nil {
		logger.Errorf("failed to bind: %v", c.Error(domain.NewBadRequest(err.Error())))
		return
	}

	ownerEntity, err := o.ownerService.OwnerByUserIDAndSchoolID(ctx, schoolID, userID)
	if err != nil {
		logger.Errorf("failed to get owner by id: %v", c.Error(err))
		return
	}

	c.JSON(http.StatusOK, response.NewOwner(ownerEntity))
}

// OwnerList endpoint for getting a list of owners.
func (o *Owner) OwnerList(c *gin.Context) {
	var (
		ctx    = c.Request.Context()
		logger = liblog.Must(ctx)
		req    request.OwnerList
	)

	if err := c.ShouldBindQuery(&req); err != nil {
		logger.Errorf("failed to bind: %v", c.Error(domain.NewBadRequest(err.Error())))
		return
	}

	logger = logger.WithFields(liblog.Fields{"request": req})
	ctx = liblog.With(ctx, logger)

	list, total, err := o.ownerService.OwnerList(
		ctx,
		domain.NewOwnerListFilter(
			domain.NewDateFilter(req.CreatedDate.DateFrom(), req.CreatedDate.DateTill()),
			domain.NewListFilter(req.SortOrder, domain.NewPagination(req.Page, req.PerPage)),
			req.OrganizationUUIDs(),
		),
	)
	if err != nil {
		logger.Errorf("failed to get owners list: %v", c.Error(err))
		return
	}

	c.JSON(http.StatusOK, response.NewOwnerList(list, response.Pagination{
		Page:    req.Page,
		PerPage: req.PerPage,
		Total:   total,
	}))
}

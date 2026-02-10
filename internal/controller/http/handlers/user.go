package handlers

import (
	"net/http"

	"github.com/google/uuid"

	"github.com/gin-gonic/gin"

	"bum-service/internal/controller/http/handlers/request"
	"bum-service/internal/controller/http/handlers/response"
	"bum-service/internal/domain"
	"bum-service/internal/service/user"
	"bum-service/pkg/liblog"
)

// User is a handler for user.
type User struct {
	userService IUserService
}

// NewUser creates a new User handler.
func NewUser(userService IUserService) *User {
	return &User{
		userService: userService,
	}
}

// UserFullInfoByID returns full user information.
func (u User) UserFullInfoByID(c *gin.Context) {
	var (
		ctx    = c.Request.Context()
		logger = liblog.Must(ctx)

		reqParam = request.GetUserIDPathVar(c)
		userID   uuid.UUID

		err error
	)

	logger = logger.WithFields(liblog.Fields{"user_id": userID})

	if userID, err = uuid.Parse(reqParam); err != nil {
		logger.Errorf("failed to bind: %v", c.Error(domain.NewBadRequest(err.Error())))
		return
	}

	userFullInfo, err := u.userService.UserFullInfoByID(ctx, userID)
	if err != nil {
		logger.Errorf("failed to get user full info by id: %v", c.Error(err))
		return
	}

	c.JSON(http.StatusOK, response.NewUser(userFullInfo))
}

// UserFullInfoFromToken returns full user information.
func (u User) UserFullInfoFromToken(c *gin.Context) {
	var (
		ctx    = c.Request.Context()
		userID = MustGetUserID(c)
		logger = liblog.Must(ctx)
	)

	logger = logger.WithFields(liblog.Fields{"user_id": userID})

	userFullInfo, err := u.userService.UserFullInfoByID(ctx, userID)
	if err != nil {
		logger.Errorf("failed to get user full info by id: %v", c.Error(err))
		return
	}

	c.JSON(http.StatusOK, response.NewUser(userFullInfo))
}

// AddUser add user handler.
func (u User) AddUser(c *gin.Context) {
	var (
		ctx    = c.Request.Context()
		logger = liblog.Must(ctx)
		req    request.AddUser
	)

	if err := c.ShouldBindJSON(&req); err != nil {
		logger.Errorf("failed to bind: %v", c.Error(domain.NewBadRequest(err.Error())))
		return
	}

	logger = logger.WithFields(liblog.Fields{"request": req.LogInfo()})
	ctx = liblog.With(ctx, logger)

	userEntity, err := u.userService.AddUser(
		ctx,
		user.AddUserArgs{
			FirstName:  req.FirstName,
			LastName:   req.LastName,
			MiddleName: req.MiddleName,
			Gender:     req.Gender,
			Phone:      req.Phone,
			Email:      req.Email,
			Password:   req.Password,
		},
	)
	if err != nil {
		logger.Errorf("failed to add user: %v", c.Error(err))
		return
	}

	c.JSON(http.StatusCreated, response.NewUser(userEntity))
}

// UserList returns a list of Users.
func (u User) UserList(c *gin.Context) {
	var (
		ctx    = c.Request.Context()
		logger = liblog.Must(ctx)
		req    request.UserList
	)

	if err := c.ShouldBindQuery(&req); err != nil {
		logger.Errorf("failed to bind: %v", c.Error(domain.NewBadRequest(err.Error())))
		return
	}

	logger = logger.WithFields(liblog.Fields{"request": req})

	filters, err := domain.NewUserListFilter(
		req.OrganizationUUIDs(),
		req.SchoolUUIDs(),

		req.UserRoles(),

		req.Emails,
		req.Phones,

		domain.NewListFilter(req.SortOrder, domain.NewPagination(req.Page, req.PerPage)),
	)
	if err != nil {
		logger.Errorf("failed to create user filter: %v", c.Error(err))
		return
	}

	userList, count, err := u.userService.UserList(ctx, filters)
	if err != nil {
		logger.Errorf("failed to get user list: %v", c.Error(err))
		return
	}

	c.JSON(http.StatusOK, response.NewUserList(userList,
		response.Pagination{
			Page:    req.Page,
			PerPage: req.PerPage,
			Total:   count,
		},
	))
}

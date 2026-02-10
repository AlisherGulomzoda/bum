package handlers

import (
	"net/http"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/golang-jwt/jwt/v5"

	"bum-service/internal/controller/http/handlers/request"
	"bum-service/internal/controller/http/handlers/response"
	"bum-service/internal/domain"
	"bum-service/pkg/liblog"
)

// Auth is auth handler.
type Auth struct {
	authService IAuthService
	userService IUserService

	jwtSecret       []byte
	accessTokenExp  time.Duration
	refreshTokenExp time.Duration
}

// NewAuth creates a new auth handler.
func NewAuth(
	authService IAuthService,
	userService IUserService,

	jwtSecret string,
	accessTokenExp time.Duration,
	refreshTokenExp time.Duration,
) *Auth {
	return &Auth{
		authService: authService,
		userService: userService,

		jwtSecret:       []byte(jwtSecret),
		accessTokenExp:  accessTokenExp,
		refreshTokenExp: refreshTokenExp,
	}
}

// UserClaims is user auth claims.
type UserClaims struct {
	UserID string `json:"user_id"`
	jwt.RegisteredClaims
}

// LoginByEmail creates a new user token by email.
func (a Auth) LoginByEmail(c *gin.Context) {
	var (
		ctx    = c.Request.Context()
		logger = liblog.Must(ctx)
		req    request.LoginByEmail
	)

	if err := c.ShouldBindJSON(&req); err != nil {
		logger.Errorf("failed to bind: %v", c.Error(domain.NewBadRequest(err.Error())))
		return
	}

	userID, err := a.authService.GetUserIDByEmailAndPassword(ctx, req.Email, req.Password)
	if err != nil {
		logger.Errorf("failed to create token: %v", c.Error(err))
		return
	}

	// If password is correct then create a new access and refresh token

	// Create a new token object, specifying signing method and the claims
	// you would like it to contain.
	accessToken, err := jwt.NewWithClaims(jwt.SigningMethodHS256,
		UserClaims{
			UserID: userID.String(),
			RegisteredClaims: jwt.RegisteredClaims{
				ExpiresAt: jwt.NewNumericDate(time.Now().Add(a.accessTokenExp)),
			},
		},
		// Sign and get the complete encoded token as a string using the secret
	).SignedString(a.jwtSecret)
	if err != nil {
		logger.Errorf("failed to create access token: %v", c.Error(err))
		return
	}

	// Create a new refresh token object, specifying signing method and the claims
	// you would like it to contain.
	refreshToken, err := jwt.NewWithClaims(jwt.SigningMethodHS256,
		UserClaims{
			UserID: userID.String(),
			RegisteredClaims: jwt.RegisteredClaims{
				ExpiresAt: jwt.NewNumericDate(time.Now().Add(a.refreshTokenExp)),
			},
		},
		// Sign and get the complete encoded token as a string using the secret
	).SignedString(a.jwtSecret)
	if err != nil {
		logger.Errorf("failed to create refresh token: %v", c.Error(err))
		return
	}

	c.JSON(http.StatusOK, response.NewUserToken(accessToken, refreshToken))
}

// RefreshToken updates user tokens.
func (a Auth) RefreshToken(c *gin.Context) {
	var (
		ctx    = c.Request.Context()
		logger = liblog.Must(ctx)
		userID = MustGetUserID(c)
	)

	_, err := a.userService.UserByID(ctx, userID)
	if err != nil {
		logger.Errorf("failed to create token: %v", c.Error(err))
		return
	}

	// If password is correct then create a new access and refresh token

	// Create a new token object, specifying signing method and the claims
	// you would like it to contain.
	accessToken, err := jwt.NewWithClaims(jwt.SigningMethodHS256,
		UserClaims{
			UserID: userID.String(),
			RegisteredClaims: jwt.RegisteredClaims{
				ExpiresAt: jwt.NewNumericDate(time.Now().Add(a.accessTokenExp)),
			},
		},
		// Sign and get the complete encoded token as a string using the secret
	).SignedString(a.jwtSecret)
	if err != nil {
		logger.Errorf("failed to create access token: %v", c.Error(err))
		return
	}

	// Create a new refresh token object, specifying signing method and the claims
	// you would like it to contain.
	refreshToken, err := jwt.NewWithClaims(jwt.SigningMethodHS256,
		UserClaims{
			UserID: userID.String(),
			RegisteredClaims: jwt.RegisteredClaims{
				ExpiresAt: jwt.NewNumericDate(time.Now().Add(a.refreshTokenExp)),
			},
		},
		// Sign and get the complete encoded token as a string using the secret
	).SignedString(a.jwtSecret)
	if err != nil {
		logger.Errorf("failed to create refresh token: %v", c.Error(err))
		return
	}

	c.JSON(http.StatusOK, response.NewUserToken(accessToken, refreshToken))
}

package handlers

import (
	"bytes"
	"context"
	"errors"
	"fmt"
	"net"
	"net/http/httputil"
	"os"
	"runtime"
	"strings"

	"github.com/gin-gonic/gin"
	"github.com/golang-jwt/jwt/v5"
	"github.com/google/uuid"

	"bum-service/internal/domain"
	"bum-service/pkg/liberror"
	"bum-service/pkg/liblog"
)

// ContextKey is context key type.
type ContextKey string

const (
	userIDContextKey ContextKey = "X-User-Id"

	tokenHeader     = "Authorization"
	tokenHeaderType = "Bearer"
)

// LoggingEndpointMiddleware middleware for logging endpoint calls and putting logger to request context.
func LoggingEndpointMiddleware(logger liblog.Logger) gin.HandlerFunc {
	return func(c *gin.Context) {
		request := c.Request

		if request.Header.Get(liblog.ReqIDHeader) == "" {
			request.Header.Set(liblog.ReqIDHeader, uuid.New().String())
		}

		ctx := liblog.WithRequest(
			liblog.With(request.Context(), logger),
			request,
		)

		c.Request = request.WithContext(ctx)

		c.Next()
	}
}

// ErrorHandlingMiddleware middleware for handling errors.
func ErrorHandlingMiddleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		c.Next()

		if c.Errors.Last() == nil {
			return
		}

		err := c.Errors.Last().Err

		logger := liblog.Must(c.Request.Context())
		if errorEncoder := liberror.ErrorEncoder(c.Request.Context(), err, c.Writer); errorEncoder != nil {
			logger.Errorf("failed to handle error: %v\n", errorEncoder)
		}

		c.Abort()
	}
}

// RecoverMiddleware middleware for handling panics.
//
//nolint:gocognit, nestif, wsl, errcheck, gocritic // its built in gin method that was copied
func RecoverMiddleware(logger liblog.Logger) gin.HandlerFunc {
	return func(c *gin.Context) {
		defer func() {
			if err := recover(); err != nil {
				// Check for a broken connection, as it is not really a
				// condition that warrants a panic stack trace.
				var brokenPipe bool
				if ne, ok := err.(*net.OpError); ok {
					var se *os.SyscallError
					if errors.As(ne, &se) {
						seStr := strings.ToLower(se.Error())
						if strings.Contains(seStr, "broken pipe") ||
							strings.Contains(seStr, "connection reset by peer") {
							brokenPipe = true

							logger.Errorf("connection is dead: %+v", se)
						}
					}
				}

				stack := stackFrame()
				httpRequest, _ := httputil.DumpRequest(c.Request, false)
				headers := strings.Split(string(httpRequest), "\r\n")
				for idx, header := range headers {
					current := strings.Split(header, ":")
					if current[0] == "Authorization" {
						headers[idx] = current[0] + ": *"
					}
				}
				headersToStr := strings.Join(headers, "\r\n")
				if brokenPipe {
					logger.Errorf("%v\n%s\n", err, headersToStr)
				} else if gin.IsDebugging() {
					logger.Errorf("[Recovery] %s panic recovered:\n%s\n%s\n", headersToStr, err, stack)
				} else {
					logger.Errorf("[Recovery] %s panic recovered:\n%s\n", err, stack)
				}

				_ = c.Error(domain.ErrInternalServerError)

				c.Abort()
			}
		}()
		c.Next()
	}
}

const (
	skipFrame = 3
)

// stackFrame returns a nicely formatted stack frame, skipping skip frames.
//
//nolint:wsl, gocritic, revive // its built in gin method that was copied
func stackFrame() []byte {
	buf := new(bytes.Buffer) // the returned data
	// As we loop, we open files and read them. These variables record the currently
	// loaded file.
	var lines [][]byte
	var lastFile string
	for i := skipFrame; ; i++ { // Skip the expected number of frames
		pc, file, line, ok := runtime.Caller(i)
		if !ok {
			break
		}
		// Print this much at least.  If we can't find the source, it won't show.
		fmt.Fprintf(buf, "%s:%d (0x%x)\n", file, line, pc)
		if file != lastFile {
			data, err := os.ReadFile(file)
			if err != nil {
				continue
			}
			lines = bytes.Split(data, []byte{'\n'})
			lastFile = file
		}
		fmt.Fprintf(buf, "\t%s: %s\n", function(pc), source(lines, line))
	}

	return buf.Bytes()
}

// source returns a space-trimmed slice of the n'th line.
//
//nolint:wsl, gocritic, nlreturn // its built in gin method that was copied
func source(lines [][]byte, n int) []byte {
	n-- // in stack trace, lines are 1-indexed but our array is 0-indexed
	if n < 0 || n >= len(lines) {
		return dunno
	}
	return bytes.TrimSpace(lines[n])
}

// function returns, if possible, the name of the function containing the PC.
//
//nolint:wsl, gocritic // its built in gin method that was copied
func function(pc uintptr) []byte {
	fn := runtime.FuncForPC(pc)
	if fn == nil {
		return dunno
	}
	name := []byte(fn.Name())
	// The name includes the path name to the package, which is unnecessary
	// since the file name is already included.  Plus, it has center dots.
	// That is, we see
	//	runtime/debug.*T·ptrmethod
	// and want
	//	*T.ptrmethod
	// Also the package path might contain dot (e.g. code.google.com/...),
	// so first eliminate the path prefix
	if lastSlash := bytes.LastIndex(name, slash); lastSlash >= 0 {
		name = name[lastSlash+1:]
	}
	if period := bytes.Index(name, dot); period >= 0 {
		name = name[period+1:]
	}
	name = bytes.Replace(name, centerDot, dot, -1)

	return name
}

//nolint:gochecknoglobals // these are built in gin variables that was copied
var (
	dunno     = []byte("???")
	centerDot = []byte("·")
	dot       = []byte(".")
	slash     = []byte("/")
)

// AuthMiddleware is middleware for checking auth token.
func (a Auth) AuthMiddleware(c *gin.Context) {
	var (
		authHeader = c.GetHeader(tokenHeader)

		logger = liblog.Must(c.Request.Context())
	)

	// Authorization header is missing in HTTP request
	if authHeader == "" {
		logger.Errorf("empty auth header: %v\n", c.Error(domain.ErrEmptyAuthHeader).Error())
		c.Abort()

		return
	}

	authTokens := strings.Split(authHeader, " ")

	// The value of authorization header is invalid
	// It should start with "Bearer ", then the token value
	if len(authTokens) != 2 || authTokens[0] != tokenHeaderType {
		logger.Errorf("invalid auth header: %v\n", c.Error(domain.ErrEmptyAuthHeader).Error())
		c.Abort()

		return
	}

	var (
		tokenString = authTokens[1]
		userClaims  *UserClaims
	)

	token, err := jwt.ParseWithClaims(tokenString, &UserClaims{}, func(_ *jwt.Token) (any, error) {
		return a.jwtSecret, nil
	})
	if err != nil {
		switch {
		case errors.Is(err, jwt.ErrTokenExpired) || errors.Is(err, jwt.ErrTokenNotValidYet):
			logger.Errorf("token is expired: %v: %v\n", err, c.Error(domain.ErrTokenIsExpired).Error())
			c.Abort()

			return
		default:
			logger.Errorf("invalid token: %v: %v\n ", err, c.Error(domain.ErrInvalidToken).Error())
			c.Abort()

			return
		}
	}

	userClaims, ok := token.Claims.(*UserClaims)
	if userClaims == nil || !ok {
		logger.Errorf("invalid token: %v\n", c.Error(domain.ErrInvalidToken).Error())
		c.Abort()

		return
	}

	userUUID, err := uuid.Parse(userClaims.UserID)
	if err != nil {
		logger.Errorf("invalid token: %v: %v\n", err, c.Error(domain.ErrInvalidToken).Error())
		c.Abort()

		return
	}

	c.Request = c.Request.WithContext(context.WithValue(c.Request.Context(), userIDContextKey, userUUID))
}

// MustGetUserID gets userID from context.
//
//nolint:forcetypeassert // it's must method so it's ok here.
func MustGetUserID(c *gin.Context) uuid.UUID {
	return c.Request.Context().Value(userIDContextKey).(uuid.UUID)
}

//nolint:funlen,revive // it's test functions.
package handlers_test

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"net/http"
	"net/http/httptest"
	"testing"
	"time"

	"github.com/brianvoe/gofakeit/v7"
	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"github.com/stretchr/testify/suite"
	"go.uber.org/mock/gomock"

	"bum-service/internal/controller/http/handlers"
	mockhandlers "bum-service/internal/controller/http/handlers/mocks"
	"bum-service/internal/controller/http/handlers/request"
	"bum-service/internal/controller/http/handlers/response"
	"bum-service/internal/domain"
	eduorganization "bum-service/internal/service/edu-organization"
	"bum-service/pkg/liberror"
	"bum-service/pkg/liblog"
	"bum-service/pkg/utils"
)

// In your Go test file.
func TestEduOrganizationSuite(t *testing.T) {
	t.Parallel()

	suite.Run(t, new(EduOrganizationSuite))
}

// EduOrganizationSuite is suite for testing educational organization handlers.
type EduOrganizationSuite struct {
	suite.Suite

	logger liblog.Logger

	timeFunc                 func() time.Time
	expectedTimeFromTimeFunc time.Time
}

func (s *EduOrganizationSuite) SetupSuite() {
	// setting out logger.
	s.logger = liblog.NewDummyLogger()

	// setting our time func to get expected time.
	s.expectedTimeFromTimeFunc = time.Now()
	s.timeFunc = func() time.Time { return s.expectedTimeFromTimeFunc }

	// setting rest mode for router.
	gin.SetMode(gin.TestMode)
}

func (*EduOrganizationSuite) SetupTest() {}

func (*EduOrganizationSuite) TearDownTest() {}

func (*EduOrganizationSuite) TearDownSuite() {}

func (s *EduOrganizationSuite) TestEduOrganization_EduOrganizationByIDPositiveCases() {
	var (
		asserting = s.Assert()
		requiring = s.Require()

		ctx = context.Background()

		// setting controller for our mocks.
		ctrl = gomock.NewController(s.Suite.T())

		// initializing mock organization service.
		eduOrganizationSvc = mockhandlers.NewMockIEduOrganizationService(ctrl)

		// setting our expected organization domain for positive cases.
		existedEduOrganizationDomain = domain.EduOrganization{
			ID:        uuid.New(),
			Name:      gofakeit.Company(),
			CreatedAt: s.expectedTimeFromTimeFunc,
			UpdatedAt: s.expectedTimeFromTimeFunc,
			DeletedAt: nil,
		}

		// setting id of existed organization.
		requestID = existedEduOrganizationDomain.ID

		// initialization of educational organization handlers.
		eduOrganizationHandlers = handlers.NewEduOrganization(eduOrganizationSvc)

		// initialization of router and registering handlers.
		router = NewRouter()

		want = struct {
			Code     int
			Response response.EduOrganization
		}{
			Code: http.StatusOK,
			Response: response.EduOrganization{
				ID:   existedEduOrganizationDomain.ID,
				Name: existedEduOrganizationDomain.Name,
				Logo: existedEduOrganizationDomain.Logo,

				CreatedAt: utils.RFC3339Time(existedEduOrganizationDomain.CreatedAt),
				UpdatedAt: utils.RFC3339Time(existedEduOrganizationDomain.UpdatedAt),
				DeletedAt: (*utils.RFC3339Time)(existedEduOrganizationDomain.DeletedAt),
			},
		}
	)

	// setting expectation for existed organization case.
	eduOrganizationSvc.EXPECT().
		EduOrganizationByID(gomock.Any(), existedEduOrganizationDomain.ID).
		Return(existedEduOrganizationDomain, nil).AnyTimes()

	// registering our handler.
	router.GET("/edu-organizations/:edu_organization_id", eduOrganizationHandlers.EduOrganizationByID)

	// creating our request.
	url := fmt.Sprintf("/edu-organizations/%s", requestID)
	r, err := http.NewRequestWithContext(ctx, http.MethodGet, url, http.NoBody)
	requiring.NoError(err)

	// creating our test response writer.
	recorder := httptest.NewRecorder()

	// serving the request.
	router.ServeHTTP(recorder, r)
	requiring.Equal(want.Code, recorder.Code)

	// passing the response.
	var resp response.EduOrganization

	s.Require().NoError(json.Unmarshal(recorder.Body.Bytes(), &resp))

	// checking the response.
	asserting.Equal(want.Response.ID, resp.ID)
	asserting.Equal(want.Response.Name, resp.Name)
	asserting.Equal(want.Response.Logo, resp.Logo)
	asserting.Equal(want.Response.CreatedAt.String(), resp.CreatedAt.String())
	asserting.Equal(want.Response.UpdatedAt.String(), resp.UpdatedAt.String())
	requiring.Equal(want.Response.DeletedAt == nil, resp.DeletedAt == nil)

	if want.Response.DeletedAt != nil {
		asserting.Equal(want.Response.DeletedAt.String(), resp.DeletedAt.String())
	}
}

func (s *EduOrganizationSuite) TestEduOrganization_EduOrganizationByIDNegativeCases() {
	s.T().Parallel()

	var (
		// setting controller for our mocks.
		ctrl = gomock.NewController(s.Suite.T())

		// initializing mock organization service.
		eduOrganizationSvc = mockhandlers.NewMockIEduOrganizationService(ctrl)

		// initializing id for internal error.
		internalErrID = uuid.New()

		// initialization of educational organization handlers.
		eduOrganizationHandlers = handlers.NewEduOrganization(eduOrganizationSvc)

		// initialization of router and registering handlers.
		router = NewRouter()
	)

	// setting expectation for not found case.
	eduOrganizationSvc.EXPECT().
		EduOrganizationByID(gomock.Any(), gomock.Not(gomock.AnyOf(internalErrID))).
		Return(domain.EduOrganization{}, domain.ErrNotFound).AnyTimes()

	// setting expectation for internal error case.
	eduOrganizationSvc.EXPECT().
		EduOrganizationByID(gomock.Any(), internalErrID).
		Return(domain.EduOrganization{}, domain.ErrInternalServerError).AnyTimes()

	// registering our handler.
	router.GET("/edu-organizations/:edu_organization_id", eduOrganizationHandlers.EduOrganizationByID)

	type args struct {
		id string
	}

	type want struct {
		Code     int
		Response *liberror.Error
	}

	tests := []struct {
		name string
		req  args
		want want
	}{
		{
			name: "a non existed identifier",
			req: args{
				id: uuid.New().String(),
			},
			want: want{
				Code:     http.StatusNotFound,
				Response: domain.ErrNotFound,
			},
		},
		{
			name: "an invalid id",
			req: args{
				id: gofakeit.Name(),
			},
			want: want{
				Code:     http.StatusBadRequest,
				Response: domain.ErrBadRequest,
			},
		},
		{
			name: "internal error",
			req: args{
				id: internalErrID.String(),
			},
			want: want{
				Code:     http.StatusInternalServerError,
				Response: domain.ErrInternalServerError,
			},
		},
	}
	for _, tt := range tests {
		tt := tt
		s.T().Run(tt.name, func(t *testing.T) {
			t.Parallel()

			var (
				asserting = assert.New(t)
				requiring = require.New(t)

				ctx = context.Background()
			)

			// creating our request.
			url := fmt.Sprintf("/edu-organizations/%s", tt.req.id)
			r, err := http.NewRequestWithContext(ctx, http.MethodGet, url, http.NoBody)
			requiring.NoError(err)

			// creating our test response writer.
			recorder := httptest.NewRecorder()

			// serving the request.
			router.ServeHTTP(recorder, r)
			requiring.Equal(tt.want.Code, recorder.Code)

			// passing the response.
			var resp liberror.Error

			requiring.NoError(json.Unmarshal(recorder.Body.Bytes(), &resp))

			// checking the response.
			asserting.Equal(tt.want.Response.Code, resp.Code)
		})
	}
}

func (s *EduOrganizationSuite) TestEduOrganization_CreateEduOrganizationPositiveCases() {
	var (
		asserting = s.Assert()
		requiring = s.Require()

		ctx = context.Background()

		// setting controller for our mocks.
		ctrl = gomock.NewController(s.Suite.T())

		// initializing mock organization service.
		eduOrganizationSvc = mockhandlers.NewMockIEduOrganizationService(ctrl)

		// initialization of educational organization handlers.
		eduOrganizationHandlers = handlers.NewEduOrganization(eduOrganizationSvc)

		// initialization of router and registering handlers.
		router = NewRouter()

		logo    = gofakeit.URL()
		reqArgs = request.CreateEduOrganizational{
			Name: gofakeit.Company(),
			Logo: &logo,
		}

		reqBytes = []byte(fmt.Sprintf(`
			{
				"name": "%s",
				"logo": "%s"
			}
		`, reqArgs.Name, *reqArgs.Logo))

		// setting our expected organization domain response for positive cases.
		expectedEduOrganizationDomain = domain.EduOrganization{
			ID:        uuid.New(),
			Name:      reqArgs.Name,
			Logo:      reqArgs.Logo,
			CreatedAt: s.expectedTimeFromTimeFunc,
			UpdatedAt: s.expectedTimeFromTimeFunc,
			DeletedAt: nil,
		}

		want = struct {
			Code     int
			Response response.EduOrganization
		}{
			Code: http.StatusCreated,
			Response: response.EduOrganization{
				ID:   expectedEduOrganizationDomain.ID,
				Name: expectedEduOrganizationDomain.Name,
				Logo: expectedEduOrganizationDomain.Logo,

				CreatedAt: utils.RFC3339Time(expectedEduOrganizationDomain.CreatedAt),
				UpdatedAt: utils.RFC3339Time(expectedEduOrganizationDomain.UpdatedAt),
				DeletedAt: (*utils.RFC3339Time)(expectedEduOrganizationDomain.DeletedAt),
			},
		}
	)

	// setting expectation for existed organization case.
	eduOrganizationSvc.EXPECT().
		CreateEduOrganization(gomock.Any(), eduorganization.CreateEduOrganizationArgs{
			Name: expectedEduOrganizationDomain.Name,
			Logo: expectedEduOrganizationDomain.Logo,
		}).
		Return(expectedEduOrganizationDomain, nil).AnyTimes()

	// registering our handler.
	router.POST("/edu-organizations", eduOrganizationHandlers.CreateEduOrganization)

	// creating our request.
	r, err := http.NewRequestWithContext(ctx, http.MethodPost, "/edu-organizations", bytes.NewReader(reqBytes))
	requiring.NoError(err)

	// creating our test response writer.
	recorder := httptest.NewRecorder()

	// serving the request.
	router.ServeHTTP(recorder, r)
	asserting.Equal(want.Code, recorder.Code)

	// passing the response.
	var resp response.EduOrganization

	s.Require().NoError(json.Unmarshal(recorder.Body.Bytes(), &resp))

	// checking the response.
	asserting.Equal(want.Response.ID, resp.ID)
	asserting.Equal(want.Response.Name, resp.Name)
	asserting.Equal(want.Response.Logo, resp.Logo)
	asserting.Equal(want.Response.CreatedAt.String(), resp.CreatedAt.String())
	asserting.Equal(want.Response.UpdatedAt.String(), resp.UpdatedAt.String())
	requiring.Equal(want.Response.DeletedAt == nil, resp.DeletedAt == nil)

	if want.Response.DeletedAt != nil {
		asserting.Equal(want.Response.DeletedAt.String(), resp.DeletedAt.String())
	}
}

func (s *EduOrganizationSuite) TestEduOrganization_CreateEduOrganizationNegativeCases() {
	s.T().Parallel()

	var (
		// setting controller for our mocks.
		ctrl = gomock.NewController(s.Suite.T())

		// initializing mock organization service.
		eduOrganizationSvc = mockhandlers.NewMockIEduOrganizationService(ctrl)

		// initializing domain for already existing domain
		logo             = gofakeit.URL()
		alreadyExistsArg = eduorganization.CreateEduOrganizationArgs{
			Name: gofakeit.Company(),
			Logo: &logo,
		}

		// initializing domain for internal error.
		internalErrArg = eduorganization.CreateEduOrganizationArgs{
			Name: gofakeit.Company(),
			Logo: &logo,
		}

		// initialization of educational organization handlers.
		eduOrganizationHandlers = handlers.NewEduOrganization(eduOrganizationSvc)

		// initialization of router and registering handlers.
		router = NewRouter()
	)

	// setting expectation for already exist organization error.
	eduOrganizationSvc.EXPECT().
		CreateEduOrganization(gomock.Any(), alreadyExistsArg).
		Return(domain.EduOrganization{}, domain.ErrEduOrganizationAlreadyExists).AnyTimes()

	// setting expectation for internal error case.
	eduOrganizationSvc.EXPECT().
		CreateEduOrganization(gomock.Any(), internalErrArg).
		Return(domain.EduOrganization{}, domain.ErrInternalServerError).AnyTimes()

	// registering our handler.
	router.POST("/edu-organizations", eduOrganizationHandlers.CreateEduOrganization)

	type args struct {
		reqBytes []byte
	}

	type want struct {
		Code     int
		Response *liberror.Error
	}

	tests := []struct {
		name string
		req  args
		want want
	}{
		{
			name: "bad request",
			req: args{
				reqBytes: []byte(`
			{
				"name": "aa",
				"logo": "bb"
		`),
			},
			want: want{
				Code:     http.StatusBadRequest,
				Response: domain.ErrBadRequest,
			},
		},
		{
			name: "already existed organization",
			req: args{
				reqBytes: []byte(fmt.Sprintf(`
				{
					"name": "%s",
					"logo": "%s"
				}`,
					alreadyExistsArg.Name, *alreadyExistsArg.Logo)),
			},
			want: want{
				Code:     http.StatusConflict,
				Response: domain.ErrEduOrganizationAlreadyExists,
			},
		},
		{
			name: "internal server error",
			req: args{
				reqBytes: []byte(fmt.Sprintf(`
				{
					"name": "%s",
					"logo": "%s"
				}`,
					internalErrArg.Name, *internalErrArg.Logo)),
			},
			want: want{
				Code:     http.StatusInternalServerError,
				Response: domain.ErrInternalServerError,
			},
		},
	}
	for _, tt := range tests {
		tt := tt
		s.T().Run(tt.name, func(t *testing.T) {
			t.Parallel()

			var (
				assertion = assert.New(t)
				requiring = require.New(t)

				ctx = context.Background()
			)

			// creating our request.
			r, err := http.NewRequestWithContext(ctx, http.MethodPost, "/edu-organizations", bytes.NewReader(tt.req.reqBytes))
			requiring.NoError(err)

			// creating our test response writer.
			recorder := httptest.NewRecorder()

			// serving the request.
			router.ServeHTTP(recorder, r)
			assertion.Equal(tt.want.Code, recorder.Code)

			// passing the response.
			var resp liberror.Error

			requiring.NoError(json.Unmarshal(recorder.Body.Bytes(), &resp))

			// checking the response.
			assertion.Equal(tt.want.Response.Code, resp.Code)
		})
	}
}

func (s *EduOrganizationSuite) TestEduOrganization_EduOrganizationListPositiveCases() {
	s.T().Parallel()

	var (
		// setting controller for our mocks.
		ctrl = gomock.NewController(s.Suite.T())

		// initializing mock organization service.
		eduOrganizationSvc = mockhandlers.NewMockIEduOrganizationService(ctrl)

		// initialization of educational organization handlers.
		eduOrganizationHandlers = handlers.NewEduOrganization(eduOrganizationSvc)

		// initialization of router and registering handlers.
		router = NewRouter()
	)

	// registering our handler.
	router.GET("/edu-organizations", eduOrganizationHandlers.EduOrganizationList)

	type args struct {
		reqBytes func(t *testing.T) map[string]string
	}

	type want struct {
		Code     int
		Response string
	}

	tests := []struct {
		name string
		req  args
		want want
	}{
		{
			name: "default pagination",
			req: args{
				reqBytes: func(t *testing.T) map[string]string { //nolint:thelper // it's ok without helper.
					var (
						id              = uuid.MustParse("351a12a6-7cdb-47e1-8105-b68d0a7b975d")
						totalCount      = 1
						orgName         = "test organization"
						expectedTimeStr = "2024-05-01T22:56:45+05:00"
					)

					expectedTime, err := time.Parse(time.RFC3339, expectedTimeStr)
					require.NoError(t, err)

					// setting expectation for default pagination case.
					eduOrganizationSvc.EXPECT().
						EduOrganizationList(gomock.Any(), domain.EduOrganizationFilters{
							ListFilter: domain.ListFilter{
								SortOrder: domain.SortOrderDesc,
								Pagination: domain.Pagination{
									Limit:  10,
									Offset: 0,
								},
							},
						}).
						Return(domain.EduOrganizations{
							domain.EduOrganization{
								ID:        id,
								Name:      orgName,
								Logo:      nil,
								CreatedAt: expectedTime,
								UpdatedAt: expectedTime,
								DeletedAt: nil,
							},
						}, totalCount, nil).AnyTimes()

					return map[string]string{}
				},
			},
			want: want{
				Code: http.StatusOK,
				Response: `
					{
					  "organizations": [
						{
						  "id": "351a12a6-7cdb-47e1-8105-b68d0a7b975d",
						  "name": "test organization",
						  "created_at": "2024-05-01T22:56:45+05:00",
						  "updated_at": "2024-05-01T22:56:45+05:00"
						}
					  ],
					  "pagination": {
						"page": 1,
						"per_page": 10,
						"total": 1
					  }
					}`,
			},
		},
	}
	for _, tt := range tests {
		tt := tt
		s.T().Run(tt.name, func(t *testing.T) {
			t.Parallel()

			var (
				assertion = assert.New(t)
				requiring = require.New(t)

				ctx = context.Background()
			)

			// creating our request.
			r, err := http.NewRequestWithContext(ctx, http.MethodGet, "/edu-organizations", http.NoBody)
			requiring.NoError(err)

			q := r.URL.Query()
			for k, v := range tt.req.reqBytes(t) {
				q.Add(k, v)
			}

			r.URL.RawQuery = q.Encode()

			// creating our test response writer.
			recorder := httptest.NewRecorder()

			// serving the request.
			router.ServeHTTP(recorder, r)
			assertion.Equal(tt.want.Code, recorder.Code)

			// checking the response.
			assertion.JSONEq(tt.want.Response, recorder.Body.String())
		})
	}
}

func (s *EduOrganizationSuite) TestEduOrganization_EduOrganizationListNegativeCases() {
	s.T().Parallel()

	var (
		// setting controller for our mocks.
		ctrl = gomock.NewController(s.Suite.T())

		// initializing mock organization service.
		eduOrganizationSvc = mockhandlers.NewMockIEduOrganizationService(ctrl)

		// initialization of educational organization handlers.
		eduOrganizationHandlers = handlers.NewEduOrganization(eduOrganizationSvc)

		// initialization of router and registering handlers.
		router = NewRouter()
	)

	// registering our handler.
	router.GET("/edu-organizations", eduOrganizationHandlers.EduOrganizationList)

	type args struct {
		reqBytes func(t *testing.T) map[string]string
	}

	type want struct {
		Code     int
		Response string
	}

	tests := []struct {
		name string
		req  args
		want want
	}{
		{
			name: "bad request",
			req: args{
				reqBytes: func(_ *testing.T) map[string]string {
					return map[string]string{"sort_order": "non_existed_order"}
				},
			},
			want: want{
				Code: http.StatusBadRequest,
				Response: `
					{
					  "code": "BAD_REQUEST"
					}`,
			},
		},
		{
			name: "internal error",
			req: args{
				reqBytes: func(_ *testing.T) map[string]string {
					var totalCount int

					// setting expectation for internal server error.
					eduOrganizationSvc.EXPECT().
						EduOrganizationList(gomock.Any(), domain.EduOrganizationFilters{
							ListFilter: domain.ListFilter{
								SortOrder: domain.SortOrderDesc,
								Pagination: domain.Pagination{
									Limit:  10,
									Offset: 0,
								},
							},
						}).
						Return(domain.EduOrganizations{}, totalCount, domain.ErrInternalServerError).AnyTimes()

					return map[string]string{}
				},
			},
			want: want{
				Code: http.StatusInternalServerError,
				Response: `
					{
					  "code": "INTERNAL_SERVER_ERROR"
					}`,
			},
		},
	}
	for _, tt := range tests {
		tt := tt
		s.T().Run(tt.name, func(t *testing.T) {
			t.Parallel()

			var (
				assertion = assert.New(t)
				requiring = require.New(t)

				ctx = context.Background()
			)

			// creating our request.
			r, err := http.NewRequestWithContext(ctx, http.MethodGet, "/edu-organizations", http.NoBody)
			requiring.NoError(err)

			q := r.URL.Query()
			for k, v := range tt.req.reqBytes(t) {
				q.Add(k, v)
			}

			r.URL.RawQuery = q.Encode()

			// creating our test response writer.
			recorder := httptest.NewRecorder()

			// serving the request.
			router.ServeHTTP(recorder, r)
			assertion.Equal(tt.want.Code, recorder.Code)

			// checking the response.
			compareErrorResponse(t, tt.want.Response, recorder.Body.String())
		})
	}
}

// NewRouter returns a new router with error handler and logger.
func NewRouter() *gin.Engine {
	var (
		router = gin.New()
		logger = liblog.NewDummyLogger()
	)

	router.Use(handlers.LoggingEndpointMiddleware(logger))
	router.Use(handlers.ErrorHandlingMiddleware())

	return router
}

func compareErrorResponse(t *testing.T, expected, actual string) {
	t.Helper()

	var expectedJSON, actualJSON liberror.Error

	if err := json.Unmarshal([]byte(expected), &expectedJSON); err != nil {
		assert.NoError(t, err)
	}

	if err := json.Unmarshal([]byte(actual), &actualJSON); err != nil {
		assert.NoError(t, err)
	}

	assert.Equal(t, expectedJSON.Code, actualJSON.Code)
}

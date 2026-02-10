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

// School is school handler.
type School struct {
	schoolService ISchoolService
}

// NewSchool creates a new school handler.
func NewSchool(schoolService ISchoolService) *School {
	return &School{
		schoolService: schoolService,
	}
}

// AddSchool adds a new school.
func (s School) AddSchool(c *gin.Context) {
	var (
		ctx    = c.Request.Context()
		logger = liblog.Must(ctx)
		req    request.AddSchool
		err    error
	)

	if err = c.ShouldBindJSON(&req); err != nil {
		logger.Errorf("failed to bind: %v", c.Error(domain.NewBadRequest(err.Error())))
		return
	}

	logger = logger.WithFields(liblog.Fields{"request": req.LogFields()})
	ctx = liblog.With(ctx, logger)

	entity, err := s.schoolService.AddSchool(
		ctx,
		school.CreateSchoolArgs{
			Name:           req.Name,
			OrganizationID: req.OrganizationUUID(),
			Location:       req.Location,
			Phone:          req.Phone,
			Email:          req.Email,
		},
	)
	if err != nil {
		logger.Errorf("failed to create a new school: %v", c.Error(err))
		return
	}

	c.JSON(http.StatusCreated, response.NewSchool(entity))
}

// SchoolByID get school by id.
func (s School) SchoolByID(c *gin.Context) { //nolint:dupl // It's ok.
	var (
		ctx      = c.Request.Context()
		logger   = liblog.Must(ctx)
		reqParam = request.GetSchoolIDPathVar(c)
		schoolID uuid.UUID
		err      error
	)

	logger = logger.WithFields(liblog.Fields{"request": liblog.Fields{"school_id": reqParam}})
	ctx = liblog.With(ctx, logger)

	if schoolID, err = uuid.Parse(reqParam); err != nil {
		logger.Errorf("failed to bind: %v", c.Error(domain.NewBadRequest(err.Error())))
		return
	}

	schoolEntity, err := s.schoolService.SchoolByID(ctx, schoolID)
	if err != nil {
		logger.Errorf("failed to get school by ID: %v", c.Error(err))
		return
	}

	c.JSON(http.StatusOK, response.NewSchool(schoolEntity))
}

// UpdateSchool updates school.
func (s School) UpdateSchool(c *gin.Context) {
	var (
		ctx      = c.Request.Context()
		logger   = liblog.Must(ctx)
		req      request.UpdateSchool
		reqParam = request.GetSchoolIDPathVar(c)
		schoolID uuid.UUID
		err      error
	)

	if err = c.ShouldBindJSON(&req); err != nil {
		logger.Errorf("failed to bind: %v", c.Error(domain.NewBadRequest(err.Error())))
		return
	}

	if schoolID, err = uuid.Parse(reqParam); err != nil {
		logger.Errorf("failed to bind: %v", c.Error(domain.NewBadRequest(err.Error())))
		return
	}

	logger = logger.WithFields(liblog.Fields{"request": req})
	ctx = liblog.With(ctx, logger)

	entity, err := s.schoolService.UpdateSchool(
		ctx,
		school.UpdateSchoolArgs{
			ID:              schoolID,
			Name:            req.Name,
			Location:        req.Location,
			Phone:           req.Phone,
			Email:           req.Email,
			GradeStandardID: req.GradeStandardID,
		},
	)
	if err != nil {
		logger.Errorf("failed to update school: %v", c.Error(err))
		return
	}

	c.JSON(http.StatusOK, response.NewSchool(entity))
}

// SchoolList get school list.
func (s School) SchoolList(c *gin.Context) {
	var (
		ctx    = c.Request.Context()
		logger = liblog.Must(ctx)
		req    request.SchoolList
		err    error
	)

	if err = c.ShouldBindQuery(&req); err != nil {
		logger.Errorf("failed to bind: %v", c.Error(domain.NewBadRequest(err.Error())))
		return
	}

	logger = logger.WithFields(liblog.Fields{"request": req})
	ctx = liblog.With(ctx, logger)

	list, total, err := s.schoolService.SchoolList(
		ctx,
		domain.NewSchoolFilters(
			domain.NewListFilter(req.SortOrder, domain.NewPagination(req.Page, req.PerPage)),
			req.Emails,
			req.Phones,
			req.OrganizationUUIDs(),
		),
	)
	if err != nil {
		logger.Errorf("failed to get school list: %v", c.Error(err))
		return
	}

	c.JSON(http.StatusOK, response.NewSchools(list, response.Pagination{
		Page:    req.Page,
		PerPage: req.PerPage,
		Total:   total,
	}))
}

// CreateGroup creates a new group.
func (s School) CreateGroup(c *gin.Context) {
	var (
		ctx             = c.Request.Context()
		logger          = liblog.Must(ctx)
		schoolIDPathVar = request.GetSchoolIDPathVar(c)
		schoolID        uuid.UUID
		req             request.CreateGroup
		err             error
	)

	logger = logger.WithFields(liblog.Fields{"request": liblog.Fields{"school_id": schoolIDPathVar}})
	ctx = liblog.With(ctx, logger)

	if schoolID, err = uuid.Parse(schoolIDPathVar); err != nil {
		logger.Errorf("failed to parse uuid: %v", c.Error(domain.NewBadRequest(err.Error())))
		return
	}

	if err = c.ShouldBindJSON(&req); err != nil {
		logger.Errorf("failed to bind: %v", c.Error(domain.NewBadRequest(err.Error())))
		return
	}

	logger = logger.WithFields(liblog.Fields{"request": req.LogFields()})
	ctx = liblog.With(ctx, logger)

	groupDomain, err := s.schoolService.CreateGroup(
		ctx,
		school.CreateGroupArgs{
			SchoolID: schoolID,
			Name:     req.Name,
			GradeID:  req.GradeID,
		},
	)
	if err != nil {
		logger.Errorf("failed to create a new group: %v", c.Error(err))
		return
	}

	c.JSON(http.StatusCreated, response.NewGroup(groupDomain))
}

// GroupByID gets group by id.
func (s School) GroupByID(c *gin.Context) {
	var (
		ctx             = c.Request.Context()
		logger          = liblog.Must(ctx)
		schoolIDPathVar = request.GetSchoolIDPathVar(c)
		groupIDPathVar  = request.GetGroupIDPathVar(c)
		groupID         uuid.UUID
		err             error
	)

	logger = logger.WithFields(liblog.Fields{
		"request": liblog.Fields{
			"school_id": schoolIDPathVar,
			"group_id":  groupIDPathVar,
		},
	})
	ctx = liblog.With(ctx, logger)

	if groupID, err = uuid.Parse(groupIDPathVar); err != nil {
		logger.Errorf("failed to parse uuid: %v", c.Error(domain.NewBadRequest(err.Error())))
		return
	}

	groupDomain, err := s.schoolService.GroupByID(ctx, groupID)
	if err != nil {
		logger.Errorf("failed to create a new group: %v", c.Error(err))
		return
	}

	c.JSON(http.StatusOK, response.NewGroup(groupDomain))
}

// UpdateGroup gets group by id.
func (s School) UpdateGroup(c *gin.Context) {
	var (
		ctx             = c.Request.Context()
		logger          = liblog.Must(ctx)
		schoolIDPathVar = request.GetSchoolIDPathVar(c)
		groupIDPathVar  = request.GetGroupIDPathVar(c)
		groupID         uuid.UUID
		schoolID        uuid.UUID
		req             request.UpdateGroup
		err             error
	)

	logger = logger.WithFields(liblog.Fields{"request": liblog.Fields{"school_id": schoolIDPathVar}})
	ctx = liblog.With(ctx, logger)

	if schoolID, err = uuid.Parse(schoolIDPathVar); err != nil {
		logger.Errorf("failed to parse uuid: %v", c.Error(domain.NewBadRequest(err.Error())))
		return
	}

	_ = schoolID

	if groupID, err = uuid.Parse(groupIDPathVar); err != nil {
		logger.Errorf("failed to parse uuid: %v", c.Error(domain.NewBadRequest(err.Error())))
		return
	}

	if err = c.ShouldBindJSON(&req); err != nil {
		logger.Errorf("failed to bind: %v", c.Error(domain.NewBadRequest(err.Error())))
		return
	}

	logger = logger.WithFields(liblog.Fields{
		"request":   req,
		"group_id":  groupIDPathVar,
		"school_id": schoolIDPathVar,
	})
	ctx = liblog.With(ctx, logger)

	groupDomain, err := s.schoolService.UpdateGroup(
		ctx,
		school.UpdateGroupArgs{
			ID:      groupID,
			Name:    req.Name,
			GradeID: req.GradeID,

			ClassTeacherID:         req.ClassTeacherID,
			ClassPresidentID:       req.ClassPresidentID,
			DeputyClassPresidentID: req.DeputyClassPresidentID,
		},
	)
	if err != nil {
		logger.Errorf("failed to update group: %v", c.Error(err))
		return
	}

	c.JSON(http.StatusOK, response.NewGroup(groupDomain))
}

// GroupList get group list.
func (s School) GroupList(c *gin.Context) {
	var (
		ctx             = c.Request.Context()
		logger          = liblog.Must(ctx)
		schoolIDPathVar = request.GetSchoolIDPathVar(c)
		schoolID        uuid.UUID
		req             request.GroupList
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

	logger = logger.WithFields(liblog.Fields{"request": req})
	ctx = liblog.With(ctx, logger)

	list, total, err := s.schoolService.GroupList(
		ctx,
		schoolID,
		domain.NewGroupFilters(),
	)
	if err != nil {
		logger.Errorf("failed to get group list: %v", c.Error(err))
		return
	}

	c.JSON(http.StatusOK, response.NewGroupsList(list, total))
}

// AddSchoolSubject adds a new school subject.
func (s School) AddSchoolSubject(c *gin.Context) {
	var (
		ctx             = c.Request.Context()
		logger          = liblog.Must(ctx)
		schoolIDPathVar = request.GetSchoolIDPathVar(c)
		schoolID        uuid.UUID
		req             request.CreateSchoolSubject
		err             error
	)

	logger = logger.WithFields(liblog.Fields{"school_id": schoolIDPathVar})
	ctx = liblog.With(ctx, logger)

	if schoolID, err = uuid.Parse(schoolIDPathVar); err != nil {
		logger.Errorf("failed to parse uuid: %v", c.Error(domain.NewBadRequest(err.Error())))
		return
	}

	if err = c.ShouldBindJSON(&req); err != nil {
		logger.Errorf("failed to bind: %v", c.Error(domain.NewBadRequest(err.Error())))
		return
	}

	logger = logger.WithFields(liblog.Fields{"request": req.LogFields()})
	ctx = liblog.With(ctx, logger)

	schoolSubjectDomain, err := s.schoolService.CreateSchoolSubject(
		ctx,
		school.CreateSchoolSubjectArgs{
			SchoolID:    schoolID,
			SubjectID:   req.SubjectID,
			Name:        req.Name,
			Description: req.Description,
		},
	)
	if err != nil {
		logger.Errorf("failed to create a new school subject: %v", c.Error(err))
		return
	}

	c.JSON(http.StatusCreated, response.NewSchoolSubject(schoolSubjectDomain))
}

// SchoolSubjectByIDAndSchoolID gets school subject by id and school id.
func (s School) SchoolSubjectByIDAndSchoolID(c *gin.Context) {
	var (
		ctx                    = c.Request.Context()
		logger                 = liblog.Must(ctx)
		schoolIDPathVar        = request.GetSchoolIDPathVar(c)
		schoolID               uuid.UUID
		schoolSubjectID        uuid.UUID
		schoolSubjectIDPathVar = request.GetSchoolSubjectIDPathVar(c)
		err                    error
	)

	logger = logger.WithFields(liblog.Fields{"school_id": schoolIDPathVar})
	ctx = liblog.With(ctx, logger)

	if schoolID, err = uuid.Parse(schoolIDPathVar); err != nil {
		logger.Errorf("failed to parse uuid: %v", c.Error(domain.NewBadRequest(err.Error())))
		return
	}

	if schoolSubjectID, err = uuid.Parse(schoolSubjectIDPathVar); err != nil {
		logger.Errorf("failed to parse uuid: %v", c.Error(domain.NewBadRequest(err.Error())))
		return
	}

	schoolSubjectDomain, err := s.schoolService.SchoolSubjectByIDAndSchoolID(ctx, schoolSubjectID, schoolID)
	if err != nil {
		logger.Errorf("failed to create a new school subject: %v", c.Error(err))
		return
	}

	c.JSON(http.StatusOK, response.NewSchoolSubject(schoolSubjectDomain))
}

// SchoolSubjectList get school subject list.
func (s School) SchoolSubjectList(c *gin.Context) {
	var (
		ctx             = c.Request.Context()
		logger          = liblog.Must(ctx)
		req             request.SchoolSubjectList
		schoolIDPathVar = request.GetSchoolIDPathVar(c)
		schoolID        uuid.UUID
		err             error
	)

	logger = logger.WithFields(liblog.Fields{"school_id": schoolIDPathVar})
	ctx = liblog.With(ctx, logger)

	if schoolID, err = uuid.Parse(schoolIDPathVar); err != nil {
		logger.Errorf("failed to parse uuid: %v", c.Error(domain.NewBadRequest(err.Error())))
		return
	}

	if err = c.ShouldBindQuery(&req); err != nil {
		logger.Errorf("failed to bind: %v", c.Error(domain.NewBadRequest(err.Error())))
		return
	}

	logger = logger.WithFields(liblog.Fields{"request": req})
	ctx = liblog.With(ctx, logger)

	list, total, err := s.schoolService.SchoolSubjectList(
		ctx,
		domain.NewSchoolSubjectFilters(
			domain.NewListFilter(req.SortOrder, domain.NewPagination(req.Page, req.PerPage)),
			schoolID,
		),
	)
	if err != nil {
		logger.Errorf("failed to get school subjects list: %v", c.Error(err))
		return
	}

	c.JSON(http.StatusOK, response.NewSchoolSubjectsList(list, response.Pagination{
		Page:    req.Page,
		PerPage: req.PerPage,
		Total:   total,
	}))
}

// GroupSubjectList get group subject list.
func (s School) GroupSubjectList(c *gin.Context) {
	var (
		ctx            = c.Request.Context()
		logger         = liblog.Must(ctx)
		groupIDPathVar = request.GetGroupIDPathVar(c)
		groupID        uuid.UUID
		err            error
	)

	logger = logger.WithFields(liblog.Fields{
		"group_id": groupIDPathVar,
	})
	ctx = liblog.With(ctx, logger)

	if groupID, err = uuid.Parse(groupIDPathVar); err != nil {
		logger.Errorf("failed to parse uuid: %v", c.Error(domain.NewBadRequest(err.Error())))
		return
	}

	list, err := s.schoolService.GroupSubjectList(ctx, groupID)
	if err != nil {
		logger.Errorf("failed to get group subjects list: %v", c.Error(err))
		return
	}

	c.JSON(http.StatusOK, response.NewGroupSubjects(list))
}

// AddGroupSubject adds a new subjects to group.
func (s School) AddGroupSubject(c *gin.Context) {
	var (
		ctx             = c.Request.Context()
		logger          = liblog.Must(ctx)
		req             request.AddGroupSubject
		schoolIDPathVar = request.GetSchoolIDPathVar(c)
		groupIDPathVar  = request.GetGroupIDPathVar(c)
		schoolID        uuid.UUID
		groupID         uuid.UUID
		err             error
	)

	if schoolID, err = uuid.Parse(schoolIDPathVar); err != nil {
		logger.Errorf("failed to parse uuid: %v", c.Error(domain.NewBadRequest(err.Error())))
		return
	}

	if groupID, err = uuid.Parse(groupIDPathVar); err != nil {
		logger.Errorf("failed to parse uuid: %v", c.Error(domain.NewBadRequest(err.Error())))
		return
	}

	if err = c.ShouldBindJSON(&req); err != nil {
		logger.Errorf("failed to bind: %v", c.Error(domain.NewBadRequest(err.Error())))
		return
	}

	logger = logger.WithFields(liblog.Fields{
		"request":   req,
		"school_id": schoolID,
		"group_id":  groupID,
	})
	ctx = liblog.With(ctx, logger)

	newGroupSubject, err := s.schoolService.AddGroupSubject(ctx, schoolID, groupID, school.AddGroupSubjectArgs{
		SchoolSubjectID: req.SchoolSubjectID,
		TeacherID:       req.TeacherID,
		Count:           req.Count,
	})
	if err != nil {
		logger.Errorf("failed to add group subject: %v", c.Error(err))
		return
	}

	c.JSON(http.StatusCreated, response.NewGroupSubject(newGroupSubject))
}

// AssignStudyPlans assigns study plana.
func (s School) AssignStudyPlans(c *gin.Context) {
	var (
		ctx                   = c.Request.Context()
		logger                = liblog.Must(ctx)
		req                   []request.AssignStudyPlan
		groupSubjectIDPathVar = request.GetGroupSubjectIDPathVar(c)
		groupSubjectID        uuid.UUID
		err                   error
	)

	if groupSubjectID, err = uuid.Parse(groupSubjectIDPathVar); err != nil {
		logger.Errorf("failed to parse uuid: %v", c.Error(domain.NewBadRequest(err.Error())))
		return
	}

	if err = c.ShouldBindJSON(&req); err != nil {
		logger.Errorf("failed to bind: %v", c.Error(domain.NewBadRequest(err.Error())))
		return
	}

	logger = logger.WithFields(liblog.Fields{
		"request":          req,
		"group_subject_id": groupSubjectID,
	})
	ctx = liblog.With(ctx, logger)

	studyPlans := make([]school.AddStudyPlanArgs, 0, len(req))

	for i, sp := range req {
		studyPlans = append(studyPlans, school.AddStudyPlanArgs{
			ID:          sp.ID,
			Title:       sp.Title,
			Description: sp.Description,
			PlanOrder:   int16(i),
		})
	}

	newStudyPlans, err := s.schoolService.AssignStudyPlans(ctx, groupSubjectID, studyPlans)
	if err != nil {
		logger.Errorf("failed to add study plan: %v", c.Error(err))
		return
	}

	c.JSON(http.StatusOK, response.NewStudyPlans(newStudyPlans))
}

// StudyPlanList get study plan list.
func (s School) StudyPlanList(c *gin.Context) {
	var (
		ctx                   = c.Request.Context()
		logger                = liblog.Must(ctx)
		groupSubjectIDPathVar = request.GetGroupSubjectIDPathVar(c)
		groupSubjectID        uuid.UUID
		err                   error
	)

	if groupSubjectID, err = uuid.Parse(groupSubjectIDPathVar); err != nil {
		logger.Errorf("failed to parse uuid: %v", c.Error(domain.NewBadRequest(err.Error())))
		return
	}

	logger = logger.WithFields(liblog.Fields{
		"group_subject_id": groupSubjectID,
	})
	ctx = liblog.With(ctx, logger)

	newStudyPlans, err := s.schoolService.StudyPlanList(ctx, groupSubjectID)
	if err != nil {
		logger.Errorf("failed to get study plan list: %v", c.Error(err))
		return
	}

	c.JSON(http.StatusOK, response.NewStudyPlans(newStudyPlans))
}

// StudyPlanChangeStatus set study plan status.
func (s School) StudyPlanChangeStatus(c *gin.Context) {
	var (
		ctx                   = c.Request.Context()
		logger                = liblog.Must(ctx)
		groupSubjectIDPathVar = request.GetGroupSubjectIDPathVar(c)
		studyPlanIDPathVar    = request.GetStudyPlanIDPathVar(c)
		StudyPlanStatus       = request.GetStudyPlanStatusPathVar(c)
		groupSubjectID        uuid.UUID
		studyPlanID           uuid.UUID
		err                   error
	)

	if groupSubjectID, err = uuid.Parse(groupSubjectIDPathVar); err != nil {
		logger.Errorf("failed to parse uuid: %v", c.Error(domain.NewBadRequest(err.Error())))
		return
	}

	if studyPlanID, err = uuid.Parse(studyPlanIDPathVar); err != nil {
		logger.Errorf("failed to parse uuid: %v", c.Error(domain.NewBadRequest(err.Error())))
		return
	}

	logger = logger.WithFields(liblog.Fields{
		"group_subject_id": groupSubjectID,
		"study_plan_id":    groupSubjectID,
	})
	ctx = liblog.With(ctx, logger)

	err = s.schoolService.StudyPlanChangeStatus(ctx, groupSubjectID, studyPlanID, StudyPlanStatus)
	if err != nil {
		logger.Errorf("failed to change study plan status: %v", c.Error(err))
		return
	}

	c.Status(http.StatusOK)
}

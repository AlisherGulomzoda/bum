package http

import (
	"time"

	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"

	"bum-service/internal/controller/http/handlers"
	"bum-service/pkg/liblog"
)

// RegisterHandlers creates a new router and registers handlers.
func RegisterHandlers(
	router *gin.Engine,
	logger liblog.Logger,

	jwtSecret string,
	accessTokenExp time.Duration,
	refreshTokenExp time.Duration,

	authService handlers.IAuthService,
	systemService handlers.ISystemService,
	eduOrganizationService handlers.IEduOrganizationService,
	ownerService handlers.IOwnerService,
	schoolService handlers.ISchoolService,
	directorService handlers.IDirectorService,
	headmasterService handlers.IHeadmasterService,
	subjectService handlers.ISubjectService,
	userService handlers.IUserService,
	teacherService handlers.ITeacherService,
	gradesService handlers.IGradesService,
	studentService handlers.IStudentService,
	lessonService handlers.ILessonService,
) error {
	router.Use(gin.Logger())
	router.Use(handlers.LoggingEndpointMiddleware(logger))
	router.Use(handlers.ErrorHandlingMiddleware())
	router.Use(handlers.RecoverMiddleware(logger))

	corsConfig := cors.DefaultConfig()
	corsConfig.AllowAllOrigins = true
	corsConfig.AddAllowHeaders("Authorization")

	router.Use(cors.New(corsConfig))

	routerV1 := router.Group("/v1")

	auth := registerAuthHandlers(routerV1, jwtSecret, accessTokenExp, refreshTokenExp, authService, userService)

	registerSystemHandlers(routerV1, systemService)

	registerEduOrganizationHandlers(routerV1, eduOrganizationService)

	registerOwnerHandlers(routerV1, ownerService)

	registerSchoolHandlers(routerV1, schoolService)

	registerDirectorHandlers(routerV1, directorService)

	registerHeadmasterHandlers(routerV1, headmasterService)

	registerSubjectHandlers(routerV1, subjectService)

	registerUserHandlers(routerV1, auth, userService)

	registerTeacherHandlers(routerV1, teacherService)

	registerStudentsHandlers(routerV1, studentService)

	registerGradesHandlers(routerV1, gradesService)

	registerLessonsHandlers(routerV1, lessonService)

	return nil
}

// registerSystemHandlers registers all system handlers.
func registerSystemHandlers(router *gin.RouterGroup, systemService handlers.ISystemService) {
	h := handlers.NewSystem(systemService)

	router.GET("/health", h.HealthCheck)
}

// registerAuthHandlers registers all Auth handlers.
func registerAuthHandlers(
	router *gin.RouterGroup,
	jwtSecret string,
	accessTokenExp time.Duration,
	refreshTokenExp time.Duration,
	authService handlers.IAuthService,
	userService handlers.IUserService,
) *handlers.Auth {
	h := handlers.NewAuth(
		authService,
		userService,

		jwtSecret,
		accessTokenExp,
		refreshTokenExp,
	)

	router.POST("/login/email", h.LoginByEmail)
	router.POST("/login/refresh", h.AuthMiddleware, h.RefreshToken)

	return h
}

func registerUserHandlers(router *gin.RouterGroup, auth *handlers.Auth, userService handlers.IUserService) {
	h := handlers.NewUser(userService)

	router.POST("/users", h.AddUser)
	router.GET("/user/full-info", auth.AuthMiddleware, h.UserFullInfoFromToken)
	router.GET("/user/:user_id/full-info", h.UserFullInfoByID)
	router.GET("/users", h.UserList)
}

// registerEduOrganizationHandlers registers all educational organization handlers.
func registerEduOrganizationHandlers(router *gin.RouterGroup, eduOrganizationService handlers.IEduOrganizationService) {
	h := handlers.NewEduOrganization(eduOrganizationService)

	router.POST("/edu-organizations", h.CreateEduOrganization)
	router.GET("/edu-organizations/:edu_organization_id", h.EduOrganizationByID)
	router.PUT("/edu-organizations/:edu_organization_id", h.UpdateEduOrganizationByID)
	router.GET("/edu-organizations", h.EduOrganizationList)
}

// registerSchoolHandlers registers all school handlers.
func registerSchoolHandlers(
	router *gin.RouterGroup,

	schoolService handlers.ISchoolService,
) {
	schoolHandlers := handlers.NewSchool(schoolService)

	// SCHOOL
	router.POST("/schools", schoolHandlers.AddSchool)
	router.GET("/schools/:school_id", schoolHandlers.SchoolByID)
	router.PUT("/schools/:school_id", schoolHandlers.UpdateSchool)
	router.GET("/schools", schoolHandlers.SchoolList)

	// GROUPS
	router.POST("/schools/:school_id/groups", schoolHandlers.CreateGroup)
	router.GET("/schools/:school_id/groups/:group_id", schoolHandlers.GroupByID)
	router.PUT("/schools/:school_id/groups/:group_id", schoolHandlers.UpdateGroup)
	router.GET("/schools/:school_id/groups", schoolHandlers.GroupList)

	// GROUPS SUBJECTS
	router.POST("/schools/:school_id/groups/:group_id/subjects", schoolHandlers.AddGroupSubject)
	router.GET("/schools/:school_id/groups/:group_id/subjects", schoolHandlers.GroupSubjectList)

	// SCHOOL SUBJECTS
	router.POST("/schools/:school_id/subjects", schoolHandlers.AddSchoolSubject)
	router.GET("/schools/:school_id/subjects/:school_subject_id", schoolHandlers.SchoolSubjectByIDAndSchoolID)
	router.GET("/schools/:school_id/subjects", schoolHandlers.SchoolSubjectList)

	// AUDITORIUMS
	router.POST("/schools/:school_id/auditoriums", schoolHandlers.CreateAuditorium)
	router.GET("/schools/:school_id/auditoriums", schoolHandlers.AuditoriumList)
	router.GET("/schools/:school_id/auditoriums/:auditorium_id", schoolHandlers.AuditoriumByIDAndSchoolID)

	// STUDY PLAN
	router.PUT("/schools/:school_id/group_subjects/:group_subject_id/study-plan", schoolHandlers.AssignStudyPlans)
	router.GET("/schools/:school_id/group_subjects/:group_subject_id/study-plan", schoolHandlers.StudyPlanList)
	router.PATCH(
		"/schools/:school_id/group_subjects/:group_subject_id/study-plan/:study_plan_id/status/:study_plan_status",
		schoolHandlers.StudyPlanChangeStatus,
	)
}

// registerOwnerHandlers registers all owner handlers.
func registerOwnerHandlers(router *gin.RouterGroup, ownerService handlers.IOwnerService) {
	h := handlers.NewOwner(ownerService)

	router.POST("/owners", h.AddOwner)
	router.GET("/owners/:owner_id", h.OwnerByID)
	router.GET("/owners/me", h.OwnerByUserIDAndSchoolID)
	router.GET("/owners", h.OwnerList)
}

// registerDirectorHandlers registers all director handlers.
func registerDirectorHandlers(router *gin.RouterGroup, directorService handlers.IDirectorService) {
	h := handlers.NewDirector(directorService)

	router.POST("/directors", h.AddDirector)
	router.GET("/directors/:director_id", h.DirectorByID)
	router.GET("/directors", h.DirectorList)
}

// registerHeadmasterHandlers registers all headmaster handlers.
func registerHeadmasterHandlers(router *gin.RouterGroup, headmasterService handlers.IHeadmasterService) {
	h := handlers.NewHeadmaster(headmasterService)

	router.POST("/headmasters", h.AddHeadmaster)
	router.GET("/headmasters/:headmaster_id", h.HeadmasterByID)
	router.GET("/headmasters", h.HeadmasterList)
}

func registerSubjectHandlers(router *gin.RouterGroup, subjectService handlers.ISubjectService) {
	h := handlers.NewSubject(subjectService)

	router.POST("/subjects", h.CreateSubject)
	router.GET("/subjects/:subject_id", h.SubjectByID)
	router.GET("/subjects", h.SubjectList)
}

func registerTeacherHandlers(router *gin.RouterGroup, teachersUseCase handlers.ITeacherService) {
	h := handlers.NewTeacher(teachersUseCase)

	router.POST("/teachers", h.AddTeacher)
	router.GET("/teachers/:teacher_id", h.TeacherByID)
	router.GET("/teachers", h.ListTeacher)
}

func registerStudentsHandlers(router *gin.RouterGroup, studentService handlers.IStudentService) {
	studentHandlers := handlers.NewStudent(studentService)

	// STUDENTS
	router.POST("/students", studentHandlers.AddStudent)
	router.GET("/students/:student_id", studentHandlers.StudentByID)
	router.GET("/students", studentHandlers.StudentList)

	// STUDENT GUARDIANS
	router.POST("/students/:student_id/guardians", studentHandlers.AssignStudentGuardian)
	router.GET("/students/:student_id/guardians", studentHandlers.StudentGuardians)
	router.GET("/students/guardians/:user_id", studentHandlers.StudentGuardianByUserID)
	router.GET("/students/guardians", studentHandlers.StudentGuardianList)
}

// registerGradesHandlers registers all grade-standard handlers.
func registerGradesHandlers(router *gin.RouterGroup, gradesUseCase handlers.IGradesService) {
	h := handlers.NewGrades(gradesUseCase)

	router.POST("/grade-standards", h.CreateGradeStandard)
	router.GET("/grade-standards/:grade_standard_id", h.GradeStandardByID)
	router.GET("/grade-standards", h.GradeStandardList)
}

// registerLessonsHandlers registers all lessons handlers.
func registerLessonsHandlers(router *gin.RouterGroup, lessonService handlers.ILessonService) {
	h := handlers.NewLesson(lessonService)

	// LESSONS
	router.PUT("/lessons", h.AssignWeekLessons)
	router.GET("/lessons", h.LessonsList)

	// STUDENT MARKS
	router.POST("lessons/marks", h.AddMark)
	router.GET("lessons/marks/:mark_id", h.MarkByID)
	router.GET("lessons/marks", h.AddMark)
}

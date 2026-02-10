package handlers

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"

	"bum-service/internal/controller/http/handlers/request"
	"bum-service/internal/controller/http/handlers/response"
	"bum-service/internal/domain"
	"bum-service/internal/service/lesson"
	"bum-service/pkg/liblog"
)

// Lesson is Lesson handler.
type Lesson struct {
	lessonService ILessonService
}

// NewLesson creates a new lesson handler.
func NewLesson(
	lessonService ILessonService,
) *Lesson {
	return &Lesson{
		lessonService: lessonService,
	}
}

// AssignWeekLessons assigns lessons for the week.
func (l *Lesson) AssignWeekLessons(c *gin.Context) {
	var (
		ctx               = c.Request.Context()
		logger            = liblog.Must(ctx)
		req               request.AssignWeekLessons
		schoolIDHeaderVar = request.GetSchoolIDHeader(c)
		schoolID          uuid.UUID
		err               error
	)

	if schoolID, err = uuid.Parse(schoolIDHeaderVar); err != nil {
		logger.Errorf("failed to parse uuid: %v", c.Error(domain.NewBadRequest(err.Error())))
		return
	}

	if err = c.ShouldBindJSON(&req); err != nil {
		logger.Errorf("failed to bind: %v", c.Error(domain.NewBadRequest(err.Error())))
		return
	}

	logger = logger.WithFields(liblog.Fields{"request": req})
	ctx = liblog.With(ctx, logger)

	addLesson, err := l.lessonService.AssignLessons(ctx, convertAssignWeekLessonsToServiceArgs(req, schoolID))
	if err != nil {
		logger.Errorf("failed to assign lessons: %v", c.Error(domain.NewBadRequest(err.Error())))
		return
	}

	c.JSON(http.StatusOK, addLesson)
}

func convertAssignWeekLessonsToServiceArgs(r request.AssignWeekLessons, schoolID uuid.UUID) lesson.AddWeekLessonsArgs {
	lessons := make([]lesson.Lesson, 0, len(r.LessonItem))

	for _, l := range r.LessonItem {
		lessons = append(lessons, lesson.Lesson{
			GroupSubjectID: l.GroupSubjectID,
			TeacherID:      l.TeacherID,
			AuditoriumID:   l.AuditoriumID,
			StartTime:      l.StartTime,
			EndTime:        l.EndTime,
			Description:    l.Description,
		})
	}

	return lesson.AddWeekLessonsArgs{
		SchoolID: schoolID,
		GroupID:  r.GroupID,
		WeekDate: r.WeekDate,
		Lessons:  lessons,
	}
}

// LessonsList get lessons list.
func (l *Lesson) LessonsList(c *gin.Context) {
	var (
		ctx    = c.Request.Context()
		logger = liblog.Must(ctx)
		req    request.LessonsList
		err    error
	)

	if err = c.ShouldBindQuery(&req); err != nil {
		logger.Errorf("failed to bind: %v", c.Error(domain.NewBadRequest(err.Error())))
		return
	}

	logger = logger.WithFields(liblog.Fields{"request": req})
	ctx = liblog.With(ctx, logger)

	list, err := l.lessonService.LessonsList(ctx,
		domain.NewLessonsListFilter(
			domain.NewDateFilter(req.Period.DateFrom(), req.Period.DateTill()),
			domain.NewListFilter(req.SortOrder, domain.NewPagination(req.Page, req.PerPage)),
			req.SchoolID,
			req.TeacherID,
			req.GroupID,
		),
	)
	if err != nil {
		logger.Errorf("failed to get lessons list: %v", c.Error(err))
		return
	}

	c.JSON(http.StatusOK, response.NewLessons(list))
}

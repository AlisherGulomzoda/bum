package handlers

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"

	"bum-service/internal/controller/http/handlers/request"
	"bum-service/internal/domain"
	"bum-service/internal/service/lesson"
	"bum-service/pkg/liblog"
)

// AddMark adds mark to student.
func (l *Lesson) AddMark(c *gin.Context) {
	var (
		ctx    = c.Request.Context()
		logger = liblog.Must(ctx)
		req    request.AddMark
		err    error
	)

	if err = c.ShouldBindJSON(&req); err != nil {
		logger.Errorf("failed to bind: %v", c.Error(domain.NewBadRequest(err.Error())))
		return
	}

	logger = logger.WithFields(liblog.Fields{"request": req})
	ctx = liblog.With(ctx, logger)

	markEntity, err := l.lessonService.AddMark(ctx, lesson.AddMarkArgs{
		LessonID:    req.LessonID,
		StudentID:   req.StudentID,
		Mark:        req.Mark,
		Description: req.Description,
	})
	if err != nil {
		logger.Errorf("failed add marks: %v", c.Error(domain.NewBadRequest(err.Error())))
		return
	}

	c.JSON(http.StatusCreated, markEntity)
}

// MarkByID returns mark by id.
func (l *Lesson) MarkByID(c *gin.Context) {
	var (
		ctx      = c.Request.Context()
		logger   = liblog.Must(ctx)
		markID   = request.GetMarkIDPathVar(c)
		markUUID uuid.UUID
		err      error
	)

	markUUID, err = uuid.Parse(markID)
	if err != nil {
		logger.Errorf("failed to parce to uuid: %v", c.Error(domain.NewBadRequest(err.Error())))
		return
	}

	logger = logger.WithFields(liblog.Fields{"mark_id": markID})
	ctx = liblog.With(ctx, logger)

	markEntity, err := l.lessonService.MarkByID(ctx, markUUID)
	if err != nil {
		logger.Errorf("failed get mark by id: %v", c.Error(domain.NewBadRequest(err.Error())))
		return
	}

	c.JSON(http.StatusCreated, markEntity)
}

// MarkList returns list of marks.
func (l *Lesson) MarkList(c *gin.Context) {
	var (
		ctx      = c.Request.Context()
		logger   = liblog.Must(ctx)
		markID   = request.GetMarkIDPathVar(c)
		markUUID uuid.UUID
		err      error
	)

	markUUID, err = uuid.Parse(markID)
	if err != nil {
		logger.Errorf("failed to parce to uuid: %v", c.Error(domain.NewBadRequest(err.Error())))
		return
	}

	logger = logger.WithFields(liblog.Fields{"mark_id": markID})
	ctx = liblog.With(ctx, logger)

	markEntity, err := l.lessonService.MarkByID(ctx, markUUID)
	if err != nil {
		logger.Errorf("failed get mark by id: %v", c.Error(domain.NewBadRequest(err.Error())))
		return
	}

	c.JSON(http.StatusCreated, markEntity)
}

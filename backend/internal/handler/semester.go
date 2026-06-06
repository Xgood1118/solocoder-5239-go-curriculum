package handler

import (
	"curriculum/internal/model"
	"curriculum/internal/service"
	"curriculum/internal/store"
	"net/http"

	"github.com/gin-gonic/gin"
)

type Handler struct {
	store *store.Store
	svc   *service.CurriculumService
}

func NewHandler(s *store.Store, svc *service.CurriculumService) *Handler {
	return &Handler{store: s, svc: svc}
}

func (h *Handler) ListSemesters(c *gin.Context) {
	semesters, err := h.store.ListSemesters()
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, semesters)
}

func (h *Handler) GetSemester(c *gin.Context) {
	id := c.Param("id")
	sem, err := h.store.GetSemester(id)
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, sem)
}

func (h *Handler) CreateSemester(c *gin.Context) {
	var sem model.Semester
	if err := c.ShouldBindJSON(&sem); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	if err := h.store.CreateSemester(&sem); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusCreated, sem)
}

func (h *Handler) UpdateSemester(c *gin.Context) {
	id := c.Param("id")
	var sem model.Semester
	if err := c.ShouldBindJSON(&sem); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	sem.ID = id
	if err := h.store.UpdateSemester(&sem); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, sem)
}

func (h *Handler) DeleteSemester(c *gin.Context) {
	id := c.Param("id")
	if err := h.store.DeleteSemester(id); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, gin.H{"message": "deleted"})
}

func (h *Handler) ListUnits(c *gin.Context) {
	semID := c.Param("semesterId")
	sem, err := h.store.GetSemester(semID)
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, sem.Units)
}

func (h *Handler) CreateUnit(c *gin.Context) {
	semID := c.Param("semesterId")
	var unit model.Unit
	if err := c.ShouldBindJSON(&unit); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	if err := h.store.AddUnit(semID, &unit); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusCreated, unit)
}

func (h *Handler) UpdateUnit(c *gin.Context) {
	semID := c.Param("semesterId")
	unitID := c.Param("unitId")
	var unit model.Unit
	if err := c.ShouldBindJSON(&unit); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	unit.ID = unitID
	if err := h.store.UpdateUnit(semID, &unit); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, unit)
}

func (h *Handler) DeleteUnit(c *gin.Context) {
	semID := c.Param("semesterId")
	unitID := c.Param("unitId")
	if err := h.store.DeleteUnit(semID, unitID); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, gin.H{"message": "deleted"})
}

func (h *Handler) SuggestUnitLessons(c *gin.Context) {
	semID := c.Param("semesterId")
	unitID := c.Param("unitId")
	suggestion, err := h.svc.SuggestUnitLessons(semID, unitID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, suggestion)
}

func (h *Handler) ListLessons(c *gin.Context) {
	semID := c.Param("semesterId")
	unitID := c.Param("unitId")
	unit, err := h.store.GetUnit(semID, unitID)
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, unit.Lessons)
}

func (h *Handler) CreateLesson(c *gin.Context) {
	semID := c.Param("semesterId")
	unitID := c.Param("unitId")
	var lesson model.Lesson
	if err := c.ShouldBindJSON(&lesson); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	if err := h.store.AddLesson(semID, unitID, &lesson); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusCreated, lesson)
}

func (h *Handler) UpdateLesson(c *gin.Context) {
	semID := c.Param("semesterId")
	unitID := c.Param("unitId")
	lessonID := c.Param("lessonId")
	var lesson model.Lesson
	if err := c.ShouldBindJSON(&lesson); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	lesson.ID = lessonID
	if err := h.store.UpdateLesson(semID, unitID, &lesson); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, lesson)
}

func (h *Handler) DeleteLesson(c *gin.Context) {
	semID := c.Param("semesterId")
	unitID := c.Param("unitId")
	lessonID := c.Param("lessonId")
	if err := h.store.DeleteLesson(semID, unitID, lessonID); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, gin.H{"message": "deleted"})
}

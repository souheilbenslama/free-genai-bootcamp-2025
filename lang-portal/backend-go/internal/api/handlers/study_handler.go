package handlers

import (
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
	"github.com/souheilbenslama/free-genai-bootcamp-2025/lang-portal/backend-go/internal/repository/sqlite"
)

type StudyHandler struct {
	repo *sqlite.StudyRepository
}

func NewStudyHandler(repo *sqlite.StudyRepository) *StudyHandler {
	return &StudyHandler{repo: repo}
}

func (h *StudyHandler) GetLastStudySession(c *gin.Context) {
	session, err := h.repo.GetLastStudySession()
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	if session == nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "no study sessions found"})
		return
	}

	c.JSON(http.StatusOK, session)
}

func (h *StudyHandler) GetStudyProgress(c *gin.Context) {
	progress, err := h.repo.GetStudyProgress()
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, progress)
}

func (h *StudyHandler) GetQuickStats(c *gin.Context) {
	stats, err := h.repo.GetQuickStats()
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, stats)
}

type StartStudySessionRequest struct {
	GroupID int `json:"group_id" binding:"required"`
}

func (h *StudyHandler) StartStudySession(c *gin.Context) {
	var req StartStudySessionRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	session, err := h.repo.CreateStudySession(req.GroupID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusCreated, session)
}

type WordReviewRequest struct {
	WordID  int  `json:"word_id" binding:"required"`
	Correct bool `json:"correct" binding:"required"`
}

func (h *StudyHandler) RecordWordReview(c *gin.Context) {
	sessionID, err := strconv.Atoi(c.Param("session_id"))
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid session ID"})
		return
	}

	var req WordReviewRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	err = h.repo.RecordWordReview(sessionID, req.WordID, req.Correct)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.Status(http.StatusNoContent)
}

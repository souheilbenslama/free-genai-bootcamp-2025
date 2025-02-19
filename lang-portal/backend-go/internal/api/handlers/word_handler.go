package handlers

import (
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
	"github.com/souheilbenslama/free-genai-bootcamp-2025/lang-portal/backend-go/internal/models"
	"github.com/souheilbenslama/free-genai-bootcamp-2025/lang-portal/backend-go/internal/repository"
)

type WordHandler struct {
	wordRepo repository.WordRepository
}

func NewWordHandler(wordRepo repository.WordRepository) *WordHandler {
	return &WordHandler{wordRepo: wordRepo}
}

func (h *WordHandler) GetWord(c *gin.Context) {
	id, err := strconv.Atoi(c.Param("id"))
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid id"})
		return
	}

	word, err := h.wordRepo.GetWord(c.Request.Context(), id)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, word)
}

func (h *WordHandler) ListWords(c *gin.Context) {
	offset, _ := strconv.Atoi(c.DefaultQuery("offset", "0"))
	limit, _ := strconv.Atoi(c.DefaultQuery("limit", "100"))

	words, err := h.wordRepo.ListWords(c.Request.Context(), offset, limit)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"items": words,
		"pagination": gin.H{
			"offset": offset,
			"limit":  limit,
		},
	})
}

func (h *WordHandler) CreateWord(c *gin.Context) {
	var word models.Word
	if err := c.ShouldBindJSON(&word); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	if err := h.wordRepo.CreateWord(c.Request.Context(), &word); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusCreated, word)
}

package routes

import (
	"github.com/gin-gonic/gin"
	"github.com/souheilbenslama/free-genai-bootcamp-2025/lang-portal/backend-go/internal/api/handlers"
)

func SetupRoutes(r *gin.Engine, wordHandler *handlers.WordHandler) {
	api := r.Group("/api")
	{
		// Word routes
		words := api.Group("/words")
		{
			words.GET("", wordHandler.ListWords)
			words.GET("/:id", wordHandler.GetWord)
			words.POST("", wordHandler.CreateWord)
		}

		// TODO: Add routes for:
		// - /api/dashboard/*
		// - /api/groups/*
		// - /api/study-sessions/*
	}
}

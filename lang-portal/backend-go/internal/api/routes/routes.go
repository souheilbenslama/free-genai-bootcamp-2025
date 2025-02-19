package routes

import (
	"github.com/gin-gonic/gin"
	"github.com/souheilbenslama/free-genai-bootcamp-2025/lang-portal/backend-go/internal/api/handlers"
)

func SetupRoutes(
	r *gin.Engine,
	wordHandler *handlers.WordHandler,
	groupHandler *handlers.GroupHandler,
	studyHandler *handlers.StudyHandler,
) {
	api := r.Group("/api")
	{
		// Word routes
		words := api.Group("/words")
		{
			words.GET("", wordHandler.ListWords)
			words.GET("/:id", wordHandler.GetWord)
			words.POST("", wordHandler.CreateWord)
			words.PUT("/:id", wordHandler.UpdateWord)
			words.DELETE("/:id", wordHandler.DeleteWord)
		}

		// Group routes
		groups := api.Group("/groups")
		{
			groups.GET("", groupHandler.GetGroups)
			groups.GET("/:id", groupHandler.GetGroup)
			groups.POST("", groupHandler.CreateGroup)
			groups.PUT("/:id", groupHandler.UpdateGroup)
			groups.DELETE("/:id", groupHandler.DeleteGroup)
			groups.POST("/:id/words", groupHandler.AddWordToGroup)
			groups.DELETE("/:id/words/:word_id", groupHandler.RemoveWordFromGroup)
		}

		// Dashboard routes
		dashboard := api.Group("/dashboard")
		{
			dashboard.GET("/last_study_session", studyHandler.GetLastStudySession)
			dashboard.GET("/study_progress", studyHandler.GetStudyProgress)
			dashboard.GET("/quick_stats", studyHandler.GetQuickStats)
		}

		// Study session routes
		study := api.Group("/study-sessions")
		{
			study.POST("", studyHandler.StartStudySession)
			study.POST("/:session_id/reviews", studyHandler.RecordWordReview)
		}
	}
}

package handlers_test

import (
	"bytes"
	"encoding/json"
	"fmt"
	"net/http"
	"net/http/httptest"

	"github.com/gin-gonic/gin"
	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"
	"github.com/souheilbenslama/free-genai-bootcamp-2025/lang-portal/backend-go/internal/api/handlers"
	"github.com/souheilbenslama/free-genai-bootcamp-2025/lang-portal/backend-go/internal/models"
	"github.com/souheilbenslama/free-genai-bootcamp-2025/lang-portal/backend-go/internal/repository/sqlite"
	"github.com/souheilbenslama/free-genai-bootcamp-2025/lang-portal/backend-go/internal/api/handlers/test"
)

var _ = Describe("GroupHandler", func() {
	var (
		router *gin.Engine
		groupHandler *handlers.GroupHandler
	)

	BeforeEach(func() {
		gin.SetMode(gin.TestMode)
		router = gin.New()
		
		// Initialize with a test database
		db := test.SetupTestDB()
		groupRepo := sqlite.NewGroupRepository(db)
		wordRepo := sqlite.NewWordRepository(db)
		groupHandler = handlers.NewGroupHandler(groupRepo)
		wordHandler := handlers.NewWordHandler(wordRepo)

		// Setup routes
		groups := router.Group("/api/groups")
		{
			groups.GET("", groupHandler.GetGroups)
			groups.GET("/:id", groupHandler.GetGroup)
			groups.POST("", groupHandler.CreateGroup)
			groups.PUT("/:id", groupHandler.UpdateGroup)
			groups.DELETE("/:id", groupHandler.DeleteGroup)
			groups.POST("/:id/words", groupHandler.AddWordToGroup)
			groups.DELETE("/:id/words/:word_id", groupHandler.RemoveWordFromGroup)
		}

		words := router.Group("/api/words")
		{
			words.POST("", wordHandler.CreateWord)
		}
	})

	Describe("POST /api/groups", func() {
		Context("when creating a new group", func() {
			It("succeeds with valid group data", func() {
				group := models.Group{
					Name:        "Basic Vocabulary",
					Description: "Essential German words for beginners",
				}

				jsonValue, err := json.Marshal(group)
				Expect(err).NotTo(HaveOccurred())

				req := httptest.NewRequest("POST", "/api/groups", bytes.NewBuffer(jsonValue))
				req.Header.Set("Content-Type", "application/json")
				
				w := httptest.NewRecorder()
				router.ServeHTTP(w, req)

				Expect(w.Code).To(Equal(http.StatusCreated))

				var response models.Group
				err = json.Unmarshal(w.Body.Bytes(), &response)
				Expect(err).NotTo(HaveOccurred())
				Expect(response.ID).NotTo(BeZero())
				Expect(response.Name).To(Equal(group.Name))
				Expect(response.Description).To(Equal(group.Description))
			})

			It("fails with missing name", func() {
				group := models.Group{
					Description: "Essential German words for beginners",
				}

				jsonValue, err := json.Marshal(group)
				Expect(err).NotTo(HaveOccurred())

				req := httptest.NewRequest("POST", "/api/groups", bytes.NewBuffer(jsonValue))
				req.Header.Set("Content-Type", "application/json")
				
				w := httptest.NewRecorder()
				router.ServeHTTP(w, req)

				Expect(w.Code).To(Equal(http.StatusBadRequest))
			})
		})
	})

	Describe("POST /api/groups/:id/words", func() {
		var (
			createdGroup models.Group
			createdWord models.Word
		)

		BeforeEach(func() {
			// Create a group and a word
			group := models.Group{
				Name:        "Test Group",
				Description: "Test Description",
			}
			word := models.Word{
				German:  "Haus",
				English: "house",
				Parts:   "das",
			}

			// Create group
			jsonValue, err := json.Marshal(group)
			Expect(err).NotTo(HaveOccurred())

			req := httptest.NewRequest("POST", "/api/groups", bytes.NewBuffer(jsonValue))
			req.Header.Set("Content-Type", "application/json")
			w := httptest.NewRecorder()
			router.ServeHTTP(w, req)
			
			err = json.Unmarshal(w.Body.Bytes(), &createdGroup)
			Expect(err).NotTo(HaveOccurred())

			// Create word
			jsonValue, err = json.Marshal(word)
			Expect(err).NotTo(HaveOccurred())

			req = httptest.NewRequest("POST", "/api/words", bytes.NewBuffer(jsonValue))
			req.Header.Set("Content-Type", "application/json")
			w = httptest.NewRecorder()
			router.ServeHTTP(w, req)
			
			err = json.Unmarshal(w.Body.Bytes(), &createdWord)
			Expect(err).NotTo(HaveOccurred())
		})

		Context("when adding a word to a group", func() {
			It("succeeds with valid association", func() {
				url := fmt.Sprintf("/api/groups/%d/words", createdGroup.ID)
				payload := map[string]int{"word_id": createdWord.ID}
				jsonValue, err := json.Marshal(payload)
				Expect(err).NotTo(HaveOccurred())
				req := httptest.NewRequest("POST", url, bytes.NewBuffer(jsonValue))
				req.Header.Set("Content-Type", "application/json")
				w := httptest.NewRecorder()
				router.ServeHTTP(w, req)

				Expect(w.Code).To(Equal(http.StatusOK))
			})

			It("fails with non-existent group", func() {
				url := fmt.Sprintf("/api/groups/%d/words", 999)
				payload := map[string]int{"word_id": createdWord.ID}
				jsonValue, err := json.Marshal(payload)
				Expect(err).NotTo(HaveOccurred())
				req := httptest.NewRequest("POST", url, bytes.NewBuffer(jsonValue))
				req.Header.Set("Content-Type", "application/json")
				w := httptest.NewRecorder()
				router.ServeHTTP(w, req)

				Expect(w.Code).To(Equal(http.StatusNotFound))
			})

			It("fails with non-existent word", func() {
				url := fmt.Sprintf("/api/groups/%d/words", createdGroup.ID)
				payload := map[string]int{"word_id": 999}
				jsonValue, err := json.Marshal(payload)
				Expect(err).NotTo(HaveOccurred())
				req := httptest.NewRequest("POST", url, bytes.NewBuffer(jsonValue))
				req.Header.Set("Content-Type", "application/json")
				w := httptest.NewRecorder()
				router.ServeHTTP(w, req)

				Expect(w.Code).To(Equal(http.StatusNotFound))
			})
		})
	})

	Describe("GET /api/groups/:id", func() {
		var createdGroup models.Group

		BeforeEach(func() {
			// Create a test group
			group := models.Group{
				Name:        "Test Group",
				Description: "Test Description",
			}

			jsonValue, err := json.Marshal(group)
			Expect(err).NotTo(HaveOccurred())

			req := httptest.NewRequest("POST", "/api/groups", bytes.NewBuffer(jsonValue))
			req.Header.Set("Content-Type", "application/json")
			w := httptest.NewRecorder()
			router.ServeHTTP(w, req)
			
			err = json.Unmarshal(w.Body.Bytes(), &createdGroup)
			Expect(err).NotTo(HaveOccurred())
		})

		Context("when getting a group", func() {
			It("succeeds with existing group ID", func() {
				req := httptest.NewRequest("GET", fmt.Sprintf("/api/groups/%d", createdGroup.ID), nil)
				w := httptest.NewRecorder()
				router.ServeHTTP(w, req)

				Expect(w.Code).To(Equal(http.StatusOK))

				var response struct {
					ID          int    `json:"id"`
					Name        string `json:"name"`
					Description string `json:"description"`
					Words       []models.Word `json:"words"`
				}
				err := json.Unmarshal(w.Body.Bytes(), &response)
				Expect(err).NotTo(HaveOccurred())
				Expect(response.ID).To(Equal(createdGroup.ID))
				Expect(response.Name).To(Equal(createdGroup.Name))
				Expect(response.Description).To(Equal(createdGroup.Description))
			})

			It("fails with non-existent group ID", func() {
				req := httptest.NewRequest("GET", "/api/groups/999", nil)
				w := httptest.NewRecorder()
				router.ServeHTTP(w, req)

				Expect(w.Code).To(Equal(http.StatusNotFound))
			})

			It("fails with invalid group ID", func() {
				req := httptest.NewRequest("GET", "/api/groups/invalid", nil)
				w := httptest.NewRecorder()
				router.ServeHTTP(w, req)

				Expect(w.Code).To(Equal(http.StatusBadRequest))
			})
		})
	})
})

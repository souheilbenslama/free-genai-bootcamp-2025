package handlers_test

import (
	"bytes"
	"encoding/json"
	"net/http"
	"net/http/httptest"

	"github.com/gin-gonic/gin"
	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"
	"github.com/souheilbenslama/free-genai-bootcamp-2025/lang-portal/backend-go/internal/api/handlers"
	"github.com/souheilbenslama/free-genai-bootcamp-2025/lang-portal/backend-go/internal/api/handlers/test"
	"github.com/souheilbenslama/free-genai-bootcamp-2025/lang-portal/backend-go/internal/models"
	"github.com/souheilbenslama/free-genai-bootcamp-2025/lang-portal/backend-go/internal/repository/sqlite"
)

var _ = Describe("WordHandler", func() {
	var (
		router *gin.Engine
		wordHandler *handlers.WordHandler
	)

	BeforeEach(func() {
		gin.SetMode(gin.TestMode)
		router = gin.New()
		
		// Initialize with a test database
		db := test.SetupTestDB()
		wordRepo := sqlite.NewWordRepository(db)
		wordHandler = handlers.NewWordHandler(wordRepo)

		// Setup routes
		words := router.Group("/api/words")
		{
			words.GET("", wordHandler.ListWords)
			words.GET("/:id", wordHandler.GetWord)
			words.POST("", wordHandler.CreateWord)
			words.PUT("/:id", wordHandler.UpdateWord)
			words.DELETE("/:id", wordHandler.DeleteWord)
		}
	})

	Describe("POST /api/words", func() {
		Context("when creating a new word", func() {
			It("succeeds with valid word data", func() {
				word := models.Word{
					German:  "Haus",
					English: "house",
					Parts:   "das",
				}

				jsonValue, err := json.Marshal(word)
				Expect(err).NotTo(HaveOccurred())

				req := httptest.NewRequest("POST", "/api/words", bytes.NewBuffer(jsonValue))
				req.Header.Set("Content-Type", "application/json")
				
				w := httptest.NewRecorder()
				router.ServeHTTP(w, req)

				Expect(w.Code).To(Equal(http.StatusOK))

				var response models.Word
				err = json.Unmarshal(w.Body.Bytes(), &response)
				Expect(err).NotTo(HaveOccurred())
				Expect(response.ID).NotTo(BeZero())
				Expect(response.German).To(Equal(word.German))
				Expect(response.English).To(Equal(word.English))
				Expect(response.Parts).To(Equal(word.Parts))
			})

			It("fails with missing required fields", func() {
				word := models.Word{
					German: "Haus",
				}

				jsonValue, err := json.Marshal(word)
				Expect(err).NotTo(HaveOccurred())

				req := httptest.NewRequest("POST", "/api/words", bytes.NewBuffer(jsonValue))
				req.Header.Set("Content-Type", "application/json")
				
				w := httptest.NewRecorder()
				router.ServeHTTP(w, req)

				Expect(w.Code).To(Equal(http.StatusBadRequest))
			})
		})
	})

	Describe("GET /api/words/:id", func() {
		var createdWord models.Word

		BeforeEach(func() {
			// Create a word to test with
			word := models.Word{
				German:  "Haus",
				English: "house",
				Parts:   "das",
			}

			jsonValue, err := json.Marshal(word)
			Expect(err).NotTo(HaveOccurred())

			req := httptest.NewRequest("POST", "/api/words", bytes.NewBuffer(jsonValue))
			req.Header.Set("Content-Type", "application/json")
			w := httptest.NewRecorder()
			router.ServeHTTP(w, req)

			err = json.Unmarshal(w.Body.Bytes(), &createdWord)
			Expect(err).NotTo(HaveOccurred())
		})

		Context("when getting a word", func() {
			It("succeeds with existing word ID", func() {
				req := httptest.NewRequest("GET", "/api/words/1", nil)
				w := httptest.NewRecorder()
				router.ServeHTTP(w, req)

				Expect(w.Code).To(Equal(http.StatusOK))

				var response models.Word
				err := json.Unmarshal(w.Body.Bytes(), &response)
				Expect(err).NotTo(HaveOccurred())
				Expect(response.ID).NotTo(BeZero())
				Expect(response.German).NotTo(BeEmpty())
				Expect(response.English).NotTo(BeEmpty())
			})

			It("fails with non-existent word ID", func() {
				req := httptest.NewRequest("GET", "/api/words/999", nil)
				w := httptest.NewRecorder()
				router.ServeHTTP(w, req)

				Expect(w.Code).To(Equal(http.StatusNotFound))
			})

			It("fails with invalid word ID", func() {
				req := httptest.NewRequest("GET", "/api/words/invalid", nil)
				w := httptest.NewRecorder()
				router.ServeHTTP(w, req)

				Expect(w.Code).To(Equal(http.StatusBadRequest))
			})
		})
	})

	Describe("GET /api/words", func() {
		BeforeEach(func() {
			// Create some test words
			words := []models.Word{
				{German: "Haus", English: "house", Parts: "das"},
				{German: "Auto", English: "car", Parts: "das"},
				{German: "Katze", English: "cat", Parts: "die"},
			}

			for _, word := range words {
				jsonValue, err := json.Marshal(word)
				Expect(err).NotTo(HaveOccurred())

				req := httptest.NewRequest("POST", "/api/words", bytes.NewBuffer(jsonValue))
				req.Header.Set("Content-Type", "application/json")
				w := httptest.NewRecorder()
				router.ServeHTTP(w, req)
				Expect(w.Code).To(Equal(http.StatusOK))
			}
		})

		Context("when listing words", func() {
			It("returns all words without pagination", func() {
				req := httptest.NewRequest("GET", "/api/words", nil)
				w := httptest.NewRecorder()
				router.ServeHTTP(w, req)

				Expect(w.Code).To(Equal(http.StatusOK))

				var response struct {
					Items []models.Word `json:"items"`
				}
				err := json.Unmarshal(w.Body.Bytes(), &response)
				Expect(err).NotTo(HaveOccurred())
				Expect(response.Items).To(HaveLen(3))

				// Verify each word has required fields
				for _, word := range response.Items {
					Expect(word.ID).NotTo(BeZero())
					Expect(word.German).NotTo(BeEmpty())
					Expect(word.English).NotTo(BeEmpty())
				}
			})

			It("returns paginated results", func() {
				req := httptest.NewRequest("GET", "/api/words?limit=2&offset=0", nil)
				w := httptest.NewRecorder()
				router.ServeHTTP(w, req)

				Expect(w.Code).To(Equal(http.StatusOK))

				var response struct {
					Items []models.Word `json:"items"`
				}
				err := json.Unmarshal(w.Body.Bytes(), &response)
				Expect(err).NotTo(HaveOccurred())
				Expect(response.Items).To(HaveLen(2))

				// Verify each word has required fields
				for _, word := range response.Items {
					Expect(word.ID).NotTo(BeZero())
					Expect(word.German).NotTo(BeEmpty())
					Expect(word.English).NotTo(BeEmpty())
				}
			})
		})
	})
})

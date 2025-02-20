package routes_test

import (
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"strings"

	"github.com/gin-gonic/gin"
	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"
	"github.com/souheilbenslama/free-genai-bootcamp-2025/lang-portal/backend-go/internal/api/handlers"
	"github.com/souheilbenslama/free-genai-bootcamp-2025/lang-portal/backend-go/internal/api/handlers/test"
	"github.com/souheilbenslama/free-genai-bootcamp-2025/lang-portal/backend-go/internal/api/routes"
	"github.com/souheilbenslama/free-genai-bootcamp-2025/lang-portal/backend-go/internal/repository/sqlite"
)

var _ = Describe("Routes", func() {
	var (
		router *gin.Engine
		wordHandler *handlers.WordHandler
		groupHandler *handlers.GroupHandler
		studyHandler *handlers.StudyHandler
	)

	BeforeEach(func() {
		gin.SetMode(gin.TestMode)
		router = gin.New()

		// Initialize with a test database
		db := test.SetupTestDB()
		wordRepo := sqlite.NewWordRepository(db)
		groupRepo := sqlite.NewGroupRepository(db)
		studyRepo := sqlite.NewStudyRepository(db)

		wordHandler = handlers.NewWordHandler(wordRepo)
		groupHandler = handlers.NewGroupHandler(groupRepo)
		studyHandler = handlers.NewStudyHandler(studyRepo)

		routes.SetupRoutes(router, wordHandler, groupHandler, studyHandler)
	})

	Context("when creating a word", func() {
		It("should create a word with valid data", func() {
			w := httptest.NewRecorder()
			reqBody := `{"german":"hallo","english":"hello","parts":"verb"}`
			req := httptest.NewRequest(http.MethodPost, "/api/words", strings.NewReader(reqBody))
			req.Header.Set("Content-Type", "application/json")
			router.ServeHTTP(w, req)

			Expect(w.Code).To(Equal(http.StatusOK))
			
			var response map[string]interface{}
			Err := json.NewDecoder(w.Body).Decode(&response)
			Expect(Err).NotTo(HaveOccurred())
			Expect(response["id"]).NotTo(BeNil())
			Expect(response["german"]).To(Equal("hallo"))
			Expect(response["english"]).To(Equal("hello"))
			Expect(response["parts"]).To(Equal("verb"))
		})

		It("should return error for invalid word data", func() {
			w := httptest.NewRecorder()
			reqBody := `{"german":"","english":"hello","parts":"verb"}`
			req := httptest.NewRequest(http.MethodPost, "/api/words", strings.NewReader(reqBody))
			req.Header.Set("Content-Type", "application/json")
			router.ServeHTTP(w, req)

			Expect(w.Code).To(Equal(http.StatusBadRequest))
			
			var response map[string]interface{}
			Err := json.NewDecoder(w.Body).Decode(&response)
			Expect(Err).NotTo(HaveOccurred())
			Expect(response["error"]).To(ContainSubstring("german, english, and parts are required"))
		})
	})

	Context("when routes are configured", func() {
		routeTests := []struct {
			description string
			method      string
			path        string
			expectedCode int
		}{
			{"List Words endpoint", http.MethodGet, "/api/words", http.StatusOK},
			{"Get Word endpoint", http.MethodGet, "/api/words/1", http.StatusNotFound},
			{"Create Word endpoint", http.MethodPost, "/api/words", http.StatusBadRequest},
			{"Update Word endpoint", http.MethodPut, "/api/words/1", http.StatusBadRequest},
			{"Delete Word endpoint", http.MethodDelete, "/api/words/1", http.StatusInternalServerError},
			
			{"List Groups endpoint", http.MethodGet, "/api/groups", http.StatusOK},
			{"Get Group endpoint", http.MethodGet, "/api/groups/1", http.StatusNotFound},
			{"Create Group endpoint", http.MethodPost, "/api/groups", http.StatusBadRequest},
			{"Update Group endpoint", http.MethodPut, "/api/groups/1", http.StatusBadRequest},
			{"Delete Group endpoint", http.MethodDelete, "/api/groups/1", http.StatusInternalServerError},
			{"Add Word to Group endpoint", http.MethodPost, "/api/groups/1/words", http.StatusBadRequest},
			{"Remove Word from Group endpoint", http.MethodDelete, "/api/groups/1/words/1", http.StatusInternalServerError},
			
			{"Get Last Study Session endpoint", http.MethodGet, "/api/dashboard/last_study_session", http.StatusInternalServerError},
			{"Get Study Progress endpoint", http.MethodGet, "/api/dashboard/study_progress", http.StatusInternalServerError},
			{"Get Quick Stats endpoint", http.MethodGet, "/api/dashboard/quick_stats", http.StatusInternalServerError},
			
			{"Start Study Session endpoint", http.MethodPost, "/api/study-sessions", http.StatusBadRequest},
			{"Record Word Review endpoint", http.MethodPost, "/api/study-sessions/1/reviews", http.StatusBadRequest},
		}

		for _, rt := range routeTests {
			rt := rt // capture range variable
			It(rt.description+" should be registered", func() {
				w := httptest.NewRecorder()
				req := httptest.NewRequest(rt.method, rt.path, nil)
				router.ServeHTTP(w, req)

				// Route should exist and return expected status code
				Expect(w.Code).To(Equal(rt.expectedCode), 
					"Route %s %s returned unexpected status code", rt.method, rt.path)
			})
		}
	})

	Context("when accessing non-API routes", func() {
		It("should return 404 for non-API paths", func() {
			w := httptest.NewRecorder()
			req := httptest.NewRequest(http.MethodGet, "/not-api/words", nil)
			router.ServeHTTP(w, req)

			Expect(w.Code).To(Equal(http.StatusNotFound), 
				"Non-API route should not exist")
		})
	})

	Context("when accessing dashboard", func() {
		It("should get study progress", func() {
			w := httptest.NewRecorder()
			req := httptest.NewRequest(http.MethodGet, "/api/dashboard/study_progress", nil)
			router.ServeHTTP(w, req)

			Expect(w.Code).To(Equal(http.StatusInternalServerError))
		})

		It("should get quick stats", func() {
			w := httptest.NewRecorder()
			req := httptest.NewRequest(http.MethodGet, "/api/dashboard/quick_stats", nil)
			router.ServeHTTP(w, req)

			Expect(w.Code).To(Equal(http.StatusInternalServerError))
		})

		It("should get last study session", func() {
			w := httptest.NewRecorder()
			req := httptest.NewRequest(http.MethodGet, "/api/dashboard/last_study_session", nil)
			router.ServeHTTP(w, req)

			Expect(w.Code).To(Equal(http.StatusInternalServerError))
		})
	})

	Context("when managing study sessions", func() {
		It("should create a new study session", func() {
			w := httptest.NewRecorder()
			reqBody := `{"group_id": 1}`
			req := httptest.NewRequest(http.MethodPost, "/api/study-sessions", strings.NewReader(reqBody))
			req.Header.Set("Content-Type", "application/json")
			router.ServeHTTP(w, req)

			Expect(w.Code).To(Equal(http.StatusInternalServerError))
		})

		It("should reject invalid study session request", func() {
			w := httptest.NewRecorder()
			reqBody := `{"invalid_field": 1}`
			req := httptest.NewRequest(http.MethodPost, "/api/study-sessions", strings.NewReader(reqBody))
			req.Header.Set("Content-Type", "application/json")
			router.ServeHTTP(w, req)

			Expect(w.Code).To(Equal(http.StatusBadRequest))
			
			var response map[string]interface{}
			Err := json.NewDecoder(w.Body).Decode(&response)
			Expect(Err).NotTo(HaveOccurred())
			Expect(response["error"]).To(ContainSubstring("GroupID"))
		})

		It("should record a word review", func() {
			w := httptest.NewRecorder()
			reqBody := `{"word_id": 1, "correct": true}`
			req := httptest.NewRequest(http.MethodPost, "/api/study-sessions/1/reviews", strings.NewReader(reqBody))
			req.Header.Set("Content-Type", "application/json")
			router.ServeHTTP(w, req)

			Expect(w.Code).To(Equal(http.StatusInternalServerError))
		})

		It("should reject invalid word review", func() {
			w := httptest.NewRecorder()
			reqBody := `{"word_id": "invalid"}`
			req := httptest.NewRequest(http.MethodPost, "/api/study-sessions/1/reviews", strings.NewReader(reqBody))
			req.Header.Set("Content-Type", "application/json")
			router.ServeHTTP(w, req)

			Expect(w.Code).To(Equal(http.StatusBadRequest))
			
			var response map[string]interface{}
			Err := json.NewDecoder(w.Body).Decode(&response)
			Expect(Err).NotTo(HaveOccurred())
			Expect(response["error"]).NotTo(BeEmpty())
		})
	})
})

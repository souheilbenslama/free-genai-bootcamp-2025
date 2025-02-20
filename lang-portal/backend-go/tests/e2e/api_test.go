package e2e

import (
	"bytes"
	"database/sql"
	"encoding/json"
	"fmt"
	"net/http"
	"os"
	"testing"
	"time"

	"github.com/gin-gonic/gin"
	_ "github.com/mattn/go-sqlite3"
	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"
	"github.com/souheilbenslama/free-genai-bootcamp-2025/lang-portal/backend-go/internal/api/handlers"
	"github.com/souheilbenslama/free-genai-bootcamp-2025/lang-portal/backend-go/internal/api/routes"
	"github.com/souheilbenslama/free-genai-bootcamp-2025/lang-portal/backend-go/internal/models"
	"github.com/souheilbenslama/free-genai-bootcamp-2025/lang-portal/backend-go/internal/repository/sqlite"
)

const testDBPath = "test.db"
const serverAddr = "localhost:8081"
const baseURL = "http://" + serverAddr

var db *sql.DB
var server *http.Server

func TestE2E(t *testing.T) {
	RegisterFailHandler(Fail)
	RunSpecs(t, "E2E API Suite")
}

var _ = BeforeSuite(func() {
	// Set up test database
	setupTestDB()

	// Start the server
	gin.SetMode(gin.TestMode)
	router := gin.New()

	wordRepo := sqlite.NewWordRepository(db)
	groupRepo := sqlite.NewGroupRepository(db)
	studyRepo := sqlite.NewStudyRepository(db)

	wordHandler := handlers.NewWordHandler(wordRepo)
	groupHandler := handlers.NewGroupHandler(groupRepo)
	studyHandler := handlers.NewStudyHandler(studyRepo)

	routes.SetupRoutes(router, wordHandler, groupHandler, studyHandler)

	server = &http.Server{
		Addr:    serverAddr,
		Handler: router,
	}

	go func() {
		if err := server.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			Fail(fmt.Sprintf("Failed to start server: %v", err))
		}
	}()

	// Wait for server to start
	time.Sleep(100 * time.Millisecond)
})

var _ = AfterSuite(func() {
	if server != nil {
		server.Close()
	}
	if db != nil {
		db.Close()
	}
	os.Remove(testDBPath)
})

func setupTestDB() {
	var err error
	db, err = sql.Open("sqlite3", testDBPath)
	Expect(err).NotTo(HaveOccurred())

	// Read and execute migration SQL
	migrationSQL, err := os.ReadFile("../../database/migrations/001_initial_schema.sql")
	Expect(err).NotTo(HaveOccurred())

	_, err = db.Exec(string(migrationSQL))
	Expect(err).NotTo(HaveOccurred())
}

var _ = Describe("API E2E Tests", func() {
	var createdWordID int
	var createdGroupID int
	var studySessionID int

	Context("Word Management Flow", func() {
		It("should create a new word", func() {
			word := models.Word{
				German:  "Apfel",
				English: "apple",
				Parts:   "noun",
			}
			body, err := json.Marshal(word)
			Expect(err).NotTo(HaveOccurred())

			resp, err := http.Post(baseURL+"/api/words", "application/json", bytes.NewReader(body))
			Expect(err).NotTo(HaveOccurred())
			Expect(resp.StatusCode).To(Equal(http.StatusOK))

			var response map[string]interface{}
			err = json.NewDecoder(resp.Body).Decode(&response)
			Expect(err).NotTo(HaveOccurred())
			Expect(response["german"]).To(Equal("Apfel"))
			createdWordID = int(response["id"].(float64))
		})

		It("should get the created word", func() {
			resp, err := http.Get(fmt.Sprintf("%s/api/words/%d", baseURL, createdWordID))
			Expect(err).NotTo(HaveOccurred())
			Expect(resp.StatusCode).To(Equal(http.StatusOK))

			var word models.Word
			err = json.NewDecoder(resp.Body).Decode(&word)
			Expect(err).NotTo(HaveOccurred())
			Expect(word.German).To(Equal("Apfel"))
		})
	})

	Context("Group Management Flow", func() {
		It("should create a new group", func() {
			group := map[string]string{
				"name":        "Fruits",
				"description": "Common fruit names",
			}
			body, err := json.Marshal(group)
			Expect(err).NotTo(HaveOccurred())

			resp, err := http.Post(baseURL+"/api/groups", "application/json", bytes.NewReader(body))
			Expect(err).NotTo(HaveOccurred())
			Expect(resp.StatusCode).To(Equal(http.StatusCreated))

			var response map[string]interface{}
			err = json.NewDecoder(resp.Body).Decode(&response)
			Expect(err).NotTo(HaveOccurred())
			Expect(response["name"]).To(Equal("Fruits"))
			createdGroupID = int(response["id"].(float64))
		})

		It("should add word to group", func() {
			url := fmt.Sprintf("%s/api/groups/%d/words", baseURL, createdGroupID)
			body := fmt.Sprintf(`{"word_id": %d}`, createdWordID)
			
			req, err := http.NewRequest(http.MethodPost, url, bytes.NewReader([]byte(body)))
			Expect(err).NotTo(HaveOccurred())
			req.Header.Set("Content-Type", "application/json")
			
			resp, err := http.DefaultClient.Do(req)
			Expect(err).NotTo(HaveOccurred())
			Expect(resp.StatusCode).To(Equal(http.StatusOK))
		})
	})

	Context("Study Session Flow", func() {
		It("should start a study session", func() {
			body := fmt.Sprintf(`{"group_id": %d}`, createdGroupID)
			resp, err := http.Post(baseURL+"/api/study-sessions", "application/json", bytes.NewReader([]byte(body)))
			Expect(err).NotTo(HaveOccurred())
			Expect(resp.StatusCode).To(Equal(http.StatusCreated))

			var response map[string]interface{}
			err = json.NewDecoder(resp.Body).Decode(&response)
			Expect(err).NotTo(HaveOccurred())
			studySessionID = int(response["id"].(float64))
		})

		It("should record word review", func() {
			url := fmt.Sprintf("%s/api/study-sessions/%d/reviews", baseURL, studySessionID)
			body := fmt.Sprintf(`{"word_id": %d, "correct": true}`, createdWordID)
			
			req, err := http.NewRequest(http.MethodPost, url, bytes.NewReader([]byte(body)))
			Expect(err).NotTo(HaveOccurred())
			req.Header.Set("Content-Type", "application/json")
			
			resp, err := http.DefaultClient.Do(req)
			Expect(err).NotTo(HaveOccurred())
			Expect(resp.StatusCode).To(Equal(http.StatusNoContent))
		})
	})

	Context("Dashboard Flow", func() {
		It("should get study progress", func() {
			resp, err := http.Get(baseURL + "/api/dashboard/study_progress")
			Expect(err).NotTo(HaveOccurred())
			Expect(resp.StatusCode).To(Equal(http.StatusOK))

			var progress models.StudyProgress
			err = json.NewDecoder(resp.Body).Decode(&progress)
			Expect(err).NotTo(HaveOccurred())
			Expect(progress.TotalWordsStudied).To(BeNumerically(">", 0))
		})

		It("should get quick stats", func() {
			resp, err := http.Get(baseURL + "/api/dashboard/quick_stats")
			Expect(err).NotTo(HaveOccurred())
			Expect(resp.StatusCode).To(Equal(http.StatusOK))

			var stats models.DashboardStats
			err = json.NewDecoder(resp.Body).Decode(&stats)
			Expect(err).NotTo(HaveOccurred())
			Expect(stats.TotalStudySessions).To(BeNumerically(">", 0))
		})

		It("should get last study session", func() {
			resp, err := http.Get(baseURL + "/api/dashboard/last_study_session")
			Expect(err).NotTo(HaveOccurred())
			Expect(resp.StatusCode).To(Equal(http.StatusOK))

			var session models.StudySession
			err = json.NewDecoder(resp.Body).Decode(&session)
			Expect(err).NotTo(HaveOccurred())
			Expect(session.ID).To(Equal(studySessionID))
		})
	})

	Context("Cleanup Flow", func() {
		It("should delete word from group", func() {
			url := fmt.Sprintf("%s/api/groups/%d/words/%d", baseURL, createdGroupID, createdWordID)
			req, err := http.NewRequest(http.MethodDelete, url, nil)
			Expect(err).NotTo(HaveOccurred())
			
			resp, err := http.DefaultClient.Do(req)
			Expect(err).NotTo(HaveOccurred())
			Expect(resp.StatusCode).To(Equal(http.StatusOK))
		})

		It("should delete group", func() {
			url := fmt.Sprintf("%s/api/groups/%d", baseURL, createdGroupID)
			req, err := http.NewRequest(http.MethodDelete, url, nil)
			Expect(err).NotTo(HaveOccurred())
			
			resp, err := http.DefaultClient.Do(req)
			Expect(err).NotTo(HaveOccurred())
			Expect(resp.StatusCode).To(Equal(http.StatusOK))
		})

		It("should delete word", func() {
			url := fmt.Sprintf("%s/api/words/%d", baseURL, createdWordID)
			req, err := http.NewRequest(http.MethodDelete, url, nil)
			Expect(err).NotTo(HaveOccurred())
			
			resp, err := http.DefaultClient.Do(req)
			Expect(err).NotTo(HaveOccurred())
			Expect(resp.StatusCode).To(Equal(http.StatusOK))
		})
	})
})

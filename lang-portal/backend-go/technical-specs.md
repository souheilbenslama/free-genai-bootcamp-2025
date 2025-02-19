# Backend server in go Technical specs


## Technical requirements

- the backend server must be written in go
- the database will be sqlite3  
- the api will be built using Gin frameworkwind
- the api will always return json 
- there will be no authentication or authorization
- every thing will for a single user    

## Database Schema

Our database will be a single sqlite database called `words.db` that will be in the root of the project folder of `backend_go`

We have the following tables:
- words - stored vocabulary words
  - id integer
  - german string
  - english string
  - parts json
- words_groups - join table for words and groups many-to-many
  - id integer
  - word_id integer
  - group_id integer
- groups - thematic groups of words
  - id integer
  - name string
- study_sessions - records of study sessions grouping word_review_items
  - id integer
  - group_id integer
  - created_at datetime
  - study_activity_id integer
- study_activities - a specific study activity, linking a study session to group
  - id integer
  - study_session_id integer
  - group_id integer
  - created_at datetime
- word_review_items - a record of word practice, determining if the word was correct or not
  - word_id integer
  - study_session_id integer
  - correct boolean
  - created_at datetime


# APi  endpoints 
- GET /api/dashboard/last_study_session
- GET /api/dashboard/study_progress
- GET /api/dashboard/quick_stats
- GET /api/words
    - pagination with 100 items per page
- GET /api/words/:id
- GET /api/groups
    - pagination with 100 items per page
- GET /api/groups/:id
- GET /api/groups/:id/words
- GET /api/study_activities
- GET /api/study_activities/:id
- GET /api/study_activities/:id/study_sessions
- GET /api/study_sessions
    - pagination with 100 items per page
- GET /api/study_sessions/:id
- GET /api/study_sessions/:id/words
- POST /api/reset_history
- POST /api/full_reset

- POST /api/study_sessions
- POST /api/study_activities
  - required params : group_id , study_activity_id
- POST /api/study_sessions/:id/words/:words_id/review
  - required params : correct

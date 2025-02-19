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

- GET /words
- GET /words/:id
- GET /groups
- GET /groups/:id
- GET /groups/:id/words
- POST /study_sessions
- POST /study_activities

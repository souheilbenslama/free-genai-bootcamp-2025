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
  - correct booleanl
  - created_at datetime


# API Endpoints

## Dashboard Endpoints

### GET /api/dashboard/last_study_session
Returns information about the most recent study session.
```json
{
  "id": 123,
  "created_at": "2025-02-19T08:07:15+01:00",
  "group_id": 1,
  "correct_count": 15,
  "wrong_count": 5
}
```

### GET /api/dashboard/study_progress
Returns study progress statistics.
```json
{
  "total_words_studied": 3,
  "total_available_words": 124,
  "mastery_percentage": 0
}
```

### GET /api/dashboard/quick_stats
Returns overview statistics.
```json
{
  "success_rate": 0.8,
  "total_study_sessions": 4,
  "total_active_groups": 3,
  "study_streak": 4
}
```

## Words Endpoints

### GET /api/words
Returns paginated list of words. Default 100 items per page.
```json
{
  "items": [
    {
      "id": 1,
      "german": "Haus",
      "english": "house",
      "parts": {"article": "das", "plural": "Häuser"},
      "correct_count": 10,
      "wrong_count": 2
    }
  ],
  "total": 1000,
  "page": 1,
  "per_page": 100
}
```

### GET /api/words/:id
Returns details of a specific word.
```json
{
  "id": 1,
  "german": "Haus",
  "english": "house",
  "correct_count": 10,
  "wrong_count": 2,
  "groups": [
    {"id": 1, "name": "Basic Vocabulary"}
  ]
}
```

## Groups Endpoints

### GET /api/groups
Returns paginated list of word groups.
```json
{
  "items": [
    {
      "id": 1,
      "name": "Basic Vocabulary",
      "word_count": 100
    }
  ],
  "total": 25,
  "page": 1,
  "per_page": 100
}
```

### GET /api/groups/:id
Returns details of a specific group.
```json
{
  "id": 1,
  "name": "Basic Vocabulary",
  "word_count": 100,
  "last_studied": "2025-02-19T08:07:15+01:00",
  "accuracy": 0.85
}
```

### GET /api/groups/:id/words
Returns words belonging to a specific group.
```json
{
  "group_id": 1,
  "group_name": "Basic Vocabulary",
  "words": [
    {
      "id": 1,
      "german": "Haus",
      "english": "house",
      "parts": {"article": "das", "plural": "Häuser"}
    }
  ]
}
```

## Study Activities Endpoints

### GET /api/study_activities
Returns list of study activities.
```json
{
  "items": [
    {
      "id": 1,
      "name": "Vocabulary Quiz",
      "thumbnail_url": "https://example.com/thumbnails/vocab-quiz.png",
      "description": "Practice your vocabulary with flashcards"
    }
  ]
}
```

### GET /api/study_activities/:id
Returns details of a specific study activity.
```json
{
  "id": 1,
  "name": "Vocabulary Quiz",
  "thumbnail_url": "https://example.com/thumbnails/vocab-quiz.png",
  "description": "Practice your vocabulary with flashcards",
  "study_sessions": [
    {
      "id": 1,
      "activity_name": "Vocabulary Quiz",
      "group_name": "Basic Vocabulary",
      "start_time": "2025-02-19T08:00:00+01:00",
      "end_time": "2025-02-19T08:15:00+01:00",
      "review_items_count": 20
    }
  ]
}
```

### GET /api/study_activities/:id/study_sessions
Returns study sessions for a specific activity.
```json
{
  "activity_id": 1,
  "sessions": [
    {
      "id": 1,
      "created_at": "2025-02-19T08:07:15+01:00",
      "words_reviewed": 20,
      "correct_count": 15
    }
  ]
}
```

## Study Sessions Endpoints

### GET /api/study_sessions
Returns paginated list of study sessions.
```json
{
  "items": [
    {
      "id": 1,
      "created_at": "2025-02-19T08:07:15+01:00",
      "group_id": 1,
      "words_reviewed": 20,
      "correct_count": 15
    }
  ],
  "total": 100,
  "page": 1,
  "per_page": 100
}
```

### GET /api/study_sessions/:id
Returns details of a specific study session.
```json
{
  "id": 1,
  "created_at": "2025-02-19T08:07:15+01:00",
  "group_id": 1,
  "activity_id": 1,
  "words_reviewed": 20,
  "correct_count": 15,
  "accuracy": 0.75
}
```

### GET /api/study_sessions/:id/words
Returns words reviewed in a specific study session.
```json
{
  "session_id": 1,
  "words": [
    {
      "word_id": 1,
      "german": "Haus",
      "english": "house",
      "parts": {"article": "das", "plural": "Häuser"},
      "correct": true,
      "reviewed_at": "2025-02-19T08:07:15+01:00"
    }
  ]
}
```

## Data Management Endpoints

### POST /api/reset_history
Resets all study history while keeping words and groups.
```json
{
  "success": true,
  "message": "Study history has been reset",
  "deleted_sessions": 50,
  "deleted_activities": 10
}
```

### POST /api/full_reset
Resets entire database including words and groups.
```json
{
  "success": true,
  "message": "Database has been reset to initial state",
  "deleted_words": 1000,
  "deleted_groups": 25,
  "deleted_sessions": 50,
  "deleted_activities": 10
}
```

### POST /api/study_sessions
Creates a new study session.
```json
{
  "id": 1,
  "created_at": "2025-02-19T08:07:15+01:00",
  "group_id": 1,
  "activity_id": 1
}
```

### POST /api/study_activities
Creates a new study activity.
Required params: group_id, study_activity_id
```json
{
  "id": 1,
  "group_id": 1,
  "created_at": "2025-02-19T08:07:15+01:00"
}
```

### POST /api/study_sessions/:id/words/:words_id/review
Records a word review in a study session.
Required params: correct
```json
{
  "success": true,
  "session_id": 1,
  "word_id": 1,
  "correct": true,
  "created_at": "2025-02-19T08:07:15+01:00"
}
```

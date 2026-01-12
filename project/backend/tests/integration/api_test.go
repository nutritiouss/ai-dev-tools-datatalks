package integration

import (
	"database/sql"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/gin-gonic/gin"
	_ "github.com/mattn/go-sqlite3"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"forum-api-wrapper/internal/api"
	"forum-api-wrapper/internal/repository"
	"forum-api-wrapper/internal/service"
)

func setupTestDB(t *testing.T) *sql.DB {
	db, err := sql.Open("sqlite3", ":memory:")
	require.NoError(t, err)

	// Run migrations
	schema := `
	CREATE TABLE forums (
		id INTEGER PRIMARY KEY AUTOINCREMENT,
		name TEXT NOT NULL,
		description TEXT,
		topic_count INTEGER DEFAULT 0,
		post_count INTEGER DEFAULT 0,
		created_at DATETIME DEFAULT CURRENT_TIMESTAMP,
		updated_at DATETIME DEFAULT CURRENT_TIMESTAMP
	);

	CREATE TABLE users (
		id INTEGER PRIMARY KEY AUTOINCREMENT,
		username TEXT NOT NULL UNIQUE,
		post_count INTEGER DEFAULT 0,
		topic_count INTEGER DEFAULT 0,
		registered_at DATETIME DEFAULT CURRENT_TIMESTAMP,
		last_active_at DATETIME
	);

	CREATE TABLE topics (
		id INTEGER PRIMARY KEY AUTOINCREMENT,
		title TEXT NOT NULL,
		forum_id INTEGER NOT NULL REFERENCES forums(id),
		author_id INTEGER NOT NULL REFERENCES users(id),
		reply_count INTEGER DEFAULT 0,
		view_count INTEGER DEFAULT 0,
		last_post_id INTEGER,
		last_post_at DATETIME,
		created_at DATETIME DEFAULT CURRENT_TIMESTAMP,
		updated_at DATETIME DEFAULT CURRENT_TIMESTAMP
	);

	CREATE TABLE posts (
		id INTEGER PRIMARY KEY AUTOINCREMENT,
		topic_id INTEGER NOT NULL REFERENCES topics(id),
		author_id INTEGER NOT NULL REFERENCES users(id),
		content TEXT NOT NULL,
		is_first_post INTEGER DEFAULT 0,
		created_at DATETIME DEFAULT CURRENT_TIMESTAMP,
		updated_at DATETIME DEFAULT CURRENT_TIMESTAMP
	);
	`

	_, err = db.Exec(schema)
	require.NoError(t, err)

	// Insert test data
	_, err = db.Exec(`
		INSERT INTO forums (id, name, description) VALUES (1, 'Test Forum', 'Test Description');
		INSERT INTO users (id, username) VALUES (1, 'testuser');
		INSERT INTO topics (id, title, forum_id, author_id, reply_count, view_count) VALUES (1, 'Test Topic', 1, 1, 0, 0);
		INSERT INTO posts (id, topic_id, author_id, content, is_first_post) VALUES (1, 1, 1, 'Test post content', 1);
	`)
	require.NoError(t, err)

	return db
}

func setupTestServer(t *testing.T) *httptest.Server {
	db := setupTestDB(t)
	repo := repository.NewRepository(db)
	svc := service.NewService(repo)
	handler := api.NewHandler(svc)

	gin.SetMode(gin.TestMode)
	router := gin.New()
	
	apiGroup := router.Group("/api")
	{
		apiGroup.GET("/health", handler.HealthCheck)
		apiGroup.GET("/forums", handler.GetForums)
		apiGroup.GET("/forums/:id", handler.GetForum)
		apiGroup.GET("/topics", handler.GetTopics)
		apiGroup.GET("/topics/:topicId", handler.GetTopic)
		apiGroup.GET("/posts", handler.GetPosts)
		apiGroup.GET("/posts/:postId", handler.GetPost)
		apiGroup.GET("/users", handler.GetUsers)
		apiGroup.GET("/users/:userId", handler.GetUser)
		apiGroup.GET("/search", handler.Search)
	}

	return httptest.NewServer(router)
}

func TestHealthCheck(t *testing.T) {
	server := setupTestServer(t)
	defer server.Close()

	resp, err := http.Get(server.URL + "/api/health")
	require.NoError(t, err)
	defer resp.Body.Close()

	assert.Equal(t, http.StatusOK, resp.StatusCode)

	var response map[string]interface{}
	err = json.NewDecoder(resp.Body).Decode(&response)
	require.NoError(t, err)
	assert.Equal(t, "ok", response["status"])
}

func TestGetForums(t *testing.T) {
	server := setupTestServer(t)
	defer server.Close()

	resp, err := http.Get(server.URL + "/api/forums?page=1&limit=20")
	require.NoError(t, err)
	defer resp.Body.Close()

	assert.Equal(t, http.StatusOK, resp.StatusCode)

	var response struct {
		Forums     []interface{} `json:"forums"`
		Pagination interface{}   `json:"pagination"`
	}
	err = json.NewDecoder(resp.Body).Decode(&response)
	require.NoError(t, err)
	assert.Greater(t, len(response.Forums), 0)
}

func TestGetForumByID(t *testing.T) {
	server := setupTestServer(t)
	defer server.Close()

	resp, err := http.Get(server.URL + "/api/forums/1")
	require.NoError(t, err)
	defer resp.Body.Close()

	assert.Equal(t, http.StatusOK, resp.StatusCode)

	var forum map[string]interface{}
	err = json.NewDecoder(resp.Body).Decode(&forum)
	require.NoError(t, err)
	assert.Equal(t, "Test Forum", forum["name"])
}

func TestGetTopics(t *testing.T) {
	server := setupTestServer(t)
	defer server.Close()

	resp, err := http.Get(server.URL + "/api/topics?page=1&limit=20")
	require.NoError(t, err)
	defer resp.Body.Close()

	assert.Equal(t, http.StatusOK, resp.StatusCode)

	var response struct {
		Topics     []interface{} `json:"topics"`
		Pagination interface{}   `json:"pagination"`
	}
	err = json.NewDecoder(resp.Body).Decode(&response)
	require.NoError(t, err)
	assert.Greater(t, len(response.Topics), 0)
}

func TestGetTopicWithPosts(t *testing.T) {
	db := setupTestDB(t)
	defer db.Close()

	repo := repository.NewRepository(db)
	svc := service.NewService(repo)

	req := httptest.NewRequest("GET", "/api/topics/1", nil)
	topicDetail, err := svc.GetTopic(req.Context(), 1, 1, 20)

	require.NoError(t, err)
	assert.NotNil(t, topicDetail)
	assert.Equal(t, "Test Topic", topicDetail.Title)
	assert.Greater(t, len(topicDetail.Posts), 0)
}

func TestSearch(t *testing.T) {
	db := setupTestDB(t)
	defer db.Close()

	repo := repository.NewRepository(db)
	svc := service.NewService(repo)

	req := httptest.NewRequest("GET", "/api/search?q=test", nil)
	results, err := svc.Search(req.Context(), "test", "all", nil, 1, 20)

	require.NoError(t, err)
	assert.NotNil(t, results)
	assert.Greater(t, results.TotalResults, 0)
}

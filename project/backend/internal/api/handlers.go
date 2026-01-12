package api

import (
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
	"forum-api-wrapper/internal/repository"
	"forum-api-wrapper/internal/service"
)

// Handler handles HTTP requests
type Handler struct {
	service *service.Service
}

// NewHandler creates a new handler instance
func NewHandler(svc *service.Service) *Handler {
	return &Handler{service: svc}
}

// HealthCheck handles GET /health
func (h *Handler) HealthCheck(c *gin.Context) {
	c.JSON(http.StatusOK, gin.H{
		"status":    "ok",
		"timestamp": "2024-01-01T00:00:00Z",
	})
}

// GetForums handles GET /forums
func (h *Handler) GetForums(c *gin.Context) {
	page, limit := parsePagination(c)

	response, err := h.service.GetForums(c.Request.Context(), page, limit)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, response)
}

// GetForum handles GET /forums/:id
func (h *Handler) GetForum(c *gin.Context) {
	id, err := strconv.Atoi(c.Param("id"))
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid forum ID"})
		return
	}

	forum, err := h.service.GetForum(c.Request.Context(), id)
	if err != nil {
		if err.Error() == "forum not found" {
			c.JSON(http.StatusNotFound, gin.H{"error": "forum not found"})
			return
		}
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, forum)
}

// GetTopics handles GET /topics
func (h *Handler) GetTopics(c *gin.Context) {
	page, limit := parsePagination(c)
	sort := c.DefaultQuery("sort", "newest")

	filter := repository.TopicFilter{Sort: sort}

	if forumIDStr := c.Query("forumId"); forumIDStr != "" {
		forumID, err := strconv.Atoi(forumIDStr)
		if err == nil {
			filter.ForumID = &forumID
		}
	}

	response, err := h.service.GetTopics(c.Request.Context(), filter, page, limit)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, response)
}

// GetTopic handles GET /topics/:id
func (h *Handler) GetTopic(c *gin.Context) {
	id, err := strconv.Atoi(c.Param("topicId"))
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid topic ID"})
		return
	}

	page, limit := parsePagination(c)

	response, err := h.service.GetTopic(c.Request.Context(), id, page, limit)
	if err != nil {
		if err.Error() == "topic not found" {
			c.JSON(http.StatusNotFound, gin.H{"error": "topic not found"})
			return
		}
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, response)
}

// GetPosts handles GET /posts
func (h *Handler) GetPosts(c *gin.Context) {
	page, limit := parsePagination(c)

	filter := repository.PostFilter{}

	if topicIDStr := c.Query("topicId"); topicIDStr != "" {
		topicID, err := strconv.Atoi(topicIDStr)
		if err == nil {
			filter.TopicID = &topicID
		}
	}

	if userIDStr := c.Query("userId"); userIDStr != "" {
		userID, err := strconv.Atoi(userIDStr)
		if err == nil {
			filter.UserID = &userID
		}
	}

	response, err := h.service.GetPosts(c.Request.Context(), filter, page, limit)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, response)
}

// GetPost handles GET /posts/:id
func (h *Handler) GetPost(c *gin.Context) {
	id, err := strconv.Atoi(c.Param("postId"))
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid post ID"})
		return
	}

	post, err := h.service.GetPost(c.Request.Context(), id)
	if err != nil {
		if err.Error() == "post not found" {
			c.JSON(http.StatusNotFound, gin.H{"error": "post not found"})
			return
		}
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, post)
}

// GetUsers handles GET /users
func (h *Handler) GetUsers(c *gin.Context) {
	page, limit := parsePagination(c)

	response, err := h.service.GetUsers(c.Request.Context(), page, limit)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, response)
}

// GetUser handles GET /users/:id
func (h *Handler) GetUser(c *gin.Context) {
	id, err := strconv.Atoi(c.Param("userId"))
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid user ID"})
		return
	}

	user, err := h.service.GetUser(c.Request.Context(), id)
	if err != nil {
		if err.Error() == "user not found" {
			c.JSON(http.StatusNotFound, gin.H{"error": "user not found"})
			return
		}
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, user)
}

// Search handles GET /search
func (h *Handler) Search(c *gin.Context) {
	query := c.Query("q")
	if query == "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "query parameter 'q' is required"})
		return
	}

	searchType := c.DefaultQuery("type", "all")
	page, limit := parsePagination(c)

	var forumID *int
	if forumIDStr := c.Query("forumId"); forumIDStr != "" {
		id, err := strconv.Atoi(forumIDStr)
		if err == nil {
			forumID = &id
		}
	}

	response, err := h.service.Search(c.Request.Context(), query, searchType, forumID, page, limit)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, response)
}

// parsePagination parses page and limit from query parameters
func parsePagination(c *gin.Context) (int, int) {
	page := 1
	limit := 20

	if pageStr := c.Query("page"); pageStr != "" {
		if p, err := strconv.Atoi(pageStr); err == nil && p > 0 {
			page = p
		}
	}

	if limitStr := c.Query("limit"); limitStr != "" {
		if l, err := strconv.Atoi(limitStr); err == nil && l > 0 && l <= 100 {
			limit = l
		}
	}

	return page, limit
}

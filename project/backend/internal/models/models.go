package models

import "time"

// Forum represents a forum category
type Forum struct {
	ID          int       `json:"id" db:"id"`
	Name        string    `json:"name" db:"name"`
	Description string    `json:"description" db:"description"`
	TopicCount  int       `json:"topicCount" db:"topic_count"`
	PostCount int       `json:"postCount" db:"post_count"`
	CreatedAt   time.Time `json:"createdAt" db:"created_at"`
	UpdatedAt   time.Time `json:"updatedAt" db:"updated_at"`
}

// Topic represents a forum topic
type Topic struct {
	ID          int       `json:"id" db:"id"`
	Title       string    `json:"title" db:"title"`
	ForumID     int       `json:"forumId" db:"forum_id"`
	ForumName   string    `json:"forumName" db:"forum_name"`
	AuthorID    int       `json:"authorId" db:"author_id"`
	AuthorName  string    `json:"authorName" db:"author_name"`
	ReplyCount  int       `json:"replyCount" db:"reply_count"`
	ViewCount   int       `json:"viewCount" db:"view_count"`
	LastPostID  *int      `json:"lastPostId,omitempty" db:"last_post_id"`
	LastPostAt  *time.Time `json:"lastPostAt,omitempty" db:"last_post_at"`
	CreatedAt   time.Time `json:"createdAt" db:"created_at"`
	UpdatedAt   time.Time `json:"updatedAt" db:"updated_at"`
}

// Post represents a forum post
type Post struct {
	ID          int       `json:"id" db:"id"`
	TopicID     int       `json:"topicId" db:"topic_id"`
	TopicTitle  string    `json:"topicTitle" db:"topic_title"`
	AuthorID    int       `json:"authorId" db:"author_id"`
	AuthorName  string    `json:"authorName" db:"author_name"`
	Content     string    `json:"content" db:"content"`
	IsFirstPost bool      `json:"isFirstPost" db:"is_first_post"`
	CreatedAt   time.Time `json:"createdAt" db:"created_at"`
	UpdatedAt   time.Time `json:"updatedAt" db:"updated_at"`
}

// User represents a forum user
type User struct {
	ID          int        `json:"id" db:"id"`
	Username    string     `json:"username" db:"username"`
	PostCount   int        `json:"postCount" db:"post_count"`
	TopicCount  int        `json:"topicCount" db:"topic_count"`
	RegisteredAt time.Time `json:"registeredAt" db:"registered_at"`
	LastActiveAt *time.Time `json:"lastActiveAt,omitempty" db:"last_active_at"`
}

// Pagination represents pagination metadata
type Pagination struct {
	Page      int  `json:"page"`
	Limit     int  `json:"limit"`
	Total     int  `json:"total"`
	TotalPages int `json:"totalPages"`
	HasNext   bool `json:"hasNext"`
	HasPrev   bool `json:"hasPrev"`
}

// CalculatePagination calculates pagination metadata
func CalculatePagination(page, limit, total int) Pagination {
	totalPages := (total + limit - 1) / limit
	if totalPages == 0 {
		totalPages = 1
	}
	
	return Pagination{
		Page:      page,
		Limit:     limit,
		Total:     total,
		TotalPages: totalPages,
		HasNext:   page < totalPages,
		HasPrev:   page > 1,
	}
}

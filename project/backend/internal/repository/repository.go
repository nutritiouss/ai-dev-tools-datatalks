package repository

import (
	"context"
	"database/sql"
	"fmt"
	"forum-api-wrapper/internal/models"
	"time"
)

// Repository defines the database operations interface
type Repository interface {
	// Forums
	GetForums(ctx context.Context, page, limit int) ([]models.Forum, int, error)
	GetForumByID(ctx context.Context, id int) (*models.Forum, error)

	// Topics
	GetTopics(ctx context.Context, filter TopicFilter, page, limit int) ([]models.Topic, int, error)
	GetTopicByID(ctx context.Context, id int) (*models.Topic, error)
	GetTopicPosts(ctx context.Context, topicID int, page, limit int) ([]models.Post, int, error)

	// Posts
	GetPosts(ctx context.Context, filter PostFilter, page, limit int) ([]models.Post, int, error)
	GetPostByID(ctx context.Context, id int) (*models.Post, error)

	// Users
	GetUsers(ctx context.Context, page, limit int) ([]models.User, int, error)
	GetUserByID(ctx context.Context, id int) (*models.User, error)

	// Search
	Search(ctx context.Context, query string, searchType string, forumID *int, page, limit int) (SearchResults, int, error)
}

// TopicFilter filters for topic queries
type TopicFilter struct {
	ForumID *int
	Sort    string
}

// PostFilter filters for post queries
type PostFilter struct {
	TopicID *int
	UserID  *int
}

// SearchResults contains search results
type SearchResults struct {
	Topics []models.Topic
	Posts  []models.Post
	Users  []models.User
}

// DBRepository implements Repository using database/sql
type DBRepository struct {
	db *sql.DB
}

// NewRepository creates a new repository instance
func NewRepository(db *sql.DB) Repository {
	return &DBRepository{db: db}
}

// GetForums retrieves forums with pagination
func (r *DBRepository) GetForums(ctx context.Context, page, limit int) ([]models.Forum, int, error) {
	offset := (page - 1) * limit

	// Get total count
	var total int
	err := r.db.QueryRowContext(ctx, "SELECT COUNT(*) FROM forums").Scan(&total)
	if err != nil {
		return nil, 0, fmt.Errorf("failed to count forums: %w", err)
	}

	// Get forums
	query := `
		SELECT id, name, description, topic_count, post_count, created_at, updated_at
		FROM forums
		ORDER BY name
		LIMIT $1 OFFSET $2
	`

	rows, err := r.db.QueryContext(ctx, query, limit, offset)
	if err != nil {
		return nil, 0, fmt.Errorf("failed to query forums: %w", err)
	}
	defer rows.Close()

	var forums []models.Forum
	for rows.Next() {
		var f models.Forum
		err := rows.Scan(&f.ID, &f.Name, &f.Description, &f.TopicCount, &f.PostCount, &f.CreatedAt, &f.UpdatedAt)
		if err != nil {
			return nil, 0, fmt.Errorf("failed to scan forum: %w", err)
		}
		forums = append(forums, f)
	}

	return forums, total, nil
}

// GetForumByID retrieves a forum by ID
func (r *DBRepository) GetForumByID(ctx context.Context, id int) (*models.Forum, error) {
	query := `
		SELECT id, name, description, topic_count, post_count, created_at, updated_at
		FROM forums
		WHERE id = $1
	`

	var f models.Forum
	err := r.db.QueryRowContext(ctx, query, id).Scan(
		&f.ID, &f.Name, &f.Description, &f.TopicCount, &f.PostCount, &f.CreatedAt, &f.UpdatedAt,
	)
	if err == sql.ErrNoRows {
		return nil, nil
	}
	if err != nil {
		return nil, fmt.Errorf("failed to get forum: %w", err)
	}

	return &f, nil
}

// GetTopics retrieves topics with filtering and pagination
func (r *DBRepository) GetTopics(ctx context.Context, filter TopicFilter, page, limit int) ([]models.Topic, int, error) {
	offset := (page - 1) * limit

	// Build query
	whereClause := "1=1"
	args := []interface{}{}
	argPos := 1

	if filter.ForumID != nil {
		whereClause += fmt.Sprintf(" AND forum_id = $%d", argPos)
		args = append(args, *filter.ForumID)
		argPos++
	}

	// Get total count
	countQuery := fmt.Sprintf("SELECT COUNT(*) FROM topics WHERE %s", whereClause)
	var total int
	err := r.db.QueryRowContext(ctx, countQuery, args...).Scan(&total)
	if err != nil {
		return nil, 0, fmt.Errorf("failed to count topics: %w", err)
	}

	// Determine sort order
	orderBy := "created_at DESC"
	switch filter.Sort {
	case "oldest":
		orderBy = "created_at ASC"
	case "most_replies":
		orderBy = "reply_count DESC"
	case "most_views":
		orderBy = "view_count DESC"
	}

	// Get topics
	query := fmt.Sprintf(`
		SELECT 
			t.id, t.title, t.forum_id, f.name as forum_name,
			t.author_id, u.username as author_name,
			t.reply_count, t.view_count,
			t.last_post_id, t.last_post_at,
			t.created_at, t.updated_at
		FROM topics t
		JOIN forums f ON t.forum_id = f.id
		JOIN users u ON t.author_id = u.id
		WHERE %s
		ORDER BY %s
		LIMIT $%d OFFSET $%d
	`, whereClause, orderBy, argPos, argPos+1)

	args = append(args, limit, offset)
	rows, err := r.db.QueryContext(ctx, query, args...)
	if err != nil {
		return nil, 0, fmt.Errorf("failed to query topics: %w", err)
	}
	defer rows.Close()

	var topics []models.Topic
	for rows.Next() {
		var t models.Topic
		var lastPostID sql.NullInt64
		var lastPostAt sql.NullTime

		err := rows.Scan(
			&t.ID, &t.Title, &t.ForumID, &t.ForumName,
			&t.AuthorID, &t.AuthorName,
			&t.ReplyCount, &t.ViewCount,
			&lastPostID, &lastPostAt,
			&t.CreatedAt, &t.UpdatedAt,
		)
		if err != nil {
			return nil, 0, fmt.Errorf("failed to scan topic: %w", err)
		}

		if lastPostID.Valid {
			id := int(lastPostID.Int64)
			t.LastPostID = &id
		}
		if lastPostAt.Valid {
			t.LastPostAt = &lastPostAt.Time
		}

		topics = append(topics, t)
	}

	return topics, total, nil
}

// GetTopicByID retrieves a topic by ID
func (r *DBRepository) GetTopicByID(ctx context.Context, id int) (*models.Topic, error) {
	query := `
		SELECT 
			t.id, t.title, t.forum_id, f.name as forum_name,
			t.author_id, u.username as author_name,
			t.reply_count, t.view_count,
			t.last_post_id, t.last_post_at,
			t.created_at, t.updated_at
		FROM topics t
		JOIN forums f ON t.forum_id = f.id
		JOIN users u ON t.author_id = u.id
		WHERE t.id = $1
	`

	var t models.Topic
	var lastPostID sql.NullInt64
	var lastPostAt sql.NullTime

	err := r.db.QueryRowContext(ctx, query, id).Scan(
		&t.ID, &t.Title, &t.ForumID, &t.ForumName,
		&t.AuthorID, &t.AuthorName,
		&t.ReplyCount, &t.ViewCount,
		&lastPostID, &lastPostAt,
		&t.CreatedAt, &t.UpdatedAt,
	)
	if err == sql.ErrNoRows {
		return nil, nil
	}
	if err != nil {
		return nil, fmt.Errorf("failed to get topic: %w", err)
	}

	if lastPostID.Valid {
		id := int(lastPostID.Int64)
		t.LastPostID = &id
	}
	if lastPostAt.Valid {
		t.LastPostAt = &lastPostAt.Time
	}

	return &t, nil
}

// GetTopicPosts retrieves posts for a topic
func (r *DBRepository) GetTopicPosts(ctx context.Context, topicID int, page, limit int) ([]models.Post, int, error) {
	offset := (page - 1) * limit

	// Get total count
	var total int
	err := r.db.QueryRowContext(ctx, "SELECT COUNT(*) FROM posts WHERE topic_id = $1", topicID).Scan(&total)
	if err != nil {
		return nil, 0, fmt.Errorf("failed to count posts: %w", err)
	}

	// Get posts
	query := `
		SELECT 
			p.id, p.topic_id, t.title as topic_title,
			p.author_id, u.username as author_name,
			p.content, p.is_first_post,
			p.created_at, p.updated_at
		FROM posts p
		JOIN topics t ON p.topic_id = t.id
		JOIN users u ON p.author_id = u.id
		WHERE p.topic_id = $1
		ORDER BY p.created_at ASC
		LIMIT $2 OFFSET $3
	`

	rows, err := r.db.QueryContext(ctx, query, topicID, limit, offset)
	if err != nil {
		return nil, 0, fmt.Errorf("failed to query posts: %w", err)
	}
	defer rows.Close()

	var posts []models.Post
	for rows.Next() {
		var p models.Post
		err := rows.Scan(
			&p.ID, &p.TopicID, &p.TopicTitle,
			&p.AuthorID, &p.AuthorName,
			&p.Content, &p.IsFirstPost,
			&p.CreatedAt, &p.UpdatedAt,
		)
		if err != nil {
			return nil, 0, fmt.Errorf("failed to scan post: %w", err)
		}
		posts = append(posts, p)
	}

	return posts, total, nil
}

// GetPosts retrieves posts with filtering and pagination
func (r *DBRepository) GetPosts(ctx context.Context, filter PostFilter, page, limit int) ([]models.Post, int, error) {
	offset := (page - 1) * limit

	whereClause := "1=1"
	args := []interface{}{}
	argPos := 1

	if filter.TopicID != nil {
		whereClause += fmt.Sprintf(" AND topic_id = $%d", argPos)
		args = append(args, *filter.TopicID)
		argPos++
	}
	if filter.UserID != nil {
		whereClause += fmt.Sprintf(" AND author_id = $%d", argPos)
		args = append(args, *filter.UserID)
		argPos++
	}

	// Get total count
	countQuery := fmt.Sprintf("SELECT COUNT(*) FROM posts WHERE %s", whereClause)
	var total int
	err := r.db.QueryRowContext(ctx, countQuery, args...).Scan(&total)
	if err != nil {
		return nil, 0, fmt.Errorf("failed to count posts: %w", err)
	}

	// Get posts
	query := fmt.Sprintf(`
		SELECT 
			p.id, p.topic_id, t.title as topic_title,
			p.author_id, u.username as author_name,
			p.content, p.is_first_post,
			p.created_at, p.updated_at
		FROM posts p
		JOIN topics t ON p.topic_id = t.id
		JOIN users u ON p.author_id = u.id
		WHERE %s
		ORDER BY p.created_at DESC
		LIMIT $%d OFFSET $%d
	`, whereClause, argPos, argPos+1)

	args = append(args, limit, offset)
	rows, err := r.db.QueryContext(ctx, query, args...)
	if err != nil {
		return nil, 0, fmt.Errorf("failed to query posts: %w", err)
	}
	defer rows.Close()

	var posts []models.Post
	for rows.Next() {
		var p models.Post
		err := rows.Scan(
			&p.ID, &p.TopicID, &p.TopicTitle,
			&p.AuthorID, &p.AuthorName,
			&p.Content, &p.IsFirstPost,
			&p.CreatedAt, &p.UpdatedAt,
		)
		if err != nil {
			return nil, 0, fmt.Errorf("failed to scan post: %w", err)
		}
		posts = append(posts, p)
	}

	return posts, total, nil
}

// GetPostByID retrieves a post by ID
func (r *DBRepository) GetPostByID(ctx context.Context, id int) (*models.Post, error) {
	query := `
		SELECT 
			p.id, p.topic_id, t.title as topic_title,
			p.author_id, u.username as author_name,
			p.content, p.is_first_post,
			p.created_at, p.updated_at
		FROM posts p
		JOIN topics t ON p.topic_id = t.id
		JOIN users u ON p.author_id = u.id
		WHERE p.id = $1
	`

	var p models.Post
	err := r.db.QueryRowContext(ctx, query, id).Scan(
		&p.ID, &p.TopicID, &p.TopicTitle,
		&p.AuthorID, &p.AuthorName,
		&p.Content, &p.IsFirstPost,
		&p.CreatedAt, &p.UpdatedAt,
	)
	if err == sql.ErrNoRows {
		return nil, nil
	}
	if err != nil {
		return nil, fmt.Errorf("failed to get post: %w", err)
	}

	return &p, nil
}

// GetUsers retrieves users with pagination
func (r *DBRepository) GetUsers(ctx context.Context, page, limit int) ([]models.User, int, error) {
	offset := (page - 1) * limit

	// Get total count
	var total int
	err := r.db.QueryRowContext(ctx, "SELECT COUNT(*) FROM users").Scan(&total)
	if err != nil {
		return nil, 0, fmt.Errorf("failed to count users: %w", err)
	}

	// Get users
	query := `
		SELECT id, username, post_count, topic_count, registered_at, last_active_at
		FROM users
		ORDER BY username
		LIMIT $1 OFFSET $2
	`

	rows, err := r.db.QueryContext(ctx, query, limit, offset)
	if err != nil {
		return nil, 0, fmt.Errorf("failed to query users: %w", err)
	}
	defer rows.Close()

	var users []models.User
	for rows.Next() {
		var u models.User
		var lastActiveAt sql.NullTime

		err := rows.Scan(
			&u.ID, &u.Username, &u.PostCount, &u.TopicCount,
			&u.RegisteredAt, &lastActiveAt,
		)
		if err != nil {
			return nil, 0, fmt.Errorf("failed to scan user: %w", err)
		}

		if lastActiveAt.Valid {
			u.LastActiveAt = &lastActiveAt.Time
		}

		users = append(users, u)
	}

	return users, total, nil
}

// GetUserByID retrieves a user by ID
func (r *DBRepository) GetUserByID(ctx context.Context, id int) (*models.User, error) {
	query := `
		SELECT id, username, post_count, topic_count, registered_at, last_active_at
		FROM users
		WHERE id = $1
	`

	var u models.User
	var lastActiveAt sql.NullTime

	err := r.db.QueryRowContext(ctx, query, id).Scan(
		&u.ID, &u.Username, &u.PostCount, &u.TopicCount,
		&u.RegisteredAt, &lastActiveAt,
	)
	if err == sql.ErrNoRows {
		return nil, nil
	}
	if err != nil {
		return nil, fmt.Errorf("failed to get user: %w", err)
	}

	if lastActiveAt.Valid {
		u.LastActiveAt = &lastActiveAt.Time
	}

	return &u, nil
}

// Search performs a full-text search across topics, posts, and users
func (r *DBRepository) Search(ctx context.Context, query string, searchType string, forumID *int, page, limit int) (SearchResults, int, error) {
	offset := (page - 1) * limit
	results := SearchResults{}

	// Build search conditions
	searchPattern := "%" + query + "%"
	args := []interface{}{searchPattern}
	argPos := 2

	forumClause := ""
	if forumID != nil {
		forumClause = fmt.Sprintf(" AND t.forum_id = $%d", argPos)
		args = append(args, *forumID)
		argPos++
	}

	var totalResults int

	// Search topics
	if searchType == "all" || searchType == "topics" {
		querySQL := fmt.Sprintf(`
			SELECT COUNT(*) FROM topics t
			WHERE (t.title ILIKE $1 OR EXISTS (
				SELECT 1 FROM posts p WHERE p.topic_id = t.id AND p.content ILIKE $1
			)) %s
		`, forumClause)
		err := r.db.QueryRowContext(ctx, querySQL, args...).Scan(&totalResults)
		if err == nil {
			searchQuery := fmt.Sprintf(`
				SELECT DISTINCT
					t.id, t.title, t.forum_id, f.name as forum_name,
					t.author_id, u.username as author_name,
					t.reply_count, t.view_count,
					t.last_post_id, t.last_post_at,
					t.created_at, t.updated_at
				FROM topics t
				JOIN forums f ON t.forum_id = f.id
				JOIN users u ON t.author_id = u.id
				WHERE (t.title ILIKE $1 OR EXISTS (
					SELECT 1 FROM posts p WHERE p.topic_id = t.id AND p.content ILIKE $1
				)) %s
				ORDER BY t.created_at DESC
				LIMIT $%d OFFSET $%d
			`, forumClause, argPos, argPos+1)

			searchArgs := append(args, limit, offset)
			rows, err := r.db.QueryContext(ctx, searchQuery, searchArgs...)
			if err == nil {
				defer rows.Close()
				for rows.Next() {
					var t models.Topic
					var lastPostID sql.NullInt64
					var lastPostAt sql.NullTime

					scanArgs := []interface{}{
						&t.ID, &t.Title, &t.ForumID, &t.ForumName,
						&t.AuthorID, &t.AuthorName,
						&t.ReplyCount, &t.ViewCount,
						&lastPostID, &lastPostAt,
						&t.CreatedAt, &t.UpdatedAt,
					}
					if err := rows.Scan(scanArgs...); err == nil {
						if lastPostID.Valid {
							id := int(lastPostID.Int64)
							t.LastPostID = &id
						}
						if lastPostAt.Valid {
							t.LastPostAt = &lastPostAt.Time
						}
						results.Topics = append(results.Topics, t)
					}
				}
			}
		}
	}

	// Search posts
	if searchType == "all" || searchType == "posts" {
		postForumClause := ""
		if forumID != nil {
			postForumClause = fmt.Sprintf(" AND t.forum_id = $%d", argPos)
		}
		querySQL := fmt.Sprintf(`
			SELECT COUNT(*) FROM posts p
			JOIN topics t ON p.topic_id = t.id
			WHERE p.content ILIKE $1 %s
		`, postForumClause)
		var postTotal int
		err := r.db.QueryRowContext(ctx, querySQL, args...).Scan(&postTotal)
		if err == nil {
			searchQuery := fmt.Sprintf(`
				SELECT 
					p.id, p.topic_id, t.title as topic_title,
					p.author_id, u.username as author_name,
					p.content, p.is_first_post,
					p.created_at, p.updated_at
				FROM posts p
				JOIN topics t ON p.topic_id = t.id
				JOIN users u ON p.author_id = u.id
				WHERE p.content ILIKE $1 %s
				ORDER BY p.created_at DESC
				LIMIT $%d OFFSET $%d
			`, postForumClause, argPos, argPos+1)

			searchArgs := append(args, limit, offset)
			rows, err := r.db.QueryContext(ctx, searchQuery, searchArgs...)
			if err == nil {
				defer rows.Close()
				for rows.Next() {
					var p models.Post
					if err := rows.Scan(
						&p.ID, &p.TopicID, &p.TopicTitle,
						&p.AuthorID, &p.AuthorName,
						&p.Content, &p.IsFirstPost,
						&p.CreatedAt, &p.UpdatedAt,
					); err == nil {
						results.Posts = append(results.Posts, p)
					}
				}
			}
		}
		totalResults += postTotal
	}

	// Search users
	if searchType == "all" || searchType == "users" {
		querySQL := "SELECT COUNT(*) FROM users WHERE username ILIKE $1"
		var userTotal int
		err := r.db.QueryRowContext(ctx, querySQL, args[0]).Scan(&userTotal)
		if err == nil {
			searchQuery := `
				SELECT id, username, post_count, topic_count, registered_at, last_active_at
				FROM users
				WHERE username ILIKE $1
				ORDER BY username
				LIMIT $2 OFFSET $3
			`
			rows, err := r.db.QueryContext(ctx, searchQuery, args[0], limit, offset)
			if err == nil {
				defer rows.Close()
				for rows.Next() {
					var u models.User
					var lastActiveAt sql.NullTime
					if err := rows.Scan(
						&u.ID, &u.Username, &u.PostCount, &u.TopicCount,
						&u.RegisteredAt, &lastActiveAt,
					); err == nil {
						if lastActiveAt.Valid {
							u.LastActiveAt = &lastActiveAt.Time
						}
						results.Users = append(results.Users, u)
					}
				}
			}
		}
		totalResults += userTotal
	}

	return results, totalResults, nil
}

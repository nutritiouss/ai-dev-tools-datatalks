package service

import (
	"context"
	"fmt"
	"forum-api-wrapper/internal/models"
	"forum-api-wrapper/internal/repository"
)

// Service provides business logic for the API
type Service struct {
	repo repository.Repository
}

// NewService creates a new service instance
func NewService(repo repository.Repository) *Service {
	return &Service{repo: repo}
}

// GetForums retrieves forums with pagination
func (s *Service) GetForums(ctx context.Context, page, limit int) (*ForumListResponse, error) {
	forums, total, err := s.repo.GetForums(ctx, page, limit)
	if err != nil {
		return nil, fmt.Errorf("failed to get forums: %w", err)
	}

	return &ForumListResponse{
		Forums:     forums,
		Pagination: models.CalculatePagination(page, limit, total),
	}, nil
}

// GetForum retrieves a forum by ID
func (s *Service) GetForum(ctx context.Context, id int) (*models.Forum, error) {
	forum, err := s.repo.GetForumByID(ctx, id)
	if err != nil {
		return nil, fmt.Errorf("failed to get forum: %w", err)
	}
	if forum == nil {
		return nil, fmt.Errorf("forum not found")
	}
	return forum, nil
}

// GetTopics retrieves topics with filtering and pagination
func (s *Service) GetTopics(ctx context.Context, filter repository.TopicFilter, page, limit int) (*TopicListResponse, error) {
	topics, total, err := s.repo.GetTopics(ctx, filter, page, limit)
	if err != nil {
		return nil, fmt.Errorf("failed to get topics: %w", err)
	}

	return &TopicListResponse{
		Topics:     topics,
		Pagination: models.CalculatePagination(page, limit, total),
	}, nil
}

// GetTopic retrieves a topic by ID with its posts
func (s *Service) GetTopic(ctx context.Context, id int, page, limit int) (*TopicDetailResponse, error) {
	topic, err := s.repo.GetTopicByID(ctx, id)
	if err != nil {
		return nil, fmt.Errorf("failed to get topic: %w", err)
	}
	if topic == nil {
		return nil, fmt.Errorf("topic not found")
	}

	posts, total, err := s.repo.GetTopicPosts(ctx, id, page, limit)
	if err != nil {
		return nil, fmt.Errorf("failed to get topic posts: %w", err)
	}

	return &TopicDetailResponse{
		Topic:          *topic,
		Posts:          posts,
		PostPagination: models.CalculatePagination(page, limit, total),
	}, nil
}

// GetPosts retrieves posts with filtering and pagination
func (s *Service) GetPosts(ctx context.Context, filter repository.PostFilter, page, limit int) (*PostListResponse, error) {
	posts, total, err := s.repo.GetPosts(ctx, filter, page, limit)
	if err != nil {
		return nil, fmt.Errorf("failed to get posts: %w", err)
	}

	return &PostListResponse{
		Posts:      posts,
		Pagination: models.CalculatePagination(page, limit, total),
	}, nil
}

// GetPost retrieves a post by ID
func (s *Service) GetPost(ctx context.Context, id int) (*models.Post, error) {
	post, err := s.repo.GetPostByID(ctx, id)
	if err != nil {
		return nil, fmt.Errorf("failed to get post: %w", err)
	}
	if post == nil {
		return nil, fmt.Errorf("post not found")
	}
	return post, nil
}

// GetUsers retrieves users with pagination
func (s *Service) GetUsers(ctx context.Context, page, limit int) (*UserListResponse, error) {
	users, total, err := s.repo.GetUsers(ctx, page, limit)
	if err != nil {
		return nil, fmt.Errorf("failed to get users: %w", err)
	}

	return &UserListResponse{
		Users:      users,
		Pagination: models.CalculatePagination(page, limit, total),
	}, nil
}

// GetUser retrieves a user by ID
func (s *Service) GetUser(ctx context.Context, id int) (*models.User, error) {
	user, err := s.repo.GetUserByID(ctx, id)
	if err != nil {
		return nil, fmt.Errorf("failed to get user: %w", err)
	}
	if user == nil {
		return nil, fmt.Errorf("user not found")
	}
	return user, nil
}

// Search performs a search across topics, posts, and users
func (s *Service) Search(ctx context.Context, query string, searchType string, forumID *int, page, limit int) (*SearchResponse, error) {
	results, total, err := s.repo.Search(ctx, query, searchType, forumID, page, limit)
	if err != nil {
		return nil, fmt.Errorf("failed to search: %w", err)
	}

	return &SearchResponse{
		Results: SearchResults{
			Topics: results.Topics,
			Posts:  results.Posts,
			Users:  results.Users,
		},
		Pagination:  models.CalculatePagination(page, limit, total),
		Query:       query,
		TotalResults: total,
	}, nil
}

// Response types
type ForumListResponse struct {
	Forums     []models.Forum `json:"forums"`
	Pagination models.Pagination `json:"pagination"`
}

type TopicListResponse struct {
	Topics     []models.Topic `json:"topics"`
	Pagination models.Pagination `json:"pagination"`
}

type TopicDetailResponse struct {
	models.Topic
	Posts          []models.Post `json:"posts"`
	PostPagination models.Pagination `json:"postPagination"`
}

type PostListResponse struct {
	Posts      []models.Post `json:"posts"`
	Pagination models.Pagination `json:"pagination"`
}

type UserListResponse struct {
	Users      []models.User `json:"users"`
	Pagination models.Pagination `json:"pagination"`
}

type SearchResults struct {
	Topics []models.Topic `json:"topics"`
	Posts  []models.Post  `json:"posts"`
	Users  []models.User   `json:"users"`
}

type SearchResponse struct {
	Results      SearchResults     `json:"results"`
	Pagination   models.Pagination `json:"pagination"`
	Query        string            `json:"query"`
	TotalResults int               `json:"totalResults"`
}

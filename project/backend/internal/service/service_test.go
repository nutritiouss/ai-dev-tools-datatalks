package service

import (
	"context"
	"forum-api-wrapper/internal/models"
	"forum-api-wrapper/internal/repository"
	"testing"
	"time"
)

// mockRepository is a mock implementation of repository.Repository
type mockRepository struct {
	forums []models.Forum
	topics []models.Topic
	posts  []models.Post
	users  []models.User
}

func (m *mockRepository) GetForums(ctx context.Context, page, limit int) ([]models.Forum, int, error) {
	return m.forums, len(m.forums), nil
}

func (m *mockRepository) GetForumByID(ctx context.Context, id int) (*models.Forum, error) {
	for _, f := range m.forums {
		if f.ID == id {
			return &f, nil
		}
	}
	return nil, nil
}

func (m *mockRepository) GetTopics(ctx context.Context, filter repository.TopicFilter, page, limit int) ([]models.Topic, int, error) {
	return m.topics, len(m.topics), nil
}

func (m *mockRepository) GetTopicByID(ctx context.Context, id int) (*models.Topic, error) {
	for _, t := range m.topics {
		if t.ID == id {
			return &t, nil
		}
	}
	return nil, nil
}

func (m *mockRepository) GetTopicPosts(ctx context.Context, topicID int, page, limit int) ([]models.Post, int, error) {
	var topicPosts []models.Post
	for _, p := range m.posts {
		if p.TopicID == topicID {
			topicPosts = append(topicPosts, p)
		}
	}
	return topicPosts, len(topicPosts), nil
}

func (m *mockRepository) GetPosts(ctx context.Context, filter repository.PostFilter, page, limit int) ([]models.Post, int, error) {
	return m.posts, len(m.posts), nil
}

func (m *mockRepository) GetPostByID(ctx context.Context, id int) (*models.Post, error) {
	for _, p := range m.posts {
		if p.ID == id {
			return &p, nil
		}
	}
	return nil, nil
}

func (m *mockRepository) GetUsers(ctx context.Context, page, limit int) ([]models.User, int, error) {
	return m.users, len(m.users), nil
}

func (m *mockRepository) GetUserByID(ctx context.Context, id int) (*models.User, error) {
	for _, u := range m.users {
		if u.ID == id {
			return &u, nil
		}
	}
	return nil, nil
}

func (m *mockRepository) Search(ctx context.Context, query string, searchType string, forumID *int, page, limit int) (repository.SearchResults, int, error) {
	return repository.SearchResults{
		Topics: m.topics,
		Posts:  m.posts,
		Users:  m.users,
	}, len(m.topics) + len(m.posts) + len(m.users), nil
}

func TestService_GetForums(t *testing.T) {
	mockRepo := &mockRepository{
		forums: []models.Forum{
			{ID: 1, Name: "Test Forum", Description: "Test Description"},
		},
	}
	svc := NewService(mockRepo)

	ctx := context.Background()
	response, err := svc.GetForums(ctx, 1, 20)

	if err != nil {
		t.Fatalf("Expected no error, got %v", err)
	}

	if len(response.Forums) != 1 {
		t.Fatalf("Expected 1 forum, got %d", len(response.Forums))
	}

	if response.Forums[0].Name != "Test Forum" {
		t.Errorf("Expected forum name 'Test Forum', got '%s'", response.Forums[0].Name)
	}
}

func TestService_GetForum_NotFound(t *testing.T) {
	mockRepo := &mockRepository{
		forums: []models.Forum{},
	}
	svc := NewService(mockRepo)

	ctx := context.Background()
	_, err := svc.GetForum(ctx, 999)

	if err == nil {
		t.Fatal("Expected error for non-existent forum")
	}

	if err.Error() != "forum not found" {
		t.Errorf("Expected 'forum not found' error, got '%s'", err.Error())
	}
}

func TestService_GetTopics(t *testing.T) {
	now := time.Now()
	mockRepo := &mockRepository{
		topics: []models.Topic{
			{
				ID:        1,
				Title:     "Test Topic",
				ForumID:   1,
				AuthorID:  1,
				CreatedAt: now,
			},
		},
	}
	svc := NewService(mockRepo)

	ctx := context.Background()
	response, err := svc.GetTopics(ctx, repository.TopicFilter{}, 1, 20)

	if err != nil {
		t.Fatalf("Expected no error, got %v", err)
	}

	if len(response.Topics) != 1 {
		t.Fatalf("Expected 1 topic, got %d", len(response.Topics))
	}

	if response.Topics[0].Title != "Test Topic" {
		t.Errorf("Expected topic title 'Test Topic', got '%s'", response.Topics[0].Title)
	}
}

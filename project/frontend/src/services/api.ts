// Centralized API client for all backend communication

const API_BASE_URL = import.meta.env.VITE_API_BASE_URL || '/api';

export interface Pagination {
  page: number;
  limit: number;
  total: number;
  totalPages: number;
  hasNext: boolean;
  hasPrev: boolean;
}

export interface Forum {
  id: number;
  name: string;
  description: string;
  topicCount: number;
  postCount: number;
  createdAt: string;
  updatedAt: string;
}

export interface Topic {
  id: number;
  title: string;
  forumId: number;
  forumName: string;
  authorId: number;
  authorName: string;
  replyCount: number;
  viewCount: number;
  lastPostId: number;
  lastPostAt: string;
  createdAt: string;
  updatedAt: string;
}

export interface Post {
  id: number;
  topicId: number;
  topicTitle: string;
  authorId: number;
  authorName: string;
  content: string;
  isFirstPost: boolean;
  createdAt: string;
  updatedAt: string;
}

export interface User {
  id: number;
  username: string;
  postCount: number;
  topicCount: number;
  registeredAt: string;
  lastActiveAt: string;
}

export interface ForumListResponse {
  forums: Forum[];
  pagination: Pagination;
}

export interface TopicListResponse {
  topics: Topic[];
  pagination: Pagination;
}

export interface TopicDetail extends Topic {
  posts: Post[];
  postPagination: Pagination;
}

export interface PostListResponse {
  posts: Post[];
  pagination: Pagination;
}

export interface UserListResponse {
  users: User[];
  pagination: Pagination;
}

export interface UserDetail extends User {
  recentTopics: Topic[];
  recentPosts: Post[];
}

export interface SearchResponse {
  results: {
    topics: Topic[];
    posts: Post[];
    users: User[];
  };
  pagination: Pagination;
  query: string;
  totalResults: number;
}

class ApiClient {
  private baseUrl: string;

  constructor(baseUrl: string = API_BASE_URL) {
    this.baseUrl = baseUrl;
  }

  private async request<T>(
    endpoint: string,
    options: RequestInit = {}
  ): Promise<T> {
    const url = `${this.baseUrl}${endpoint}`;
    const response = await fetch(url, {
      ...options,
      headers: {
        'Content-Type': 'application/json',
        ...options.headers,
      },
    });

    if (!response.ok) {
      const error = await response.json().catch(() => ({
        error: `HTTP ${response.status}: ${response.statusText}`,
      }));
      throw new Error(error.error || 'An error occurred');
    }

    return response.json();
  }

  // Health check
  async healthCheck(): Promise<{ status: string; timestamp: string }> {
    return this.request('/health');
  }

  // Forums
  async getForums(page: number = 1, limit: number = 20): Promise<ForumListResponse> {
    return this.request(`/forums?page=${page}&limit=${limit}`);
  }

  async getForum(forumId: number): Promise<Forum> {
    return this.request(`/forums/${forumId}`);
  }

  // Topics
  async getTopics(params: {
    forumId?: number;
    page?: number;
    limit?: number;
    sort?: 'newest' | 'oldest' | 'most_replies' | 'most_views';
  } = {}): Promise<TopicListResponse> {
    const queryParams = new URLSearchParams();
    if (params.forumId) queryParams.append('forumId', params.forumId.toString());
    if (params.page) queryParams.append('page', params.page.toString());
    if (params.limit) queryParams.append('limit', params.limit.toString());
    if (params.sort) queryParams.append('sort', params.sort);
    
    return this.request(`/topics?${queryParams.toString()}`);
  }

  async getTopic(topicId: number, page: number = 1, limit: number = 20): Promise<TopicDetail> {
    return this.request(`/topics/${topicId}?page=${page}&limit=${limit}`);
  }

  // Posts
  async getPosts(params: {
    topicId?: number;
    userId?: number;
    page?: number;
    limit?: number;
  } = {}): Promise<PostListResponse> {
    const queryParams = new URLSearchParams();
    if (params.topicId) queryParams.append('topicId', params.topicId.toString());
    if (params.userId) queryParams.append('userId', params.userId.toString());
    if (params.page) queryParams.append('page', params.page.toString());
    if (params.limit) queryParams.append('limit', params.limit.toString());
    
    return this.request(`/posts?${queryParams.toString()}`);
  }

  async getPost(postId: number): Promise<Post> {
    return this.request(`/posts/${postId}`);
  }

  // Users
  async getUsers(page: number = 1, limit: number = 20): Promise<UserListResponse> {
    return this.request(`/users?page=${page}&limit=${limit}`);
  }

  async getUser(userId: number): Promise<UserDetail> {
    return this.request(`/users/${userId}`);
  }

  // Search
  async search(params: {
    q: string;
    type?: 'all' | 'topics' | 'posts' | 'users';
    forumId?: number;
    page?: number;
    limit?: number;
  }): Promise<SearchResponse> {
    const queryParams = new URLSearchParams();
    queryParams.append('q', params.q);
    if (params.type) queryParams.append('type', params.type);
    if (params.forumId) queryParams.append('forumId', params.forumId.toString());
    if (params.page) queryParams.append('page', params.page.toString());
    if (params.limit) queryParams.append('limit', params.limit.toString());
    
    return this.request(`/search?${queryParams.toString()}`);
  }
}

export const apiClient = new ApiClient();

import { useState, useEffect } from 'react';
import { useSearchParams, Link } from 'react-router-dom';
import { apiClient, SearchResponse, Topic, Post, User } from '../services/api';
import './SearchResults.css';

export function SearchResults() {
  const [searchParams] = useSearchParams();
  const query = searchParams.get('q') || '';
  const [results, setResults] = useState<SearchResponse | null>(null);
  const [loading, setLoading] = useState(true);
  const [error, setError] = useState<string | null>(null);
  const [activeTab, setActiveTab] = useState<'all' | 'topics' | 'posts' | 'users'>('all');

  useEffect(() => {
    if (query) {
      performSearch();
    }
  }, [query]);

  const performSearch = async () => {
    try {
      setLoading(true);
      setError(null);
      const data = await apiClient.search({
        q: query,
        type: activeTab === 'all' ? undefined : activeTab,
      });
      setResults(data);
    } catch (err) {
      setError(err instanceof Error ? err.message : 'Failed to search');
    } finally {
      setLoading(false);
    }
  };

  const formatDate = (dateString: string) => {
    return new Date(dateString).toLocaleDateString('en-US', {
      year: 'numeric',
      month: 'short',
      day: 'numeric',
    });
  };

  if (loading) {
    return <div className="loading">Searching...</div>;
  }

  if (error) {
    return <div className="error">Error: {error}</div>;
  }

  if (!results) {
    return null;
  }

  const hasResults = 
    results.results.topics.length > 0 ||
    results.results.posts.length > 0 ||
    results.results.users.length > 0;

  return (
    <div className="search-results">
      <h2>Search Results for "{query}"</h2>
      <div className="results-summary">
        Found {results.totalResults} result{results.totalResults !== 1 ? 's' : ''}
      </div>

      <div className="tabs">
        <button
          className={activeTab === 'all' ? 'active' : ''}
          onClick={() => setActiveTab('all')}
        >
          All ({results.totalResults})
        </button>
        <button
          className={activeTab === 'topics' ? 'active' : ''}
          onClick={() => setActiveTab('topics')}
        >
          Topics ({results.results.topics.length})
        </button>
        <button
          className={activeTab === 'posts' ? 'active' : ''}
          onClick={() => setActiveTab('posts')}
        >
          Posts ({results.results.posts.length})
        </button>
        <button
          className={activeTab === 'users' ? 'active' : ''}
          onClick={() => setActiveTab('users')}
        >
          Users ({results.results.users.length})
        </button>
      </div>

      {!hasResults ? (
        <div className="empty">No results found</div>
      ) : (
        <div className="results-content">
          {(activeTab === 'all' || activeTab === 'topics') && results.results.topics.length > 0 && (
            <div className="results-section">
              <h3>Topics</h3>
              {results.results.topics.map((topic: Topic) => (
                <div key={topic.id} className="result-item">
                  <Link to={`/topics/${topic.id}`} className="result-title">
                    {topic.title}
                  </Link>
                  <div className="result-meta">
                    <span>{topic.forumName}</span>
                    <span>by {topic.authorName}</span>
                    <span>{formatDate(topic.createdAt)}</span>
                  </div>
                </div>
              ))}
            </div>
          )}

          {(activeTab === 'all' || activeTab === 'posts') && results.results.posts.length > 0 && (
            <div className="results-section">
              <h3>Posts</h3>
              {results.results.posts.map((post: Post) => (
                <div key={post.id} className="result-item">
                  <Link to={`/topics/${post.topicId}`} className="result-title">
                    {post.topicTitle}
                  </Link>
                  <div className="result-meta">
                    <span>by {post.authorName}</span>
                    <span>{formatDate(post.createdAt)}</span>
                  </div>
                  <div className="result-preview" dangerouslySetInnerHTML={{ __html: post.content.substring(0, 200) + '...' }} />
                </div>
              ))}
            </div>
          )}

          {(activeTab === 'all' || activeTab === 'users') && results.results.users.length > 0 && (
            <div className="results-section">
              <h3>Users</h3>
              {results.results.users.map((user: User) => (
                <div key={user.id} className="result-item">
                  <Link to={`/users/${user.id}`} className="result-title">
                    {user.username}
                  </Link>
                  <div className="result-meta">
                    <span>{user.postCount} posts</span>
                    <span>{user.topicCount} topics</span>
                    <span>Joined {formatDate(user.registeredAt)}</span>
                  </div>
                </div>
              ))}
            </div>
          )}
        </div>
      )}
    </div>
  );
}

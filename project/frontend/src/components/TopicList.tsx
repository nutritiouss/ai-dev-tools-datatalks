import { useState, useEffect } from 'react';
import { Link } from 'react-router-dom';
import { apiClient, Topic, Pagination } from '../services/api';
import './TopicList.css';

interface TopicListProps {
  forumId?: number;
  sort?: 'newest' | 'oldest' | 'most_replies' | 'most_views';
}

export function TopicList({ forumId, sort = 'newest' }: TopicListProps) {
  const [topics, setTopics] = useState<Topic[]>([]);
  const [pagination, setPagination] = useState<Pagination | null>(null);
  const [loading, setLoading] = useState(true);
  const [error, setError] = useState<string | null>(null);
  const [currentPage, setCurrentPage] = useState(1);

  useEffect(() => {
    loadTopics();
  }, [forumId, sort, currentPage]);

  const loadTopics = async () => {
    try {
      setLoading(true);
      setError(null);
      const response = await apiClient.getTopics({
        forumId,
        page: currentPage,
        limit: 20,
        sort,
      });
      setTopics(response.topics);
      setPagination(response.pagination);
    } catch (err) {
      setError(err instanceof Error ? err.message : 'Failed to load topics');
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

  if (loading && topics.length === 0) {
    return <div className="loading">Loading topics...</div>;
  }

  if (error) {
    return <div className="error">Error: {error}</div>;
  }

  return (
    <div className="topic-list">
      <div className="topic-list-header">
        <h2>Topics</h2>
        {pagination && (
          <div className="pagination-info">
            Page {pagination.page} of {pagination.totalPages} ({pagination.total} total)
          </div>
        )}
      </div>

      {topics.length === 0 ? (
        <div className="empty">No topics found</div>
      ) : (
        <>
          <div className="topics">
            {topics.map((topic) => (
              <div key={topic.id} className="topic-item">
                <div className="topic-main">
                  <Link to={`/topics/${topic.id}`} className="topic-title">
                    {topic.title}
                  </Link>
                  <div className="topic-meta">
                    <span className="forum-name">{topic.forumName}</span>
                    <span className="author">by {topic.authorName}</span>
                    <span className="date">{formatDate(topic.createdAt)}</span>
                  </div>
                </div>
                <div className="topic-stats">
                  <div className="stat">
                    <span className="stat-value">{topic.replyCount}</span>
                    <span className="stat-label">replies</span>
                  </div>
                  <div className="stat">
                    <span className="stat-value">{topic.viewCount}</span>
                    <span className="stat-label">views</span>
                  </div>
                  {topic.lastPostAt && (
                    <div className="stat last-post">
                      <span className="stat-label">Last post:</span>
                      <span className="stat-value">{formatDate(topic.lastPostAt)}</span>
                    </div>
                  )}
                </div>
              </div>
            ))}
          </div>

          {pagination && (
            <div className="pagination">
              <button
                onClick={() => setCurrentPage((p) => Math.max(1, p - 1))}
                disabled={!pagination.hasPrev}
              >
                Previous
              </button>
              <span>
                Page {pagination.page} of {pagination.totalPages}
              </span>
              <button
                onClick={() => setCurrentPage((p) => p + 1)}
                disabled={!pagination.hasNext}
              >
                Next
              </button>
            </div>
          )}
        </>
      )}
    </div>
  );
}

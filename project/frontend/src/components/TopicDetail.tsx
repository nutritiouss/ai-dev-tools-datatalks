import { useState, useEffect } from 'react';
import { useParams, Link } from 'react-router-dom';
import { apiClient, TopicDetail, Post } from '../services/api';
import './TopicDetail.css';

export function TopicDetail() {
  const { topicId } = useParams<{ topicId: string }>();
  const [topic, setTopic] = useState<TopicDetail | null>(null);
  const [loading, setLoading] = useState(true);
  const [error, setError] = useState<string | null>(null);
  const [currentPage, setCurrentPage] = useState(1);

  useEffect(() => {
    if (topicId) {
      loadTopic();
    }
  }, [topicId, currentPage]);

  const loadTopic = async () => {
    if (!topicId) return;
    
    try {
      setLoading(true);
      setError(null);
      const data = await apiClient.getTopic(parseInt(topicId), currentPage, 20);
      setTopic(data);
    } catch (err) {
      setError(err instanceof Error ? err.message : 'Failed to load topic');
    } finally {
      setLoading(false);
    }
  };

  const formatDate = (dateString: string) => {
    return new Date(dateString).toLocaleString('en-US', {
      year: 'numeric',
      month: 'short',
      day: 'numeric',
      hour: '2-digit',
      minute: '2-digit',
    });
  };

  if (loading && !topic) {
    return <div className="loading">Loading topic...</div>;
  }

  if (error) {
    return <div className="error">Error: {error}</div>;
  }

  if (!topic) {
    return <div className="empty">Topic not found</div>;
  }

  return (
    <div className="topic-detail">
      <div className="topic-header">
        <Link to="/" className="back-link">← Back to topics</Link>
        <h1>{topic.title}</h1>
        <div className="topic-info">
          <span>Forum: {topic.forumName}</span>
          <span>Author: {topic.authorName}</span>
          <span>Created: {formatDate(topic.createdAt)}</span>
          <span>{topic.replyCount} replies</span>
          <span>{topic.viewCount} views</span>
        </div>
      </div>

      <div className="posts">
        {topic.posts.map((post) => (
          <div key={post.id} className={`post ${post.isFirstPost ? 'first-post' : ''}`}>
            <div className="post-header">
              <div className="post-author">
                <strong>{post.authorName}</strong>
                <span className="post-date">{formatDate(post.createdAt)}</span>
              </div>
            </div>
            <div 
              className="post-content"
              dangerouslySetInnerHTML={{ __html: post.content }}
            />
          </div>
        ))}
      </div>

      {topic.postPagination && (
        <div className="pagination">
          <button
            onClick={() => setCurrentPage((p) => Math.max(1, p - 1))}
            disabled={!topic.postPagination.hasPrev}
          >
            Previous
          </button>
          <span>
            Page {topic.postPagination.page} of {topic.postPagination.totalPages}
          </span>
          <button
            onClick={() => setCurrentPage((p) => p + 1)}
            disabled={!topic.postPagination.hasNext}
          >
            Next
          </button>
        </div>
      )}
    </div>
  );
}

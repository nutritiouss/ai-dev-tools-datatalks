import { useState } from 'react'
import { useNavigate } from 'react-router-dom'
import './SessionManager.css'

function SessionManager() {
  const [sessionId, setSessionId] = useState('')
  const navigate = useNavigate()

  const createSession = async () => {
    try {
      const response = await fetch('/api/sessions', {
        method: 'POST',
        headers: {
          'Content-Type': 'application/json',
        },
      })
      const data = await response.json()
      navigate(`/session/${data.sessionId}`)
    } catch (error) {
      console.error('Failed to create session:', error)
    }
  }

  const joinSession = () => {
    if (sessionId.trim()) {
      navigate(`/session/${sessionId.trim()}`)
    }
  }

  return (
    <div className="session-manager">
      <div className="session-manager-container">
        <h1>Coding Interview Platform</h1>
        <p className="subtitle">Collaborative code editor for technical interviews</p>
        
        <div className="session-actions">
          <button onClick={createSession} className="btn btn-primary">
            Create New Session
          </button>
          
          <div className="divider">or</div>
          
          <div className="join-session">
            <input
              type="text"
              placeholder="Enter session ID"
              value={sessionId}
              onChange={(e) => setSessionId(e.target.value)}
              onKeyPress={(e) => e.key === 'Enter' && joinSession()}
              className="session-input"
            />
            <button onClick={joinSession} className="btn btn-secondary">
              Join Session
            </button>
          </div>
        </div>
      </div>
    </div>
  )
}

export default SessionManager








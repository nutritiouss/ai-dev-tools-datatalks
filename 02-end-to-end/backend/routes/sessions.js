import express from 'express'
import { v4 as uuidv4 } from 'uuid'

const router = express.Router()

// In-memory storage for sessions
const sessions = new Map()

// Create a new session
router.post('/sessions', (req, res) => {
  const sessionId = uuidv4()
  sessions.set(sessionId, {
    sessionId,
    code: '// Write your code here\n',
    language: 'javascript',
    createdAt: new Date().toISOString(),
  })
  res.json({ sessionId })
})

// Get session info
router.get('/sessions/:id', (req, res) => {
  const { id } = req.params
  const session = sessions.get(id)
  if (!session) {
    return res.status(404).json({ error: 'Session not found' })
  }
  res.json(session)
})

export default router


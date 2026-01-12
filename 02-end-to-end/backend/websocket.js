import { v4 as uuidv4 } from 'uuid'

// In-memory storage for sessions
const sessions = new Map()

// Track users per session
const sessionUsers = new Map()

export function setupWebSocket(io) {
  io.on('connection', (socket) => {
    console.log('Client connected:', socket.id)

    socket.on('join-session', ({ sessionId }) => {
      // Initialize session if it doesn't exist
      if (!sessions.has(sessionId)) {
        sessions.set(sessionId, {
          sessionId,
          code: '// Write your code here\n',
          language: 'javascript',
          createdAt: new Date().toISOString(),
        })
      }

      // Join the socket room for this session
      socket.join(sessionId)

      // Track user in session
      if (!sessionUsers.has(sessionId)) {
        sessionUsers.set(sessionId, new Set())
      }
      sessionUsers.get(sessionId).add(socket.id)

      // Send current session state to the new user
      const session = sessions.get(sessionId)
      socket.emit('session-joined', {
        code: session.code,
        language: session.language,
      })

      // Broadcast updated user count
      const users = Array.from(sessionUsers.get(sessionId))
      io.to(sessionId).emit('users-update', users)

      console.log(`User ${socket.id} joined session ${sessionId}`)
    })

    socket.on('code-change', ({ sessionId, code }) => {
      // Update session code
      if (sessions.has(sessionId)) {
        const session = sessions.get(sessionId)
        session.code = code
        sessions.set(sessionId, session)

        // Broadcast to all other users in the session
        socket.to(sessionId).emit('code-change', {
          sessionId,
          code,
        })
      }
    })

    socket.on('disconnect', () => {
      console.log('Client disconnected:', socket.id)

      // Remove user from all sessions
      for (const [sessionId, users] of sessionUsers.entries()) {
        if (users.has(socket.id)) {
          users.delete(socket.id)
          // Broadcast updated user count
          const usersArray = Array.from(users)
          io.to(sessionId).emit('users-update', usersArray)

          // Clean up empty sessions after a delay
          if (users.size === 0) {
            setTimeout(() => {
              if (sessionUsers.get(sessionId)?.size === 0) {
                sessionUsers.delete(sessionId)
                sessions.delete(sessionId)
                console.log(`Cleaned up empty session ${sessionId}`)
              }
            }, 300000) // 5 minutes
          }
        }
      }
    })
  })
}








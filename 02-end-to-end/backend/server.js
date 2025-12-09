import express from 'express'
import { createServer } from 'http'
import { Server } from 'socket.io'
import cors from 'cors'
import { fileURLToPath } from 'url'
import { dirname, join } from 'path'
import sessionRoutes from './routes/sessions.js'
import { setupWebSocket } from './websocket.js'

const __filename = fileURLToPath(import.meta.url)
const __dirname = dirname(__filename)

const app = express()
const httpServer = createServer(app)

// Determine CORS origin based on environment
const corsOrigin = process.env.NODE_ENV === 'production' 
  ? '*' 
  : 'http://localhost:3000'

const io = new Server(httpServer, {
  cors: {
    origin: corsOrigin,
    methods: ['GET', 'POST'],
  },
})

app.use(cors())
app.use(express.json())

// API routes
app.use('/api', sessionRoutes)

// Health check
app.get('/health', (req, res) => {
  res.json({ status: 'ok' })
})

// Setup WebSocket
setupWebSocket(io)

// Serve static files in production
if (process.env.NODE_ENV === 'production') {
  app.use(express.static(join(__dirname, '../public')))
  
  // Serve index.html for all non-API routes
  app.get('*', (req, res) => {
    if (!req.path.startsWith('/api') && !req.path.startsWith('/socket.io')) {
      res.sendFile(join(__dirname, '../public/index.html'))
    }
  })
}

const PORT = process.env.PORT || 5000

httpServer.listen(PORT, () => {
  console.log(`Server running on port ${PORT}`)
})


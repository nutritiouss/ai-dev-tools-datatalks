import request from 'supertest'
import { createServer } from 'http'
import express from 'express'
import { Server } from 'socket.io'
import { Client } from 'socket.io-client'
import cors from 'cors'
import sessionRoutes from '../routes/sessions.js'
import { setupWebSocket } from '../websocket.js'

describe('Integration Tests', () => {
  let app, httpServer, io

  beforeAll((done) => {
    app = express()
    httpServer = createServer(app)
    io = new Server(httpServer, {
      cors: {
        origin: '*',
        methods: ['GET', 'POST'],
      },
    })

    app.use(cors())
    app.use(express.json())
    app.use('/api', sessionRoutes)
    setupWebSocket(io)

    httpServer.listen(0, () => {
      done()
    })
  })

  afterAll((done) => {
    io.close()
    httpServer.close(done)
  })

  describe('API Endpoints', () => {
    test('POST /api/sessions should create a new session', async () => {
      const response = await request(app)
        .post('/api/sessions')
        .expect(200)

      expect(response.body).toHaveProperty('sessionId')
      expect(response.body.sessionId).toBeTruthy()
    })

    test('GET /api/sessions/:id should return session info', async () => {
      // First create a session
      const createResponse = await request(app)
        .post('/api/sessions')
        .expect(200)

      const sessionId = createResponse.body.sessionId

      // Then get it
      const getResponse = await request(app)
        .get(`/api/sessions/${sessionId}`)
        .expect(200)

      expect(getResponse.body).toHaveProperty('sessionId', sessionId)
      expect(getResponse.body).toHaveProperty('code')
    })

    test('GET /api/sessions/:id should return 404 for non-existent session', async () => {
      await request(app)
        .get('/api/sessions/non-existent-id')
        .expect(404)
    })
  })

  describe('WebSocket Integration', () => {
    let client1, client2
    const port = httpServer.address().port

    afterEach((done) => {
      if (client1) client1.close()
      if (client2) client2.close()
      done()
    })

    test('Clients can connect and join a session', (done) => {
      client1 = new Client(`http://localhost:${port}`)
      
      client1.on('connect', () => {
        client1.emit('join-session', { sessionId: 'test-session' })
      })

      client1.on('session-joined', (data) => {
        expect(data).toHaveProperty('code')
        done()
      })
    })

    test('Code changes are broadcast to all clients in a session', (done) => {
      const sessionId = 'test-session-2'
      let codeChangeReceived = false

      client1 = new Client(`http://localhost:${port}`)
      client2 = new Client(`http://localhost:${port}`)

      client1.on('connect', () => {
        client1.emit('join-session', { sessionId })
      })

      client2.on('connect', () => {
        client2.emit('join-session', { sessionId })
      })

      client2.on('session-joined', () => {
        // Wait a bit for both clients to be ready
        setTimeout(() => {
          client1.emit('code-change', {
            sessionId,
            code: 'const test = 123;',
          })
        }, 100)
      })

      client2.on('code-change', (data) => {
        expect(data.sessionId).toBe(sessionId)
        expect(data.code).toBe('const test = 123;')
        codeChangeReceived = true
        done()
      })

      // Timeout fallback
      setTimeout(() => {
        if (!codeChangeReceived) {
          done(new Error('Code change not received'))
        }
      }, 2000)
    })

    test('User count updates when clients join/leave', (done) => {
      const sessionId = 'test-session-3'
      let userUpdates = []

      client1 = new Client(`http://localhost:${port}`)
      client2 = new Client(`http://localhost:${port}`)

      client1.on('connect', () => {
        client1.emit('join-session', { sessionId })
      })

      client1.on('users-update', (users) => {
        userUpdates.push(users.length)
        
        if (userUpdates.length === 1 && users.length === 1) {
          // First client joined, now join second
          client2.on('connect', () => {
            client2.emit('join-session', { sessionId })
          })
        } else if (userUpdates.length === 2 && users.length === 2) {
          // Both clients joined
          expect(userUpdates).toEqual([1, 2])
          done()
        }
      })

      client2.on('connect', () => {
        client2.emit('join-session', { sessionId })
      })
    })
  })
})








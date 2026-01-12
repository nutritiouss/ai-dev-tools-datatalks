# Coding Interview Platform

A real-time collaborative coding interview platform built with React, Express.js, and WebSockets. This application allows interviewers to create coding sessions and share them with candidates, enabling real-time code editing, syntax highlighting, and browser-based code execution.

## Features

- **Real-time Collaboration**: Multiple users can edit code simultaneously with live updates
- **Session Management**: Create and share unique session links
- **Syntax Highlighting**: Support for JavaScript and Python with Monaco Editor
- **Code Execution**: Execute code safely in the browser using Pyodide (Python WASM) and native JavaScript
- **User Tracking**: See how many users are connected to each session

## Tech Stack

- **Frontend**: React 18 + Vite
- **Backend**: Express.js
- **WebSocket**: Socket.io
- **Code Editor**: Monaco Editor (VS Code editor)
- **Syntax Highlighting**: Monaco Editor (built-in)
- **Code Execution**: Pyodide (Python WASM), native JavaScript eval
- **Testing**: Jest, Supertest
- **Containerization**: Docker
- **Base Image**: node:20-alpine

## Project Structure

```
02-end-to-end/
├── frontend/          # React + Vite frontend
│   ├── src/
│   │   ├── components/
│   │   │   ├── CodeEditor.jsx
│   │   │   └── SessionManager.jsx
│   │   ├── App.jsx
│   │   └── main.jsx
│   └── package.json
├── backend/           # Express.js backend
│   ├── routes/
│   │   └── sessions.js
│   ├── tests/
│   │   └── integration.test.js
│   ├── server.js
│   ├── websocket.js
│   └── package.json
├── Dockerfile         # Production Docker image
├── package.json       # Root package with concurrently
└── README.md
```

## Setup Instructions

### Prerequisites

- Node.js 18+ installed
- npm or yarn package manager

### Installation

1. Install root dependencies:
```bash
npm install
```

2. Install frontend dependencies:
```bash
cd frontend
npm install
cd ..
```

3. Install backend dependencies:
```bash
cd backend
npm install
cd ..
```

## Running the Application

### Development Mode (Both Frontend and Backend)

From the root directory, run:

```bash
npm run dev
```

This uses `concurrently` to run both the frontend (port 3000) and backend (port 5000) simultaneously.

The command in `package.json` is:
```json
"dev": "concurrently \"npm run dev --prefix backend\" \"npm run dev --prefix frontend\""
```

### Running Separately

**Frontend only:**
```bash
cd frontend
npm run dev
```

**Backend only:**
```bash
cd backend
npm run dev
```

## Testing

### Run All Tests

From the root directory:
```bash
npm test
```

### Run Integration Tests

Integration tests check the interaction between client and server, including:
- API endpoint functionality
- WebSocket connections
- Code synchronization between multiple clients
- User count updates

**Backend integration tests:**
```bash
cd backend
npm test
```

The terminal command for executing tests is:
```bash
npm test
```

### Frontend Tests

```bash
cd frontend
npm test
```

## Building for Production

Build the frontend:
```bash
cd frontend
npm run build
```

## Docker

### Build Docker Image

```bash
docker build -t coding-interview-platform .
```

### Run Docker Container

```bash
docker run -p 5000:5000 coding-interview-platform
```

The application will be available at `http://localhost:5000`

### Docker Base Image

The Dockerfile uses `node:20-alpine` as the base image for both build and production stages.

## Deployment

The application can be deployed to various cloud platforms:

- **Render**: Recommended for easy deployment
- **Railway**: Alternative option
- **Heroku**: Traditional option
- **AWS/GCP/Azure**: For enterprise deployments

For deployment, the Dockerfile creates a single container with both frontend and backend. The frontend is built and served as static files by the Express server.

### Deployment Service Used

**Render** (recommended choice)

To deploy to Render:
1. Connect your GitHub repository
2. Select the Docker deployment option
3. Render will automatically build and deploy using the Dockerfile

## Usage

1. **Create a Session**: Click "Create New Session" to generate a unique session ID
2. **Share the Link**: Copy the session link and share it with candidates
3. **Collaborate**: Multiple users can join and edit code in real-time
4. **Select Language**: Choose between JavaScript and Python
5. **Execute Code**: Click "Run Code" to execute the code in the browser
6. **View Output**: See the execution results in the output panel

## API Endpoints

- `POST /api/sessions` - Create a new coding session
- `GET /api/sessions/:id` - Get session information
- `GET /health` - Health check endpoint

## WebSocket Events

### Client to Server:
- `join-session` - Join a coding session
- `code-change` - Broadcast code changes

### Server to Client:
- `session-joined` - Confirmation of joining a session
- `code-change` - Receive code updates from other users
- `users-update` - Updated list of connected users

## Homework Answers

### Question 1: Initial Implementation Prompt

The initial prompt given to AI was:

```
Create a collaborative coding interview platform with the following features:
- Frontend: React + Vite application
- Backend: Express.js server with WebSocket support
- Real-time code collaboration using WebSockets
- Support for creating and joining coding sessions via unique session IDs
- Code editor with syntax highlighting for JavaScript and Python
- Browser-based code execution (Python via WASM, JavaScript natively)
- User tracking to show connected users per session

Implement both frontend and backend in one go. Use Socket.io for WebSocket communication.
Store sessions in memory for now (can be replaced with database later).
```

### Question 2: Integration Test Command

The terminal command for executing integration tests is:

```bash
npm test
```

This runs tests in both frontend and backend. For backend integration tests specifically:

```bash
cd backend && npm test
```

### Question 3: npm dev Command

The command in `package.json` for running both frontend and backend concurrently:

```json
"dev": "concurrently \"npm run dev --prefix backend\" \"npm run dev --prefix frontend\""
```

This uses the `concurrently` package to run both services simultaneously.

### Question 4: Syntax Highlighting Library

**Monaco Editor** (`@monaco-editor/react`)

Monaco Editor is the same editor that powers VS Code. It provides:
- Built-in syntax highlighting for 100+ languages
- IntelliSense (code completion)
- Error detection
- Multi-cursor editing
- And many other VS Code features

### Question 5: WASM Library for Python

**Pyodide** (`pyodide`)

Pyodide compiles Python to WebAssembly (WASM) and runs it in the browser. It's loaded from CDN:
```javascript
const pyodide = await loadPyodide({
  indexURL: 'https://cdn.jsdelivr.net/pyodide/v0.26.1/full/',
})
```

Pyodide allows executing Python code safely in the browser without requiring a backend Python interpreter.

### Question 6: Docker Base Image

**node:20-alpine**

The Dockerfile uses `node:20-alpine` as the base image for both the build stage and production stage. Alpine Linux is chosen for its small size, making the final Docker image more lightweight.

### Question 7: Deployment Service

**Render**

Render is recommended for deployment because:
- Easy GitHub integration
- Automatic deployments on push
- Free tier available
- Simple Docker deployment
- Built-in SSL certificates
- Good documentation and support

## Development Notes

- The application uses in-memory storage for sessions (sessions are lost on server restart)
- For production, consider adding a database (PostgreSQL, MongoDB, etc.)
- Code execution is limited to browser environment for security
- WebSocket connections are managed per session room
- Empty sessions are cleaned up after 5 minutes of inactivity

## License

This project is created for educational purposes as part of the AI Dev Tools Zoomcamp course.








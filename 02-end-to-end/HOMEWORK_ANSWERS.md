# Homework 2: End-To-End Project - Answers

## Question 1: Initial prompt (1 point)

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

## Question 2: Integration tests (1 point)

```
npm test
```

## Question 3: Running server & client (1 point)

```
concurrently "npm run dev --prefix backend" "npm run dev --prefix frontend"
```

Or the full command from package.json:
```json
"dev": "concurrently \"npm run dev --prefix backend\" \"npm run dev --prefix frontend\""
```

## Question 4: Syntax highlight (1 point)

```
Monaco Editor
```

Package name: `@monaco-editor/react`

## Question 5: Python WASM library (1 point)

```
Pyodide
```

Package name: `pyodide`

## Question 6: Base image (1 point)

```
node:20-alpine
```

## Question 7: Deployment (1 point)

```
Render
```

---

## GitHub Repository Link

If you need to submit a GitHub link to the folder:

```
https://github.com/YOUR_USERNAME/YOUR_REPO/tree/main/02-end-to-end
```

Replace `YOUR_USERNAME` and `YOUR_REPO` with your actual GitHub username and repository name.


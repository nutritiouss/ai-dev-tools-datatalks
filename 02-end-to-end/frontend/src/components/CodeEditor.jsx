import { useState, useEffect, useRef } from 'react'
import { useParams, useNavigate } from 'react-router-dom'
import { io } from 'socket.io-client'
import Editor from '@monaco-editor/react'
import { loadPyodide } from 'pyodide'
import './CodeEditor.css'

function CodeEditor() {
  const { sessionId } = useParams()
  const navigate = useNavigate()
  const [code, setCode] = useState('// Write your code here\n')
  const [language, setLanguage] = useState('javascript')
  const [output, setOutput] = useState('')
  const [connectedUsers, setConnectedUsers] = useState([])
  const [isExecuting, setIsExecuting] = useState(false)
  const socketRef = useRef(null)
  const pyodideRef = useRef(null)

  useEffect(() => {
    // Initialize Pyodide for Python execution
    const initPyodide = async () => {
      try {
        const pyodide = await loadPyodide({
          indexURL: 'https://cdn.jsdelivr.net/pyodide/v0.26.1/full/',
        })
        pyodideRef.current = pyodide
      } catch (error) {
        console.error('Failed to load Pyodide:', error)
      }
    }
    initPyodide()

    // Connect to WebSocket server
    socketRef.current = io('http://localhost:5000', {
      query: { sessionId },
    })

    socketRef.current.on('connect', () => {
      console.log('Connected to server')
    })

    socketRef.current.on('code-change', (data) => {
      if (data.sessionId === sessionId) {
        setCode(data.code)
      }
    })

    socketRef.current.on('users-update', (users) => {
      setConnectedUsers(users)
    })

    socketRef.current.on('session-joined', (data) => {
      if (data.code) {
        setCode(data.code)
      }
    })

    // Request current session state
    socketRef.current.emit('join-session', { sessionId })

    return () => {
      if (socketRef.current) {
        socketRef.current.disconnect()
      }
    }
  }, [sessionId])

  const handleCodeChange = (value) => {
    setCode(value || '')
    if (socketRef.current) {
      socketRef.current.emit('code-change', {
        sessionId,
        code: value || '',
      })
    }
  }

  const executeCode = async () => {
    setIsExecuting(true)
    setOutput('')

    try {
      if (language === 'python') {
        if (!pyodideRef.current) {
          setOutput('Error: Pyodide not loaded yet. Please wait...')
          setIsExecuting(false)
          return
        }

        try {
          const result = await pyodideRef.current.runPythonAsync(code)
          setOutput(String(result || 'Code executed successfully (no output)'))
        } catch (error) {
          setOutput(`Error: ${error.message}`)
        }
      } else if (language === 'javascript') {
        // Execute JavaScript in browser
        try {
          const originalConsoleLog = console.log
          let logs = []
          console.log = (...args) => {
            logs.push(args.map(arg => 
              typeof arg === 'object' ? JSON.stringify(arg, null, 2) : String(arg)
            ).join(' '))
          }

          const result = eval(code)
          console.log = originalConsoleLog

          let outputText = logs.length > 0 ? logs.join('\n') : ''
          if (result !== undefined) {
            outputText += (outputText ? '\n' : '') + String(result)
          }
          setOutput(outputText || 'Code executed successfully (no output)')
        } catch (error) {
          setOutput(`Error: ${error.message}`)
        }
      }
    } catch (error) {
      setOutput(`Error: ${error.message}`)
    } finally {
      setIsExecuting(false)
    }
  }

  const copySessionLink = () => {
    const link = `${window.location.origin}/session/${sessionId}`
    navigator.clipboard.writeText(link)
    alert('Session link copied to clipboard!')
  }

  return (
    <div className="code-editor-container">
      <div className="editor-header">
        <div className="header-left">
          <button onClick={() => navigate('/')} className="btn-back">
            ← Back
          </button>
          <div className="session-info">
            <span className="session-id">Session: {sessionId}</span>
            <button onClick={copySessionLink} className="btn-link">
              Copy Link
            </button>
          </div>
        </div>
        <div className="header-right">
          <select
            value={language}
            onChange={(e) => setLanguage(e.target.value)}
            className="language-select"
          >
            <option value="javascript">JavaScript</option>
            <option value="python">Python</option>
          </select>
          <button
            onClick={executeCode}
            disabled={isExecuting}
            className="btn-run"
          >
            {isExecuting ? 'Running...' : 'Run Code'}
          </button>
          <div className="users-count">
            👥 {connectedUsers.length} user{connectedUsers.length !== 1 ? 's' : ''}
          </div>
        </div>
      </div>

      <div className="editor-content">
        <div className="editor-panel">
          <Editor
            height="100%"
            language={language}
            value={code}
            onChange={handleCodeChange}
            theme="vs-dark"
            options={{
              minimap: { enabled: true },
              fontSize: 14,
              wordWrap: 'on',
              automaticLayout: true,
            }}
          />
        </div>
        <div className="output-panel">
          <div className="output-header">Output</div>
          <pre className="output-content">{output || 'Output will appear here...'}</pre>
        </div>
      </div>
    </div>
  )
}

export default CodeEditor








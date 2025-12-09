import { render, screen } from '@testing-library/react'
import { BrowserRouter } from 'react-router-dom'
import SessionManager from '../SessionManager'

describe('SessionManager', () => {
  it('renders the main heading', () => {
    render(
      <BrowserRouter>
        <SessionManager />
      </BrowserRouter>
    )
    expect(screen.getByText('Coding Interview Platform')).toBeInTheDocument()
  })

  it('renders create session button', () => {
    render(
      <BrowserRouter>
        <SessionManager />
      </BrowserRouter>
    )
    expect(screen.getByText('Create New Session')).toBeInTheDocument()
  })

  it('renders join session input and button', () => {
    render(
      <BrowserRouter>
        <SessionManager />
      </BrowserRouter>
    )
    expect(screen.getByPlaceholderText('Enter session ID')).toBeInTheDocument()
    expect(screen.getByText('Join Session')).toBeInTheDocument()
  })
})


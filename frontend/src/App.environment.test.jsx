import { render, screen } from '@testing-library/react'
import { describe, it, expect, beforeEach, afterEach, vi } from 'vitest'
import App from './App'

describe('App Environment Integration', () => {
  let originalWindowEnv

  beforeEach(() => {
    // Save original window.env
    originalWindowEnv = window.env
  })

  afterEach(() => {
    // Restore original window.env
    window.env = originalWindowEnv
  })

  it('renders with environment configuration available', () => {
    render(<App />)
    
    // The app should render without crashing
    expect(screen.getByText('PlayHub')).toBeInTheDocument()
    expect(screen.getByText('Your Gaming Hub - Queue, Play, Trade')).toBeInTheDocument()
    
    // Environment should be available
    expect(window.env).toBeDefined()
    expect(window.env.REACT_APP_ENV).toBeDefined()
    expect(window.env.REACT_APP_API_BASE_URL).toBeDefined()
  })

  it('handles missing environment gracefully', () => {
    // Remove window.env
    delete window.env
    
    // App should still render
    render(<App />)
    expect(screen.getByText('PlayHub')).toBeInTheDocument()
  })

  it('can access environment variables in component', () => {
    // Mock a component that uses environment variables
    const TestComponent = () => {
      const apiUrl = window.env?.REACT_APP_API_BASE_URL || 'http://localhost:8081'
      const environment = window.env?.REACT_APP_ENV || 'development'
      
      return (
        <div>
          <span data-testid="api-url">{apiUrl}</span>
          <span data-testid="environment">{environment}</span>
        </div>
      )
    }
    
    render(<TestComponent />)
    
    expect(screen.getByTestId('api-url')).toBeInTheDocument()
    expect(screen.getByTestId('environment')).toBeInTheDocument()
  })

  it('environment variables are accessible during render', () => {
    let capturedEnv = null
    
    const TestComponent = () => {
      capturedEnv = window.env
      return <div>Test</div>
    }
    
    render(<TestComponent />)
    
    expect(capturedEnv).toBeDefined()
    expect(capturedEnv.REACT_APP_ENV).toBeDefined()
    expect(capturedEnv.REACT_APP_API_BASE_URL).toBeDefined()
  })
})

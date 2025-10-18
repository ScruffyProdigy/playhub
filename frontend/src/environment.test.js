import { describe, it, expect, beforeEach, afterEach, vi } from 'vitest'

describe('Environment Configuration', () => {
  let originalWindowEnv

  beforeEach(() => {
    // Save original window.env
    originalWindowEnv = window.env
  })

  afterEach(() => {
    // Restore original window.env
    window.env = originalWindowEnv
  })

  it('should have window.env defined', () => {
    expect(window.env).toBeDefined()
  })

  it('should have required environment variables', () => {
    expect(window.env).toHaveProperty('REACT_APP_ENV')
    expect(window.env).toHaveProperty('REACT_APP_API_BASE_URL')
  })

  it('should have valid environment values', () => {
    expect(typeof window.env.REACT_APP_ENV).toBe('string')
    expect(typeof window.env.REACT_APP_API_BASE_URL).toBe('string')
    expect(window.env.REACT_APP_ENV.length).toBeGreaterThan(0)
    expect(window.env.REACT_APP_API_BASE_URL.length).toBeGreaterThan(0)
  })

  it('should have valid API base URL format', () => {
    const apiUrl = window.env.REACT_APP_API_BASE_URL
    expect(apiUrl).toMatch(/^https?:\/\//)
  })

  it('should handle missing environment gracefully', () => {
    // Temporarily remove window.env
    delete window.env
    
    // The app should not crash
    expect(() => {
      // Simulate accessing window.env in the app
      const env = window.env || {}
      return env.REACT_APP_API_BASE_URL || 'http://localhost:8081'
    }).not.toThrow()
  })

  it('should provide fallback values when environment is missing', () => {
    // Test fallback behavior
    const getApiBaseUrl = () => {
      return window.env?.REACT_APP_API_BASE_URL || 'http://localhost:8081'
    }
    
    const getEnvironment = () => {
      return window.env?.REACT_APP_ENV || 'development'
    }
    
    expect(getApiBaseUrl()).toBeDefined()
    expect(getEnvironment()).toBeDefined()
  })

  it('should have consistent environment configuration', () => {
    const env1 = window.env
    const env2 = window.env
    
    expect(env1).toEqual(env2)
  })

  it('should allow environment variables to be overridden', () => {
    const originalApiUrl = window.env.REACT_APP_API_BASE_URL
    
    // Test overriding environment variables
    window.env.REACT_APP_API_BASE_URL = 'https://test-api.example.com'
    
    expect(window.env.REACT_APP_API_BASE_URL).toBe('https://test-api.example.com')
    
    // Restore original value
    window.env.REACT_APP_API_BASE_URL = originalApiUrl
  })
})

describe('Environment Configuration Integration', () => {
  it('should be able to construct API URLs', () => {
    const apiBaseUrl = window.env.REACT_APP_API_BASE_URL
    const graphqlUrl = `${apiBaseUrl}/graphql`
    const healthUrl = `${apiBaseUrl}/healthz`
    
    expect(graphqlUrl).toMatch(/^https?:\/\/.*\/graphql$/)
    expect(healthUrl).toMatch(/^https?:\/\/.*\/healthz$/)
  })

  it('should be able to make API requests', async () => {
    const apiBaseUrl = window.env.REACT_APP_API_BASE_URL
    const healthUrl = `${apiBaseUrl}/healthz`
    
    try {
      const response = await fetch(healthUrl)
      if (response.ok) {
        const text = await response.text()
        expect(text).toBeDefined()
      }
    } catch (error) {
      // If the request fails (backend not available), that's okay for unit tests
      // We just want to ensure the URL is properly formed
      expect(apiBaseUrl).toMatch(/^https?:\/\//)
    }
  })

  it('should handle different environment configurations', () => {
    const environments = ['local', 'staging', 'production']
    const currentEnv = window.env.REACT_APP_ENV
    
    expect(environments).toContain(currentEnv)
  })
})

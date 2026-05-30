import { describe, it, expect } from 'vitest'
import { getApiBaseUrl, getGraphQLUrl, getGraphQLWsUrl } from './env'

describe('env helpers', () => {
  it('uses relative GraphQL path when API base is empty', () => {
    window.env = { REACT_APP_API_BASE_URL: '' }
    expect(getApiBaseUrl()).toBe('')
    expect(getGraphQLUrl()).toBe('/graphql')
  })

  it('builds GraphQL URL from configured API base', () => {
    window.env = { REACT_APP_API_BASE_URL: 'http://localhost:8080/' }
    expect(getApiBaseUrl()).toBe('http://localhost:8080')
    expect(getGraphQLUrl()).toBe('http://localhost:8080/graphql')
  })

  it('uses a direct backend WebSocket URL for subscriptions', () => {
    window.env = { REACT_APP_API_BASE_URL: 'http://localhost:8080' }
    expect(getGraphQLWsUrl()).toBe('ws://localhost:8080/graphql')
  })

  it('defaults WebSocket to local backend when API base is empty', () => {
    window.env = { REACT_APP_API_BASE_URL: '' }
    expect(getGraphQLWsUrl()).toBe('ws://127.0.0.1:8080/graphql')
  })
})

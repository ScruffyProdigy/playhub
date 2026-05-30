import '@testing-library/jest-dom'

// Mock window.env for testing
Object.defineProperty(window, 'env', {
  writable: true,
  configurable: true,
  value: {
    REACT_APP_ENV: 'test',
    REACT_APP_API_BASE_URL: '',
  },
})

// Mock window.matchMedia
Object.defineProperty(window, 'matchMedia', {
  writable: true,
  value: vi.fn().mockImplementation(query => ({
    matches: false,
    media: query,
    onchange: null,
    addListener: vi.fn(), // deprecated
    removeListener: vi.fn(), // deprecated
    addEventListener: vi.fn(),
    removeEventListener: vi.fn(),
    dispatchEvent: vi.fn(),
  })),
})

// Mock IntersectionObserver
global.IntersectionObserver = vi.fn().mockImplementation(() => ({
  observe: vi.fn(),
  unobserve: vi.fn(),
  disconnect: vi.fn(),
}))

// Mock ResizeObserver
global.ResizeObserver = vi.fn().mockImplementation(() => ({
  observe: vi.fn(),
  unobserve: vi.fn(),
  disconnect: vi.fn(),
}))

export function mockGraphQLResponse(data) {
  global.fetch = vi.fn().mockResolvedValue({
    ok: true,
    json: async () => ({ data }),
  })
}

export function mockUnauthenticatedSession() {
  mockGraphQLResponse({ me: null })
}

export function mockAuthenticatedSession(user = {
  id: 'user-1',
  email: 'player@example.com',
  displayName: 'player',
  createdAt: '2026-01-01T00:00:00Z',
}) {
  mockGraphQLResponse({ me: user })
}

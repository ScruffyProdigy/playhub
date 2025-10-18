# Frontend Testing Guide

This document describes the comprehensive testing setup for the React frontend application.

## Testing Stack

- **Unit & Integration Tests**: Vitest + React Testing Library
- **E2E Tests**: Playwright
- **Test Runner**: Vitest (fast, Vite-native)
- **Assertions**: Vitest + Jest DOM matchers
- **Mocking**: Vitest built-in mocking

## Test Types

### 1. Unit Tests (`src/*.test.jsx`)

Test individual components in isolation:

```bash
npm run test:run
```

**Coverage:**
- Component rendering
- Props handling
- State management
- Event handlers
- Conditional rendering

### 2. Integration Tests (`src/integration/*.test.jsx`)

Test component interactions and user journeys:

```bash
npm run test:run
```

**Coverage:**
- User interactions
- Component communication
- State persistence
- Error handling
- Performance

### 3. E2E Tests (`tests/e2e/*.spec.js`)

Test the complete application in real browsers:

```bash
npm run test:e2e
```

**Coverage:**
- Full user workflows
- Cross-browser compatibility
- Accessibility
- Performance
- Mobile responsiveness

## Running Tests

### Development

```bash
# Run unit tests in watch mode
npm run test:watch

# Run unit tests with UI
npm run test:ui

# Run all tests (unit + E2E)
npm run test:all
```

### CI/CD

```bash
# Run unit tests with coverage
npm run test:coverage

# Run E2E tests
npm run test:e2e

# Run E2E tests with UI
npm run test:e2e:ui
```

## Test Structure

```
frontend/
├── src/
│   ├── test/
│   │   └── setup.js              # Test configuration
│   ├── App.test.jsx              # Unit tests
│   ├── main.test.jsx             # Entry point tests
│   └── integration/
│       └── App.integration.test.jsx  # Integration tests
├── tests/
│   └── e2e/
│       └── app.spec.js           # E2E tests
├── playwright.config.js          # Playwright configuration
└── vite.config.js               # Vitest configuration
```

## Writing Tests

### Unit Test Example

```jsx
import { describe, it, expect } from 'vitest'
import { render, screen } from '@testing-library/react'
import userEvent from '@testing-library/user-event'
import App from './App'

describe('App Component', () => {
  it('increments counter when clicked', async () => {
    const user = userEvent.setup()
    render(<App />)
    
    const button = screen.getByRole('button', { name: /count is 0/i })
    await user.click(button)
    
    expect(screen.getByText('count is 1')).toBeInTheDocument()
  })
})
```

### Integration Test Example

```jsx
describe('User Journey: Counter Interaction', () => {
  it('allows user to increment counter multiple times', async () => {
    const user = userEvent.setup()
    render(<App />)
    
    const button = screen.getByRole('button')
    
    // Simulate user clicking multiple times
    await user.click(button)
    await user.click(button)
    await user.click(button)
    
    expect(screen.getByText('count is 3')).toBeInTheDocument()
  })
})
```

### E2E Test Example

```jsx
import { test, expect } from '@playwright/test'

test('increments counter when clicked', async ({ page }) => {
  await page.goto('/')
  
  const button = page.getByRole('button', { name: 'count is 0' })
  await button.click()
  
  await expect(page.getByRole('button', { name: 'count is 1' })).toBeVisible()
})
```

## Test Configuration

### Vitest Configuration (`vite.config.js`)

```js
export default defineConfig({
  test: {
    globals: true,
    environment: 'jsdom',
    setupFiles: ['./src/test/setup.js'],
  },
})
```

### Playwright Configuration (`playwright.config.js`)

```js
export default defineConfig({
  testDir: './tests/e2e',
  use: {
    baseURL: 'http://localhost:5173',
  },
  webServer: {
    command: 'npm run dev',
    url: 'http://localhost:5173',
    reuseExistingServer: !process.env.CI,
  },
})
```

## Best Practices

### 1. Test Structure

- **Arrange**: Set up test data and render components
- **Act**: Perform user actions or trigger events
- **Assert**: Verify expected outcomes

### 2. Naming Conventions

- Test files: `*.test.jsx` or `*.spec.js`
- Test descriptions: Use "should" or "when" statements
- Test groups: Group related tests with `describe`

### 3. Accessibility Testing

```jsx
// Test accessibility attributes
expect(button).toHaveAccessibleName(/count is \d+/)
expect(button).toBeFocused()

// Test keyboard navigation
await page.keyboard.press('Tab')
await page.keyboard.press('Enter')
```

### 4. Performance Testing

```jsx
// Test response times
const startTime = Date.now()
await button.click()
const responseTime = Date.now() - startTime
expect(responseTime).toBeLessThan(100)
```

### 5. Error Handling

```jsx
// Test error states
expect(() => render(<Component />)).not.toThrow()

// Test error boundaries
expect(screen.getByText('Something went wrong')).toBeInTheDocument()
```

## CI/CD Integration

### GitHub Actions

The `.github/workflows/frontend-tests.yml` workflow runs:

1. **Unit Tests**: Linting, unit tests, coverage
2. **E2E Tests**: Cross-browser testing with Playwright
3. **Artifacts**: Test reports and screenshots

### Local Development

```bash
# Pre-commit hook (add to .git/hooks/pre-commit)
#!/bin/bash
cd frontend
npm run test:run
npm run lint
```

## Debugging Tests

### Unit Tests

```bash
# Run specific test file
npm run test App.test.jsx

# Run with debug output
npm run test -- --reporter=verbose

# Run in watch mode
npm run test:watch
```

### E2E Tests

```bash
# Run with browser UI
npm run test:e2e:ui

# Run in headed mode
npm run test:e2e:headed

# Run specific test
npx playwright test app.spec.js
```

### Debug Mode

```bash
# Debug unit tests
npm run test -- --inspect-brk

# Debug E2E tests
npx playwright test --debug
```

## Coverage Reports

```bash
# Generate coverage report
npm run test:coverage

# View coverage in browser
open coverage/index.html
```

## Common Issues

### 1. Async Operations

```jsx
// Use waitFor for async operations
await waitFor(() => {
  expect(screen.getByText('Updated text')).toBeInTheDocument()
})
```

### 2. Mocking

```jsx
// Mock external dependencies
vi.mock('./api', () => ({
  fetchData: vi.fn(() => Promise.resolve({ data: 'test' }))
}))
```

### 3. Cleanup

```jsx
// Clean up after tests
afterEach(() => {
  cleanup()
  vi.clearAllMocks()
})
```

## Performance Testing

### Load Testing

```jsx
// Test with many interactions
for (let i = 0; i < 100; i++) {
  await user.click(button)
}
expect(screen.getByText('count is 100')).toBeInTheDocument()
```

### Memory Testing

```jsx
// Test for memory leaks
const { unmount } = render(<Component />)
// Perform many operations
unmount()
// Check for cleanup
```

## Accessibility Testing

### Screen Reader Testing

```jsx
// Test with screen reader
expect(button).toHaveAccessibleName('Increment counter')
expect(button).toHaveAttribute('aria-label')
```

### Keyboard Navigation

```jsx
// Test keyboard accessibility
await user.tab()
expect(button).toBeFocused()
await user.keyboard('{Enter}')
```

## Mobile Testing

### Responsive Design

```jsx
// Test mobile viewport
await page.setViewportSize({ width: 375, height: 667 })
await page.goto('/')
// Test mobile interactions
```

## Continuous Integration

### Pre-commit Hooks

```bash
# Install husky for git hooks
npm install --save-dev husky

# Add pre-commit hook
npx husky add .husky/pre-commit "cd frontend && npm run test:run"
```

### Branch Protection

Configure GitHub branch protection rules:
- Require status checks: `unit-tests`, `e2e-tests`
- Require up-to-date branches
- Dismiss stale reviews

## Monitoring

### Test Metrics

- Test execution time
- Coverage percentage
- Flaky test detection
- Performance regression detection

### Alerts

Set up alerts for:
- Test failures
- Coverage drops
- Performance regressions
- Flaky tests

## Resources

- [Vitest Documentation](https://vitest.dev/)
- [React Testing Library](https://testing-library.com/docs/react-testing-library/intro/)
- [Playwright Documentation](https://playwright.dev/)
- [Jest DOM Matchers](https://github.com/testing-library/jest-dom)
- [User Event](https://testing-library.com/docs/user-event/intro/)

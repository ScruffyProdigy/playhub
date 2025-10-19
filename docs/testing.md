# Testing Guide

This guide covers the testing strategies and how to run tests for PlayHub.

## Testing Status

### ✅ Implemented
- **Unit Tests**: Component and function testing with Vitest
- **Integration Tests**: Component interaction testing
- **E2E Tests**: Full user workflow testing with Playwright
- **Environment Configuration Tests**: Runtime environment validation
- **Backend Tests**: GraphQL resolver testing with mock data
- **Drift Detection**: gqlgen code generation validation

### 🚧 In Development
- **Database Tests**: Integration tests with real database
- **Authentication Tests**: JWT and user management testing
- **Performance Tests**: Load testing and benchmarking

### 📋 Planned
- **API Contract Tests**: Schema validation and API testing
- **Security Tests**: Vulnerability and penetration testing
- **Chaos Engineering**: Failure scenario testing

## Testing Philosophy

PlayHub uses a comprehensive testing strategy with multiple layers:

- **Unit Tests**: Test individual functions and components
- **Integration Tests**: Test component interactions
- **End-to-End Tests**: Test complete user workflows
- **Performance Tests**: Benchmark critical paths

## Backend Testing

### Test Structure
```
backend/
├── graph/
│   ├── healthz_test.go          # Basic functionality tests
│   ├── resolvers_test.go        # GraphQL resolver tests
│   ├── benchmark_test.go        # Performance benchmarks
│   ├── gqlgen_drift_test.go     # Code generation drift detection
│   └── drift_demo_test.go       # Drift detection demo
```

### Running Backend Tests

```bash
# All tests
cd backend && go test ./...

# Specific test file
cd backend && go test ./graph -v

# Run benchmarks
cd backend && go test -bench=. ./graph

# Run with coverage
cd backend && go test -cover ./...
```

### Test Types

1. **Health Check Tests** (`healthz_test.go`)
   - Basic API functionality
   - GraphQL query execution
   - Direct resolver calls

2. **Resolver Tests** (`resolvers_test.go`)
   - GraphQL query/mutation testing
   - Error handling
   - Pagination
   - Data validation

3. **Benchmark Tests** (`benchmark_test.go`)
   - Performance measurement
   - Load testing
   - Memory usage

4. **Drift Detection** (`gqlgen_drift_test.go`)
   - Ensures generated code is up-to-date
   - Prevents "forgot to run generate" errors

## Frontend Testing

### Test Structure
```
frontend/
├── src/
│   ├── App.test.jsx                    # Component unit tests
│   ├── App.environment.test.jsx        # Environment integration tests
│   ├── environment.test.js             # Environment configuration tests
│   ├── main.test.jsx                   # Entry point tests
│   └── integration/
│       └── App.integration.test.jsx    # Integration tests
└── tests/
    └── e2e/
        ├── app.spec.js                 # End-to-end tests
        └── environment-config.spec.js  # Environment configuration E2E tests
```

### Running Frontend Tests

```bash
# Unit and integration tests
cd frontend && npm run test:run

# Watch mode
cd frontend && npm run test:watch

# With coverage
cd frontend && npm run test:coverage

# E2E tests
cd frontend && npm run test:e2e

# All tests
cd frontend && npm run test:all
```

### Test Types

1. **Unit Tests** (`*.test.jsx`)
   - Component rendering
   - User interactions
   - State management
   - Props validation

2. **Integration Tests** (`integration/*.test.jsx`)
   - Component interactions
   - User journeys
   - Performance testing
   - Error handling

3. **E2E Tests** (`tests/e2e/*.spec.js`)
   - Full user workflows
   - Cross-browser testing
   - Accessibility testing
   - Performance testing

4. **Environment Configuration Tests**
   - **Unit Tests** (`environment.test.js`)
     - Environment variable validation
     - API URL construction
     - Fallback behavior testing
   - **Integration Tests** (`App.environment.test.jsx`)
     - Component environment integration
     - Runtime environment access
     - Error handling
   - **E2E Tests** (`environment-config.spec.js`)
     - `window.env` loading verification
     - API connectivity testing
     - Runtime environment injection
     - Cross-environment validation

## Test Configuration

### Vitest Configuration
```javascript
// vite.config.js
export default defineConfig({
  test: {
    globals: true,
    environment: 'jsdom',
    setupFiles: ['./src/test/setup.js'],
    exclude: ['**/node_modules/**', '**/dist/**', '**/tests/e2e/**'],
  },
})
```

### Playwright Configuration
```javascript
// playwright.config.js
export default defineConfig({
  testDir: './tests/e2e',
  fullyParallel: true,
  projects: [
    { name: 'chromium', use: { ...devices['Desktop Chrome'] } },
    { name: 'firefox', use: { ...devices['Desktop Firefox'] } },
    { name: 'webkit', use: { ...devices['Desktop Safari'] } },
  ],
})
```

## CI/CD Testing

### GitHub Actions
The project includes automated testing in CI/CD:

- **Backend Tests**: Run on every push/PR
- **Frontend Tests**: Unit, integration, and E2E tests
- **Drift Detection**: Ensures generated code is current
- **Cross-browser Testing**: Multiple browser environments

### Test Scripts
```bash
# Run all tests (from project root)
./scripts/test.sh

# Backend only
./scripts/test-backend.sh

# Frontend only
./scripts/test-frontend.sh
```

## Best Practices

### Writing Tests

1. **Test Structure**
   - Arrange: Set up test data
   - Act: Execute the function/component
   - Assert: Verify the results

2. **Naming Conventions**
   - `TestFunctionName` for Go tests
   - `describe('Component', () => { it('should do something', () => {}) })` for JS tests

3. **Test Data**
   - Use factories for test data
   - Keep tests independent
   - Clean up after tests

### Test Maintenance

1. **Keep Tests Fast**
   - Mock external dependencies
   - Use in-memory databases
   - Parallel test execution

2. **Test Coverage**
   - Aim for >80% coverage
   - Focus on critical paths
   - Don't test implementation details

3. **Flaky Tests**
   - Fix flaky tests immediately
   - Use proper waits in E2E tests
   - Retry mechanisms for network calls

## Debugging Tests

### Backend
```bash
# Verbose output
go test -v ./graph

# Run specific test
go test -run TestHealthz ./graph

# Debug with delve
dlv test ./graph
```

### Frontend
```bash
# Debug mode
npm run test:ui

# Watch mode with debugging
npm run test:watch

# E2E debugging
npm run test:e2e:headed
```

## Performance Testing

### Backend Benchmarks
```bash
cd backend && go test -bench=. -benchmem ./graph
```

### Frontend Performance
- Lighthouse CI integration
- Bundle size monitoring
- Runtime performance tests

## Troubleshooting

### Common Issues

1. **Tests timing out**
   - Increase timeout values
   - Check for infinite loops
   - Verify async operations complete

2. **Flaky E2E tests**
   - Add proper waits
   - Use stable selectors
   - Check for race conditions

3. **Test environment issues**
   - Verify dependencies installed
   - Check environment variables
   - Clear caches if needed

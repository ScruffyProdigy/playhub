# Development Guide

This guide will help you set up and run JoinQuest locally for development.

**Context:** JoinQuest is the shared lobby and integration platform described in **[Product vision](vision.md)**. This repo is the server + player UI; partner games (e.g. `demo-game-rps`) integrate via the handoff protocol.

## Current Development Status

### ✅ Ready for Development
- **Environment Configuration**: Docker-based runtime environment injection
- **GraphQL API**: Catalog matchmaking, auth, PostgreSQL-backed resolvers
- **Frontend Foundation**: React app with testing infrastructure
- **Kubernetes Deployment**: Multi-environment deployment scripts
- **Testing Suite**: Comprehensive test coverage
- **Database Integration**: PostgreSQL connection management
- **CI/CD Pipeline**: GitHub Actions workflows with proper linting and testing

### 🚧 In Active Development
- **Digital goods**: Player purchase and trade flows

### 📋 Next Steps (not in current branch)
- **E2E partner path**: [end-to-end-partner-checklist.md](./end-to-end-partner-checklist.md)
- **`seatTemplate` + LFG**: [seat-templates-and-matchmaking.md](./seat-templates-and-matchmaking.md) — manifest still flat `seats[]` in code
- **Post-game GraphQL**: [match-lifecycle-callbacks.md](./match-lifecycle-callbacks.md) (`reportPlayerFinished`, `reportMatchResult`)

## Prerequisites

- **Go 1.25+** - [Download](https://golang.org/dl/)
- **Node.js 20+** - [Download](https://nodejs.org/)
- **Docker & Docker Compose** - [Download](https://www.docker.com/)
- **Git** - [Download](https://git-scm.com/)

## Quick Setup

1. **Clone the repository**
   ```bash
   git clone https://github.com/scruffyprodigy/playhub.git
   cd playhub
   ```

2. **Run the setup script**
   ```bash
   ./scripts/setup.sh
   ```

3. **Start development servers**
   ```bash
   ./scripts/dev.sh
   ```

## Manual Setup

### Backend Setup

1. **Navigate to backend directory**
   ```bash
   cd backend
   ```

2. **Install Go dependencies**
   ```bash
   go mod download
   ```

3. **Generate GraphQL code**
   ```bash
   go run github.com/99designs/gqlgen@v0.17.81 generate
   ```

4. **Run tests**
   ```bash
   go test ./...
   ```

5. **Start the server**
   ```bash
   go run server.go
   ```

The backend will be available at `http://localhost:8080`

### Frontend Setup

1. **Navigate to frontend directory**
   ```bash
   cd frontend
   ```

2. **Install dependencies**
   ```bash
   npm install
   ```

3. **Run tests**
   ```bash
   npm run test:run
   ```

4. **Start development server**
   ```bash
   npm run dev
   ```

The frontend will be available at `http://localhost:5173`

## Development Workflow

### Making Changes

1. **Backend Changes**
   - Modify GraphQL schema in `backend/graph/schema/`
   - Update resolvers in `backend/graph/`
   - Run `go run github.com/99designs/gqlgen@v0.17.81 generate` after schema changes
   - Test with `go test ./...`

2. **Frontend Changes**
   - Modify components in `frontend/src/`
   - Update tests as needed
   - Test with `npm run test:run`

### Code Generation

The backend uses gqlgen for GraphQL code generation:

```bash
cd backend
go run github.com/99designs/gqlgen@v0.17.81 generate
```

### Testing

- **Backend tests**: `./scripts/test-backend.sh` (uses isolated `playhub_test` DB; does not touch dev `playhub`)
- **Frontend unit tests**: `cd frontend && npm run test:run`
- **Frontend E2E tests**: `cd frontend && npm run test:e2e` (uses dev `playhub` via running backend)
- **All tests**: `./scripts/test.sh`
- **Test DB only**: `./scripts/db.sh test-migrate` then `export DATABASE_URL="$(./scripts/db.sh test-url)"`

### Linting

- **Backend**: `cd backend && go vet ./...`
- **Frontend**: `cd frontend && npm run lint`

## Environment Configuration

JoinQuest uses a Docker-based environment configuration system that allows the same Docker image to work across different environments (local, staging, production).

### Environment Variables

#### Backend
- `PORT` - Server port (default: 8080)
- `DATABASE_URL` - PostgreSQL connection string
- `MAGIC_LINK_BASE_URL` - Prefix for email sign-in links (see `.env.example`)
- `SESSION_COOKIE_NAME` / `SESSION_COOKIE_SECURE` - Session cookie settings
- `CORS_ALLOWED_ORIGINS` - Allowed browser origins for GraphQL
- `JWT_PRIVATE_KEY_PEM` / `JWKS_KID` - JWT signing for session cookies and seat tokens (see `.env.example`)

#### Frontend (Runtime Injection)
- `REACT_APP_ENV` - Environment identifier (local, staging, production)
- `REACT_APP_API_BASE_URL` - Backend API URL for the current environment

### Environment-Specific Configurations

#### Local Development
- **Frontend**: `http://localhost:5173`
- **GraphQL**: `http://localhost:8080/graphql`
- **Start**: `./scripts/dev.sh`

#### Production (JoinQuest)
- **Public URL**: `https://joinquest.cc`
- **GraphQL**: `https://joinquest.cc/graphql`
- **Deploy**: `./scripts/deploy-joinquest.sh`

### How Environment Configuration Works

1. **Docker Entrypoint**: The frontend Docker image includes a script that generates `env.js` at runtime
2. **Kubernetes ConfigMaps**: Each environment has its own ConfigMap with environment-specific values
3. **Runtime Injection**: When the container starts, it reads environment variables and creates the `env.js` file
4. **Frontend Access**: The frontend loads `/env.js` and accesses variables via `window.env`

For detailed information, see [Environment Configuration Guide](environment-configuration.md).

## Recent Updates

### GitHub Workflow Fixes (Latest)

The CI/CD pipeline has been updated with the following improvements:

- **ESLint Configuration**: Fixed linting issues with proper environment configuration
- **Test Environment Setup**: Added proper mocking for `window.env` in test environment
- **GraphQL Code Generation**: Ensured generated code is up to date with schema
- **Database Integration**: Added PostgreSQL connection management

All workflows now pass successfully:
- ✅ Frontend Tests (linting, unit tests, E2E tests)
- ✅ Environment Configuration Tests
- ✅ GraphQL Drift Detection

## Troubleshooting

### Common Issues

1. **GraphQL generation fails**
   - Ensure you're in the backend directory
   - Check that all schema files are valid GraphQL
   - Run `make generate` to regenerate code

2. **Frontend tests fail**
   - Clear node_modules: `rm -rf node_modules && npm install`
   - Check that all dependencies are installed
   - Ensure `window.env` is properly mocked in test setup

3. **Linting errors**
   - Run `npm run lint` to check for issues
   - ESLint is configured for Node.js, browser, and test environments

4. **Port conflicts**
   - Backend default: 8080
   - Frontend default: 5173
   - Change ports in respective config files if needed

### Getting Help

- Check the [API Documentation](api.md)
- Review [Architecture Overview](architecture.md)
- Open an issue on GitHub

# Development Guide

This guide will help you set up and run PlayHub locally for development.

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

- **Backend tests**: `cd backend && go test ./...`
- **Frontend unit tests**: `cd frontend && npm run test:run`
- **Frontend E2E tests**: `cd frontend && npm run test:e2e`
- **All tests**: `./scripts/test.sh`

### Linting

- **Backend**: `cd backend && go vet ./...`
- **Frontend**: `cd frontend && npm run lint`

## Environment Variables

### Backend
- `PORT` - Server port (default: 8080)
- `DATABASE_URL` - PostgreSQL connection string
- `JWT_SECRET` - JWT signing secret

### Frontend
- `VITE_API_URL` - Backend API URL (default: http://localhost:8080)

## Troubleshooting

### Common Issues

1. **GraphQL generation fails**
   - Ensure you're in the backend directory
   - Check that all schema files are valid GraphQL

2. **Frontend tests fail**
   - Clear node_modules: `rm -rf node_modules && npm install`
   - Check that all dependencies are installed

3. **Port conflicts**
   - Backend default: 8080
   - Frontend default: 5173
   - Change ports in respective config files if needed

### Getting Help

- Check the [API Documentation](api.md)
- Review [Architecture Overview](architecture.md)
- Open an issue on GitHub

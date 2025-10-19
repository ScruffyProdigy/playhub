# PlayHub

PlayHub is a gaming lobby platform that connects players to games and enables trading of digital goods and currency.

## Features

### ✅ Implemented
- **Environment Configuration**: Docker-based runtime environment injection
- **GraphQL API**: Complete GraphQL schema with mock resolvers
- **Frontend Foundation**: React application with testing infrastructure
- **Kubernetes Deployment**: Multi-environment deployment scripts
- **Testing Suite**: Comprehensive unit, integration, and E2E tests
- **CI/CD Pipeline**: GitHub Actions workflows for testing and deployment
- **Database Integration**: PostgreSQL setup with connection management
- **Linting & Code Quality**: ESLint configuration with proper test environment setup

### 🚧 In Development
- **User Authentication**: JWT-based authentication system
- **Game Management**: CRUD operations for games
- **Queue System**: Player queuing and matchmaking

### 📋 Planned
- **Game Integration**: Connection to 3rd party games
- **Digital Trading**: Currency and digital goods trading system
- **Real-time Updates**: WebSocket subscriptions for live updates
- **Payment Processing**: Integration with payment providers
- **Analytics Dashboard**: Usage and performance metrics

## Architecture

This project consists of:

- **Backend**: Go-based GraphQL API with gqlgen
- **Frontend**: React + Vite application
- **Database**: PostgreSQL with connection management

## Getting Started

### Prerequisites

- Go 1.25+
- Node.js 20+
- Docker (optional)

### Quick Setup

```bash
git clone https://github.com/scruffyprodigy/playhub.git
cd playhub
./scripts/setup.sh
./scripts/dev.sh
```

### Manual Setup

See [Development Guide](docs/development.md) for detailed setup instructions.

## Development

### Running Tests

```bash
# All tests
./scripts/test.sh

# Backend only
./scripts/test-backend.sh

# Frontend only
./scripts/test-frontend.sh

# With E2E tests
./scripts/test-frontend.sh --e2e
```

### Development Servers

```bash
# Start both frontend and backend
./scripts/dev.sh
```

### Code Generation

The backend uses gqlgen for GraphQL code generation:

```bash
cd backend
go run github.com/99designs/gqlgen@v0.17.81 generate
```

### Deployment

Deploy to different environments:

```bash
# Local development (minikube)
./deploy-local.sh

# Staging environment
./deploy-staging.sh

# Production environment
./deploy-production.sh
```

See [Environment Configuration](docs/environment-configuration.md) for detailed deployment instructions.

## Documentation

- **[Development Guide](docs/development.md)** - Setup and development workflow
- **[Architecture Overview](docs/architecture.md)** - System design and components
- **[API Documentation](docs/api.md)** - GraphQL API reference
- **[Testing Guide](docs/testing.md)** - Testing strategies and running tests
- **[Environment Configuration](docs/environment-configuration.md)** - Environment setup for different deployments
- **[Contributing Guide](docs/contributing.md)** - How to contribute to the project

## Contributing

See [Contributing Guide](docs/contributing.md) for detailed information.

## License

[Add your license here]

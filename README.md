# PlayHub

PlayHub is a gaming lobby platform that connects players to games and enables trading of digital goods and currency.

## Features

- **Game Queuing**: Join queues for your favorite games
- **Game Integration**: Seamless connection to 3rd party games
- **Digital Trading**: Trade currency for digital goods and vice versa
- **User Management**: Secure user accounts and profiles

## Architecture

This project consists of:

- **Backend**: Go-based GraphQL API with gqlgen
- **Frontend**: React + Vite application
- **Database**: (To be implemented)

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

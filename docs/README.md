# PlayHub Documentation

This directory contains comprehensive documentation for the PlayHub gaming lobby platform.

## 📚 Documentation Structure

- **[Development Guide](development.md)** - How to set up and run the project locally
- **[Architecture Overview](architecture.md)** - System design and component relationships
- **[API Documentation](api.md)** - GraphQL API reference and examples
- **[Environment Configuration](environment-configuration.md)** - Environment setup and configuration system
- **[Contributing Guide](contributing.md)** - How to contribute to the project
- **[Testing Guide](testing.md)** - Testing strategies and running tests

## 🚀 Quick Start

1. **Prerequisites**: Go 1.25+, Node.js 20+, Docker
2. **Clone**: `git clone https://github.com/scruffyprodigy/playhub.git`
3. **Setup**: Run `./scripts/setup.sh` from the project root
4. **Start**: Run `./scripts/dev.sh` to start both frontend and backend

## 📁 Project Structure

```
playhub/
├── backend/          # Go GraphQL API
├── frontend/         # React + Vite application
├── k8s/             # Kubernetes configurations
├── docs/            # This documentation
├── scripts/         # Shared development scripts
└── .github/         # CI/CD workflows
```

## 🔗 External Links

- [GraphQL Schema](backend/graph/schema/)
- [Frontend Components](frontend/src/)
- [Kubernetes Configs](k8s/)
- [GitHub Actions](.github/workflows/)

# Contributing to PlayHub

Thank you for your interest in contributing to PlayHub! This guide will help you get started.

## Getting Started

1. **Fork the repository** on GitHub
2. **Clone your fork** locally
3. **Set up the development environment** using our setup script
4. **Create a feature branch** for your changes
5. **Make your changes** and test them
6. **Submit a pull request**

## Development Setup

### Quick Setup
```bash
git clone https://github.com/your-username/playhub.git
cd playhub
./scripts/setup.sh
```

### Manual Setup
See [Development Guide](development.md) for detailed setup instructions.

## Development Workflow

### 1. Create a Feature Branch
```bash
git checkout -b feature/your-feature-name
```

### 2. Make Your Changes
- Write clean, readable code
- Follow existing code style
- Add tests for new functionality
- Update documentation as needed

### 3. Test Your Changes
```bash
# Run all tests
./scripts/test.sh

# Run specific test suites
./scripts/test-backend.sh
./scripts/test-frontend.sh --e2e
```

### 4. Commit Your Changes
```bash
git add .
git commit -m "feat: add new feature description"
```

### 5. Push and Create Pull Request
```bash
git push origin feature/your-feature-name
```

Then create a pull request on GitHub.

## Code Style

### Backend (Go)
- Follow [Effective Go](https://golang.org/doc/effective_go.html)
- Use `gofmt` for formatting
- Run `go vet` for static analysis
- Write tests for all public functions

### Frontend (JavaScript/React)
- Use ESLint configuration provided
- Follow React best practices
- Write tests for components
- Use TypeScript when possible

### General
- Use meaningful variable and function names
- Write clear comments for complex logic
- Keep functions small and focused
- Follow the existing project structure

## Testing

### Backend Testing
- Write unit tests for all resolvers
- Test error conditions
- Include benchmark tests for performance-critical code
- Ensure drift detection tests pass

### Frontend Testing
- Write unit tests for components
- Test user interactions
- Include integration tests for user journeys
- Run E2E tests for critical paths

### Test Coverage
- Aim for >80% test coverage
- Focus on critical business logic
- Don't test implementation details

## Pull Request Guidelines

### Before Submitting
- [ ] All tests pass
- [ ] Code follows style guidelines
- [ ] Documentation is updated
- [ ] No sensitive data is committed
- [ ] Commit messages are clear

### Pull Request Template
```markdown
## Description
Brief description of changes

## Type of Change
- [ ] Bug fix
- [ ] New feature
- [ ] Breaking change
- [ ] Documentation update

## Testing
- [ ] Unit tests pass
- [ ] Integration tests pass
- [ ] E2E tests pass
- [ ] Manual testing completed

## Checklist
- [ ] Code follows style guidelines
- [ ] Self-review completed
- [ ] Documentation updated
- [ ] No sensitive data committed
```

## Issue Guidelines

### Bug Reports
When reporting bugs, please include:
- Clear description of the issue
- Steps to reproduce
- Expected vs actual behavior
- Environment details (OS, browser, etc.)
- Screenshots if applicable

### Feature Requests
When requesting features, please include:
- Clear description of the feature
- Use case and motivation
- Proposed implementation (if you have ideas)
- Any relevant examples

## Code Review Process

### For Contributors
- Address all review comments
- Keep PRs focused and small
- Respond to feedback promptly
- Be open to suggestions

### For Reviewers
- Be constructive and helpful
- Focus on code quality and correctness
- Test the changes locally when possible
- Approve when ready

## Release Process

### Versioning
We use [Semantic Versioning](https://semver.org/):
- **MAJOR**: Breaking changes
- **MINOR**: New features (backward compatible)
- **PATCH**: Bug fixes (backward compatible)

### Release Notes
Include clear release notes for each version:
- New features
- Bug fixes
- Breaking changes
- Migration instructions (if needed)

## Community Guidelines

### Be Respectful
- Use welcoming and inclusive language
- Be respectful of differing viewpoints
- Accept constructive criticism gracefully
- Focus on what's best for the community

### Be Collaborative
- Help others when you can
- Share knowledge and experience
- Be patient with newcomers
- Work together toward common goals

## Getting Help

### Documentation
- Check existing documentation first
- Look for similar issues or PRs
- Read the codebase for examples

### Communication
- Use GitHub issues for bug reports and feature requests
- Use GitHub discussions for questions and ideas
- Be specific and provide context

### Mentorship
- New contributors are welcome
- Ask questions if you're unsure
- We're here to help you succeed

## Recognition

Contributors will be recognized in:
- CONTRIBUTORS.md file
- Release notes
- Project documentation

Thank you for contributing to PlayHub! 🚀

# Contributing to SPIRIT

Thank you for your interest in contributing to SPIRIT! This document provides guidelines for contributing to the project.

## How to Contribute

### Reporting Issues

1. Check if the issue already exists in the [issue tracker](https://github.com/TheOrionAI/spirit/issues)
2. If not, create a new issue with:
   - Clear description of the problem
   - Steps to reproduce
   - Expected vs actual behavior
   - Environment details (OS, version)

### Suggesting Features

Open an issue with the "enhancement" label describing:
- The feature you'd like to see
- Why it would be useful
- Any implementation ideas

### Code Contributions

1. **Fork** the repository
2. **Clone** your fork locally
3. **Create a branch** for your changes: `git checkout -b feature/my-feature`
4. **Make changes** following our coding standards
5. **Test** your changes locally
6. **Commit** with clear messages: `git commit -m "Add feature X"`
7. **Push** to your fork: `git push origin feature/my-feature`
8. **Open a Pull Request** with:
   - Clear title and description
   - Reference any related issues
   - Screenshots if applicable

## Development Setup

```bash
# Clone your fork
git clone https://github.com/YOUR_USERNAME/spirit.git
cd spirit

# Install Go dependencies
go mod download

# Build
go build -o bin/spirit ./cmd/spirit

# Test
./bin/spirit --version
```

## Coding Standards

- **Go**: Follow standard Go conventions (`gofmt`, `golint`)
- **Comments**: Explain "why", not "what"
- **Tests**: Add tests for new functionality
- **Documentation**: Update README.md if needed

## Commit Message Format

```
type(scope): subject

body (optional)

footer (optional)
```

Types:
- `feat:` New feature
- `fix:` Bug fix
- `docs:` Documentation changes
- `test:` Test changes
- `refactor:` Code restructuring

Example:
```
feat(sync): add verbose flag for detailed output

Adds --verbose flag to spirit sync that lists all tracked
and skipped files during the sync operation.

Closes #42
```

## Pull Request Process

1. Ensure your PR passes all checks
2. Update documentation if needed
3. Link to related issues
4. Request review from maintainers
5. Address review comments
6. PR will be merged by maintainers

## Questions?

Join the discussion:
- Twitter: [@my_self_orion](https://x.com/my_self_orion)
- GitHub: [TheOrionAI/spirit/discussions](https://github.com/TheOrionAI/spirit/discussions)

## Security

If you discover a security vulnerability, please email us directly instead of opening a public issue.

---

Thank you for making SPIRIT better! 🌌

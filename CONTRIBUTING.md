# Contributing to Ship 🚢

Thank you for your interest in contributing to Ship! This document provides guidelines and instructions for contributing.

## Getting Started

### Prerequisites

- Go 1.24.10 or higher
- Git
- Make (optional but recommended)

### Setup Development Environment

1. Fork the repository on GitHub

2. Clone your fork:
```bash
git clone https://github.com/YOUR_USERNAME/ship.git
cd ship
```

3. Add upstream remote:
```bash
git remote add upstream https://github.com/KambojRajan/ship.git
```

4. Install dependencies:
```bash
make deps
```

5. Build and test:
```bash
make build
make test
```

## Development Workflow

### 1. Create a Branch

```bash
git checkout -b feature/your-feature-name
# or
git checkout -b fix/your-bug-fix
```

Branch naming conventions:
- `feature/` - New features
- `fix/` - Bug fixes
- `docs/` - Documentation updates
- `refactor/` - Code refactoring
- `test/` - Test additions or modifications

### 2. Make Your Changes

Write clean, readable code following Go best practices:

```bash
# Format your code
make fmt

# Run linting
make vet

# Run tests frequently
make test
```

### 3. Test Your Changes

```bash
# Run all tests
make test

# Run specific tests
go test ./tests/command_tests/add_test.go -v

# Run with coverage
make test-cover
```

### 4. Commit Your Changes

Write clear, descriptive commit messages:

```bash
git add .
git commit -m "feat: add branch listing command"
```

Commit message format:
- `feat:` - New feature
- `fix:` - Bug fix
- `docs:` - Documentation changes
- `test:` - Test changes
- `refactor:` - Code refactoring
- `chore:` - Maintenance tasks

### 5. Keep Your Fork Updated

```bash
git fetch upstream
git rebase upstream/main
```

### 6. Push and Create Pull Request

```bash
git push origin feature/your-feature-name
```

Then create a Pull Request on GitHub.

## Code Style Guidelines

### Go Code Style

- Follow [Effective Go](https://golang.org/doc/effective_go.html)
- Use `gofmt` for formatting (automatically done by `make fmt`)
- Write clear, self-documenting code
- Add comments for exported functions and types
- Keep functions small and focused

### Example:

```go
// CalculateHash computes the SHA-1 hash of the given content
// and returns it as a hexadecimal string.
func CalculateHash(content []byte) string {
    hash := sha1.New()
    hash.Write(content)
    return hex.EncodeToString(hash.Sum(nil))
}
```

## Testing Guidelines

### Writing Tests

1. Place tests in the `tests/` directory
2. Use the `_test.go` suffix
3. Use descriptive test names
4. Test both success and failure cases

### Example Test:

```go
func TestAddCommand(t *testing.T) {
    // Setup
    testDir := setupTestRepo(t)
    defer cleanupTestRepo(t, testDir)
    
    // Test
    err := commands.Add(testDir)
    
    // Assert
    assert.NoError(t, err)
    assert.FileExists(t, filepath.Join(testDir, ".ship", "index"))
}
```

### Running Tests

```bash
# Run all tests
make test

# Run specific test file
go test ./tests/command_tests/add_test.go -v

# Run with coverage
make test-cover

# Run with race detection
go test -race ./tests/...
```

## Project Structure

```
ship/
├── cmd/              # CLI command definitions (Cobra)
├── commands/         # Command implementations
├── core/            
│   ├── common/      # Common utilities
│   ├── Entities/    # Data structures (Blob, Tree, Commit, etc.)
│   └── utils/       # Helper utilities
├── tests/
│   ├── command_tests/  # Command integration tests
│   └── helpers/        # Test utilities
├── main.go          # Application entry point
├── Makefile         # Build automation
└── install.sh       # Installation script
```

## Adding New Features

### Adding a New Command

1. Create command definition in `cmd/`:
```go
// cmd/your-command.go
var yourCmd = &cobra.Command{
    Use:   "your-command",
    Short: "Brief description",
    Run: func(cmd *cobra.Command, args []string) {
        // Call implementation
    },
}
```

2. Create implementation in `commands/`:
```go
// commands/your-command.go
func YourCommand(args []string) error {
    // Implementation
}
```

3. Add tests in `tests/command_tests/`:
```go
// tests/command_tests/your-command_test.go
func TestYourCommand(t *testing.T) {
    // Test implementation
}
```

4. Update documentation in `README.md`

### Adding a New Entity

1. Create entity in `core/entities/`
2. Implement required methods
3. Add serialization/deserialization
4. Write unit tests

## Pull Request Process

1. **Update Documentation**: Ensure README.md and relevant docs are updated
2. **Add Tests**: All new features should have tests
3. **Run Tests**: Ensure all tests pass with `make test`
4. **Format Code**: Run `make fmt` and `make vet`
5. **Write Clear PR Description**: Explain what and why
6. **Link Issues**: Reference related issues with `Fixes #123`

## Pull Request Checklist

- [ ] Code follows the style guidelines
- [ ] Self-review of code completed
- [ ] Comments added for complex code
- [ ] Documentation updated
- [ ] Tests added/updated and passing
- [ ] No new warnings generated
- [ ] Commit messages are clear

## Reporting Bugs

### Before Submitting

1. Check existing issues
2. Try the latest version
3. Collect information about your environment

### Bug Report Template

```markdown
**Describe the bug**
A clear description of the bug.

**To Reproduce**
Steps to reproduce:
1. Run `ship init .`
2. Run `ship add ...`
3. See error

**Expected behavior**
What you expected to happen.

**Environment**
- OS: [e.g., macOS 13.0]
- Go version: [e.g., 1.24.10]
- Ship version: [e.g., 1.0.0]

**Additional context**
Any other relevant information.
```

## Feature Requests

We welcome feature requests! Please:

1. Check if the feature already exists or is requested
2. Clearly describe the feature and its use case
3. Explain why it would be useful to most users

## Questions?

- Open an issue with the `question` label
- Check existing documentation
- Review closed issues for similar questions

## Code of Conduct

- Be respectful and inclusive
- Focus on constructive feedback
- Help others learn and grow

## License

By contributing, you agree that your contributions will be licensed under the same license as the project.

---

Thank you for contributing to Ship! 🚢

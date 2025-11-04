# git-cc

A terminal-based conventional commit tool for git with an interactive TUI interface.

<a href="https://asciinema.org/a/R8xrOATDrnmO2hDNrdhGGllgH" target="_blank">
    <img src="./docs/demo.svg" alt="DEMO"/>
</a>

## Features

- **Interactive TUI**: User-friendly terminal interface powered by Bubble Tea
- **Conventional Commits**: Enforces conventional commit standards
- **Smart Error Handling**: Clean error display for git operations and pre-commit hooks
- **Git Validation**: Checks for git repository and staged files before starting
- **Error Recovery**: Retry functionality when commits fail
- **Cross-platform**: Works on Linux, macOS, and Windows

## Installation

### Go Install (Recommended)

```bash
go install github.com/denysvitali/git-cc@latest
```

### Pre-built Binaries

Download the appropriate binary from the [Releases](https://github.com/denysvitali/git-cc/releases) page.

### Docker

```bash
# Pull the image
docker pull ghcr.io/denysvitali/git-cc:latest

# Run with git repository mounted
docker run --rm -it -v $(pwd):/repo -w /repo ghcr.io/denysvitali/git-cc:latest
```

### Package Managers (Coming Soon)

We're working on distribution to package managers:

- **Homebrew (macOS)**: `brew install git-cc`
- **Scoop (Windows)**: `scoop install git-cc`
- **Snap (Linux)**: `snap install git-cc`

Check the [Issues](https://github.com/denysvitali/git-cc/issues) for progress or to help set up these distribution channels!

### From Source

```bash
git clone https://github.com/denysvitali/git-cc
cd git-cc
make install
```

## Usage

1. Stage your files with `git add`
2. Run `git cc`
3. Select the type of change from the list
4. Optionally enter a scope (press Enter to skip)
5. Enter your commit message
6. Press Enter to commit

### Keyboard Shortcuts

- `↑/↓` or `j/k`: Navigate through commit types
- `Enter`: Select and move to next step / Commit
- `Ctrl+C` or `q`: Quit the application
- `r`: Retry after a failed commit
- `/`: Start filtering commit types
- `Esc`: Clear filter

## Development

### Prerequisites

- Go 1.21 or later
- Git

### Setup

```bash
# Clone the repository
git clone https://github.com/denysvitali/git-cc
cd git-cc

# Set up development environment
make dev-setup
```

### Common Tasks

```bash
# Build the application
make build

# Run tests
make test

# Run tests with coverage
make test-coverage

# Format code
make fmt

# Run linter
make lint

# Run all quality checks
make quality

# Build for multiple platforms
make build-all
```

### Project Structure

```
git-cc/
├── main.go                 # Application entry point
├── git/                    # Git operations package
│   ├── git.go             # Git command handling
│   └── git_test.go        # Git package tests
├── ui/                     # User interface package
│   ├── model.go           # TUI model and logic
│   └── model_test.go      # UI package tests
├── integration_test.go     # Integration tests
├── .github/workflows/      # GitHub Actions CI/CD
├── .golangci.yml          # Golangci-lint configuration
├── Makefile               # Development tasks
└── README.md              # This file
```

### Testing

The project includes comprehensive test coverage:

- **Unit Tests**: Test individual functions and methods
- **Integration Tests**: Test the full application workflow
- **Git Operations Tests**: Test git command handling in isolation

```bash
# Run all tests
make test

# Run only unit tests
go test ./...

# Run only integration tests
make integration-test

# Generate coverage report
make test-coverage
```

### Contributing

1. Fork the repository
2. Create a feature branch (`git checkout -b feature/amazing-feature`)
3. Make your changes
4. Run tests and ensure all quality checks pass (`make quality`)
5. Commit your changes using conventional commits
6. Push to the branch (`git push origin feature/amazing-feature`)
7. Open a Pull Request

## Code Quality

This project uses several tools to maintain code quality:

- **golangci-lint**: Comprehensive Go linting
- **go vet**: Static analysis for potential issues
- **go test**: Unit and integration testing
- **GitHub Actions**: CI/CD pipeline with automated testing
- **Pre-commit hooks**: Ensure code quality before commits

### CI/CD Pipeline

The GitHub Actions workflow includes:

- Multi-version Go testing (1.21, 1.22, 1.23)
- Linting with golangci-lint
- Unit and integration testing with coverage
- Security scanning with gosec
- Multi-platform builds

## Conventional Commits

This tool helps you create conventional commits following the [Conventional Commits specification](https://www.conventionalcommits.org/en/v1.0.0/).

### Commit Types

- `feat`: A new feature
- `fix`: A bug fix
- `docs`: Documentation only changes
- `style`: Changes that do not affect the meaning of the code
- `refactor`: A code change that neither fixes a bug nor adds a feature
- `perf`: A code change that improves performance
- `test`: Adding missing tests or correcting existing tests
- `build`: Changes that affect the build system or external dependencies
- `ci`: Changes to CI configuration files and scripts
- `chore`: Other changes that don't modify src or test files

### Format

```
<type>[optional scope]: <description>

[optional body]

[optional footer(s)]
```

Examples:
- `feat: add user authentication`
- `fix(auth): resolve login timeout issue`
- `docs: update API documentation`

## Error Handling

git-cc provides smart error handling for common git issues:

- **Pre-commit Hook Failures**: Clean display of hook output with retry option
- **Merge Conflicts**: Clear indication of conflict resolution needed
- **Empty Repository**: Helpful message when no files are staged
- **Not a Git Repository**: Validation before starting the TUI

## License

This project is licensed under the MIT License - see the [LICENSE](LICENSE) file for details.

## References

- [Conventional Commits](https://www.conventionalcommits.org/en/v1.0.0/)
- [Bubble Tea](https://github.com/charmbracelet/bubbletea)
- [Charm Bubbles](https://github.com/charmbracelet/bubbles)

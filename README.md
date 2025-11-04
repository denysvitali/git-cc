# git-cc

Interactive TUI for conventional commits.

![Demo](docs/demo.svg)

## Installation

```bash
go install github.com/denysvitali/git-cc@latest
```

Or download from [releases](https://github.com/denysvitali/git-cc/releases).

## Usage

1. Stage files: `git add .`
2. Run: `git cc`
3. Select type, write message
4. Press Enter to commit

### Controls
- `↑/↓` or `j/k`: Navigate
- `Enter`: Select/Commit
- `Ctrl+C` or `q`: Quit
- `r`: Retry after failure

## Conventional Commits

Follows the [Conventional Commits specification](https://www.conventionalcommits.org).

### Types
- `feat`: New feature
- `fix`: Bug fix
- `docs`: Documentation
- `style`: Formatting
- `refactor`: Code refactoring
- `perf`: Performance
- `test`: Tests
- `build`: Build system
- `ci`: CI configuration
- `chore`: Other changes

Format: `<type>[optional scope]: <description>`

## Development

```bash
git clone https://github.com/denysvitali/git-cc
cd git-cc
make install
```

## License

MIT
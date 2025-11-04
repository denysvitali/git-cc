# Release Process with Goreleaser

This document describes how git-cc uses Goreleaser for automated releases.

## Overview

git-cc uses Goreleaser to automate the entire release process, including:

- Multi-platform binary builds
- Package creation (archives, Docker images)
- Distribution to multiple platforms (Homebrew, Scoop, Snap, Docker)
- GitHub release automation
- Changelog generation

## Configuration

The Goreleaser configuration is in `.goreleaser.yml` and includes:

### Build Configuration

- **Platforms**: Linux, Windows, macOS (amd64/arm64)
- **Optimization**: Static binaries with build metadata injection
- **Version Info**: Git tag, commit hash, build date, and builder

### Distribution Channels

1. **GitHub Releases**
   - Automatic release creation on git tags
   - Multi-platform asset uploads
   - Changelog generation
   - Checksums for integrity verification

2. **Homebrew (macOS)**
   - Automatic formula generation
   - Test validation
   - Installation via `brew install git-cc`

3. **Scoop (Windows)**
   - Windows package manager support
   - Automatic manifest updates
   - Installation via `scoop install git-cc`

4. **Snap (Linux)**
   - Cross-distribution Snap package
   - Confined environment for security
   - Installation via `snap install git-cc`

5. **Docker Images**
   - Multi-architecture builds
   - GitHub Container Registry publishing
   - Semantic version tagging

## Release Process

### Prerequisites

1. **Git Tag**: Create a new git tag
   ```bash
   git tag v1.0.0
   git push origin v1.0.0
   ```

2. **GitHub Token**: Set up `GITHUB_TOKEN` environment variable
   ```bash
   export GITHUB_TOKEN=your_github_token
   ```

### Automated Release

The GitHub Actions workflow (`.github/workflows/release.yml`) automatically:

1. Triggers on git tags (`v*`) or manual workflow dispatch
2. Runs all quality checks (tests, linting, security scans)
3. Executes Goreleaser release process
4. Publishes to all distribution channels

### Local Release Testing

You can test the release process locally:

```bash
# Test configuration
make release-check

# Build snapshot
make release-snapshot

# Test full release process (no publishing)
make release-test
```

## Development Workflow

### Pre-release Preparation

```bash
# Run all checks
make ci-check

# Test build with version info
VERSION=1.0.0-rc1 make build

# Test Goreleaser locally
make release-test
```

### Creating a Release

1. **Finalize Changes**
   ```bash
   make clean
   make ci-check
   ```

2. **Tag Release**
   ```bash
   git tag -a v1.0.0 -m "Release version 1.0.0"
   git push origin v1.0.0
   ```

3. **Monitor Release**
   - GitHub Actions will automatically handle the release
   - Check the Actions tab for progress
   - Verify distribution channels updated

## Distribution Channels

### Homebrew

```bash
# Install
brew tap denysvitali/homebrew-tap
brew install git-cc

# Update
brew upgrade git-cc
```

### Scoop

```bash
# Install
scoop bucket add denysvitali-scoop https://github.com/denysvitali/scoop-bucket
scoop install git-cc

# Update
scoop update git-cc
```

### Snap

```bash
# Install
sudo snap install git-cc

# Update
sudo snap refresh git-cc
```

### Docker

```bash
# Pull image
docker pull ghcr.io/denysvitali/git-cc:latest

# Run
docker run --rm -it ghcr.io/denysvitali/git-cc:latest --version

# With git repository mounted
docker run --rm -it -v $(pwd):/repo -w /repo ghcr.io/denysvitali/git-cc:latest
```

## Changelog Generation

Goreleaser automatically generates changelogs based on commit messages:

- **Features**: Commits starting with `feat:`
- **Bug fixes**: Commits starting with `fix:`
- **Documentation**: Commits starting with `docs:`
- **Other**: All other commits

The changelog is organized into groups and excludes merge commits and chores.

## Version Information

Built binaries include version information:

```bash
$ git-cc --version
git-cc 1.0.0
  Commit: abc123
  Built: 2024-01-15T10:30:00Z
  Built by: goreleaser
```

## Troubleshooting

### Build Failures

1. **Missing Tags**: Ensure you have created a git tag
   ```bash
   git tag -a v1.0.0 -m "Release v1.0.0"
   ```

2. **GitHub Token Issues**: Verify token has necessary permissions
   - `repo` scope for releases
   - `workflow` scope for GitHub Actions

3. **Build Errors**: Check logs in GitHub Actions for detailed error messages

### Distribution Issues

1. **Homebrew**: Formula updates take time to propagate
   - Allow 30 minutes for formula updates
   - Check brew tap repository

2. **Scoop**: Manual updates may be required
   - Update scoop bucket manually if needed

3. **Docker**: Check image tags and platform support
   - Verify multi-platform builds worked correctly
   - Check GitHub Container Registry permissions

## Customization

To modify the release process:

1. **Edit `.goreleaser.yml`**
2. **Test locally**: `make release-test`
3. **Update CI workflow** if needed

### Adding New Distribution Channels

Add new sections to `.goreleaser.yml` for:
- Package managers
- Container registries
- Cloud storage providers

### Modifying Build Options

Update the `builds` section to:
- Add new platforms
- Change build flags
- Modify binary names

## Security Considerations

- **Binary Verification**: Always verify checksums when downloading
- **Container Security**: Use official Docker images from registries
- **Token Security**: Store GitHub tokens securely in CI/CD
- **Code Review**: Review `.goreleaser.yml` changes carefully

## Rollback Process

If a release needs to be rolled back:

1. **Delete Tag**: Remove the git tag
   ```bash
   git tag -d v1.0.0
   git push origin :refs/tags/v1.0.0
   ```

2. **New Release**: Create a new version tag
   ```bash
   git tag -a v1.0.1 -m "Rollback fix"
   git push origin v1.0.1
   ```

3. **Communication**: Announce the rollback and new version

## Monitoring

Monitor releases through:

- **GitHub**: Release page shows assets and statistics
- **Homebrew**: Formula install statistics
- **Docker**: Image pull counts on GitHub Container Registry
- **Scoop**: Download statistics if available

## Support

For release-related issues:

1. Check existing GitHub issues
2. Create new issue with detailed description
3. Include logs and error messages
4. Mention the platform and version affected
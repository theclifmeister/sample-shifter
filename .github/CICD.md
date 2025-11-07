# CI/CD Pipeline Documentation

This repository uses GitHub Actions for Continuous Integration and Continuous Deployment.

## Workflows

### Test and Lint (`test.yml`)

**Triggers:**
- On every pull request to `main` branch
- On every push to `main` branch

**What it does:**
1. **Test Job**: Runs tests on multiple platforms and Go versions
   - Platforms: Ubuntu, macOS, Windows
   - Go versions: 1.22, 1.23
   - Runs tests with race detection and generates coverage reports
   - Uploads coverage to Codecov (optional)

2. **Lint Job**: Checks code quality
   - Runs `go fmt` to ensure code is properly formatted
   - Runs `go vet` to catch common mistakes

3. **Build Job**: Verifies the binary builds correctly
   - Builds the project
   - Tests that the binary runs with `--help` flag

**Status**: Check the badge in README.md to see the current status.

### Release (`release.yml`)

**Triggers:**
- On push of version tags matching pattern `v*.*.*` (e.g., `v1.0.0`, `v2.1.3`)

**What it does:**
1. Runs all tests to ensure code quality
2. Builds binaries for multiple platforms:
   - Linux (amd64, arm64)
   - macOS (amd64, arm64)
   - Windows (amd64, arm64)
3. Creates compressed archives:
   - `.tar.gz` for Linux and macOS
   - `.zip` for Windows
4. Generates changelog from git commits
5. Creates a GitHub release with all binaries attached

## Creating a Release

To create a new release:

1. Ensure all changes are committed and pushed to `main`
2. Create and push a version tag:
   ```bash
   git tag v1.0.0
   git push origin v1.0.0
   ```
3. The release workflow will automatically:
   - Build binaries for all platforms
   - Create a GitHub release
   - Upload all binaries to the release

## Download Releases

Users can download pre-compiled binaries from the [Releases](https://github.com/theclifmeister/sample-shifter/releases) page.

Available platforms:
- `sample-shifter-linux-amd64.tar.gz` - Linux 64-bit (Intel/AMD)
- `sample-shifter-linux-arm64.tar.gz` - Linux 64-bit (ARM)
- `sample-shifter-darwin-amd64.tar.gz` - macOS (Intel)
- `sample-shifter-darwin-arm64.tar.gz` - macOS (Apple Silicon)
- `sample-shifter-windows-amd64.zip` - Windows 64-bit (Intel/AMD)
- `sample-shifter-windows-arm64.zip` - Windows 64-bit (ARM)

## Development Workflow

1. Create a feature branch
2. Make changes
3. Push to GitHub and create a PR
4. CI will automatically run tests and linting
5. Once tests pass and PR is approved, merge to `main`
6. After merging, create a release tag to trigger the release workflow

## Troubleshooting

### Tests fail on PR
- Check the Actions tab for detailed logs
- Run tests locally: `go test -v ./...`
- Ensure code is formatted: `go fmt ./...`
- Check for issues: `go vet ./...`

### Release fails
- Ensure the tag follows the `v*.*.*` pattern
- Check that all tests pass before creating the tag
- Review the Actions tab for detailed error messages

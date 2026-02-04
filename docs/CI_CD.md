# CI/CD and Releases

This document describes the continuous integration and release process for MyGoMetrics.

## Overview

MyGoMetrics uses GitHub Actions for continuous integration and automated releases. The CI/CD pipeline ensures code quality through automated testing and provides reproducible builds via container images.

## CI Pipeline

The CI pipeline is defined in `.github/workflows/ci.yml` and runs on every push and pull request.

### Automated Checks

Every push and pull request triggers:

- **Go vet** - Static analysis for common errors
- **Go test** - Unit tests with race detection
- **Docker build** - Validates the Dockerfile builds successfully

These checks ensure code quality and prevent regressions before code is merged.

### Trigger Conditions

The workflow triggers only on changes to relevant paths:
- `.github/` - CI/CD configuration
- `cmd/` - Application entry point
- `internal/` - Internal packages
- `go.mod` - Go dependencies
- `Dockerfile` - Container build
- `helm/` - Kubernetes deployment

This reduces unnecessary pipeline runs and optimizes build times.

## Releases

### Automated Release Process

When a version tag is pushed (e.g., `v0.7.0`), the CI pipeline automatically:

1. Runs all tests and linting
2. Builds a Docker image for `linux/amd64`
3. Pushes the image to GitHub Container Registry (ghcr.io)

This ensures all releases are reproducible and tested.

### Container Image Registry

**Image Location:**

Container images are published to GitHub Container Registry:
```
ghcr.io/<owner>/mygometrics:<tag>
```

Where:
- `<owner>` is the GitHub repository owner (e.g., `sheliakhin-golang-portfolio`)
- `<tag>` matches the Git tag (e.g., `v0.7.0`)

**Authentication:**

The CI pipeline uses the built-in `GITHUB_TOKEN` for authentication. No manual secret configuration is required.

### Using Released Images

**Pulling Released Images:**

```bash
docker pull ghcr.io/sheliakhin-golang-portfolio/mygometrics:v0.7.0
```

**Using with Docker:**

```bash
docker run --rm -p 9000:9000 \
  ghcr.io/sheliakhin-golang-portfolio/mygometrics:v0.7.0
```

**Using with Helm:**

The Helm chart defaults to using images from ghcr.io. Update the `image.tag` value to use a specific release:

```bash
helm install mygometrics ./helm/mygometrics \
  --set image.tag=v0.7.0
```

Or in `values.yaml`:

```yaml
image:
  registry: ghcr.io
  repository: sheliakhin-golang-portfolio/mygometrics
  tag: "v0.7.0"
```

## Creating a Release

To create a new release, follow these steps:

### 1. Update Version Numbers

Update version numbers in relevant files:
- `helm/mygometrics/Chart.yaml` - Update `version` field
- Any other version references as needed

### 2. Update Changelog

Document all changes in `CHANGELOG.md` following the existing format.

### 3. Commit Changes

```bash
git add .
git commit -m "Release v0.X.0"
```

### 4. Create and Push Tag

```bash
git tag v0.X.0
git push origin main
git push origin v0.X.0
```

### 5. Verify Release

The CI pipeline will automatically:
- Run all tests
- Build the Docker image
- Push to ghcr.io

Monitor the GitHub Actions workflow to ensure the release completes successfully.

## CI/CD Configuration

The GitHub Actions workflow configuration is located at:
```
.github/workflows/ci.yml
```

### Key Configuration Details

- **Go Version**: Automatically detected from `go.mod`
- **Working Directory**: Repository root (`MyGoMetrics`)
- **Docker Build Context**: Repository root (`.`)
- **Dockerfile Location**: `./Dockerfile`
- **Target Platform**: `linux/amd64` (single architecture)

### Future Enhancements

**Multi-Architecture Support:**

Currently, images are built for `linux/amd64` only. Future enhancements may include:
- Multi-architecture builds (`linux/amd64`, `linux/arm64`)
- Docker Buildx with manifest lists
- Automatic platform selection for ARM-based systems

For more details on CI/CD architecture decisions, see [DECISIONS.md](./DECISIONS.md) Section 14.

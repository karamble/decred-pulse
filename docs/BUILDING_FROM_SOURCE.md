# Building dcrd from Source - Implementation Guide

This document explains why and how we build dcrd from the official source repository instead of using the outdated Docker Hub image.

## Problem: Outdated Docker Image

The official `decred/dcrd:latest` image on Docker Hub:
- ❌ **Last updated**: Over 6 years ago
- ❌ **Security**: Missing years of security patches
- ❌ **Features**: Missing latest features and improvements
- ❌ **Unmaintained**: No longer receiving updates

**Source**: The image is deprecated and should not be used for production.

## Solution: Build from Official Source

We build dcrd directly from the official GitHub repository:
- ✅ **Source**: https://github.com/decred/dcrd
- ✅ **Docker Guide**: https://github.com/decred/dcrd/tree/master/contrib/docker
- ✅ **Latest Code**: Always builds from current source
- ✅ **Security**: Includes all latest security patches
- ✅ **Version Control**: Choose any version/tag/branch

## Implementation

### Directory Structure

```
umbrel-decred-dash-main/
├── dcrd/
│   ├── Dockerfile          # Multi-stage build from source
│   └── README.md           # Build documentation
├── docker compose.yml      # Uses build: ./dcrd
├── dcrd.conf              # Configuration file
└── env.example            # Includes DCRD_VERSION option
```

### Dockerfile Overview

Located in `dcrd/Dockerfile`:

```dockerfile
# Build Stage (golang:1.21-alpine)
- Clone official dcrd repository
- Build dcrd binary
- Build dcrctl binary

# Runtime Stage (alpine:latest)
- Minimal Alpine base (~30-40MB)
- Only compiled binaries
- Runs as non-root user (dcrd)
- Secure and minimal
```

**Key Features**:
- Multi-stage build (no build artifacts in final image)
- Configurable version via `DCRD_VERSION` build arg
- Security: runs as non-root user
- Size: ~30-40MB final image (vs ~500MB build stage)

### Docker Compose Configuration

```yaml
dcrd:
  build:
    context: ./dcrd
    args:
      DCRD_VERSION: ${DCRD_VERSION:-master}
  volumes:
    - dcrd-data:/home/dcrd/.dcrd  # Note: /home/dcrd not /root
```

**Changes from previous version**:
- `image: decred/dcrd:latest` → `build: ./dcrd`
- `/root/.dcrd` → `/home/dcrd/.dcrd` (non-root user)
- Added `DCRD_VERSION` build argument

## Usage

### Basic Usage

```bash
# First time (builds dcrd)
docker compose up -d
# Takes 5-10 minutes on first build

# Subsequent starts (uses cached build)
docker compose up -d
# Instant start
```

### Build Specific Version

**Option 1: Environment Variable**

Edit `.env`:
```env
DCRD_VERSION=release-v2.0.6
```

Build:
```bash
docker compose build dcrd
docker compose up -d dcrd
```

**Option 2: Command Line**

```bash
DCRD_VERSION=release-v2.0.6 docker compose build dcrd
docker compose up -d dcrd
```

**Option 3: Makefile**

```bash
make build-dcrd VERSION=release-v2.0.6
```

### Available Versions

Browse at: https://github.com/decred/dcrd/tags

Examples:
- `master` - Latest development (default)
- `release-v2.0.6` - Stable release v2.0.6
- `release-v2.0.5` - Stable release v2.0.5
- `release-v2.0.4` - Stable release v2.0.4
- `release-v1.8.1` - Older stable release

Or use specific commit hash:
```bash
DCRD_VERSION=abc123def456 docker compose build dcrd
```

## Build Process

### Timeline

1. **Clone Repository** (30 seconds)
   - Downloads dcrd source from GitHub
   
2. **Download Dependencies** (2-3 minutes)
   - Downloads Go modules and dependencies
   
3. **Compile** (2-4 minutes)
   - Builds dcrd binary
   - Builds dcrctl binary
   
4. **Create Runtime Image** (30 seconds)
   - Copies binaries to minimal Alpine image
   - Sets up non-root user

**Total**: 5-10 minutes (first build)  
**Cached**: Instant (subsequent builds)

### Build Output

```
[+] Building 428.3s (16/16) FINISHED
 => [internal] load build definition from Dockerfile
 => [internal] load .dockerignore
 => [build 2/4] RUN apk add --no-cache git
 => [build 3/4] RUN git clone --depth 1 --branch master https://github.com/decred/dcrd.git
 => [build 4/4] RUN go install . ./cmd/dcrctl
 => [stage-1 2/5] RUN apk add --no-cache ca-certificates tzdata
 => [stage-1 3/5] COPY --from=builder /go/bin/dcrd /usr/local/bin/
 => [stage-1 4/5] COPY --from=builder /go/bin/dcrctl /usr/local/bin/
 => exporting to image
```

## Updating dcrd

### Update to Latest

```bash
# Rebuild from latest source
docker compose build --no-cache dcrd
docker compose up -d dcrd

# Or use Makefile
make update-dcrd
```

**Important**: Blockchain data persists in the volume, so no need to re-sync!

### Update to Specific Version

```bash
# Set version in .env
echo "DCRD_VERSION=release-v2.0.7" >> .env

# Rebuild
docker compose build --no-cache dcrd
docker compose up -d dcrd
```

## Comparison: Old vs New

| Aspect | Old (Docker Hub Image) | New (Build from Source) |
|--------|------------------------|-------------------------|
| **Image Source** | Docker Hub | Official GitHub |
| **Last Updated** | 6+ years ago | Current (you control) |
| **Version** | Very old, fixed | Latest or any version |
| **Security** | Outdated | Latest patches |
| **Customization** | None | Full control |
| **Setup Time** | Instant pull | 5-10 min first build |
| **Subsequent Starts** | Instant | Instant (cached) |
| **Image Size** | Unknown | ~30-40MB |
| **User** | root | dcrd (non-root) |
| **Maintenance** | Abandoned | Active official repo |
| **Trust** | Uncertain | Official source code |

## Benefits

### Security
- ✅ Latest security patches
- ✅ Runs as non-root user
- ✅ Minimal attack surface
- ✅ Build from trusted source

### Flexibility
- ✅ Choose any version
- ✅ Build from any branch
- ✅ Build from any commit
- ✅ Easy updates

### Maintenance
- ✅ Active development
- ✅ Community support
- ✅ Regular updates
- ✅ Bug fixes

### Transparency
- ✅ Know exactly what's running
- ✅ Audit source code
- ✅ Build logs visible
- ✅ Reproducible builds

## Makefile Commands

We added several commands for easier management:

```bash
make update-dcrd              # Update dcrd to latest
make build-dcrd VERSION=v2.0.6  # Build specific version
make update                   # Update all services
```

Full list: Run `make help`

## Troubleshooting

### Build Fails

**Problem**: Git clone fails
```
Error: failed to fetch repository
```

**Solution**: Check internet connection, GitHub may be down, try again

---

**Problem**: Out of disk space
```
Error: no space left on device
```

**Solution**: 
```bash
docker system df              # Check usage
docker system prune -a        # Clean up
docker image prune            # Remove old images
```

---

**Problem**: Build is very slow
```
Building... (taking forever)
```

**Solution**: 
- First build is always slow (5-10 minutes)
- Subsequent builds use cache (instant)
- Ensure good internet connection
- Check CPU/memory availability

### Runtime Issues

**Problem**: Container exits immediately
```
dcrd exited with code 1
```

**Solution**: Check logs
```bash
docker compose logs dcrd
```

Common issues:
- Invalid RPC credentials
- Port already in use
- Corrupted data directory

---

**Problem**: Can't connect to RPC
```
Error: connection refused
```

**Solution**: 
- Wait for dcrd to start (check logs)
- Verify credentials match in .env
- Check healthcheck status: `docker compose ps`

## Performance

### Build Performance

**Hardware Impact**:
- **CPU**: High during build (compiling Go code)
- **Memory**: ~2-4GB during build
- **Disk**: ~500MB during build, ~40MB final
- **Network**: ~200MB download (source + dependencies)

**Optimization**:
- Use Docker build cache (don't use --no-cache unnecessarily)
- Multi-core CPU helps (Go compiles in parallel)
- SSD recommended for faster builds

### Runtime Performance

No difference from pre-built image:
- Same dcrd binary
- Same resource usage
- Same performance characteristics

## CI/CD Integration

### Automated Builds

**GitHub Actions** example:

```yaml
name: Build dcrd
on: [push]
jobs:
  build:
    runs-on: ubuntu-latest
    steps:
      - uses: actions/checkout@v2
      - name: Build dcrd
        run: docker compose build dcrd
```

**GitLab CI** example:

```yaml
build:
  image: docker:latest
  services:
    - docker:dind
  script:
    - docker compose build dcrd
```

## Production Recommendations

### Version Pinning

For production, pin to specific stable releases:

```env
# .env
DCRD_VERSION=release-v2.0.6  # Pin to stable release
```

**Don't use `master` in production** - it's development code.

### Regular Updates

1. Monitor releases: https://github.com/decred/dcrd/releases
2. Test new versions in staging
3. Update production after testing
4. Keep blockchain data backed up

### Build Caching

For faster deployments:
- Use CI/CD to pre-build images
- Push to private registry
- Pull pre-built images in production

## Resources

### Documentation
- **dcrd Repository**: https://github.com/decred/dcrd
- **Docker Contrib**: https://github.com/decred/dcrd/tree/master/contrib/docker
- **Releases**: https://github.com/decred/dcrd/releases
- **Our Build Docs**: `dcrd/README.md`

### Support
- **dcrd Issues**: https://github.com/decred/dcrd/issues
- **Decred Matrix**: https://chat.decred.org/
- **Decred Discord**: https://discord.gg/decred

## License

dcrd is licensed under the ISC License.
See: https://github.com/decred/dcrd/blob/master/LICENSE


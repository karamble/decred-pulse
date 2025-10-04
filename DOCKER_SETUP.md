# Docker Setup Guide

Complete guide to running the Decred Dashboard with a full dcrd node using Docker Compose.

## Overview

The Docker Compose stack includes:

1. **dcrd** - Full Decred node (from official [decred/dcrd](https://github.com/decred/dcrd) image)
2. **backend** - Go API server
3. **frontend** - React dashboard UI

All services are connected via a Docker bridge network with proper health checks and dependencies.

## Architecture

```
┌─────────────────┐
│   Frontend      │  Port 3000
│   (Nginx)       │
└────────┬────────┘
         │
         ▼ HTTP
┌─────────────────┐
│   Backend       │  Port 8080
│   (Go API)      │
└────────┬────────┘
         │
         ▼ RPC
┌─────────────────┐
│   dcrd Node     │  Ports 9108 (P2P), 9109 (RPC)
│   (Full Node)   │
└─────────────────┘
```

## Prerequisites

- Docker 20.10+
- Docker Compose 2.0+
- At least 10GB free disk space (for blockchain data)
- 2GB+ RAM recommended

## Quick Start

### 1. Configure Credentials

```bash
# Copy the example environment file
cp .env.example .env

# Edit the file with your secure password
nano .env
```

Set a strong password:
```env
DCRD_RPC_USER=decred
DCRD_RPC_PASS=your_secure_password_here
```

### 2. Start the Stack

```bash
# Build and start all services
docker compose up -d

# View logs
docker compose logs -f

# View specific service logs
docker compose logs -f dcrd
docker compose logs -f backend
docker compose logs -f frontend
```

### 3. Access the Dashboard

Open your browser to: **http://localhost:3000**

## Initial Sync

The dcrd node needs to sync the entire Decred blockchain on first run:

- **Mainnet**: ~8-10GB (may take several hours)
- **Testnet**: ~1-2GB (faster sync)

Monitor sync progress:

```bash
# Watch dcrd logs
docker compose logs -f dcrd

# Check node status via dcrctl
docker exec decred-dcrd-node dcrctl \
  --rpcuser=decred \
  --rpcpass=your_password \
  --rpcserver=127.0.0.1:9109 \
  --notls \
  getinfo
```

The dashboard will show "syncing" status with progress bar during initial sync.

## Service Details

### dcrd Service

**Build**: Built from source via `./dcrd/Dockerfile`

**Source**: https://github.com/decred/dcrd (official Decred repository)

**Why Build from Source?**: The official Docker Hub image is outdated (6+ years old). Building from source ensures latest version with security updates.

**Configuration**:
- P2P Port: `9108` (for blockchain sync)
- RPC Port: `9109` (for API access)
- Data Volume: `dcrd-data` (persistent blockchain storage)
- Config File: `./dcrd.conf` (mounted read-only)
- Build Version: Configurable via `DCRD_VERSION` env var (default: `master`)

**Health Check**: 
- Runs every 30 seconds
- Uses `dcrctl getinfo` to verify node is responsive
- Backend waits for healthy status before starting

**Key Features**:
- `--notls`: Disabled TLS for local container network
- `--txindex`: Full transaction index for complete data access
- `--rpclisten=0.0.0.0:9109`: Listen on all interfaces within container

### Backend Service

**Build**: `./backend/Dockerfile`

**Configuration**:
- API Port: `8080`
- Connects to: `dcrd:9109` (via Docker network)
- Depends on: dcrd (waits for healthy status)

### Frontend Service

**Build**: `./frontend/Dockerfile`

**Configuration**:
- Web Port: `3000` (mapped to 80 internally)
- API URL: `http://localhost:8080/api`
- Depends on: backend

## Docker Commands

### Start/Stop

```bash
# Start all services
docker compose up -d

# Stop all services
docker compose down

# Stop and remove volumes (WARNING: deletes blockchain data)
docker compose down -v

# Restart a specific service
docker compose restart dcrd
docker compose restart backend
docker compose restart frontend
```

### Viewing Logs

```bash
# All services
docker compose logs -f

# Specific service
docker compose logs -f dcrd
docker compose logs -f backend
docker compose logs -f frontend

# Last 100 lines
docker compose logs --tail=100 dcrd
```

### Service Status

```bash
# Check running services
docker compose ps

# Check resource usage
docker stats

# Check specific container
docker inspect decred-dcrd-node
```

### Accessing Containers

```bash
# Execute commands in dcrd container
docker exec -it decred-dcrd-node /bin/sh

# Use dcrctl to query node
docker exec decred-dcrd-node dcrctl \
  --rpcuser=decred \
  --rpcpass=your_password \
  --rpcserver=127.0.0.1:9109 \
  --notls \
  getblockcount

# Access backend container
docker exec -it decred-dashboard-backend /bin/sh

# Access frontend container
docker exec -it decred-dashboard-frontend /bin/sh
```

## Data Persistence

Blockchain data is stored in a Docker volume: `dcrd-data`

```bash
# Inspect volume
docker volume inspect umbrel-decred-dash-main_dcrd-data

# Backup volume
docker run --rm -v umbrel-decred-dash-main_dcrd-data:/data \
  -v $(pwd):/backup alpine \
  tar czf /backup/dcrd-backup.tar.gz -C /data .

# Restore volume
docker run --rm -v umbrel-decred-dash-main_dcrd-data:/data \
  -v $(pwd):/backup alpine \
  tar xzf /backup/dcrd-backup.tar.gz -C /data
```

## Network Configuration

All services communicate via the `decred-network` bridge network:

```bash
# Inspect network
docker network inspect umbrel-decred-dash-main_decred-network

# View connected containers
docker network inspect umbrel-decred-dash-main_decred-network | grep Name
```

Service DNS names within the network:
- `dcrd` → dcrd node
- `backend` → API server
- `frontend` → Web UI

## Using Testnet

To run on testnet instead of mainnet:

1. **Edit `.env`**:
   ```env
   DCRD_RPC_USER=decred
   DCRD_RPC_PASS=your_password
   DCRD_TESTNET=1
   ```

2. **Update `docker compose.yml`**:
   ```yaml
   dcrd:
     command: >
       --testnet
       --rpcuser=${DCRD_RPC_USER:-decred}
       --rpcpass=${DCRD_RPC_PASS:-decredpass}
       --rpclisten=0.0.0.0:9109
       --notls
       --txindex
   ```

3. **Restart**:
   ```bash
   docker compose down
   docker compose up -d
   ```

## Troubleshooting

### dcrd won't start

```bash
# Check logs
docker compose logs dcrd

# Common issues:
# 1. Port 9108 or 9109 already in use
# 2. Insufficient disk space
# 3. Corrupted blockchain data

# Solution: Check ports
netstat -tulpn | grep 910

# Solution: Clean data and restart
docker compose down -v
docker compose up -d
```

### Backend can't connect to dcrd

```bash
# Check dcrd health
docker compose ps dcrd

# Verify RPC credentials
docker exec decred-dcrd-node dcrctl \
  --rpcuser=decred \
  --rpcpass=your_password \
  --rpcserver=127.0.0.1:9109 \
  --notls \
  getinfo

# Check backend logs
docker compose logs backend

# Restart backend
docker compose restart backend
```

### Slow blockchain sync

```bash
# Monitor sync progress
docker compose logs -f dcrd | grep -i "sync\|block"

# Check peer connections
docker exec decred-dcrd-node dcrctl \
  --rpcuser=decred \
  --rpcpass=your_password \
  --rpcserver=127.0.0.1:9109 \
  --notls \
  getpeerinfo

# Typical sync times:
# - Mainnet: 4-8 hours (depending on hardware/network)
# - Testnet: 30-60 minutes
```

### Frontend can't reach backend

```bash
# Check all services are running
docker compose ps

# Verify network connectivity
docker exec decred-dashboard-frontend ping -c 3 backend

# Check CORS settings in backend logs
docker compose logs backend | grep -i cors

# Restart frontend
docker compose restart frontend
```

## Production Deployment

For production use:

1. **Use TLS for dcrd**:
   - Generate certificates
   - Update `--notls` to use `--rpccert` and `--rpckey`
   - Update backend to use TLS

2. **Secure RPC credentials**:
   - Use strong passwords
   - Use Docker secrets instead of environment variables

3. **Enable TLS for frontend/backend**:
   - Add reverse proxy (Nginx/Traefik)
   - Use Let's Encrypt certificates

4. **Resource limits**:
   ```yaml
   dcrd:
     deploy:
       resources:
         limits:
           cpus: '2'
           memory: 4G
   ```

5. **Backups**:
   - Regular volume backups
   - Monitor disk usage
   - Set up alerts

## Building dcrd from Source

The dcrd service is built from the official Decred source repository, not from an outdated Docker image.

### Why Build from Source?

The official `decred/dcrd:latest` Docker Hub image hasn't been updated in over 6 years. Building from source ensures:
- ✅ Latest security patches
- ✅ Latest features
- ✅ Ability to choose specific versions
- ✅ Build from trusted official source

### Build Process

**First time** (automatic when you run `docker compose up`):

```bash
docker compose up -d
# This automatically builds dcrd from source (takes 5-10 minutes)
```

**Manual build**:

```bash
# Build all services
docker compose build

# Build only dcrd
docker compose build dcrd

# Clean build (no cache)
docker compose build --no-cache dcrd
```

### Build Specific Version

**Using master branch (default - latest development)**:
```bash
docker compose build dcrd
```

**Using specific release tag**:

Edit `.env`:
```env
DCRD_VERSION=release-v2.0.6
```

Then rebuild:
```bash
docker compose build dcrd
docker compose up -d dcrd
```

**Available versions**: See https://github.com/decred/dcrd/tags

### Build Time

- **First build**: 5-10 minutes (downloads source, compiles Go code)
- **Subsequent builds**: Faster (uses Docker cache)
- **Final image size**: ~30-40MB (minimal Alpine + binaries)

See `dcrd/README.md` for detailed build documentation.

## Updating

### Update dcrd to Latest

```bash
# Rebuild dcrd from latest source
docker compose build --no-cache dcrd
docker compose up -d dcrd
```

**Note**: Blockchain data persists, so no need to re-sync.

### Update Frontend/Backend

```bash
# Rebuild custom images
docker compose build --no-cache backend frontend

# Restart with new images
docker compose up -d

# Clean up old images
docker image prune -a
```

## Performance Tuning

### dcrd Optimization

Add to `docker compose.yml` command:

```yaml
command: >
  --rpcuser=${DCRD_RPC_USER:-decred}
  --rpcpass=${DCRD_RPC_PASS:-decredpass}
  --rpclisten=0.0.0.0:9109
  --notls
  --txindex
  --maxpeers=125
  --blockmaxsize=393216
```

### Backend Optimization

Set environment variables:

```yaml
backend:
  environment:
    - PORT=8080
    - DCRD_RPC_HOST=dcrd
    - DCRD_RPC_PORT=9109
    - GOMAXPROCS=2
```

## Monitoring

### Prometheus/Grafana (Optional)

Add to `docker compose.yml`:

```yaml
prometheus:
  image: prom/prometheus
  volumes:
    - ./prometheus.yml:/etc/prometheus/prometheus.yml
  ports:
    - "9090:9090"

grafana:
  image: grafana/grafana
  ports:
    - "3001:3000"
  depends_on:
    - prometheus
```

## Support

- **dcrd Issues**: https://github.com/decred/dcrd/issues
- **Dashboard Issues**: (your repository)
- **Decred Community**: https://decred.org/community/

## License

ISC License - Part of the Decred community projects.


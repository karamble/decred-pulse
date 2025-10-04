# Decred Pulse - Quick Start Guide

Get up and running in 3 minutes! ⚡

## 🚀 One-Command Setup

```bash
# Clone and start
git clone https://github.com/karamble/decred-pulse.git && cd decred-pulse
cp env.example .env && nano .env  # Set your password
docker compose up -d
```

**Done!** Dashboard available at: http://localhost:3000

## 📋 Prerequisites

- Docker & Docker Compose installed
- 10GB+ free disk space
- 2GB+ RAM

## 🔧 Initial Configuration

### 1. Set RPC Credentials

Edit `.env`:
```env
DCRD_RPC_USER=decred
DCRD_RPC_PASS=change_me_to_secure_password
```

### 2. Start Services

```bash
docker compose up -d
```

This starts:
- **dcrd** (full Decred node - built from source)
- **backend** (Go API)
- **frontend** (React UI)

**⏱️ First Run**: Takes 5-10 minutes to build dcrd from source  
**⚡ Subsequent Runs**: Uses cached build (instant)

### 3. Monitor Sync Progress

First run downloads blockchain (~8-10GB for mainnet):

```bash
# Watch sync logs
docker compose logs -f dcrd

# Check sync status
make sync-status
```

Typical sync time: **4-8 hours** for mainnet

## 📊 Access Dashboard

Open browser: **http://localhost:3000**

The dashboard shows:
- Real-time sync progress
- Node status
- Network metrics
- Connected peers
- Blockchain info

## 🎮 Makefile Commands

**Easy management** with make:

### Basic Operations
```bash
make help                # Show all available commands
make setup               # Initial setup (copy env.example to .env)
make start               # Start all services (dcrd + backend + frontend)
make stop                # Stop all services
make restart             # Restart all services
make status              # Show status of all services
```

### Logging
```bash
make logs                # View logs from all services
make logs-dcrd           # View dcrd logs only
make logs-backend        # View backend logs only
make logs-frontend       # View frontend logs only
```

### dcrd Node Management
```bash
make sync-status         # Check blockchain sync status
make peers               # Show connected peers
make dcrctl CMD="..."    # Run dcrctl command (e.g., make dcrctl CMD="getblockcount")
```

### Building & Updates
```bash
make build               # Build/rebuild all images
make build-dcrd          # Build dcrd specific version (e.g., make build-dcrd VERSION=release-v2.0.6)
make update              # Update all services to latest
make update-dcrd         # Update dcrd to latest from source
```

### Development
```bash
make dev-backend         # Run backend in development mode (outside Docker)
make dev-frontend        # Run frontend in development mode (outside Docker)
make install-backend     # Install backend dependencies
make install-frontend    # Install frontend dependencies
```

### Container Access
```bash
make shell-dcrd          # Open shell in dcrd container
make shell-backend       # Open shell in backend container
```

### Data Management
```bash
make backup              # Backup blockchain data
make restore             # Restore blockchain data (e.g., make restore BACKUP=backups/dcrd-backup-xxx.tar.gz)
make clean               # Stop and remove everything (WARNING: deletes blockchain data!)
```

## 🐳 Docker Compose Commands (Direct)

If you prefer using `docker compose` commands directly instead of `make`:

### Build Specific dcrd Version

**Method 1: Using environment variable**

```bash
# Set version in .env file
echo "DCRD_VERSION=release-v2.0.6" >> .env

# Build dcrd
docker compose build dcrd

# Start the service
docker compose up -d dcrd
```

**Method 2: Inline environment variable**

```bash
# Build specific version directly
DCRD_VERSION=release-v2.0.6 docker compose build dcrd

# Start the service
docker compose up -d dcrd
```

**Method 3: Using build args**

```bash
# Build with build argument
docker compose build --build-arg DCRD_VERSION=release-v2.0.6 dcrd

# Start the service
docker compose up -d dcrd
```

**Available versions**: Browse at https://github.com/decred/dcrd/tags

Examples:
- `master` - Latest development (default)
- `release-v2.0.6` - Stable release v2.0.6
- `release-v2.0.5` - Stable release v2.0.5
- `release-v2.0.4` - Stable release v2.0.4
- `abc123def456` - Specific commit hash

### Update dcrd to Latest

**Update from latest source:**

```bash
# Rebuild dcrd from latest master branch
docker compose build --no-cache dcrd

# Restart dcrd with new build
docker compose up -d dcrd
```

**Update to specific newer version:**

```bash
# Set new version
export DCRD_VERSION=release-v2.0.7

# Rebuild
docker compose build --no-cache dcrd

# Restart
docker compose up -d dcrd
```

### Update All Services

**Rebuild all services:**

```bash
# Rebuild everything
docker compose build --no-cache

# Restart all services
docker compose up -d
```

**Rebuild individual services:**

```bash
# Backend only
docker compose build --no-cache backend
docker compose up -d backend

# Frontend only
docker compose build --no-cache frontend
docker compose up -d frontend
```

### Check Build Status

```bash
# View images
docker images | grep decred

# Check build logs
docker compose logs --tail=100 dcrd

# Verify dcrd version (after build)
docker exec decred-pulse-dcrd dcrd --version
```

### Force Clean Build

If you encounter build issues:

```bash
# Remove old images
docker rmi $(docker images -q 'decred-pulse*')

# Clean build cache
docker builder prune -a

# Rebuild from scratch
docker compose build --no-cache dcrd
docker compose up -d dcrd
```

## 🔍 Common Tasks

### Check if Services are Running

```bash
docker compose ps
```

Should show 3 services running:
- `decred-pulse-dcrd`
- `decred-pulse-backend`
- `decred-pulse-frontend`

### View Logs

```bash
# All logs
docker compose logs -f

# Specific service
docker compose logs -f dcrd
docker compose logs -f backend
docker compose logs -f frontend
```

### Check Node Status

```bash
make status
```

Or manually:
```bash
docker exec decred-pulse-dcrd dcrctl \
  --rpcuser=decred \
  --rpcpass=your_password \
  --rpcserver=127.0.0.1:9109 \
  --rpccert=/certs/rpc.cert \
  getinfo
```

### Restart Services

```bash
# All services
docker compose restart

# Specific service
docker compose restart dcrd
docker compose restart backend
```

## 📶 Ports

- **3000** - Dashboard UI (frontend)
- **8080** - API backend
- **9108** - dcrd P2P (blockchain sync)
- **9109** - dcrd RPC (backend connection)

## ⚠️ Troubleshooting

### Dashboard shows "RPC not connected"

1. Check dcrd is running: `docker compose ps dcrd`
2. Check logs: `docker compose logs dcrd`
3. Wait for blockchain sync (check progress: `make sync-status`)

### Slow blockchain sync

Normal! First sync takes time:
- Mainnet: 4-8 hours
- Testnet: 30-60 minutes

Check progress:
```bash
make sync-status
docker compose logs -f dcrd | grep -i sync
```

### Port already in use

Another service using ports 3000, 8080, 9108, or 9109?

Check:
```bash
netstat -tulpn | grep -E "3000|8080|9108|9109"
```

Stop conflicting service or change ports in `docker compose.yml`.

### Out of disk space

Blockchain data requires:
- Mainnet: ~10GB
- Testnet: ~2GB

Check available space:
```bash
df -h
docker system df
```

## 🧹 Cleanup

### Stop Services

```bash
docker compose down
```

### Remove Everything (including blockchain data)

```bash
make clean
# or
docker compose down -v
```

**Warning**: This deletes the entire blockchain! You'll need to re-sync.

## 🔄 Updates

### Update dcrd to Latest

```bash
# Rebuild from latest source
docker compose build --no-cache dcrd
docker compose up -d dcrd
```

Or simply:
```bash
make update-dcrd
```

### Build Specific dcrd Version

Edit `.env`:
```env
DCRD_VERSION=release-v2.0.6
```

Then rebuild:
```bash
docker compose build dcrd
docker compose up -d dcrd
```

Available versions: https://github.com/decred/dcrd/tags

### Update All Services

```bash
# Rebuild all
docker compose build --no-cache

# Restart
docker compose up -d
```

Or simply:
```bash
make update
docker compose restart
```

## 💾 Backup Blockchain Data

```bash
# Create backup
make backup

# Creates: backups/dcrd-backup-YYYYMMDD-HHMMSS.tar.gz
```

Restore:
```bash
make restore BACKUP=backups/dcrd-backup-xxx.tar.gz
```

## 🧪 Development Mode

Run frontend/backend locally (without Docker) for development:

```bash
# Backend (requires Go 1.21+)
cd backend
export DCRD_RPC_HOST=localhost
export DCRD_RPC_PORT=9109
export DCRD_RPC_USER=decred
export DCRD_RPC_PASS=your_password
go run main.go

# Frontend (requires Node.js 18+)
cd frontend
npm install
npm run dev
```

## 📚 More Documentation

- **Full Docker Guide**: [DOCKER_SETUP.md](DOCKER_SETUP.md)
- **Main README**: [README.md](README.md)
- **Backend Docs**: [backend/README.md](backend/README.md)
- **Frontend Docs**: [frontend/README.md](frontend/README.md)

## 🆘 Getting Help

**Logs are your friend!**

```bash
# When something goes wrong, check logs:
docker compose logs -f dcrd      # Node issues
docker compose logs -f backend   # API issues
docker compose logs -f frontend  # UI issues
```

**Common solutions**:
1. Restart: `docker compose restart`
2. Check credentials in `.env`
3. Wait for blockchain sync
4. Check disk space: `df -h`

## ✅ Success Checklist

- [ ] Docker & Docker Compose installed
- [ ] `.env` file created with secure password
- [ ] Services started: `docker compose up -d`
- [ ] All containers running: `docker compose ps`
- [ ] dcrd syncing: `make logs-dcrd`
- [ ] Dashboard accessible: http://localhost:3000
- [ ] Real data showing (after sync completes)

## 🎉 Next Steps

Once sync completes:
- Explore the dashboard metrics
- Monitor your node performance
- Check connected peers
- View blockchain information
- Track network statistics

---

**Happy monitoring!** 🚀

For issues, check [DOCKER_SETUP.md](DOCKER_SETUP.md) troubleshooting section.


# Decred Pulse

A modern, full-stack dashboard for monitoring Decred (dcrd) node performance and network status in real-time.

![License](https://img.shields.io/badge/license-ISC-blue.svg)
![Go](https://img.shields.io/badge/Go-1.21+-00ADD8?logo=go)
![React](https://img.shields.io/badge/React-18+-61DAFB?logo=react)
![TypeScript](https://img.shields.io/badge/TypeScript-5+-3178C6?logo=typescript)

## Features

### 🎯 Core Functionality
- **Real-time Monitoring**: Live data from dcrd node via RPC
- **Beautiful UI**: Modern dark theme with gradient effects
- **Comprehensive Metrics**: Node status, blockchain info, network stats, peer connections
- **Auto-refresh**: Dashboard updates every 30 seconds
- **Dynamic RPC Configuration**: Connect to dcrd node through the UI

### 📊 Dashboard Metrics
- **Node Status**: Sync status, progress, and version
- **Blockchain Data**: Block height, difficulty, chain size, block time
- **Network Info**: Peer count, network hashrate
- **Supply Information**: Circulating, staked, mixed supply, exchange rate
- **Treasury Data**: Self-funded treasury size
- **Connected Peers**: Real-time peer list with latency

## Architecture

```
┌─────────────┐      REST API      ┌─────────────┐      RPC      ┌─────────┐
│   React     │ ◄────────────────► │  Go Backend │ ◄───────────► │  dcrd   │
│  Frontend   │     JSON/HTTP      │     API     │   JSON-RPC    │  Node   │
└─────────────┘                    └─────────────┘               └─────────┘
```

### Frontend (React + TypeScript)
- **Framework**: React 18 with TypeScript
- **Build Tool**: Vite for fast development
- **Styling**: Tailwind CSS
- **HTTP Client**: Axios
- **Icons**: Lucide React

### Backend (Go)
- **Language**: Go 1.21+
- **RPC Client**: dcrd rpcclient v8
- **Router**: Gorilla Mux
- **CORS**: rs/cors middleware

## Prerequisites

### Backend
- Go 1.21 or later
- Access to a dcrd node with RPC enabled

### Frontend
- Node.js 18+ (or Bun)
- Modern web browser

### dcrd Node
```bash
# Example dcrd.conf
rpcuser=your_username
rpcpass=your_password
rpclisten=127.0.0.1:9109
```

## Quick Start

### Option A: Docker Compose (Recommended - Includes dcrd)

**Easiest way** - runs everything including a full dcrd node:

```bash
# 1. Clone repository
git clone <repository-url>
cd decred-pulse

# 2. Set up credentials
cp env.example .env
# Edit .env with your RPC password
nano .env

# 3. Start everything (dcrd + backend + frontend)
docker compose up -d

# 4. Monitor initial blockchain sync (first run only)
docker compose logs -f dcrd

# 5. Open browser
# http://localhost:3000
```

**Using Makefile** (even easier):
```bash
make setup    # Create .env file
make start    # Start all services
make status   # Check status
make logs     # View logs
```

See [DOCKER_SETUP.md](docs/DOCKER_SETUP.md) for complete Docker documentation.

### Option B: Manual Setup (External dcrd)

If you already have dcrd running elsewhere:

#### 1. Clone Repository

```bash
git clone <repository-url>
cd umbrel-decred-dash-main
```

#### 2. Start Backend

```bash
cd backend

# Install dependencies
go mod download

# Set environment variables
export DCRD_RPC_HOST=localhost
export DCRD_RPC_PORT=9109
export DCRD_RPC_USER=your_username
export DCRD_RPC_PASS=your_password

# Run backend
go run main.go
```

Backend will start on `http://localhost:8080`

#### 3. Start Frontend

```bash
cd frontend

# Install dependencies
npm install
# or
bun install

# Run development server
npm run dev
# or
bun run dev
```

Frontend will start on `http://localhost:3000`

#### 4. Open Browser

Navigate to `http://localhost:3000`

## Configuration

### Backend Configuration

Environment variables:

```bash
# Server
PORT=8080

# Decred RPC
DCRD_RPC_HOST=localhost
DCRD_RPC_PORT=9109
DCRD_RPC_USER=your_username
DCRD_RPC_PASS=your_password
```

See `backend/config.example` for more details.

### Frontend Configuration

Optional `.env` file:

```bash
VITE_API_URL=http://localhost:8080/api
```

## API Endpoints

### Health Check
```
GET /api/health
```

### Dashboard Data (All-in-one)
```
GET /api/dashboard
```
Returns complete dashboard data.

### Node Status
```
GET /api/node/status
```

### Blockchain Information
```
GET /api/blockchain/info
```

### Network Peers
```
GET /api/network/peers
```

### Connect RPC
```
POST /api/connect
Content-Type: application/json

{
  "host": "localhost",
  "port": "9109",
  "username": "user",
  "password": "pass"
}
```

## Building for Production

### Backend

```bash
cd backend
go build -o decred-dashboard-backend
./decred-dashboard-backend
```

### Frontend

```bash
cd frontend
npm run build
# Output in dist/
```

Deploy the `dist/` folder to any static hosting service.

## Docker Compose Stack (Recommended)

The project includes a complete Docker Compose orchestration with **dcrd full node built from source** using the official [decred/dcrd repository](https://github.com/decred/dcrd/tree/master/contrib/docker).

### Stack Components

```
┌─────────────────┐
│   Frontend      │  Port 3000 (Nginx + React)
│   (React UI)    │
└────────┬────────┘
         │ HTTP REST API
         ▼
┌─────────────────┐
│   Backend       │  Port 8080 (Go API)
│   (Go + Mux)    │
└────────┬────────┘
         │ JSON-RPC
         ▼
┌─────────────────┐
│   dcrd Node     │  Ports 9108 (P2P), 9109 (RPC)
│   (Full Node)   │  Built from source (latest)
└─────────────────┘
```

### Features

- ✅ **Built from Source**: Builds dcrd from official GitHub repository (not outdated Docker image)
- ✅ **Always Latest**: Uses current source code with latest security patches
- ✅ **Version Control**: Choose specific release tags or use master branch
- ✅ **Persistent Storage**: Blockchain data stored in Docker volume
- ✅ **Health Checks**: Backend waits for dcrd to be ready
- ✅ **Auto-restart**: All services restart on failure
- ✅ **Network Isolation**: Services communicate via private Docker network
- ✅ **Easy Management**: Makefile commands for common operations

### Quick Commands

```bash
# Start everything (one command!)
docker compose up -d
# Note: First run builds dcrd from source (takes 5-10 minutes)

# Or use Makefile
make start        # Start all services
make logs         # View logs
make logs-dcrd    # View dcrd sync progress
make status       # Check status and blockchain info
make stop         # Stop all services

# Advanced
make sync-status  # Check blockchain sync percentage
make peers        # Show connected peers
make backup       # Backup blockchain data
make dcrctl CMD="getblockcount"  # Run dcrctl commands
```

### Build dcrd from Source

The stack builds dcrd from the official repository on first run:

```bash
# Build specific version
DCRD_VERSION=release-v2.0.6 docker compose build dcrd

# Update to latest
docker compose build --no-cache dcrd
docker compose up -d dcrd
```

See `dcrd/README.md` for detailed build documentation.

### Files

- `docker compose.yml` - Main orchestration file
- `dcrd.conf` - dcrd configuration (mounted into container)
- `env.example` - Environment variables template
- `Makefile` - Helper commands
- `docs/DOCKER_SETUP.md` - Complete Docker documentation
- `docs/QUICKSTART.md` - Quick start guide
- `docs/BUILDING_FROM_SOURCE.md` - Building dcrd from source

### Initial Sync

On first run, dcrd needs to sync the blockchain:
- **Mainnet**: ~8-10GB, takes 4-8 hours
- **Testnet**: ~1-2GB, takes 30-60 minutes

Monitor progress:
```bash
docker compose logs -f dcrd
# or
make logs-dcrd
```

The dashboard will show sync status and progress bar during initial sync.

### Environment Variables

Create `.env` file (from `env.example`):

```env
DCRD_RPC_USER=decred
DCRD_RPC_PASS=your_secure_password_here
```

Both dcrd and backend use these credentials automatically.

### Volume Management

Blockchain data persists in Docker volume: `dcrd-data`

```bash
# Backup
make backup

# Clean everything (WARNING: deletes blockchain!)
make clean
```

See [DOCKER_SETUP.md](DOCKER_SETUP.md) for complete Docker documentation including:
- Detailed architecture
- Troubleshooting guide
- Production deployment
- Performance tuning
- Monitoring setup

## Project Structure

```
umbrel-decred-dash-main/
├── backend/              # Go backend API
│   ├── main.go          # Main server file
│   ├── go.mod           # Go dependencies
│   ├── config.example   # Configuration example
│   └── README.md        # Backend documentation
│
├── frontend/            # React frontend
│   ├── src/
│   │   ├── components/  # React components
│   │   ├── services/    # API service layer
│   │   ├── App.tsx      # Main app component
│   │   └── main.tsx     # Entry point
│   ├── package.json
│   ├── vite.config.ts
│   └── README.md        # Frontend documentation
│
└── README.md           # This file
```

## Development

### Backend Development

```bash
cd backend
go run main.go
```

### Frontend Development

```bash
cd frontend
npm run dev
```

Hot reload is enabled for both frontend and backend during development.

## Troubleshooting

### Backend won't connect to dcrd

1. Check dcrd is running: `dcrd --version`
2. Verify RPC credentials in dcrd.conf
3. Check firewall allows port 9109
4. Test RPC connection: `curl -u user:pass http://localhost:9109`

### Frontend can't reach backend

1. Verify backend is running on port 8080
2. Check CORS settings in backend
3. Verify `VITE_API_URL` is correct
4. Check browser console for errors

### Data not updating

1. Check backend logs for RPC errors
2. Verify dcrd node is synced
3. Check auto-refresh interval in frontend

## Performance

- **Backend**: Minimal CPU/memory usage, ~10-20MB RAM
- **Frontend**: Optimized React build, < 500KB gzipped
- **RPC Calls**: Cached where appropriate, rate-limited to avoid overwhelming dcrd

## Security

- **RPC Credentials**: Never exposed to frontend, stored in backend environment
- **CORS**: Configurable for production environments
- **TLS**: Support for TLS connections to dcrd
- **Input Validation**: All API inputs validated

## Contributing

Contributions are welcome! Please:

1. Fork the repository
2. Create a feature branch
3. Make your changes
4. Test thoroughly
5. Submit a pull request

## License

ISC License - Part of the Decred community projects.

## Acknowledgments

- Decred development team for dcrd and rpcclient
- React and Go communities
- Tailwind CSS for the amazing styling framework

## Support

For issues and questions:
- GitHub Issues
- Decred Matrix channels
- Decred Discord

## Roadmap

- [ ] WebSocket support for real-time updates
- [ ] Historical data charts
- [ ] Transaction monitoring
- [ ] Ticket pool statistics
- [ ] Mobile app version
- [ ] Multi-node support
- [ ] Alert notifications

---

**Made with ❤️ for the Decred community**

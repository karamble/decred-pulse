# Installation Guide

Complete installation guide for Decred Pulse. Choose the installation method that best suits your needs.

## 📋 Prerequisites

Before installing Decred Pulse, ensure you have:

### Hardware Requirements

**Minimum**:
- 2 CPU cores
- 4 GB RAM
- 15 GB free disk space
- Internet connection

**Recommended**:
- 4+ CPU cores
- 8 GB RAM
- 20+ GB free disk space (SSD preferred)
- Stable internet (10+ Mbps)

### Software Requirements

Choose your installation method:

#### Option A: Docker Installation (Recommended)
- **Docker**: Version 20.10+
- **Docker Compose**: Version 2.0+
- **Operating System**: Linux, macOS, or Windows

#### Option B: Manual Installation
- **Go**: Version 1.21+
- **Node.js**: Version 18+
- **Git**: Latest version
- **dcrd**: Running instance
- **dcrwallet**: Running instance (optional)

---

## 🚀 Option A: Docker Installation (Recommended)

The easiest and most reliable way to run Decred Pulse.

### Step 1: Install Docker

**Linux (Ubuntu/Debian)**:
```bash
# Update package index
sudo apt update

# Install Docker
sudo apt install -y docker.io docker-compose

# Add user to docker group
sudo usermod -aG docker $USER

# Log out and back in for group changes to take effect
```

**macOS**:
```bash
# Install Docker Desktop from:
# https://www.docker.com/products/docker-desktop

# Or via Homebrew:
brew install --cask docker
```

**Windows**:
```powershell
# Install Docker Desktop from:
# https://www.docker.com/products/docker-desktop

# Requires WSL2 (Windows Subsystem for Linux)
```

**Verify Installation**:
```bash
docker --version
docker compose version
```

Expected output:
```
Docker version 24.0.5, build...
Docker Compose version v2.20.0
```

---

### Step 2: Clone Repository

```bash
# Clone the repository
git clone https://github.com/karamble/decred-pulse.git

# Enter directory
cd decred-pulse
```

---

### Step 3: Configure Environment

```bash
# Create .env file from example
cp env.example .env

# Edit with your credentials
nano .env
```

**Required changes**:
```bash
# Set secure RPC passwords
DCRD_RPC_PASS=your_secure_password_here
DCRWALLET_RPC_PASS=your_secure_wallet_password_here
```

**Generate secure passwords**:
```bash
# Random 32-character password
openssl rand -base64 32
```

---

### Step 4: Start Services

```bash
# Start all services
docker compose up -d

# Or using Makefile
make start
```

**What happens**:
1. Builds dcrd from source (5-10 minutes, first time only)
2. Builds dcrwallet from source
3. Builds backend and frontend
4. Starts all containers
5. Begins blockchain sync

**Monitor startup**:
```bash
# View logs
docker compose logs -f

# Check status
docker compose ps
```

---

### Step 5: Access Dashboard

**Open browser**: http://localhost:3000

**Initial view**:
- Node Dashboard shows sync status
- "Syncing" status with progress bar
- Peer connections establishing

**Note**: First sync takes 4-8 hours for mainnet

---

### Step 6: Verify Installation

```bash
# Check all services are running
make status

# Or
docker compose ps
```

**Expected output**:
```
NAME                      STATUS      PORTS
decred-pulse-dcrd         Up (healthy)  9108-9109
decred-pulse-dcrwallet    Up (healthy)  9110
decred-pulse-backend      Up (healthy)  8080
decred-pulse-frontend     Up            3000
```

---

## 🔧 Option B: Manual Installation

For development or custom setups.

### Step 1: Install Dependencies

**Go (Backend)**:
```bash
# Download from https://go.dev/dl/
# Or via package manager

# Ubuntu/Debian
sudo apt install golang-go

# macOS
brew install go

# Verify
go version  # Should be 1.21+
```

**Node.js (Frontend)**:
```bash
# Download from https://nodejs.org/
# Or via package manager

# Ubuntu/Debian
curl -fsSL https://deb.nodesource.com/setup_18.x | sudo -E bash -
sudo apt install -y nodejs

# macOS
brew install node

# Verify
node --version  # Should be 18+
npm --version
```

---

### Step 2: Install dcrd

```bash
# Clone dcrd
git clone https://github.com/decred/dcrd.git
cd dcrd

# Build
go install .

# Verify
dcrd --version
```

**Configure dcrd**:
```bash
# Create config directory
mkdir -p ~/.dcrd

# Create config file
nano ~/.dcrd/dcrd.conf
```

**Minimal dcrd.conf**:
```ini
rpcuser=your_rpc_username
rpcpass=your_rpc_password
rpclisten=127.0.0.1:9109
txindex=1
```

**Start dcrd**:
```bash
dcrd
```

---

### Step 3: Install dcrwallet (Optional)

```bash
# Clone dcrwallet
git clone https://github.com/decred/dcrwallet.git
cd dcrwallet

# Build
go install .

# Verify
dcrwallet --version
```

**Configure dcrwallet**:
```bash
# Create config directory
mkdir -p ~/.dcrwallet

# Create config file
nano ~/.dcrwallet/dcrwallet.conf
```

**Minimal dcrwallet.conf**:
```ini
username=your_wallet_username
password=your_wallet_password
rpclisten=127.0.0.1:9110
```

**Start dcrwallet**:
```bash
dcrwallet
```

---

### Step 4: Install Decred Pulse Backend

```bash
# Clone repository
git clone https://github.com/karamble/decred-pulse.git
cd decred-pulse/backend

# Install dependencies
go mod download

# Set environment variables
export DCRD_RPC_HOST=localhost
export DCRD_RPC_PORT=9109
export DCRD_RPC_USER=your_rpc_username
export DCRD_RPC_PASS=your_rpc_password

# If using wallet
export DCRWALLET_RPC_HOST=localhost
export DCRWALLET_RPC_PORT=9110
export DCRWALLET_RPC_USER=your_wallet_username
export DCRWALLET_RPC_PASS=your_wallet_password

# Start backend
go run main.go
```

**Backend should start on**: http://localhost:8080

---

### Step 5: Install Decred Pulse Frontend

```bash
# In new terminal, navigate to frontend
cd decred-pulse/frontend

# Install dependencies
npm install

# Start development server
npm run dev
```

**Frontend should start on**: http://localhost:5173

---

### Step 6: Verify Installation

**Test backend**:
```bash
curl http://localhost:8080/api/health
```

Expected: `{"status":"healthy"}`

**Test frontend**: Open http://localhost:5173 in browser

**Test node connection**: Dashboard should show node data

---

## 🐧 Linux-Specific Installation

### Ubuntu/Debian

```bash
# Update system
sudo apt update && sudo apt upgrade -y

# Install Docker
sudo apt install -y docker.io docker-compose git

# Add user to docker group
sudo usermod -aG docker $USER
newgrp docker

# Clone and start
git clone https://github.com/karamble/decred-pulse.git
cd decred-pulse
cp env.example .env
nano .env  # Edit passwords
make start
```

### Fedora/RHEL

```bash
# Install Docker
sudo dnf install -y docker docker-compose git

# Start Docker service
sudo systemctl start docker
sudo systemctl enable docker

# Add user to docker group
sudo usermod -aG docker $USER
newgrp docker

# Clone and start
git clone https://github.com/karamble/decred-pulse.git
cd decred-pulse
cp env.example .env
nano .env  # Edit passwords
make start
```

### Arch Linux

```bash
# Install Docker
sudo pacman -S docker docker-compose git

# Start Docker service
sudo systemctl start docker.service
sudo systemctl enable docker.service

# Add user to docker group
sudo usermod -aG docker $USER
newgrp docker

# Clone and start
git clone https://github.com/karamble/decred-pulse.git
cd decred-pulse
cp env.example .env
nano .env  # Edit passwords
make start
```

---

## 🍎 macOS-Specific Installation

```bash
# Install Homebrew (if not installed)
/bin/bash -c "$(curl -fsSL https://raw.githubusercontent.com/Homebrew/install/HEAD/install.sh)"

# Install Docker Desktop
brew install --cask docker

# Start Docker Desktop from Applications

# Install Git (if not installed)
brew install git

# Clone repository
git clone https://github.com/karamble/decred-pulse.git
cd decred-pulse

# Setup and start
cp env.example .env
nano .env  # Edit passwords
make start
```

---

## 🪟 Windows-Specific Installation

### Windows with WSL2 (Recommended)

```powershell
# Install WSL2
wsl --install

# Restart computer

# In WSL2 terminal:
# Install Docker Desktop for Windows with WSL2 backend

# Clone repository
git clone https://github.com/karamble/decred-pulse.git
cd decred-pulse

# Setup and start
cp env.example .env
nano .env  # Edit passwords
make start
```

### Windows with Docker Desktop

```powershell
# Install Docker Desktop from:
# https://www.docker.com/products/docker-desktop

# Install Git from:
# https://git-scm.com/download/win

# Clone repository (Git Bash)
git clone https://github.com/karamble/decred-pulse.git
cd decred-pulse

# Create .env file
copy env.example .env
notepad .env  # Edit passwords

# Start (PowerShell or Git Bash)
docker compose up -d
```

---

## 🧪 Testnet Installation

For testing without real DCR:

```bash
# Clone repository
git clone https://github.com/karamble/decred-pulse.git
cd decred-pulse

# Create .env file
cp env.example .env

# Edit .env and enable testnet
nano .env
```

**Add/uncomment**:
```bash
DCRD_TESTNET=1
```

**Start services**:
```bash
make start
```

**Testnet benefits**:
- ✅ Faster sync (~30-60 minutes)
- ✅ Smaller size (~1-2 GB)
- ✅ Free testnet coins
- ✅ Safe for experimentation

---

## 🔧 Post-Installation Configuration

### Customize Ports

Edit `docker-compose.yml`:

```yaml
services:
  frontend:
    ports:
      - "8000:80"  # Change 3000 to 8000
  
  backend:
    ports:
      - "8081:8080"  # Change 8080 to 8081
```

### Enable Transaction Indexing

Edit `.env`:
```bash
DCRD_EXTRA_ARGS=--txindex
```

Restart:
```bash
docker compose restart dcrd
```

### Adjust Gap Limit

Edit `.env`:
```bash
DCRWALLET_GAP_LIMIT=500  # Increase for older wallets
```

Restart:
```bash
docker compose restart dcrwallet
```

---

## 🐛 Troubleshooting Installation

### Docker Installation Issues

**Problem**: "Cannot connect to Docker daemon"

**Solution**:
```bash
# Start Docker service
sudo systemctl start docker

# Or on macOS, start Docker Desktop app
```

---

**Problem**: "Permission denied"

**Solution**:
```bash
# Add user to docker group
sudo usermod -aG docker $USER

# Log out and back in, or:
newgrp docker
```

---

**Problem**: "Port already in use"

**Solution**:
```bash
# Find process using port
sudo lsof -i :3000

# Kill process or change port in docker-compose.yml
```

---

### Build Issues

**Problem**: "dcrd build failed"

**Solution**:
```bash
# Clean rebuild
docker compose down
docker compose build --no-cache dcrd
docker compose up -d
```

---

**Problem**: "Out of disk space"

**Solution**:
```bash
# Clean Docker
docker system prune -a

# Check disk space
df -h
```

---

### Network Issues

**Problem**: "Cannot pull images"

**Solution**:
```bash
# Check internet connection
ping google.com

# Check DNS
cat /etc/resolv.conf

# Retry
docker compose pull
```

---

## ✅ Verify Installation

### Quick Health Check

```bash
# All services running?
docker compose ps

# Backend healthy?
curl http://localhost:8080/api/health

# Frontend accessible?
curl http://localhost:3000

# dcrd syncing?
make sync-status
```

### Full Verification

```bash
# Node status
make status

# View logs
make logs

# Check peers
make peers

# Dashboard accessible
# Open http://localhost:3000 in browser
```

---

## 🎓 Next Steps

After successful installation:

1. **[First Steps](first-steps.md)** - What to do next
2. **[Configuration](../setup/configuration.md)** - Customize settings
3. **[Node Dashboard](../features/node-dashboard.md)** - Monitor your node
4. **[Wallet Setup](../wallet-setup.md)** - Configure wallet

---

## 📚 Additional Resources

- **[Quick Start](../quickstart.md)** - 3-minute setup
- **[Docker Setup](../docker-setup.md)** - Detailed Docker guide
- **[Troubleshooting](../guides/troubleshooting.md)** - Common issues
- **[CLI Commands](../reference/cli-commands.md)** - Command reference

---

**Questions?** Check the [FAQ](../reference/faq.md) or [Troubleshooting Guide](../guides/troubleshooting.md)


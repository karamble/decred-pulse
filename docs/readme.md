# Decred Pulse Documentation

Welcome to the **Decred Pulse** documentation! This guide will help you set up, configure, and use your Decred node and wallet dashboard.

## 📖 Documentation Index

### 🚀 Getting Started

Start here if you're new to Decred Pulse:

- **[Installation Guide](getting-started/installation.md)** - Detailed installation instructions
- **[First Steps](getting-started/first-steps.md)** - What to do after installation
- **Quick Start** - Get up and running in 3 minutes *(see Installation Guide)*

### ⚙️ Setup & Configuration

Detailed setup guides for different components:

- **[Configuration Guide](setup/configuration.md)** - Environment variables and config files
- **Docker Setup** - Complete Docker Compose guide *(see Installation Guide)*
- **Building from Source** - Build dcrd and dcrwallet from source *(see Installation Guide)*
- **Wallet Setup** - Configure and connect your wallet *(see Wallet Operations)*

### 🎯 Features

Learn about the dashboard features:

- **[Node Dashboard](features/node-dashboard.md)** - Monitor your Decred node
- **[Wallet Dashboard](features/wallet-dashboard.md)** - Manage accounts, balances, and transactions
- **[Staking Guide](features/staking-guide.md)** - Ticket purchasing and staking information
- **Transaction History** - View and manage transactions *(see Wallet Dashboard)*
- **Block Explorer** - Mini block explorer features *(future)*

### 📚 User Guides

Step-by-step guides for common tasks:

- **[Wallet Operations](guides/wallet-operations.md)** - Import xpub, rescan, sync progress
- **[Backup & Restore](guides/backup-restore.md)** - Backup blockchain and wallet data
- **[Troubleshooting](guides/troubleshooting.md)** - Common issues and solutions
- **[Monitoring Setup](deployment/monitoring-setup.md)** - Production monitoring (see Deployment section)

### 🔌 API Documentation

API reference for developers:

- **[API Reference](api/api-reference.md)** - Complete API endpoint documentation
- **Node Endpoints** - Node-specific API endpoints *(see API Reference)*
- **Wallet Endpoints** - Wallet-specific API endpoints *(see API Reference)*

### 💻 Development

Contributing and development guides:

- **[Architecture](development/architecture.md)** - System design and architecture
- **Development Setup** - Local development environment *(planned)*
- **Backend Guide** - Backend development *(planned)*
- **Frontend Guide** - Frontend development *(planned)*
- **Contributing** - How to contribute *(planned)*

### 🚀 Deployment

Production deployment guides:

- **[Production Deployment](deployment/production.md)** - Deploy to production
- **[Security Best Practices](deployment/security.md)** - Security guidelines
- **[Performance Tuning](deployment/performance.md)** - Optimize performance
- **[Monitoring Setup](deployment/monitoring-setup.md)** - Production monitoring

### 📋 Reference

Quick reference materials:

- **[CLI Commands](reference/cli-commands.md)** - Makefile and Docker Compose commands
- **Environment Variables** - All environment variables *(see Configuration Guide)*

---

## 🎯 Quick Navigation by Role

### I'm a Node Operator
1. [Installation Guide](getting-started/installation.md) → Install and setup
2. [Node Dashboard](features/node-dashboard.md) → Monitor your node
3. [First Steps](getting-started/first-steps.md) → Initial setup
4. [Backup & Restore](guides/backup-restore.md) → Data protection

### I'm Managing a Wallet
1. [Installation Guide](getting-started/installation.md) → Setup wallet
2. [Wallet Dashboard](features/wallet-dashboard.md) → View balances
3. [Wallet Operations](guides/wallet-operations.md) → Import xpub, rescan

### I'm Staking Tickets
1. [Staking Guide](features/staking-guide.md) → Learn about staking
2. [Wallet Dashboard](features/wallet-dashboard.md) → View your tickets
3. [Wallet Operations](guides/wallet-operations.md) → Manage tickets

### I'm a Developer
1. [Architecture](development/architecture.md) → Understand the system
2. [API Reference](api/api-reference.md) → Use the API
3. [Configuration Guide](setup/configuration.md) → Configure the app

### I'm Deploying to Production
1. [Production Deployment](deployment/production.md) → Deploy guide
2. [Security](deployment/security.md) → Secure your deployment
3. [Monitoring Setup](deployment/monitoring-setup.md) → Production monitoring

---

## 📦 What is Decred Pulse?

**Decred Pulse** is a modern, full-stack dashboard for monitoring Decred nodes and wallets in real-time. It provides:

### Key Features

✅ **Node Monitoring**
- Real-time node status and sync progress
- Blockchain information (height, difficulty, size)
- Network statistics (peers, hashrate)
- Mempool activity

✅ **Wallet Management**
- Account balances (spendable, locked, immature)
- Transaction history with progressive loading
- Watch-only wallet support (xpub import)
- Wallet rescan with progress tracking

✅ **Staking Information**
- Ticket pool statistics
- Your tickets (immature, live, voted, revoked)
- Stake difficulty and price
- Mempool ticket activity

✅ **Modern UI**
- Beautiful dark theme with gradients
- Real-time auto-refresh
- Responsive design
- Smooth animations

### Architecture

```
┌─────────────┐      REST API      ┌─────────────┐      RPC      ┌─────────┐
│   React     │ ◄────────────────► │  Go Backend │ ◄───────────► │  dcrd   │
│  Frontend   │     JSON/HTTP      │     API     │   JSON-RPC    │  Node   │
└─────────────┘                    └─────────────┘               └─────────┘
                                           │
                                           │ RPC
                                           ▼
                                    ┌─────────────┐
                                    │ dcrwallet   │
                                    │   (RPC)     │
                                    └─────────────┘
```

### Technology Stack

**Frontend:**
- React 18 + TypeScript
- Vite (build tool)
- Tailwind CSS
- Axios (HTTP client)
- Lucide React (icons)

**Backend:**
- Go 1.21+
- Gorilla Mux (router)
- dcrd rpcclient
- CORS middleware

**Infrastructure:**
- Docker + Docker Compose
- Nginx (frontend serving)
- dcrd (full node)
- dcrwallet (wallet RPC)

---

## 🆘 Getting Help

### Documentation
- Browse the [documentation index](#-documentation-index) above
- Review [Troubleshooting](guides/troubleshooting.md)
- Check specific feature guides for FAQs

### Community
- **GitHub Issues**: Report bugs and request features
- **Decred Discord**: Get community support
- **Decred Matrix**: Chat with developers

### Logs
When something goes wrong, check the logs:

```bash
# All services
docker compose logs -f

# Specific services
docker compose logs -f dcrd
docker compose logs -f backend
docker compose logs -f frontend
```

---

## 🔗 Quick Links

- **Main README**: [README.md](../README.md)
- **Project Repository**: [GitHub](https://github.com/karamble/decred-pulse)
- **License**: [ISC License](../LICENSE)

---

**Made with ❤️ for the Decred community**


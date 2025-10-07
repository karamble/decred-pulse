# Decred Pulse

A modern, self-hosted dashboard for monitoring your Decred node, wallet, and blockchain in real-time.

![License](https://img.shields.io/badge/license-ISC-blue.svg)
![Go](https://img.shields.io/badge/Go-1.21+-00ADD8?logo=go)
![React](https://img.shields.io/badge/React-18+-61DAFB?logo=react)
![TypeScript](https://img.shields.io/badge/TypeScript-5+-3178C6?logo=typescript)

## What is Decred Pulse?

Decred Pulse is a comprehensive dashboard that provides:
- **Node Dashboard**: Monitor your dcrd node performance, blockchain sync status, network peers, and mempool activity
- **Wallet Dashboard**: Track your wallet balances, transactions, staking tickets, and ticket pool statistics
- **Block Explorer**: Browse blocks, transactions, and addresses directly from your node

All data comes from your local dcrd and dcrwallet instances via RPC - no third-party services required.

## Quick Start

### Prerequisites
- Docker and Docker Compose
- 10GB+ free disk space for blockchain data

### Launch

```bash
# 1. Clone the repository
git clone https://github.com/karamble/decred-pulse.git
cd decred-pulse

# 2. Set up environment
cp env.example .env
# Edit .env with your preferred RPC password
nano .env

# 3. Start all services
docker compose up -d

# 4. Access the dashboard
# Open http://localhost:3000 in your browser
```

The first run will sync the blockchain (takes 4-8 hours for mainnet). Monitor progress with:
```bash
docker compose logs -f dcrd
```

### Using Makefile

```bash
make start       # Start all services
make stop        # Stop all services
make logs        # View logs
make status      # Check status
```

For more commands: `make help`

## Documentation

Complete documentation is available in the [`docs/`](docs/) folder:

📚 **[Documentation Index](docs/readme.md)** - Start here

### Quick Links

**Getting Started**
- [First Steps](docs/getting-started/first-steps.md) - What to do after installation
- [Quick Start Guide](docs/getting-started/quick-start.md) - Detailed setup instructions
- [Environment Setup](docs/getting-started/environment-setup.md) - Configuration options

**Guides**
- [Wallet Operations](docs/guides/wallet-operations.md) - Import xpub, rescan, sync monitoring
- [Backup & Restore](docs/guides/backup-restore.md) - Protect your blockchain data
- [Troubleshooting](docs/guides/troubleshooting.md) - Common issues and solutions

**Features**
- [Node Dashboard](docs/features/node-dashboard.md) - Monitor your dcrd node
- [Wallet Dashboard](docs/features/wallet-dashboard.md) - Track balances and staking
- [Block Explorer](docs/features/block-explorer.md) - Browse blocks and transactions

**Deployment**
- [Docker Deployment](docs/deployment/docker-deployment.md) - Production setup
- [Monitoring Setup](docs/deployment/monitoring-setup.md) - Health checks and alerts

**Reference**
- [CLI Commands](docs/reference/cli-commands.md) - Makefile and Docker commands
- [Configuration](docs/setup/configuration.md) - All configuration options

## Support

For issues and questions:
- [GitHub Issues](https://github.com/karamble/decred-pulse/issues)
- [Decred Matrix](https://chat.decred.org)
- [Decred Discord](https://discord.gg/decred)

## License

ISC License - Part of the Decred community projects.

---

**Made with ❤️ for the Decred community**

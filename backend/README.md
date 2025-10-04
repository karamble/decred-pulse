# Decred Pulse - Backend

Go backend API server that connects to dcrd via RPC and provides REST endpoints for the dashboard frontend.

## Features

- **RPC Client**: Connects to dcrd node via RPC protocol
- **REST API**: Provides JSON endpoints for dashboard data
- **Real-time Data**: Fetches live blockchain, network, and peer information
- **CORS Enabled**: Ready for frontend integration
- **Environment Configuration**: Flexible configuration via environment variables

## Prerequisites

- Go 1.21 or later
- Access to a dcrd node with RPC enabled
- RPC credentials (username and password)

## Installation

```bash
cd backend

# Download dependencies
go mod download

# Build the application
go build -o decred-dashboard-backend
```

## Configuration

Configure the backend using environment variables:

```bash
export PORT=8080
export DCRD_RPC_HOST=localhost
export DCRD_RPC_PORT=9109
export DCRD_RPC_USER=your_username
export DCRD_RPC_PASS=your_password
```

Or copy the example config:
```bash
cp config.example .env
# Edit .env with your values
source .env
```

## Running

```bash
# Run directly
go run main.go

# Or run the built binary
./decred-dashboard-backend
```

The server will start on `http://localhost:8080` (or the port you configured).

## API Endpoints

### Health Check
```
GET /api/health
```
Returns server health status and RPC connection state.

### Dashboard Data (All-in-one)
```
GET /api/dashboard
```
Returns complete dashboard data including node status, blockchain info, network info, peers, and supply info.

### Node Status
```
GET /api/node/status
```
Returns node status and sync progress.

### Blockchain Info
```
GET /api/blockchain/info
```
Returns blockchain information (block height, difficulty, etc.).

### Network Peers
```
GET /api/network/peers
```
Returns list of connected peers with latency information.

### Connect RPC
```
POST /api/connect
Content-Type: application/json

{
  "host": "localhost",
  "port": "9109",
  "username": "your_username",
  "password": "your_password"
}
```
Dynamically connect to dcrd RPC (useful if not configured at startup).

## Response Examples

### Dashboard Data Response
```json
{
  "nodeStatus": {
    "status": "running",
    "syncProgress": 100,
    "version": "v2.0.6"
  },
  "blockchainInfo": {
    "blockHeight": 856234,
    "blockHash": "000000000000000001a2b3c4...",
    "difficulty": 452000000000,
    "chainSize": 8800000000,
    "blockTime": "5m 12s"
  },
  "networkInfo": {
    "peerCount": 8,
    "hashrate": "452.00 PH/s",
    "networkHashPS": 452000000000000000
  },
  "peers": [
    {
      "id": 1,
      "address": "45.32.176.123:9108",
      "protocol": "TCP",
      "latency": "45ms"
    }
  ],
  "supplyInfo": {
    "circulatingSupply": "17.05M",
    "stakedSupply": "10.15M",
    "stakedPercent": 59.4,
    "exchangeRate": "$17.70",
    "treasurySize": "861.6K DCR",
    "mixedPercent": "62%"
  },
  "lastUpdate": "2025-10-04T19:45:00Z"
}
```

## Development

```bash
# Run with auto-reload (using air or similar)
go run main.go

# Run tests
go test ./...

# Format code
go fmt ./...
```

## Docker (Optional)

```dockerfile
# Build
docker build -t decred-dashboard-backend .

# Run
docker run -p 8080:8080 \
  -e DCRD_RPC_HOST=dcrd \
  -e DCRD_RPC_PORT=9109 \
  -e DCRD_RPC_USER=user \
  -e DCRD_RPC_PASS=pass \
  decred-dashboard-backend
```

## Dependencies

- `github.com/decred/dcrd/rpcclient/v8` - Decred RPC client
- `github.com/gorilla/mux` - HTTP router
- `github.com/rs/cors` - CORS middleware

## License

Part of the Decred Dashboard project.


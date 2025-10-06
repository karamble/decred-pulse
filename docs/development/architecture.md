# System Architecture

Technical overview of Decred Pulse architecture, component design, data flow, and integration patterns.

## 📐 High-Level Architecture

```
┌──────────────────────────────────────────────────────────────────┐
│                          User's Browser                          │
│                     (React SPA - Port 3000)                      │
└──────────────────────┬───────────────────────────────────────────┘
                       │ HTTP/JSON (REST API)
                       │
┌──────────────────────▼───────────────────────────────────────────┐
│                      Go Backend API                              │
│                       (Port 8080)                                │
│  ┌──────────────┬──────────────┬──────────────┬───────────────┐ │
│  │   Handlers   │   Services   │    Types     │   RPC Client  │ │
│  └──────────────┴──────────────┴──────────────┴───────────────┘ │
└──────────────────────┬──────────────────────┬────────────────────┘
                       │ JSON-RPC              │ JSON-RPC
                       │                       │
          ┌────────────▼────────────┐  ┌──────▼───────────┐
          │     dcrd Node           │  │   dcrwallet      │
          │   (Port 9108/9109)      │  │   (Port 9110)    │
          │ ┌─────────────────────┐ │  │ ┌──────────────┐ │
          │ │  Blockchain Data    │ │  │ │ Wallet DB    │ │
          │ │   (~10 GB Docker    │ │  │ │              │ │
          │ │      Volume)        │ │  │ └──────────────┘ │
          │ └─────────────────────┘ │  └──────────────────┘
          │                         │
          │  P2P Network (Port 9108)│
          └────────────┬────────────┘
                       │
                       │ Decred P2P Protocol
                       │
          ┌────────────▼────────────┐
          │   Decred Network        │
          │   (Global P2P)          │
          └─────────────────────────┘
```

---

## 🎯 Component Overview

### Frontend (React + TypeScript)

**Purpose**: User interface for monitoring and managing Decred node and wallet

**Technology Stack**:
- **Framework**: React 18
- **Language**: TypeScript
- **Build Tool**: Vite
- **Styling**: Tailwind CSS
- **HTTP Client**: Axios
- **Icons**: Lucide React
- **Routing**: React Router DOM

**Architecture Pattern**: Component-based with service layer

**Location**: `/frontend/`

---

### Backend (Go API)

**Purpose**: Bridge between frontend and Decred RPC services

**Technology Stack**:
- **Language**: Go 1.21+
- **Router**: Gorilla Mux
- **RPC Client**: dcrd rpcclient v8
- **CORS**: rs/cors
- **Concurrency**: Goroutines and channels

**Architecture Pattern**: Layered architecture (Handlers → Services → RPC)

**Location**: `/backend/`

---

### dcrd (Decred Node)

**Purpose**: Full Decred blockchain node

**Functionality**:
- Blockchain synchronization
- P2P networking
- Block validation
- Transaction relay
- RPC interface

**Built From**: Official dcrd source (GitHub)

**Version**: Configurable (default: master, recommended: release tags)

---

### dcrwallet (Decred Wallet)

**Purpose**: HD wallet for managing DCR and tickets

**Functionality**:
- HD wallet management
- Transaction creation/signing
- Ticket purchasing/voting
- Address generation
- Balance tracking

**Built From**: Official dcrwallet source (GitHub)

---

## 🔄 Data Flow

### Node Dashboard Data Flow

```
┌─────────────┐
│   Browser   │
└──────┬──────┘
       │ 1. GET /api/dashboard
       │
┌──────▼──────────────────────────────────────┐
│            Backend API                      │
│                                             │
│  2. GetDashboardDataHandler                 │
│         │                                   │
│         ├─► 3. FetchNodeStatus()           │
│         │       └─► dcrd.GetInfo()         │
│         │                                   │
│         ├─► 4. FetchBlockchainInfo()       │
│         │       └─► dcrd.GetBlockchainInfo()│
│         │                                   │
│         ├─► 5. FetchNetworkPeers()         │
│         │       └─► dcrd.GetPeerInfo()     │
│         │                                   │
│         ├─► 6. FetchMempoolInfo()          │
│         │       └─► dcrd.GetRawMempool()   │
│         │                                   │
│         └─► 7. FetchSupplyInfo()           │
│                 └─► dcrd.GetCoinSupply()   │
│                                             │
│  8. Aggregate all data                     │
│  9. Return JSON response                   │
└──────┬──────────────────────────────────────┘
       │
       │ 10. Response (JSON)
       │
┌──────▼──────┐
│   Browser   │
│             │
│ 11. Update  │
│     UI      │
└─────────────┘
```

---

### Wallet Dashboard Data Flow

```
┌─────────────┐
│   Browser   │
└──────┬──────┘
       │ 1. GET /api/wallet/dashboard
       │
┌──────▼──────────────────────────────────────┐
│            Backend API                      │
│                                             │
│  2. GetWalletDashboardHandler              │
│         │                                   │
│         ├─► 3. FetchWalletStatus()         │
│         │       └─► wallet.WalletInfo()    │
│         │                                   │
│         ├─► 4. FetchAccountInfo()          │
│         │       └─► wallet.GetBalance()    │
│         │                                   │
│         ├─► 5. FetchAllAccounts()          │
│         │       └─► wallet.GetBalance()    │
│         │                                   │
│         └─► 6. FetchWalletStakingInfo()    │
│                 ├─► wallet.GetStakeInfo()  │
│                 ├─► dcrd.GetStakeDifficulty()│
│                 └─► dcrd.EstimateStakeDiff()│
│                                             │
│  7. Aggregate all data                     │
│  8. Return JSON response                   │
└──────┬──────────────────────────────────────┘
       │
       │ 9. Response (JSON)
       │
┌──────▼──────┐
│   Browser   │
│             │
│ 10. Update  │
│     UI      │
└─────────────┘
```

---

### Xpub Import & Rescan Flow

```
┌─────────────┐
│   Browser   │
└──────┬──────┘
       │ 1. POST /api/wallet/importxpub
       │    Body: {xpub, gapLimit}
       │
┌──────▼──────────────────────────────────────┐
│            Backend API                      │
│                                             │
│  2. ImportXpubHandler                       │
│         │                                   │
│         ├─► 3. Validate input              │
│         │                                   │
│         ├─► 4. wallet.ImportXpub()         │
│         │       (RPC to dcrwallet)          │
│         │                                   │
│         └─► 5. Trigger rescan              │
└──────┬──────────────────────────────────────┘
       │
       │ 6. Response: {status: "success"}
       │
┌──────▼──────┐
│   Browser   │
│             │
│ 7. Show     │
│    progress │
│             │
│ 8. Poll:    │
│    GET /api/│
│    wallet/  │
│    sync-    │
│    progress │
└──────┬──────┘
       │ (Every 2s)
       │
┌──────▼──────────────────────────────────────┐
│            Backend API                      │
│                                             │
│  9. GetSyncProgressHandler                  │
│         │                                   │
│         ├─► 10. Read dcrwallet.log         │
│         │        (Last 500 lines)           │
│         │                                   │
│         ├─► 11. Parse rescan messages      │
│         │        Extract: progress %,       │
│         │        current block, total       │
│         │                                   │
│         ├─► 12. Check timestamp            │
│         │        (< 2 min = active)         │
│         │                                   │
│         └─► 13. Return progress            │
└──────┬──────────────────────────────────────┘
       │
       │ 14. Response: {isRescanning, progress, ...}
       │
┌──────▼──────┐
│   Browser   │
│             │
│ 15. Update  │
│     progress│
│     bar     │
│             │
│ 16. Repeat  │
│     until   │
│     complete│
└─────────────┘
```

---

## 📦 Backend Layer Architecture

### Layer 1: Handlers (`backend/handlers/`)

**Responsibility**: HTTP request handling and response formatting

**Files**:
- `node.go` - Node/dcrd endpoints
- `wallet.go` - Wallet/dcrwallet endpoints

**Functions**:
- Parse HTTP requests
- Validate input
- Call service layer
- Format JSON responses
- Handle errors

**Example**:
```go
func GetDashboardDataHandler(w http.ResponseWriter, r *http.Request) {
    // Call service layer
    data, err := services.FetchDashboardData()
    if err != nil {
        http.Error(w, err.Error(), http.StatusInternalServerError)
        return
    }
    
    // Return JSON
    w.Header().Set("Content-Type", "application/json")
    json.NewEncoder(w).Encode(data)
}
```

---

### Layer 2: Services (`backend/services/`)

**Responsibility**: Business logic and RPC orchestration

**Files**:
- `node.go` - Node data fetching and processing
- `wallet.go` - Wallet data fetching and processing

**Functions**:
- Make RPC calls
- Process/transform data
- Aggregate multiple RPC responses
- Handle concurrency
- Error handling

**Example**:
```go
func FetchDashboardData() (*types.DashboardData, error) {
    ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
    defer cancel()
    
    // Concurrent RPC calls
    results := make(chan result, 5)
    
    go fetchNodeStatus(ctx, results)
    go fetchBlockchainInfo(ctx, results)
    go fetchPeers(ctx, results)
    // ... more goroutines
    
    // Aggregate results
    return aggregateResults(results), nil
}
```

---

### Layer 3: Types (`backend/types/`)

**Responsibility**: Data structure definitions

**Files**:
- `node.go` - Node-related types
- `wallet.go` - Wallet-related types

**Structures**:
```go
type DashboardData struct {
    NodeStatus     NodeStatus
    BlockchainInfo BlockchainInfo
    NetworkInfo    NetworkInfo
    Peers          []Peer
    MempoolInfo    MempoolInfo
    SupplyInfo     SupplyInfo
    StakingInfo    StakingInfo
    LastUpdate     time.Time
}
```

---

### Layer 4: RPC Client (`backend/rpc/`)

**Responsibility**: RPC connection management

**File**: `client.go`

**Functions**:
- Initialize RPC connections
- Maintain connection state
- Handle reconnection
- Provide RPC client instances

**Global Variables**:
```go
var (
    NodeClient   *rpcclient.Client  // dcrd RPC
    WalletClient *rpcclient.Client  // dcrwallet RPC
)
```

---

### Layer 5: Utilities (`backend/utils/`)

**Responsibility**: Helper functions

**File**: `formatters.go`

**Functions**:
- Format DCR amounts
- Format byte sizes
- Format time durations
- Parse configuration

---

## 🎨 Frontend Architecture

### Component Structure

```
src/
├── App.tsx                   # Main app component, routing
├── main.tsx                  # Entry point
│
├── components/               # Reusable UI components
│   ├── NodeStatus.tsx       # Node sync status card
│   ├── MetricCard.tsx       # Generic metric display
│   ├── BlockchainInfo.tsx   # Blockchain data card
│   ├── PeersList.tsx        # Peer connections list
│   ├── StakingStats.tsx     # Staking statistics
│   ├── MempoolActivity.tsx  # Mempool transaction breakdown
│   ├── AccountInfo.tsx      # Wallet account summary
│   ├── AccountsList.tsx     # Detailed accounts list
│   ├── TransactionHistory.tsx # Wallet transactions
│   ├── TicketPoolInfo.tsx   # Network ticket pool
│   ├── MyTicketsInfo.tsx    # Personal tickets
│   ├── ImportXpubModal.tsx  # Xpub import dialog
│   ├── Header.tsx           # App header/navigation
│   ├── WalletStatus.tsx     # Wallet connection status
│   └── RPCConnection.tsx    # RPC connection form
│
├── pages/                    # Page-level components
│   ├── NodeDashboard.tsx    # Node monitoring page
│   └── WalletDashboard.tsx  # Wallet management page
│
├── services/                 # API integration layer
│   └── api.ts               # Axios client, API functions
│
└── index.css                 # Global styles (Tailwind)
```

---

### State Management

**Pattern**: Component-level state with React hooks

**No global state library**: Keep it simple, use props and local state

**Data fetching**: `useEffect` + `useState` pattern

**Example**:
```typescript
const [data, setData] = useState<DashboardData | null>(null);
const [loading, setLoading] = useState(true);
const [error, setError] = useState<string | null>(null);

useEffect(() => {
  const fetchData = async () => {
    try {
      const response = await getDashboardData();
      setData(response);
    } catch (err) {
      setError(err.message);
    } finally {
      setLoading(false);
    }
  };
  
  fetchData();
  const interval = setInterval(fetchData, 30000); // Auto-refresh
  return () => clearInterval(interval);
}, []);
```

---

### API Service Layer

**File**: `frontend/src/services/api.ts`

**Purpose**: Centralized API communication

**Pattern**: Axios instance with typed responses

**Example**:
```typescript
const api = axios.create({
  baseURL: 'http://localhost:8080/api',
  timeout: 10000,
});

export const getDashboardData = async (): Promise<DashboardData> => {
  const response = await api.get<DashboardData>('/dashboard');
  return response.data;
};

export const importXpub = async (xpub: string, gapLimit: number) => {
  const response = await api.post('/wallet/importxpub', {
    xpub,
    gapLimit,
  });
  return response.data;
};
```

---

## 🔌 Communication Protocols

### Frontend ↔ Backend

**Protocol**: HTTP/REST

**Format**: JSON

**Method**: Axios HTTP requests

**Endpoints**: `/api/*`

**Authentication**: None (localhost only)

**CORS**: Enabled for local development

---

### Backend ↔ dcrd

**Protocol**: JSON-RPC over HTTP/HTTPS

**Format**: JSON-RPC 2.0

**Port**: 9109 (RPC)

**Authentication**: Username + Password (RPC credentials)

**TLS**: Self-signed certificate

**Client**: `github.com/decred/dcrd/rpcclient`

---

### Backend ↔ dcrwallet

**Protocol**: JSON-RPC over HTTP/HTTPS

**Format**: JSON-RPC 2.0

**Port**: 9110 (Wallet RPC)

**Authentication**: Username + Password (separate credentials)

**TLS**: Self-signed certificate

**Client**: `github.com/decred/dcrd/rpcclient` (wallet mode)

---

### dcrd ↔ Decred Network

**Protocol**: Decred P2P wire protocol

**Port**: 9108 (P2P)

**Format**: Binary protocol messages

**Purpose**: Blockchain sync, transaction relay, block propagation

---

## 🐳 Docker Architecture

### Container Orchestration

**Tool**: Docker Compose

**Network**: Bridge network (`decred-pulse_decred-network`)

**Volumes**: 
- `dcrd-data` - Persistent blockchain (~10 GB)
- `dcrwallet-data` - Wallet database
- `certs` - Shared RPC certificates

---

### Container Dependencies

```
┌─────────────┐
│   frontend  │ (No dependencies, but needs backend)
└─────────────┘

┌─────────────┐
│   backend   │ (Depends on: dcrd, dcrwallet health)
└──────┬──────┘
       │
   ┌───▼────┐
   │  dcrd  │ (Independent, but backend waits for health)
   └────────┘

   ┌─────────────┐
   │  dcrwallet  │ (Depends on: dcrd)
   └─────────────┘
```

**Health Checks**:
- dcrd: RPC `getblockcount` response
- dcrwallet: RPC `walletinfo` response
- backend: HTTP `/api/health` response
- frontend: Nginx response

---

### Build Process

**Multi-stage builds**:
1. **Builder stage**: Compile from source
2. **Runtime stage**: Minimal image with binary

**dcrd build**:
```dockerfile
# Stage 1: Build from source
FROM golang:1.21-alpine AS builder
RUN git clone --depth 1 --branch ${DCRD_VERSION} https://github.com/decred/dcrd.git
WORKDIR /go/src/github.com/decred/dcrd
RUN go install .

# Stage 2: Runtime
FROM alpine:latest
COPY --from=builder /go/bin/dcrd /usr/local/bin/
# ... setup user, directories, entrypoint
```

---

## 🔒 Security Architecture

### Credential Management

**RPC Credentials**:
- Stored in `.env` file (not committed)
- Passed as environment variables to containers
- Never exposed to frontend
- Self-signed TLS certificates

**Certificate Handling**:
- Generated on first run
- Stored in Docker volume
- Shared between dcrd, dcrwallet, backend
- Backend skips verification (local only)

---

### Network Isolation

**Docker Network**:
- Private bridge network
- Containers communicate internally
- Only necessary ports exposed to host

**Port Exposure**:
- 3000: Frontend (public)
- 8080: Backend API (public)
- 9108: dcrd P2P (public, for peers)
- 9109: dcrd RPC (localhost only via Docker)
- 9110: dcrwallet RPC (localhost only via Docker)

---

### Frontend Security

**No sensitive data in frontend**:
- RPC credentials never sent to browser
- All RPC calls proxied through backend
- CORS restricted in production

---

## ⚡ Performance Considerations

### Backend Optimization

**Concurrent RPC Calls**:
```go
// Fetch multiple data sources concurrently
ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
defer cancel()

var (
    nodeStatus NodeStatus
    blockchain BlockchainInfo
    peers      []Peer
)

g, gctx := errgroup.WithContext(ctx)

g.Go(func() error {
    nodeStatus, err = fetchNodeStatus(gctx)
    return err
})

g.Go(func() error {
    blockchain, err = fetchBlockchainInfo(gctx)
    return err
})

g.Go(func() error {
    peers, err = fetchPeers(gctx)
    return err
})

if err := g.Wait(); err != nil {
    return nil, err
}
```

**Timeouts**:
- RPC calls: 10 seconds
- HTTP handlers: 30 seconds
- Wallet operations: 60 seconds

---

### Frontend Optimization

**Code Splitting**: Vite handles automatically

**Lazy Loading**: Components loaded on demand

**Memoization**: React.memo for expensive renders

**Debouncing**: For search/input fields

**Caching**: Axios response caching (if needed)

---

### Database Optimization

**dcrd**:
- LevelDB for blockchain storage
- Configurable cache (`dbcache`)
- Transaction indexing optional

**dcrwallet**:
- BoltDB for wallet storage
- Address caching
- Transaction indexing

---

## 📊 Monitoring & Observability

### Logging

**Backend**: Structured logging with Go's `log` package

**dcrd**: Configurable log levels (info, debug, trace)

**dcrwallet**: Separate wallet logs

**Frontend**: Browser console + network inspector

---

### Health Checks

**Backend**: `/api/health` endpoint

**dcrd**: RPC `getinfo` call

**dcrwallet**: RPC `walletinfo` call

**Docker**: Built-in health check commands

---

### Metrics

**Node Metrics**:
- Block height
- Peer count
- Mempool size
- Sync progress

**Wallet Metrics**:
- Balance
- Transaction count
- Ticket status
- Rescan progress

---

## 🚀 Deployment Architecture

### Development

```
Local Machine
├── Frontend: http://localhost:5173 (Vite dev server)
├── Backend: http://localhost:8080 (Go binary)
├── dcrd: Docker container
└── dcrwallet: Docker container
```

---

### Docker Compose (Recommended)

```
Host Machine
├── dcrd: Docker container
├── dcrwallet: Docker container
├── backend: Docker container
└── frontend: Docker container (Nginx)
```

**Access**: `http://localhost:3000`

---

### Production (Future)

```
Server
├── Nginx (Reverse Proxy + SSL)
│   ├── Frontend static files
│   └── Proxy to Backend API
├── Backend (Systemd service)
├── dcrd (Systemd service)
└── dcrwallet (Systemd service)
```

**Access**: `https://your-domain.com`

---

## 📚 Technology Choices

### Why Go for Backend?

- ✅ Native dcrd RPC client library
- ✅ Excellent concurrency (goroutines)
- ✅ Fast compilation and execution
- ✅ Strong typing
- ✅ Low memory footprint

### Why React for Frontend?

- ✅ Component-based architecture
- ✅ Large ecosystem
- ✅ TypeScript support
- ✅ Fast development
- ✅ Virtual DOM performance

### Why Docker Compose?

- ✅ Simple orchestration
- ✅ Reproducible environments
- ✅ Easy dependency management
- ✅ Cross-platform compatibility
- ✅ Development/production parity

### Why Tailwind CSS?

- ✅ Utility-first approach
- ✅ Rapid prototyping
- ✅ Consistent design system
- ✅ Small bundle size (purged)
- ✅ Great documentation

---

## 📚 Related Documentation

- **[Development Setup](development-setup.md)** - Local development guide
- **[Backend Guide](backend-guide.md)** - Backend development
- **[Frontend Guide](frontend-guide.md)** - Frontend development
- **[API Reference](../api/api-reference.md)** - API documentation

---

**Questions?** Check the [Development Setup](development-setup.md) or ask in the [Decred Community](https://decred.org/community/)


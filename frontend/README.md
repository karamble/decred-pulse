# Decred Dashboard - Frontend

Modern React dashboard for monitoring Decred dcrd node performance and network status.

## Features

- **Real-time Data**: Auto-refreshes every 30 seconds
- **Beautiful UI**: Modern dark theme with gradient effects and smooth animations
- **Responsive Design**: Works on all screen sizes
- **Live Metrics**: Node status, blockchain info, network stats, peer connections
- **RPC Configuration**: Connect to dcrd node via browser interface

## Tech Stack

- **React 18** with TypeScript
- **Vite** for fast development and optimized builds
- **Tailwind CSS** for styling
- **Axios** for API communication
- **Lucide React** for beautiful icons

## Prerequisites

- Node.js 18+ or Bun
- Running backend API (see `../backend`)

## Installation

```bash
cd frontend

# Using npm
npm install

# Or using bun (faster)
bun install
```

## Configuration

Create a `.env` file (optional):

```bash
VITE_API_URL=http://localhost:8080/api
```

If not set, defaults to `http://localhost:8080/api`.

## Development

```bash
# Using npm
npm run dev

# Or using bun
bun run dev
```

The app will be available at `http://localhost:3000`.

The Vite dev server is configured with a proxy to forward `/api` requests to the backend at `http://localhost:8080`.

## Building for Production

```bash
# Using npm
npm run build

# Or using bun
bun run build
```

The built files will be in the `dist/` directory.

## Preview Production Build

```bash
npm run preview
# or
bun run preview
```

## Project Structure

```
frontend/
├── src/
│   ├── components/         # React components
│   │   ├── NodeStatus.tsx
│   │   ├── MetricCard.tsx
│   │   ├── BlockchainInfo.tsx
│   │   ├── PeersList.tsx
│   │   └── RPCConnection.tsx
│   ├── services/           # API service layer
│   │   └── api.ts
│   ├── App.tsx            # Main application component
│   ├── main.tsx           # Application entry point
│   └── index.css          # Global styles
├── index.html
├── package.json
├── vite.config.ts
├── tailwind.config.js
└── tsconfig.json
```

## Features Overview

### Dashboard Components

1. **Node Status Card**
   - Shows node sync status (running, syncing, stopped)
   - Displays sync progress with animated progress bar
   - Shows node version

2. **Metric Cards** (8 total)
   - Circulating Supply
   - Network Peers
   - Block Height
   - Network Hashrate
   - Treasury Size
   - Supply Staked
   - Supply Mixed
   - Exchange Rate

3. **Blockchain Information Panel**
   - Latest block number
   - Block time
   - Chain size
   - Network difficulty

4. **RPC Connection Form**
   - Connect to dcrd node dynamically
   - Configure host, port, username, password
   - Shows connection status

5. **Peers List**
   - Shows all connected peers
   - Displays peer address, protocol, and latency
   - Scrollable list

### API Integration

The frontend communicates with the Go backend via REST API:

- `GET /api/dashboard` - Fetch all dashboard data
- `GET /api/node/status` - Get node status
- `GET /api/blockchain/info` - Get blockchain information
- `GET /api/network/peers` - Get connected peers
- `POST /api/connect` - Connect to RPC

## Environment Variables

- `VITE_API_URL` - Backend API URL (default: `http://localhost:8080/api`)

## Deployment

### Static Hosting

After building, deploy the `dist/` folder to any static hosting service:

- Netlify
- Vercel
- GitHub Pages
- AWS S3 + CloudFront
- Nginx

### Docker

```dockerfile
FROM node:18-alpine AS builder
WORKDIR /app
COPY package*.json ./
RUN npm install
COPY . .
RUN npm run build

FROM nginx:alpine
COPY --from=builder /app/dist /usr/share/nginx/html
EXPOSE 80
CMD ["nginx", "-g", "daemon off;"]
```

## Customization

### Colors

Edit `tailwind.config.js` to customize the color scheme:

```js
colors: {
  primary: { DEFAULT: 'hsl(217 91% 60%)' },
  // ... other colors
}
```

### API Endpoint

Change the backend URL in `src/services/api.ts` or via the `VITE_API_URL` environment variable.

### Auto-refresh Interval

Edit the interval in `src/App.tsx`:

```typescript
const interval = setInterval(fetchData, 30000); // 30 seconds
```

## Troubleshooting

### CORS Errors

Make sure the backend has CORS enabled for your frontend origin.

### API Not Found

Check that:
1. Backend is running on port 8080
2. `VITE_API_URL` is set correctly
3. Vite proxy is configured in `vite.config.ts`

### Build Errors

Clear cache and reinstall:
```bash
rm -rf node_modules dist
npm install
npm run build
```

## License

Part of the Decred Dashboard project.


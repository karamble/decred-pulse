// Copyright (c) 2015-2025 The Decred developers
// Use of this source code is governed by an ISC
// license that can be found in the LICENSE file.

import { useEffect, useState } from 'react';
import { 
  Users, Layers, TrendingUp, Coins, Wallet, 
  Lock, DollarSign, Shield 
} from 'lucide-react';
import { NodeStatus } from './components/NodeStatus';
import { MetricCard } from './components/MetricCard';
import { BlockchainInfo } from './components/BlockchainInfo';
import { PeersList } from './components/PeersList';
import { StakingStats } from './components/StakingStats';
import { MempoolActivity } from './components/MempoolActivity';
import { getDashboardData, DashboardData } from './services/api';

function App() {
  const [data, setData] = useState<DashboardData | null>(null);
  const [loading, setLoading] = useState(true);
  const [error, setError] = useState<string | null>(null);

  const fetchData = async () => {
    try {
      const dashboardData = await getDashboardData();
      setData(dashboardData);
      setError(null);
    } catch (err: any) {
      console.error('Error fetching dashboard data:', err);
      if (err.response?.status === 503) {
        setError('RPC client not connected. Please configure the connection below.');
      } else {
        setError(err.message || 'Failed to fetch data');
      }
    } finally {
      setLoading(false);
    }
  };

  useEffect(() => {
    fetchData();
    // Auto-refresh every 30 seconds
    const interval = setInterval(fetchData, 30000);
    return () => clearInterval(interval);
  }, []);

  return (
    <div className="min-h-screen bg-background p-6">
      <div className="max-w-7xl mx-auto space-y-6">
        {/* Header */}
        <div className="flex items-center justify-between mb-8 animate-fade-in">
          <div className="flex items-center gap-4">
            <div className="h-16 w-16 rounded-xl animate-pulse-glow flex items-center justify-center text-2xl font-bold bg-gradient-primary shadow-glow-primary">
              DCR
            </div>
            <div>
              <h1 className="text-4xl font-bold bg-gradient-primary bg-clip-text text-transparent mb-2">
                Decred Node
              </h1>
              <p className="text-muted-foreground">Monitor your dcrd node performance and network status</p>
            </div>
          </div>
          <div className="flex items-center gap-4">
            <a
              href="https://nodes.jholdstock.uk/"
              target="_blank"
              rel="noopener noreferrer"
              className="px-4 py-3 rounded-lg bg-primary/10 border border-primary/20 hover:bg-primary/20 transition-all duration-300 hover:shadow-glow-primary flex items-center gap-2"
            >
              <div className="h-10 w-10 rounded-lg flex items-center justify-center font-bold text-sm bg-gradient-primary">
                DCR
              </div>
              <span className="text-primary font-semibold">Node Mapper</span>
            </a>
            <div className="px-4 py-2 rounded-lg bg-primary/10 border border-primary/20">
              <p className="text-sm text-muted-foreground">Version</p>
              <p className="text-lg font-semibold text-primary">
                {data?.nodeStatus?.version || 'Loading...'}
              </p>
            </div>
          </div>
        </div>

        {/* Error Message */}
        {error && (
          <div className="p-4 rounded-lg bg-red-500/10 border border-red-500/20 animate-fade-in">
            <p className="text-red-500 font-medium">{error}</p>
          </div>
        )}

        {/* Loading State */}
        {loading && !data && (
          <div className="text-center py-12">
            <div className="inline-block animate-spin rounded-full h-12 w-12 border-4 border-primary border-t-transparent"></div>
            <p className="mt-4 text-muted-foreground">Loading dashboard data...</p>
          </div>
        )}

        {/* Node Status */}
        {data && (
          <NodeStatus 
            status={data.nodeStatus.status as any} 
            syncProgress={data.nodeStatus.syncProgress}
            version={data.nodeStatus.version}
            syncMessage={data.nodeStatus.syncMessage}
          />
        )}

        {/* Metrics Grid */}
        <div className="grid grid-cols-1 md:grid-cols-2 lg:grid-cols-4 gap-6">
          <MetricCard
            title="Circulating Supply"
            value={data?.supplyInfo?.circulatingSupply || 'Loading...'}
            subtitle="DCR in circulation"
            icon={Coins}
            trend={{ value: "Max 21M", isPositive: true }}
          />
          <MetricCard
            title="Network Peers"
            value={data?.networkInfo?.peerCount || 0}
            subtitle="Connected nodes"
            icon={Users}
          />
          <MetricCard
            title="Block Height"
            value={data?.blockchainInfo?.blockHeight?.toLocaleString() || 'Loading...'}
            subtitle="Latest block"
            icon={Layers}
          />
          <MetricCard
            title="Network Hashrate"
            value={data?.networkInfo?.hashrate || 'Loading...'}
            subtitle="Total network power"
            icon={TrendingUp}
          />
        </div>

        {/* Additional Metrics Grid */}
        <div className="grid grid-cols-1 md:grid-cols-2 lg:grid-cols-4 gap-6">
          <MetricCard
            title="Treasury Size"
            value={data?.supplyInfo?.treasurySize || 'Loading...'}
            subtitle="Self-funded from block reward"
            icon={Wallet}
          />
          <MetricCard
            title="Supply Staked"
            value={data?.supplyInfo?.stakedSupply || 'Loading...'}
            subtitle="Stakeholders Rule"
            icon={Lock}
          />
          <MetricCard
            title="Supply Mixed"
            value={data?.supplyInfo?.mixedPercent || 'Loading...'}
            subtitle="Privacy enhanced"
            icon={Shield}
            trend={{ value: "CoinShuffle++", isPositive: true }}
          />
          <MetricCard
            title="Exchange Rate"
            value={data?.supplyInfo?.exchangeRate || 'Loading...'}
            subtitle="USD per DCR"
            icon={DollarSign}
          />
        </div>

        {/* Details Grid */}
        <div className="grid grid-cols-1 lg:grid-cols-2 gap-6">
          <BlockchainInfo data={data?.blockchainInfo} />
          <StakingStats data={data?.stakingInfo} />
        </div>

        {/* Additional Grid */}
        <div className="grid grid-cols-1 lg:grid-cols-2 gap-6">
          <MempoolActivity data={data?.mempoolInfo} />
        </div>

        {/* Peers List */}
        <PeersList peers={data?.peers} />

        {/* Last Update */}
        {data && (
          <div className="text-center text-sm text-muted-foreground animate-fade-in">
            Last updated: {new Date(data.lastUpdate).toLocaleString()}
          </div>
        )}
      </div>
    </div>
  );
}

export default App;


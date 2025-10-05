// Copyright (c) 2015-2025 The Decred developers
// Use of this source code is governed by an ISC
// license that can be found in the LICENSE file.

import { useEffect, useState } from 'react';
import { Plus, RefreshCw } from 'lucide-react';
import { WalletStatus } from '../components/WalletStatus';
import { AccountInfo } from '../components/AccountInfo';
import { AccountsList } from '../components/AccountsList';
import { ImportXpubModal } from '../components/ImportXpubModal';
import { getWalletDashboard, WalletDashboardData, triggerRescan } from '../services/api';

export const WalletDashboard = () => {
  const [data, setData] = useState<WalletDashboardData | null>(null);
  const [loading, setLoading] = useState(true);
  const [error, setError] = useState<string | null>(null);
  const [showImportModal, setShowImportModal] = useState(false);
  const [rescanning, setRescanning] = useState(false);

  const fetchData = async () => {
    try {
      const walletData = await getWalletDashboard();
      setData(walletData);
      setError(null);
      
      // Clear rescanning state if rescan is no longer active
      if (rescanning && !walletData.walletStatus.rescanInProgress) {
        setRescanning(false);
      }
    } catch (err: any) {
      console.error('Error fetching wallet data:', err);
      
      // During rescan, keep showing last data if we have it
      // With increased timeouts, we should now get status updates even during rescan
      if (err.code === 'ECONNABORTED' || err.message?.includes('timeout') || err.response?.status === 408) {
        if (!data) {
          setError('Initializing wallet status. This may take a moment during rescan.');
        }
        // Don't clear existing data if we have it - progress bar will keep showing
      } else if (err.response?.status === 503) {
        setError('Wallet RPC not connected. Please ensure dcrwallet is running.');
      } else {
        setError(err.message || 'Failed to fetch wallet data');
      }
    } finally {
      setLoading(false);
    }
  };

  useEffect(() => {
    fetchData();
    // Auto-refresh every 10 seconds to show rescan progress
    const interval = setInterval(fetchData, 10000);
    return () => clearInterval(interval);
  }, []);

  const handleImportSuccess = (rescanEnabled: boolean) => {
    // If rescan is enabled, set rescanning state and wait for logs
    if (rescanEnabled) {
      setRescanning(true);
      // Wait 2 seconds for rescan to start writing to logs
      setTimeout(() => {
        fetchData();
      }, 2000);
    } else {
      // No rescan, just refresh data immediately
      fetchData();
    }
  };

  const handleRescan = async () => {
    if (rescanning) return;
    
    if (!confirm('This will rescan the entire blockchain from block 0. This may take 30+ minutes. Continue?')) {
      return;
    }

    try {
      setRescanning(true);
      await triggerRescan();
      setError(null);
      
      // Immediately fetch status to show progress bar
      // Keep rescanning state true - it will be cleared by polling when rescan completes
      setTimeout(() => {
        fetchData();
      }, 2000); // Wait 2 seconds for rescan to start writing to logs
    } catch (err: any) {
      console.error('Error triggering rescan:', err);
      setError(err.response?.data?.error || err.message || 'Failed to trigger rescan');
      setRescanning(false); // Only clear on error
    }
  };

  return (
    <div className="space-y-6">
      {/* Action Buttons */}
      <div className="flex justify-end gap-3">
        <button
          onClick={handleRescan}
          disabled={rescanning || data?.walletStatus.status === 'no_wallet'}
          className="px-6 py-3 rounded-lg bg-muted/20 text-foreground font-semibold hover:bg-muted/30 transition-all flex items-center gap-2 disabled:opacity-50 disabled:cursor-not-allowed"
        >
          <RefreshCw className={`h-5 w-5 ${rescanning ? 'animate-spin' : ''}`} />
          {rescanning ? 'Rescanning...' : 'Rescan'}
        </button>
        <button
          onClick={() => setShowImportModal(true)}
          className="px-6 py-3 rounded-lg bg-gradient-primary text-white font-semibold hover:shadow-glow-primary transition-all flex items-center gap-2"
        >
          <Plus className="h-5 w-5" />
          Add X-Pub
        </button>
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
          <p className="mt-4 text-muted-foreground">Loading wallet data...</p>
        </div>
      )}

      {/* Wallet Status */}
      {data && (
        <WalletStatus
          status={data.walletStatus.status as any}
          syncProgress={data.walletStatus.syncProgress}
          version={data.walletStatus.version}
          syncMessage={data.walletStatus.syncMessage}
          unlocked={data.walletStatus.unlocked}
        />
      )}

      {/* Account Info */}
      {data && data.walletStatus.status !== 'no_wallet' && (
        <AccountInfo
          accountName={data.accountInfo.accountName}
          totalBalance={data.accountInfo.totalBalance}
          spendableBalance={data.accountInfo.spendableBalance}
          immatureBalance={data.accountInfo.immatureBalance}
          unconfirmedBalance={data.accountInfo.unconfirmedBalance}
        />
      )}

      {/* Accounts List */}
      {data && data.walletStatus.status !== 'no_wallet' && data.accounts && (
        <AccountsList accounts={data.accounts} />
      )}

      {/* Last Update */}
      {data && (
        <div className="text-center text-sm text-muted-foreground animate-fade-in">
          Last updated: {new Date(data.lastUpdate).toLocaleString()}
        </div>
      )}

      {/* Import Xpub Modal */}
      <ImportXpubModal
        isOpen={showImportModal}
        onClose={() => setShowImportModal(false)}
        onSuccess={handleImportSuccess}
      />
    </div>
  );
};


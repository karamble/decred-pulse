// Copyright (c) 2015-2025 The Decred developers
// Use of this source code is governed by an ISC
// license that can be found in the LICENSE file.

import { useEffect, useState } from 'react';
import { Plus, RefreshCw } from 'lucide-react';
import { WalletStatus } from '../components/WalletStatus';
import { AccountInfo } from '../components/AccountInfo';
import { AccountsList } from '../components/AccountsList';
import { ImportXpubModal } from '../components/ImportXpubModal';
import { SyncProgressBar } from '../components/SyncProgressBar';
import { TicketPoolInfo } from '../components/TicketPoolInfo';
import { MyTicketsInfo } from '../components/MyTicketsInfo';
import { getWalletDashboard, WalletDashboardData, triggerRescan, getSyncProgress, SyncProgressData } from '../services/api';

export const WalletDashboard = () => {
  const [data, setData] = useState<WalletDashboardData | null>(null);
  const [loading, setLoading] = useState(true);
  const [error, setError] = useState<string | null>(null);
  const [showImportModal, setShowImportModal] = useState(false);
  const [rescanning, setRescanning] = useState(false);
  const [syncProgress, setSyncProgress] = useState<SyncProgressData | null>(null);
  const [showSyncProgress, setShowSyncProgress] = useState(false);

  const fetchSyncProgress = async () => {
    try {
      const progress = await getSyncProgress();
      setSyncProgress(progress);
      
      // Check if rescan is active (backend returns false when logs are stale/complete)
      if (!progress.isRescanning) {
        // Rescan is complete or stale - stop polling
        if (showSyncProgress || rescanning) {
          console.log('Rescan no longer active (stale or complete) - stopping sync progress polling');
          const wasRescanning = showSyncProgress;
          setShowSyncProgress(false);
          setRescanning(false);
          
          // Resume normal wallet data polling
          if (wasRescanning) {
            console.log('Resuming normal wallet data polling');
            setError(null);
            setLoading(false);
            fetchData(true);
          }
        }
        return; // Stop processing - rescan is not active
      }
      
      // Rescan is active
      if (progress.progress < 99) {
        // Still rescanning - show progress bar
        setShowSyncProgress(true);
        setRescanning(true);
      } else {
        // Progress reached 99% - rescan complete
        if (showSyncProgress || rescanning) {
          console.log('Rescan completed (99%) - stopping sync progress polling');
          const wasRescanning = showSyncProgress;
          setShowSyncProgress(false);
          setRescanning(false);
          
          // Resume normal wallet data polling
          if (wasRescanning) {
            console.log('Resuming normal wallet data polling');
            setError(null);
            setLoading(false);
            fetchData(true);
          }
        }
      }
    } catch (err: any) {
      console.error('Error fetching sync progress:', err);
      // Don't show error for sync progress polling - it's optional
    }
  };

  const fetchData = async (force = false) => {
    // Skip wallet data fetching during rescan - only poll sync progress
    // Unless force=true (when explicitly called after rescan completes)
    if (showSyncProgress && !force) {
      console.log('Skipping wallet data fetch - rescan in progress');
      return;
    }

    try {
      const walletData = await getWalletDashboard();
      setData(walletData);
      setError(null);
      
      // Clear rescanning state if rescan is no longer active
      if (rescanning && !walletData.walletStatus.rescanInProgress) {
        setRescanning(false);
        setShowSyncProgress(false);
      }
    } catch (err: any) {
      console.error('Error fetching wallet data:', err);
      
      // Handle errors appropriately
      if (err.code === 'ECONNABORTED' || err.message?.includes('timeout') || err.response?.status === 408) {
        if (!data) {
          setError('Initializing wallet status. This may take a moment.');
        }
        // Don't clear existing data if we have it
      } else if (err.response?.status === 503) {
        setError('Wallet RPC not connected. Please ensure dcrwallet is running.');
      } else {
        setError(err.message || 'Failed to fetch wallet data');
      }
    } finally {
      setLoading(false);
    }
  };

  // Initial data fetch - check for active rescan first
  useEffect(() => {
    const initialize = async () => {
      // First check if there's an active rescan
      try {
        const progress = await getSyncProgress();
        if (progress.isRescanning && progress.progress < 99) {
          // Rescan is active - show progress bar and skip wallet data fetch
          console.log('Active rescan detected on load - showing progress bar');
          setSyncProgress(progress);
          setShowSyncProgress(true);
          setRescanning(true);
          setLoading(false);
          setError(null); // Clear any errors
          return; // Don't fetch wallet data
        }
      } catch (err) {
        console.log('No active rescan detected, proceeding with normal fetch');
      }
      
      // No active rescan - fetch wallet data normally
      fetchData();
    };
    
    initialize();
  }, []);

  // Poll wallet data only when NOT rescanning
  useEffect(() => {
    if (!showSyncProgress) {
      // Auto-refresh every 10 seconds when not rescanning
      const interval = setInterval(fetchData, 10000);
      return () => clearInterval(interval);
    }
  }, [showSyncProgress]);

  // Poll sync progress more frequently when rescanning
  useEffect(() => {
    if (rescanning || showSyncProgress) {
      // Start polling immediately
      fetchSyncProgress();
      
      // Poll every 3 seconds during rescan
      const syncInterval = setInterval(fetchSyncProgress, 3000);
      return () => clearInterval(syncInterval);
    }
  }, [rescanning, showSyncProgress]);

  const handleImportSuccess = (rescanEnabled: boolean) => {
    // If rescan is enabled, set rescanning state and show progress bar
    if (rescanEnabled) {
      setRescanning(true);
      setShowSyncProgress(true);
      setError(null); // Clear any existing errors
      // Wait 2 seconds for rescan to start writing to logs, then start polling
      setTimeout(() => {
        fetchSyncProgress();
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
      setShowSyncProgress(true); // Show progress bar immediately
      await triggerRescan();
      setError(null);
      
      // Start polling sync progress after a brief delay for logs to start
      setTimeout(() => {
        fetchSyncProgress();
      }, 2000); // Wait 2 seconds for rescan to start writing to logs
    } catch (err: any) {
      console.error('Error triggering rescan:', err);
      setError(err.response?.data?.error || err.message || 'Failed to trigger rescan');
      setRescanning(false); // Only clear on error
      setShowSyncProgress(false);
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

      {/* Sync Progress Bar - shown during rescan */}
      {showSyncProgress && syncProgress && (
        <SyncProgressBar
          progress={syncProgress.progress}
          scanHeight={syncProgress.scanHeight}
          chainHeight={syncProgress.chainHeight}
          message={syncProgress.message}
        />
      )}

      {/* Error Message - hide if sync progress bar is showing */}
      {error && !showSyncProgress && (
        <div className="p-4 rounded-lg bg-red-500/10 border border-red-500/20 animate-fade-in">
          <p className="text-red-500 font-medium">{error}</p>
        </div>
      )}

      {/* Loading State */}
      {loading && !data && !showSyncProgress && (
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

      {/* Staking Information */}
      {data && data.walletStatus.status !== 'no_wallet' && data.stakingInfo && (
        <>
          {/* Ticket Pool & Difficulty Info */}
          <TicketPoolInfo
            poolSize={data.stakingInfo.poolSize}
            currentDifficulty={data.stakingInfo.currentDifficulty}
            nextDifficulty={data.stakingInfo.nextDifficulty}
            estimatedMin={data.stakingInfo.estimatedMin}
            estimatedMax={data.stakingInfo.estimatedMax}
            estimatedExpected={data.stakingInfo.estimatedExpected}
            allMempoolTix={data.stakingInfo.allMempoolTix}
          />

          {/* My Tickets Info */}
          <MyTicketsInfo
            ownMempoolTix={data.stakingInfo.ownMempoolTix}
            immature={data.stakingInfo.immature}
            unspent={data.stakingInfo.unspent}
            voted={data.stakingInfo.voted}
            revoked={data.stakingInfo.revoked}
            unspentExpired={data.stakingInfo.unspentExpired}
            totalSubsidy={data.stakingInfo.totalSubsidy}
          />
        </>
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


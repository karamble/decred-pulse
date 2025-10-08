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
import { TransactionHistory } from '../components/TransactionHistory';
import { getWalletDashboard, WalletDashboardData, triggerRescan, getSyncProgress, streamRescanProgress, SyncProgressData } from '../services/api';

export const WalletDashboard = () => {
  const [data, setData] = useState<WalletDashboardData | null>(null);
  const [loading, setLoading] = useState(true);
  const [error, setError] = useState<string | null>(null);
  const [showImportModal, setShowImportModal] = useState(false);
  const [syncProgress, setSyncProgress] = useState<SyncProgressData | null>(null);
  const [showSyncProgress, setShowSyncProgress] = useState(false);

  const fetchData = async () => {
    try {
      const walletData = await getWalletDashboard();
      setData(walletData);
      setError(null);
      
      // Check if wallet is syncing and show unified progress bar
      if (walletData.walletStatus.status === 'syncing' && !showSyncProgress) {
        console.log('Wallet syncing detected - activating progress bar stream');
        setSyncProgress({
          progress: walletData.walletStatus.syncProgress || 0,
          scanHeight: walletData.walletStatus.syncHeight || 0,
          chainHeight: 1016874, // Will be updated by WebSocket
          message: walletData.walletStatus.syncMessage || 'Connecting to sync stream...',
          isRescanning: true,
        });
        setShowSyncProgress(true);
      }
    } catch (err: any) {
      console.error('Error fetching wallet data:', err);
      
      // Handle errors appropriately
      if (err.code === 'ECONNABORTED' || err.message?.includes('timeout') || err.response?.status === 408) {
        if (!data) {
          setError('Initializing wallet status. This may take a moment.');
        }
      } else if (err.response?.status === 503) {
        setError('Wallet RPC not connected. Please ensure dcrwallet is running.');
      } else {
        setError(err.message || 'Failed to fetch wallet data');
      }
    } finally {
      setLoading(false);
    }
  };

  // Initial load - check for active rescan first
  useEffect(() => {
    const initialize = async () => {
      try {
        const progress = await getSyncProgress();
        if (progress.isRescanning && progress.progress < 100) {
          // Active rescan detected - show progress bar only
          console.log('Active rescan detected on load - showing progress bar');
          setSyncProgress(progress);
          setShowSyncProgress(true);
          setLoading(false);
          // Still fetch wallet data for the status card, but in background
          fetchData();
          return;
        }
      } catch (err) {
        console.log('No active rescan, loading wallet data normally');
      }
      
      // No active rescan - fetch wallet data normally
      fetchData();
    };
    
    initialize();
  }, []);

  // Auto-refresh wallet data every 10 seconds (but NOT during rescan)
  useEffect(() => {
    if (!showSyncProgress) {
      const interval = setInterval(fetchData, 10000);
      return () => clearInterval(interval);
    }
  }, [showSyncProgress]);

  // WebSocket streaming for rescan progress (replaces polling)
  useEffect(() => {
    if (showSyncProgress) {
      console.log('Starting WebSocket stream for rescan progress');
      
      const cleanup = streamRescanProgress(
        // onProgress callback
        (progress) => {
          console.log('Received progress update:', progress);
          setSyncProgress(progress);
          
          // Only hide progress bar when we get explicit completion signal
          // Don't close on temporary "not rescanning" states
          if (progress.message === "Rescan complete" && progress.progress >= 99) {
            console.log('Rescan completed - hiding progress bar');
            setShowSyncProgress(false);
            fetchData(); // Refresh wallet data to show updated balances
          }
        },
        // onError callback
        (error) => {
          console.error('WebSocket error:', error);
          setError('Failed to stream rescan progress. Please refresh the page.');
        },
        // onClose callback
        () => {
          console.log('WebSocket stream closed');
          // When WebSocket closes, refresh data and hide progress bar
          setShowSyncProgress(false);
          fetchData();
        }
      );

      // Cleanup WebSocket connection when component unmounts or showSyncProgress changes
      return cleanup;
    }
  }, [showSyncProgress]);

  const handleImportSuccess = () => {
    // Show progress bar after xpub import (WebSocket will auto-connect)
    setShowSyncProgress(true);
    setError(null);
    // Set initial placeholder sync progress
    setSyncProgress({
      isRescanning: true,
      scanHeight: 0,
      chainHeight: 1,
      progress: 0,
      message: 'Connecting to rescan stream...'
    });
  };

  const handleRescan = async () => {
    if (showSyncProgress) return; // Already rescanning
    
    if (!confirm('This will rescan the entire blockchain from block 0. This may take 30+ minutes. Continue?')) {
      return;
    }

    try {
      // Show progress bar immediately (WebSocket will auto-connect)
      setShowSyncProgress(true);
      setSyncProgress({
        isRescanning: true,
        scanHeight: 0,
        chainHeight: 1,
        progress: 0,
        message: 'Connecting to rescan stream...'
      });
      
      await triggerRescan();
      setError(null);
    } catch (err: any) {
      console.error('Error triggering rescan:', err);
      setError(err.response?.data?.error || err.message || 'Failed to trigger rescan');
      setShowSyncProgress(false);
    }
  };

  return (
    <div className="space-y-6">
      {/* Action Buttons */}
      <div className="flex justify-end gap-3">
        <button
          onClick={handleRescan}
          disabled={showSyncProgress || data?.walletStatus.status === 'no_wallet'}
          className="px-6 py-3 rounded-lg bg-muted/20 text-foreground font-semibold hover:bg-muted/30 transition-all flex items-center gap-2 disabled:opacity-50 disabled:cursor-not-allowed"
        >
          <RefreshCw className={`h-5 w-5 ${showSyncProgress ? 'animate-spin' : ''}`} />
          {showSyncProgress ? 'Rescanning...' : 'Rescan'}
        </button>
        <button
          onClick={() => setShowImportModal(true)}
          className="px-6 py-3 rounded-lg bg-gradient-primary text-white font-semibold transition-all flex items-center gap-2"
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

      {/* Error Message - hide during rescan */}
      {error && !showSyncProgress && (
        <div className="p-4 rounded-lg bg-red-500/10 border border-red-500/20 animate-fade-in">
          <p className="text-red-500 font-medium">{error}</p>
        </div>
      )}

      {/* Loading State - hide during rescan */}
      {loading && !data && !showSyncProgress && (
        <div className="text-center py-12">
          <div className="inline-block animate-spin rounded-full h-12 w-12 border-4 border-primary border-t-transparent"></div>
          <p className="mt-4 text-muted-foreground">Loading wallet data...</p>
        </div>
      )}

      {/* Wallet Status - always visible, but hides when unified progress bar is shown */}
      {data && !showSyncProgress && (
        <WalletStatus
          status={data.walletStatus.status as any}
          version={data.walletStatus.version}
          unlocked={data.walletStatus.unlocked}
        />
      )}

      {/* Hide wallet data cards during rescan to prevent RPC flooding */}
      {!showSyncProgress && (
        <>

          {/* Left Column: Account Balance + Accounts | Right Column: Transaction History */}
          {data && data.walletStatus.status !== 'no_wallet' && (
        <div className="grid grid-cols-1 lg:grid-cols-2 gap-6 animate-fade-in">
          {/* Left Column */}
          <div className="space-y-6">
            {/* Account Info */}
            <AccountInfo
              accountName={data.accountInfo.accountName}
              totalBalance={data.accountInfo.totalBalance}
              spendableBalance={data.accountInfo.spendableBalance}
              immatureBalance={data.accountInfo.immatureBalance}
              unconfirmedBalance={data.accountInfo.unconfirmedBalance}
              lockedByTickets={data.accountInfo.lockedByTickets}
              cumulativeTotal={data.accountInfo.cumulativeTotal}
              totalSpendable={data.accountInfo.totalSpendable}
              totalLockedByTickets={data.accountInfo.totalLockedByTickets}
            />

            {/* Accounts List */}
            {data.accounts && (
              <AccountsList accounts={data.accounts} />
            )}
          </div>

          {/* Right Column: Transaction History */}
          {!loading && !showSyncProgress && (
            <TransactionHistory />
          )}
        </div>
      )}

          {/* Row 2: Ticket Pool | My Tickets */}
          {data && data.walletStatus.status !== 'no_wallet' && data.stakingInfo && (
            <div className="grid grid-cols-1 lg:grid-cols-2 gap-6 animate-fade-in">
              {/* Ticket Pool & Difficulty Info */}
              <TicketPoolInfo
                poolSize={data.stakingInfo.poolSize}
                currentDifficulty={data.stakingInfo.currentDifficulty}
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
            </div>
          )}

          {/* Last Update */}
          {data && (
            <div className="text-center text-sm text-muted-foreground animate-fade-in">
              Last updated: {new Date(data.lastUpdate).toLocaleString()}
            </div>
          )}
        </>
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


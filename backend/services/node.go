// Copyright (c) 2015-2025 The Decred developers
// Use of this source code is governed by an ISC
// license that can be found in the LICENSE file.

package services

import (
	"context"
	"fmt"
	"log"
	"math"
	"sync"
	"time"

	"decred-pulse-backend/rpc"
	"decred-pulse-backend/types"
	"decred-pulse-backend/utils"
)

var (
	// Track previous sync values to calculate delta
	prevHeaders int64
	prevBlocks  int64
	syncMutex   sync.Mutex
)

func FetchDashboardData() (*types.DashboardData, error) {
	nodeStatus, err := FetchNodeStatus()
	if err != nil {
		return nil, err
	}

	blockchainInfo, err := FetchBlockchainInfo()
	if err != nil {
		return nil, err
	}

	networkInfo, err := FetchNetworkInfo()
	if err != nil {
		return nil, err
	}

	peers, err := FetchPeers()
	if err != nil {
		return nil, err
	}

	supplyInfo, err := FetchSupplyInfo()
	if err != nil {
		return nil, err
	}

	stakingInfo, err := FetchStakingInfo()
	if err != nil {
		return nil, err
	}

	mempoolInfo, err := FetchMempoolInfo()
	if err != nil {
		return nil, err
	}

	return &types.DashboardData{
		NodeStatus:     *nodeStatus,
		BlockchainInfo: *blockchainInfo,
		NetworkInfo:    *networkInfo,
		Peers:          peers,
		SupplyInfo:     *supplyInfo,
		StakingInfo:    *stakingInfo,
		MempoolInfo:    *mempoolInfo,
		LastUpdate:     time.Now(),
	}, nil
}

func FetchNodeStatus() (*types.NodeStatus, error) {
	ctx := context.Background()

	// Get version info using version command
	versionInfo, err := rpc.DcrdClient.Version(ctx)
	if err != nil {
		return nil, fmt.Errorf("failed to get version: %v", err)
	}

	// Get blockchain info for accurate sync status
	chainInfo, err := rpc.DcrdClient.GetBlockChainInfo(ctx)
	if err != nil {
		return nil, err
	}

	// Debug logging
	log.Printf("Blockchain sync status - InitialBlockDownload: %v, Blocks: %d, Headers: %d, SyncHeight: %d, VerificationProgress: %f",
		chainInfo.InitialBlockDownload, chainInfo.Blocks, chainInfo.Headers, chainInfo.SyncHeight, chainInfo.VerificationProgress)

	// Calculate sync progress based on actual blockchain sync
	var syncProgress float64
	var status string
	var syncPhase string
	var syncMessage string

	// Thread-safe access to previous values
	syncMutex.Lock()
	currentHeaders := chainInfo.Headers
	currentBlocks := chainInfo.Blocks
	deltaHeaders := currentHeaders - prevHeaders
	deltaBlocks := currentBlocks - prevBlocks
	prevHeaders = currentHeaders
	prevBlocks = currentBlocks
	syncMutex.Unlock()

	if chainInfo.InitialBlockDownload {
		// Node is still syncing
		status = "syncing"

		// Determine sync phase: headers or blocks
		if chainInfo.Blocks == 0 && chainInfo.Headers > 0 {
			// Syncing headers
			syncPhase = "headers"
			if chainInfo.SyncHeight > 0 {
				syncProgress = (float64(chainInfo.Headers) / float64(chainInfo.SyncHeight)) * 100
			}
			if deltaHeaders > 0 {
				syncMessage = fmt.Sprintf("Processed %s headers in the last 30 seconds", utils.FormatNumber(deltaHeaders))
			} else {
				syncMessage = "Syncing headers..."
			}
		} else if chainInfo.Blocks > 0 {
			// Syncing blocks
			syncPhase = "blocks"
			if chainInfo.SyncHeight > 0 {
				syncProgress = (float64(chainInfo.Blocks) / float64(chainInfo.SyncHeight)) * 100
			}
			if deltaBlocks > 0 {
				syncMessage = fmt.Sprintf("Processed %s blocks in the last 30 seconds", utils.FormatNumber(deltaBlocks))
			} else {
				syncMessage = "Syncing blocks..."
			}
		} else {
			// Initial state
			syncPhase = "starting"
			syncMessage = "Starting sync..."
			syncProgress = 0
		}

		// Use verification progress if available and more accurate
		if chainInfo.VerificationProgress > 0 {
			syncProgress = chainInfo.VerificationProgress * 100
		}
	} else {
		// Node is fully synced
		status = "running"
		syncProgress = 100.0
		syncPhase = "synced"
		syncMessage = "Fully synced"
	}

	// Ensure progress is between 0 and 100
	if syncProgress > 100 {
		syncProgress = 100
	} else if syncProgress < 0 {
		syncProgress = 0
	}

	return &types.NodeStatus{
		Status:       status,
		SyncProgress: syncProgress,
		Version:      fmt.Sprintf("%d.%d.%d", versionInfo["dcrd"].Major, versionInfo["dcrd"].Minor, versionInfo["dcrd"].Patch),
		SyncPhase:    syncPhase,
		SyncMessage:  syncMessage,
	}, nil
}

func FetchBlockchainInfo() (*types.BlockchainInfo, error) {
	ctx := context.Background()
	info, err := rpc.DcrdClient.GetBlockChainInfo(ctx)
	if err != nil {
		return nil, err
	}

	bestBlockHash, err := rpc.DcrdClient.GetBestBlockHash(ctx)
	if err != nil {
		return nil, err
	}

	blockHeader, err := rpc.DcrdClient.GetBlockHeader(ctx, bestBlockHash)
	if err != nil {
		return nil, err
	}

	// Estimate block time (time since last block)
	timeSinceBlock := time.Since(blockHeader.Timestamp)
	blockTime := fmt.Sprintf("%dm %ds", int(timeSinceBlock.Minutes()), int(timeSinceBlock.Seconds())%60)

	return &types.BlockchainInfo{
		BlockHeight: info.Blocks,
		BlockHash:   bestBlockHash.String(),
		Difficulty:  float64(info.Difficulty),
		ChainSize:   0, // Would need to calculate from disk usage
		BlockTime:   blockTime,
	}, nil
}

func FetchNetworkInfo() (*types.NetworkInfo, error) {
	ctx := context.Background()

	// Get peer count
	peerCount := 0
	peerInfo, err := rpc.DcrdClient.GetPeerInfo(ctx)
	if err == nil {
		peerCount = len(peerInfo)
	}

	// Get network difficulty and calculate hashrate - direct RPC method
	hashrateStr := "N/A"
	networkHashPS := float64(0)

	difficulty, err := rpc.DcrdClient.GetDifficulty(ctx)
	if err == nil && difficulty > 0 {
		// Calculate network hashrate from difficulty
		// Formula: hashrate = difficulty * 2^32 / target_block_time
		// For Decred, target block time is 5 minutes = 300 seconds
		networkHashPS = difficulty * math.Pow(2, 32) / 300
		hashrateStr = utils.FormatHashrate(networkHashPS)
	}

	return &types.NetworkInfo{
		PeerCount:     peerCount,
		Hashrate:      hashrateStr,
		NetworkHashPS: networkHashPS,
	}, nil
}

func FetchPeers() ([]types.Peer, error) {
	ctx := context.Background()
	peerInfo, err := rpc.DcrdClient.GetPeerInfo(ctx)
	if err != nil {
		return nil, err
	}

	peers := make([]types.Peer, 0, len(peerInfo))
	now := time.Now().Unix()

	for i, p := range peerInfo {
		// Convert pingtime from microseconds to milliseconds
		latency := fmt.Sprintf("%.0fms", float64(p.PingTime)/1000)

		// Calculate connection time
		connDuration := now - p.ConnTime
		connTime := utils.FormatDuration(connDuration)

		// Calculate total traffic in MB
		totalBytes := p.BytesSent + p.BytesRecv
		traffic := utils.FormatTraffic(totalBytes)

		// Extract version from subver string (e.g., "/dcrwire:1.0.0/dcrd:2.0.5/" -> "2.0.5")
		version := utils.ExtractDcrdVersion(p.SubVer)

		peers = append(peers, types.Peer{
			ID:         i + 1,
			Address:    p.Addr,
			Protocol:   "TCP",
			Latency:    latency,
			ConnTime:   connTime,
			Traffic:    traffic,
			Version:    version,
			IsSyncNode: p.SyncNode,
		})
	}

	return peers, nil
}

// formatDuration formats a duration in seconds to a human-readable string

func FetchSupplyInfo() (*types.SupplyInfo, error) {
	ctx := context.Background()

	// Get real circulating supply from dcrd - direct RPC method
	circulatingSupply := "N/A"
	stakedSupply := "N/A"
	stakedPercent := float64(0)
	treasuryBalance := "N/A"

	// Check if node is fully synced before calling TicketPoolValue
	chainInfo, err := rpc.DcrdClient.GetBlockChainInfo(ctx)
	isSynced := err == nil && !chainInfo.InitialBlockDownload

	coinSupply, err := rpc.DcrdClient.GetCoinSupply(ctx)
	if err == nil && coinSupply > 0 {
		// Convert atoms to DCR and format with commas
		coinSupplyDCR := coinSupply.ToCoin()
		circulatingSupply = utils.FormatDCRAmount(coinSupplyDCR)

		// Calculate staked supply from ticket pool
		// Only call GetTicketPoolValue if node is fully synced to avoid nil pointer panic during initial sync
		if isSynced {
			ticketPoolValue, err := rpc.DcrdClient.GetTicketPoolValue(ctx)
			if err == nil && ticketPoolValue > 0 {
				lockedDCR := ticketPoolValue.ToCoin()
				stakedSupply = utils.FormatDCRAmount(lockedDCR)

				if coinSupplyDCR > 0 {
					stakedPercent = (lockedDCR / coinSupplyDCR) * 100
				}
			}
		}
	}

	// Get treasury balance - direct RPC method
	// Pass nil for hash (gets latest) and false for verbose
	treasuryBalanceResult, err := rpc.DcrdClient.GetTreasuryBalance(ctx, nil, false)
	if err == nil && treasuryBalanceResult.Balance > 0 {
		// Balance is in atoms (uint64), convert to DCR by dividing by 1e8
		treasuryBalanceDCR := float64(treasuryBalanceResult.Balance) / 1e8
		// Format with 2 decimal places and commas
		treasuryBalance = utils.FormatDCRAmountWithDecimals(treasuryBalanceDCR, 2)
	}

	return &types.SupplyInfo{
		CirculatingSupply: circulatingSupply,
		StakedSupply:      stakedSupply,
		StakedPercent:     stakedPercent,
		ExchangeRate:      "N/A", // Requires external API
		TreasurySize:      treasuryBalance,
		MixedPercent:      "N/A", // Requires mixer statistics
	}, nil
}

func FetchStakingInfo() (*types.StakingInfo, error) {
	ctx := context.Background()

	// Check if node is fully synced before calling TicketPoolValue
	chainInfo, err := rpc.DcrdClient.GetBlockChainInfo(ctx)
	isSynced := err == nil && !chainInfo.InitialBlockDownload

	// Get stake difficulty (ticket price) - direct RPC method
	stakeDiff, err := rpc.DcrdClient.GetStakeDifficulty(ctx)
	if err != nil {
		return nil, fmt.Errorf("failed to get stake difficulty: %v", err)
	}
	ticketPrice := stakeDiff.CurrentStakeDifficulty

	// Get live tickets from pool - direct RPC method
	// LiveTickets returns []*chainhash.Hash directly
	liveTickets, err := rpc.DcrdClient.LiveTickets(ctx)
	poolSize := uint32(0)
	if err == nil && liveTickets != nil {
		// Count the actual number of live tickets
		poolSize = uint32(len(liveTickets))
	}

	// Get ticket pool value (total locked DCR) - direct RPC method
	// Returns dcrutil.Amount which needs to be converted to float64 DCR
	// Only call GetTicketPoolValue if node is fully synced to avoid nil pointer panic during initial sync
	lockedDCR := float64(0)
	if isSynced {
		poolValue, err := rpc.DcrdClient.GetTicketPoolValue(ctx)
		if err == nil {
			lockedDCR = poolValue.ToCoin()
		}
	}

	// Get total coin supply for participation rate calculation - direct RPC method
	// Returns dcrutil.Amount which needs to be converted to float64 DCR
	participationRate := float64(0)
	coinSupply, err := rpc.DcrdClient.GetCoinSupply(ctx)
	if err == nil && coinSupply > 0 {
		// Calculate participation rate as percentage of total supply
		coinSupplyDCR := coinSupply.ToCoin()
		participationRate = (lockedDCR / coinSupplyDCR) * 100
	}

	return &types.StakingInfo{
		TicketPrice:       ticketPrice,
		PoolSize:          poolSize,
		LockedDCR:         lockedDCR,
		ParticipationRate: participationRate,
		AllMempoolTix:     0, // Would need livetickets/missedtickets commands
		Immature:          0, // Not available from dcrd alone
		Live:              poolSize,
		Voted:             0, // Would need block analysis
		Missed:            0, // Would need missedtickets command
		Revoked:           0, // Would need block analysis
	}, nil
}

func FetchMempoolInfo() (*types.MempoolInfo, error) {
	ctx := context.Background()

	// Get raw mempool - this is a standard method available in rpcclient
	hashes, err := rpc.DcrdClient.GetRawMempool(ctx, "false")
	if err != nil {
		// If mempool query fails (e.g., during sync), return empty mempool
		return &types.MempoolInfo{
			Size:           0,
			Bytes:          0,
			TxCount:        0,
			TotalFee:       0,
			AverageFeeRate: 0,
		}, nil
	}

	// Count transactions - this is the only real data we have
	txCount := len(hashes)

	// We don't have access to actual size and fees without querying each transaction
	// Show 0 instead of estimates
	return &types.MempoolInfo{
		Size:           uint64(txCount),
		Bytes:          0,
		TxCount:        txCount,
		TotalFee:       0,
		AverageFeeRate: 0,
	}, nil
}

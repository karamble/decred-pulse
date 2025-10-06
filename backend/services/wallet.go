// Copyright (c) 2015-2025 The Decred developers
// Use of this source code is governed by an ISC
// license that can be found in the LICENSE file.

package services

import (
	"context"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"os"
	"regexp"
	"strconv"
	"strings"
	"sync"
	"time"

	"decred-pulse-backend/rpc"
	"decred-pulse-backend/types"
	"decred-pulse-backend/utils"
)

var (
	// Track wallet sync
	prevWalletHeight int64
	walletSyncMutex  sync.Mutex
)

func ParseWalletLogsForRescan() (bool, int64, error) {
	// Path to wallet log file (mounted from dcrwallet-data volume)
	logPath := "/wallet-data/logs/mainnet/dcrwallet.log"

	// Check if log file exists
	if _, err := os.Stat(logPath); os.IsNotExist(err) {
		return false, 0, nil // No log file yet, wallet might be starting
	}

	// Read the log file
	data, err := ioutil.ReadFile(logPath)
	if err != nil {
		return false, 0, fmt.Errorf("failed to read wallet log: %w", err)
	}

	// Parse logs for rescan messages with timestamps: "2025-10-05 15:23:16.672 [INF] WLLT: Rescanning block range [414000, 415999]..."
	// Regex to extract timestamp and block range
	rescanRegex := regexp.MustCompile(`^(\d{4}-\d{2}-\d{2} \d{2}:\d{2}:\d{2}\.\d{3}).*Rescanning block range \[(\d+), (\d+)\]`)

	lines := strings.Split(string(data), "\n")

	// Check the most recent rescan message (from the end, last ~100 lines for performance)
	startIndex := len(lines) - 100
	if startIndex < 0 {
		startIndex = 0
	}

	now := time.Now()
	maxAge := 2 * time.Minute // Consider rescan messages older than 2 minutes as stale

	for i := len(lines) - 1; i >= startIndex; i-- {
		line := lines[i]
		if matches := rescanRegex.FindStringSubmatch(line); matches != nil {
			// Parse the timestamp from the log line
			timestampStr := matches[1]
			logTime, err := time.Parse("2006-01-02 15:04:05.000", timestampStr)
			if err != nil {
				log.Printf("Warning: Failed to parse log timestamp '%s': %v", timestampStr, err)
				continue
			}

			// Check if the log entry is recent enough to be considered an active rescan
			age := now.Sub(logTime)
			if age > maxAge {
				log.Printf("Rescan message found but is stale (age: %v, max: %v) - considering rescan as inactive", age, maxAge)
				return false, 0, nil
			}

			// Extract the end block of the range being scanned
			endBlock, err := strconv.ParseInt(matches[3], 10, 64)
			if err == nil {
				log.Printf("Active rescan detected: block %d (log age: %v)", endBlock, age)
				return true, endBlock, nil
			}
		}
	}

	return false, 0, nil
}

func FetchWalletStatus() (*types.WalletStatus, error) {
	// Use a longer timeout for wallet status to handle rescan scenarios
	// During rescan, RPC calls can be slow but should still respond
	ctx, cancel := context.WithTimeout(context.Background(), 15*time.Second)
	defer cancel()

	// Get wallet info using getinfo
	walletInfo, err := rpc.WalletClient.GetInfo(ctx)
	if err != nil {
		// If RPC fails, check if it's because of an active rescan
		isRescanning, logScanHeight, logErr := ParseWalletLogsForRescan()
		if logErr == nil && isRescanning {
			// Wallet is busy rescanning, use log data
			if rpc.DcrdClient != nil {
				chainHeight, err := rpc.DcrdClient.GetBlockCount(ctx)
				if err == nil {
					syncProgress := (float64(logScanHeight) / float64(chainHeight)) * 100
					return &types.WalletStatus{
						Status:           "syncing",
						SyncProgress:     syncProgress,
						SyncHeight:       logScanHeight,
						BestBlockHash:    "",
						Version:          "unknown",
						Unlocked:         true,
						RescanInProgress: true,
						SyncMessage:      fmt.Sprintf("Rescanning... %d/%d blocks (%.1f%%)", logScanHeight, chainHeight, syncProgress),
					}, nil
				}
			}
		}

		return &types.WalletStatus{
			Status:      "no_wallet",
			SyncMessage: fmt.Sprintf("Wallet not available: %v", err),
		}, nil
	}

	// Determine wallet status
	status := "synced"
	syncProgress := 100.0
	syncMessage := "Fully synced"
	rescanInProgress := false
	var syncHeight int64 = 0
	bestBlockHash := ""

	// Get best block from wallet
	bestHash, bestHeight, err := rpc.WalletClient.GetBestBlock(ctx)
	if err == nil {
		syncHeight = bestHeight
		bestBlockHash = bestHash.String()

		// Get current block count from dcrd for comparison
		if rpc.DcrdClient != nil {
			chainHeight, err := rpc.DcrdClient.GetBlockCount(ctx)
			if err == nil {
				walletHeight := bestHeight

				// Calculate sync progress
				// During rescan, wallet height changes rapidly through already-synced blocks
				// Allow a buffer of 2 blocks to account for chain growth during sync
				blocksBehind := chainHeight - walletHeight
				if blocksBehind > 2 {
					status = "syncing"
					syncProgress = (float64(walletHeight) / float64(chainHeight)) * 100

					// Calculate delta for sync message
					walletSyncMutex.Lock()
					deltaHeight := walletHeight - prevWalletHeight
					prevWalletHeight = walletHeight
					walletSyncMutex.Unlock()

					if deltaHeight > 0 {
						// If scanning more than 100 blocks per check (30s), likely a rescan
						if deltaHeight > 100 {
							syncMessage = fmt.Sprintf("Rescanning... %d/%d blocks (%.1f%%)", walletHeight, chainHeight, syncProgress)
							rescanInProgress = true
						} else {
							syncMessage = fmt.Sprintf("Syncing... scanned %s blocks recently", utils.FormatNumber(deltaHeight))
						}
					} else {
						syncMessage = fmt.Sprintf("Syncing... %d/%d blocks", walletHeight, chainHeight)
					}
					rescanInProgress = true
				} else {
					// Even if heights match, check if we were recently rescanning
					walletSyncMutex.Lock()
					deltaHeight := walletHeight - prevWalletHeight
					if deltaHeight > 100 {
						// Just finished a fast rescan
						syncMessage = "Rescan completed, wallet fully synced"
					}
					prevWalletHeight = walletHeight
					walletSyncMutex.Unlock()
				}
			}
		}
	} else {
		status = "disconnected"
		syncMessage = "Wallet not connected to dcrd"
	}

	// Check Docker logs for active rescan (especially for manual rescans via rescan button or xpub import)
	// This is more reliable than RPC during intensive rescans
	isRescanning, logScanHeight, logErr := ParseWalletLogsForRescan()
	if logErr == nil && isRescanning {
		// Get chain height for progress calculation
		if rpc.DcrdClient != nil {
			chainHeight, err := rpc.DcrdClient.GetBlockCount(ctx)
			if err == nil {
				// Only override status if the wallet is actually behind
				// Allow a small buffer of 2 blocks to account for chain growth during sync
				blocksBehind := chainHeight - logScanHeight
				if blocksBehind > 2 {
					status = "syncing"
					rescanInProgress = true
					syncHeight = logScanHeight
					syncProgress = (float64(logScanHeight) / float64(chainHeight)) * 100
					syncMessage = fmt.Sprintf("Rescanning... %d/%d blocks (%.1f%%)", logScanHeight, chainHeight, syncProgress)
					log.Printf("Rescan detected from logs: block %d / %d (%.1f%%)", logScanHeight, chainHeight, syncProgress)
				} else {
					log.Printf("Rescan message found in logs but wallet is at chain tip (%d/%d, %d blocks behind) - rescan complete", logScanHeight, chainHeight, blocksBehind)
				}
			}
		}
	}

	return &types.WalletStatus{
		Status:           status,
		SyncProgress:     syncProgress,
		SyncHeight:       syncHeight,
		BestBlockHash:    bestBlockHash,
		Version:          fmt.Sprintf("%d.%d.%d", walletInfo.Version, 0, 0),
		Unlocked:         true, // Assume unlocked since we auto-unlock
		RescanInProgress: rescanInProgress,
		SyncMessage:      syncMessage,
	}, nil
}

func FetchWalletDashboardData() (*types.WalletDashboardData, error) {
	ctx := context.Background()
	return FetchWalletDashboardDataWithContext(ctx)
}

func FetchWalletDashboardDataWithContext(ctx context.Context) (*types.WalletDashboardData, error) {
	walletStatus, err := FetchWalletStatus()
	if err != nil {
		return nil, err
	}

	accountInfo := &types.AccountInfo{}
	accounts := []types.AccountInfo{}
	var stakingInfo *types.WalletStakingInfo

	// Fetch data with timeout protection - use channels to respect context
	type accountResult struct {
		data *types.AccountInfo
		err  error
	}
	type accountsResult struct {
		data []types.AccountInfo
		err  error
	}
	type stakingResult struct {
		data *types.WalletStakingInfo
		err  error
	}

	accountChan := make(chan accountResult, 1)
	accountsChan := make(chan accountsResult, 1)
	stakingChan := make(chan stakingResult, 1)

	go func() {
		info, err := FetchAccountInfoWithContext(ctx)
		accountChan <- accountResult{info, err}
	}()

	go func() {
		accts, err := FetchAllAccounts(ctx)
		accountsChan <- accountsResult{accts, err}
	}()

	go func() {
		staking, err := FetchWalletStakingInfo(ctx)
		stakingChan <- stakingResult{staking, err}
	}()

	select {
	case res := <-accountChan:
		if res.err != nil {
			log.Printf("Warning: Failed to fetch account info: %v", res.err)
		} else {
			accountInfo = res.data
		}
	case <-ctx.Done():
		log.Printf("Warning: Account info fetch cancelled: %v", ctx.Err())
	}

	select {
	case res := <-accountsChan:
		if res.err != nil {
			log.Printf("Warning: Failed to fetch accounts: %v", res.err)
		} else {
			accounts = res.data
		}
	case <-ctx.Done():
		log.Printf("Warning: Accounts fetch cancelled: %v", ctx.Err())
	}

	select {
	case res := <-stakingChan:
		if res.err != nil {
			log.Printf("Warning: Failed to fetch staking info: %v", res.err)
			// Staking info is optional - continue without it
		} else {
			stakingInfo = res.data
		}
	case <-ctx.Done():
		log.Printf("Warning: Staking info fetch cancelled: %v", ctx.Err())
	}

	return &types.WalletDashboardData{
		WalletStatus: *walletStatus,
		AccountInfo:  *accountInfo,
		Accounts:     accounts,
		StakingInfo:  stakingInfo,
		LastUpdate:   time.Now(),
	}, nil
}

func FetchAccountInfo() (*types.AccountInfo, error) {
	return FetchAccountInfoWithContext(context.Background())
}

func FetchAccountInfoWithContext(ctx context.Context) (*types.AccountInfo, error) {
	// Get balance using getbalance (no arguments for all accounts)
	result, err := rpc.WalletClient.RawRequest(ctx, "getbalance", []json.RawMessage{})
	if err != nil {
		log.Printf("Warning: Failed to get balance: %v", err)
		return &types.AccountInfo{
			AccountName:        "Total",
			TotalBalance:       0,
			SpendableBalance:   0,
			ImmatureBalance:    0,
			UnconfirmedBalance: 0,
			LockedByTickets:    0,
			AccountNumber:      0,
		}, nil
	}

	// Parse the full balance response structure
	// getbalance returns: {
	//   "balances":[{account info},...],
	//   "blockhash":"...",
	//   "totallockedbytickets": X,
	//   "totalspendable": Y,
	//   "cumulativetotal": Z
	// }
	type AccountBalance struct {
		AccountName             string  `json:"accountname"`
		ImmatureCoinbaseRewards float64 `json:"immaturecoinbaserewards"`
		ImmatureStakeGeneration float64 `json:"immaturestakegeneration"`
		LockedByTickets         float64 `json:"lockedbytickets"`
		Spendable               float64 `json:"spendable"`
		Total                   float64 `json:"total"`
		Unconfirmed             float64 `json:"unconfirmed"`
		VotingAuthority         float64 `json:"votingauthority"`
	}
	type BalanceResponse struct {
		Balances             []AccountBalance `json:"balances"`
		BlockHash            string           `json:"blockhash"`
		TotalLockedByTickets float64          `json:"totallockedbytickets"`
		TotalSpendable       float64          `json:"totalspendable"`
		CumulativeTotal      float64          `json:"cumulativetotal"`
	}

	var balanceResp BalanceResponse
	if err := json.Unmarshal(result, &balanceResp); err != nil {
		log.Printf("Warning: Failed to unmarshal balance response: %v", err)
		return &types.AccountInfo{
			AccountName:        "Total",
			TotalBalance:       0,
			SpendableBalance:   0,
			ImmatureBalance:    0,
			UnconfirmedBalance: 0,
			LockedByTickets:    0,
			AccountNumber:      0,
		}, nil
	}

	// Sum immature and unconfirmed balances across all accounts
	immature := 0.0
	unconfirmed := 0.0
	lockedByTickets := 0.0
	votingAuthority := 0.0

	for _, acct := range balanceResp.Balances {
		immature += acct.ImmatureCoinbaseRewards + acct.ImmatureStakeGeneration
		unconfirmed += acct.Unconfirmed
		lockedByTickets += acct.LockedByTickets
		votingAuthority += acct.VotingAuthority
	}

	// Return wallet-wide totals with granular breakdown
	return &types.AccountInfo{
		AccountName:        "Total",
		TotalBalance:       balanceResp.CumulativeTotal,
		SpendableBalance:   balanceResp.TotalSpendable,
		ImmatureBalance:    immature,
		UnconfirmedBalance: unconfirmed,
		LockedByTickets:    balanceResp.TotalLockedByTickets,
		VotingAuthority:    votingAuthority,
		AccountNumber:      0,
		// Wallet-wide totals
		CumulativeTotal:      balanceResp.CumulativeTotal,
		TotalSpendable:       balanceResp.TotalSpendable,
		TotalLockedByTickets: balanceResp.TotalLockedByTickets,
	}, nil
}

func FetchAllAccounts(ctx context.Context) ([]types.AccountInfo, error) {
	// Get all accounts and their balances using getbalance RPC
	result, err := rpc.WalletClient.RawRequest(ctx, "getbalance", []json.RawMessage{})
	if err != nil {
		log.Printf("Warning: Failed to get accounts: %v", err)
		return []types.AccountInfo{}, nil
	}

	// Parse the balance response structure
	type AccountBalance struct {
		AccountName             string  `json:"accountname"`
		ImmatureCoinbaseRewards float64 `json:"immaturecoinbaserewards"`
		ImmatureStakeGeneration float64 `json:"immaturestakegeneration"`
		LockedByTickets         float64 `json:"lockedbytickets"`
		Spendable               float64 `json:"spendable"`
		Total                   float64 `json:"total"`
		Unconfirmed             float64 `json:"unconfirmed"`
		VotingAuthority         float64 `json:"votingauthority"`
		AccountNumber           uint32  `json:"accountnumber"`
	}
	type BalanceResponse struct {
		Balances  []AccountBalance `json:"balances"`
		BlockHash string           `json:"blockhash"`
	}

	var balanceResp BalanceResponse
	if err := json.Unmarshal(result, &balanceResp); err != nil {
		log.Printf("Warning: Failed to unmarshal accounts: %v", err)
		return []types.AccountInfo{}, nil
	}

	accounts := make([]types.AccountInfo, 0, len(balanceResp.Balances))
	for _, acct := range balanceResp.Balances {
		accounts = append(accounts, types.AccountInfo{
			AccountName:             acct.AccountName,
			TotalBalance:            acct.Total,
			SpendableBalance:        acct.Spendable,
			ImmatureBalance:         acct.ImmatureCoinbaseRewards + acct.ImmatureStakeGeneration,
			UnconfirmedBalance:      acct.Unconfirmed,
			LockedByTickets:         acct.LockedByTickets,
			VotingAuthority:         acct.VotingAuthority,
			ImmatureCoinbaseRewards: acct.ImmatureCoinbaseRewards,
			ImmatureStakeGeneration: acct.ImmatureStakeGeneration,
			AccountNumber:           0, // Account numbers not reliably available from RPC
		})
	}

	return accounts, nil
}

// Old FetchTransactions functions removed - replaced by ListTransactions

func FetchAddresses() ([]types.Address, error) {
	return FetchAddressesWithContext(context.Background())
}

func FetchAddressesWithContext(ctx context.Context) ([]types.Address, error) {
	// List addresses via raw RPC - only return addresses with funds (not empty)
	// This prevents returning 40k+ empty addresses
	result, err := rpc.WalletClient.RawRequest(ctx, "listreceivedbyaddress", []json.RawMessage{
		json.RawMessage(`0`),     // minconf
		json.RawMessage(`false`), // include empty = false (only show addresses with funds)
	})
	if err != nil {
		log.Printf("Warning: Failed to list addresses: %v", err)
		return []types.Address{}, nil
	}

	// Parse the result
	var rawAddrList []map[string]interface{}
	if err := json.Unmarshal(result, &rawAddrList); err != nil {
		log.Printf("Warning: Failed to unmarshal addresses: %v", err)
		return []types.Address{}, nil
	}

	// Limit to 100 addresses max to prevent huge payloads
	maxAddresses := 100
	if len(rawAddrList) > maxAddresses {
		log.Printf("Warning: Wallet has %d addresses with funds, limiting to %d", len(rawAddrList), maxAddresses)
		rawAddrList = rawAddrList[:maxAddresses]
	}

	addresses := make([]types.Address, 0, len(rawAddrList))
	for _, addr := range rawAddrList {
		address, _ := addr["address"].(string)
		account, _ := addr["account"].(string)
		amount, _ := addr["amount"].(float64)

		addresses = append(addresses, types.Address{
			Address: address,
			Account: account,
			Used:    amount > 0, // Has received funds
			Path:    "",         // Would need to query separately
		})
	}

	log.Printf("Returning %d addresses with funds", len(addresses))
	return addresses, nil
}

func FetchWalletStakingInfo(ctx context.Context) (*types.WalletStakingInfo, error) {
	stakingInfo := &types.WalletStakingInfo{}

	// Fetch getstakeinfo
	stakeInfoResult, err := rpc.WalletClient.RawRequest(ctx, "getstakeinfo", []json.RawMessage{})
	if err != nil {
		log.Printf("Warning: Failed to get stake info: %v", err)
		return nil, err
	}

	type StakeInfoResponse struct {
		BlockHeight    int64   `json:"blockheight"`
		Difficulty     float64 `json:"difficulty"`
		TotalSubsidy   float64 `json:"totalsubsidy"`
		OwnMempoolTix  int32   `json:"ownmempooltix"`
		Immature       int32   `json:"immature"`
		Unspent        int32   `json:"unspent"`
		Voted          int32   `json:"voted"`
		Revoked        int32   `json:"revoked"`
		UnspentExpired int32   `json:"unspentexpired"`
		PoolSize       int32   `json:"poolsize"`
		AllMempoolTix  int32   `json:"allmempooltix"`
	}

	var stakeInfo StakeInfoResponse
	if err := json.Unmarshal(stakeInfoResult, &stakeInfo); err != nil {
		log.Printf("Warning: Failed to unmarshal stake info: %v", err)
		return nil, err
	}

	stakingInfo.BlockHeight = stakeInfo.BlockHeight
	stakingInfo.Difficulty = stakeInfo.Difficulty
	stakingInfo.TotalSubsidy = stakeInfo.TotalSubsidy
	stakingInfo.OwnMempoolTix = stakeInfo.OwnMempoolTix
	stakingInfo.Immature = stakeInfo.Immature
	stakingInfo.Unspent = stakeInfo.Unspent
	stakingInfo.Voted = stakeInfo.Voted
	stakingInfo.Revoked = stakeInfo.Revoked
	stakingInfo.UnspentExpired = stakeInfo.UnspentExpired
	stakingInfo.PoolSize = stakeInfo.PoolSize
	stakingInfo.AllMempoolTix = stakeInfo.AllMempoolTix

	// Fetch estimatestakediff
	estimateResult, err := rpc.WalletClient.RawRequest(ctx, "estimatestakediff", []json.RawMessage{})
	if err != nil {
		log.Printf("Warning: Failed to estimate stake diff: %v", err)
	} else {
		type EstimateResponse struct {
			Min      float64 `json:"min"`
			Max      float64 `json:"max"`
			Expected float64 `json:"expected"`
		}
		var estimate EstimateResponse
		if err := json.Unmarshal(estimateResult, &estimate); err == nil {
			stakingInfo.EstimatedMin = estimate.Min
			stakingInfo.EstimatedMax = estimate.Max
			stakingInfo.EstimatedExpected = estimate.Expected
		}
	}

	// Fetch getstakedifficulty
	difficultyResult, err := rpc.WalletClient.RawRequest(ctx, "getstakedifficulty", []json.RawMessage{})
	if err != nil {
		log.Printf("Warning: Failed to get stake difficulty: %v", err)
	} else {
		type DifficultyResponse struct {
			Current float64 `json:"current"`
			Next    float64 `json:"next"`
		}
		var difficulty DifficultyResponse
		if err := json.Unmarshal(difficultyResult, &difficulty); err == nil {
			stakingInfo.CurrentDifficulty = difficulty.Current
			stakingInfo.NextDifficulty = difficulty.Next
		}
	}

	return stakingInfo, nil
}

// ListTransactions fetches recent wallet transactions
func ListTransactions(ctx context.Context, count, from int) (*types.TransactionListResponse, error) {
	// Default parameters
	if count <= 0 {
		count = 50 // Default to 50 transactions
	}
	if count > 200 {
		count = 200 // Cap at 200 for performance
	}

	// Call listtransactions RPC with parameters
	result, err := rpc.WalletClient.RawRequest(ctx, "listtransactions", []json.RawMessage{
		json.RawMessage(`"*"`),                    // account (all accounts)
		json.RawMessage(fmt.Sprintf("%d", count)), // count
		json.RawMessage(fmt.Sprintf("%d", from)),  // from (skip)
		json.RawMessage("false"),                  // includewatchonly
	})
	if err != nil {
		return nil, fmt.Errorf("failed to list transactions: %w", err)
	}

	// Parse the response
	var rpcTransactions []struct {
		Account         string   `json:"account"`
		Address         string   `json:"address"`
		Amount          float64  `json:"amount"`
		BlockHash       string   `json:"blockhash"`
		BlockTime       int64    `json:"blocktime"`
		Category        string   `json:"category"`
		Confirmations   int64    `json:"confirmations"`
		Fee             float64  `json:"fee"`
		Generated       bool     `json:"generated"`
		Time            int64    `json:"time"`
		TimeReceived    int64    `json:"timereceived"`
		TxID            string   `json:"txid"`
		TxType          string   `json:"txtype"`
		Vout            uint32   `json:"vout"`
		WalletConflicts []string `json:"walletconflicts"`
	}

	if err := json.Unmarshal(result, &rpcTransactions); err != nil {
		return nil, fmt.Errorf("failed to unmarshal transactions: %w", err)
	}

	// Convert to our transaction type
	transactions := make([]types.Transaction, 0, len(rpcTransactions))
	seenTxIDs := make(map[string]bool) // Track unique transactions to avoid duplicates

	for _, rpcTx := range rpcTransactions {
		// Create a unique key for this transaction (txid + vout)
		uniqueKey := fmt.Sprintf("%s-%d", rpcTx.TxID, rpcTx.Vout)

		// Skip if we've already processed this exact output
		if seenTxIDs[uniqueKey] {
			continue
		}
		seenTxIDs[uniqueKey] = true

		tx := types.Transaction{
			TxID:          rpcTx.TxID,
			Amount:        rpcTx.Amount,
			Fee:           rpcTx.Fee,
			Confirmations: rpcTx.Confirmations,
			BlockHash:     rpcTx.BlockHash,
			BlockTime:     rpcTx.BlockTime,
			Time:          time.Unix(rpcTx.Time, 0),
			Category:      rpcTx.Category,
			TxType:        rpcTx.TxType,
			Address:       rpcTx.Address,
			Account:       rpcTx.Account,
			Vout:          rpcTx.Vout,
			Generated:     rpcTx.Generated,
		}

		transactions = append(transactions, tx)
	}

	return &types.TransactionListResponse{
		Transactions: transactions,
		Total:        len(transactions),
	}, nil
}

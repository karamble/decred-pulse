// Copyright (c) 2015-2025 The Decred developers
// Use of this source code is governed by an ISC
// license that can be found in the LICENSE file.

package handlers

import (
	"context"
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"strings"
	"time"

	"decred-pulse-backend/rpc"
	"decred-pulse-backend/services"
	"decred-pulse-backend/types"
)

// GetWalletStatusHandler handles requests for wallet status
func GetWalletStatusHandler(w http.ResponseWriter, r *http.Request) {
	if rpc.WalletClient == nil {
		http.Error(w, "Wallet RPC client not initialized", http.StatusServiceUnavailable)
		return
	}

	status, err := services.FetchWalletStatus()
	if err != nil {
		log.Printf("Error fetching wallet status: %v", err)
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(status)
}

// GetWalletDashboardHandler handles requests for complete wallet dashboard data
func GetWalletDashboardHandler(w http.ResponseWriter, r *http.Request) {
	if rpc.WalletClient == nil {
		http.Error(w, "Wallet RPC client not initialized", http.StatusServiceUnavailable)
		return
	}

	// Create a context with timeout to prevent hanging on slow RPC calls
	// Use 20 seconds to accommodate wallet rescans which can slow down RPC responses
	ctx, cancel := context.WithTimeout(r.Context(), 20*time.Second)
	defer cancel()

	// Use a channel to handle the async fetch with timeout
	type result struct {
		data *types.WalletDashboardData
		err  error
	}
	resultChan := make(chan result, 1)

	go func() {
		data, err := services.FetchWalletDashboardDataWithContext(ctx)
		resultChan <- result{data, err}
	}()

	select {
	case res := <-resultChan:
		if res.err != nil {
			log.Printf("Error fetching wallet dashboard data: %v", res.err)
			// Return partial data if available
			if res.data != nil {
				w.Header().Set("Content-Type", "application/json")
				json.NewEncoder(w).Encode(res.data)
				return
			}
			http.Error(w, res.err.Error(), http.StatusInternalServerError)
			return
		}
		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode(res.data)
	case <-ctx.Done():
		log.Printf("Wallet dashboard request timed out")
		http.Error(w, "Wallet dashboard request timed out - wallet may be rescanning", http.StatusRequestTimeout)
	}
}

// ImportXpubHandler handles xpub import requests
func ImportXpubHandler(w http.ResponseWriter, r *http.Request) {
	if rpc.WalletClient == nil {
		http.Error(w, "Wallet RPC client not initialized", http.StatusServiceUnavailable)
		return
	}

	var req types.ImportXpubRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		http.Error(w, "Invalid request body", http.StatusBadRequest)
		return
	}

	// Validate xpub format (Decred mainnet xpubs start with "dpub")
	if !strings.HasPrefix(req.Xpub, "dpub") && !strings.HasPrefix(req.Xpub, "tpub") {
		response := types.ImportXpubResponse{
			Success: false,
			Message: "Invalid xpub format. Decred mainnet xpubs must start with 'dpub'",
		}
		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode(response)
		return
	}

	// Set default account name if not provided
	accountName := req.AccountName
	if accountName == "" {
		accountName = "imported"
	}

	// Import the xpub asynchronously
	// Note: importxpub automatically triggers a blockchain rescan which can take several minutes
	// We run it in a goroutine and return immediately so the frontend doesn't timeout
	log.Printf("Starting xpub import for account: %s", accountName)

	// Start the import in a goroutine
	go func() {
		// Build params: importxpub "account" "xpub" rescanfrom
		// Start from block 0 for full rescan (or higher if user specifies)
		rescanFrom := 0
		if req.Rescan {
			rescanFrom = 0 // Full rescan from genesis if explicitly requested
		}

		params := []json.RawMessage{
			json.RawMessage(fmt.Sprintf(`"%s"`, accountName)),
			json.RawMessage(fmt.Sprintf(`"%s"`, req.Xpub)),
			json.RawMessage(fmt.Sprintf(`%d`, rescanFrom)), // CRITICAL: rescanfrom parameter
		}

		log.Printf("Importing xpub for account '%s' with rescan from block %d", accountName, rescanFrom)
		result, err := rpc.WalletClient.RawRequest(context.Background(), "importxpub", params)
		if err != nil {
			log.Printf("Failed to import xpub: %v", err)
			return
		}
		log.Printf("Xpub import completed: %v", string(result))
		log.Printf("Rescan initiated automatically by importxpub")
	}()

	// Return immediately - the frontend will poll wallet status to track rescan progress
	response := types.ImportXpubResponse{
		Success: true,
		Message: fmt.Sprintf("Xpub import started for account '%s'. Full blockchain rescan will discover your transactions. This typically takes 5-10 minutes.", accountName),
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(response)
}

// RescanWalletHandler handles wallet rescan requests
func RescanWalletHandler(w http.ResponseWriter, r *http.Request) {
	if rpc.WalletClient == nil {
		http.Error(w, "Wallet RPC client not initialized", http.StatusServiceUnavailable)
		return
	}

	var req types.RescanRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		// Default to full rescan from genesis
		req.BeginHeight = 0
	}

	// Start rescan in a goroutine - it's a long-running operation
	// Frontend will poll wallet status to track progress via log parsing
	log.Printf("Starting wallet rescan from block %d", req.BeginHeight)

	go func() {
		ctx := context.Background()
		rescanParams := []json.RawMessage{json.RawMessage(fmt.Sprintf("%d", req.BeginHeight))}
		_, err := rpc.WalletClient.RawRequest(ctx, "rescanwallet", rescanParams)
		if err != nil {
			log.Printf("Rescan completed with error: %v", err)
		} else {
			log.Printf("Rescan completed successfully")
		}
	}()

	// Return immediately so frontend can start polling for progress
	response := types.RescanResponse{
		Success: true,
		Message: fmt.Sprintf("Wallet rescan started from block %d. This may take 30+ minutes.", req.BeginHeight),
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(response)
}

// GetSyncProgressHandler handles requests for wallet sync progress from log files
func GetSyncProgressHandler(w http.ResponseWriter, r *http.Request) {
	// Get sync progress from log file parsing
	isRescanning, scanHeight, err := services.ParseWalletLogsForRescan()
	if err != nil {
		log.Printf("Error parsing wallet logs: %v", err)
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	// Get chain height from dcrd for progress calculation
	var progress float64 = 100.0
	var chainHeight int64 = 0
	var message string = "No active rescan"

	if isRescanning && rpc.DcrdClient != nil {
		ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
		defer cancel()

		height, err := rpc.DcrdClient.GetBlockCount(ctx)
		if err == nil {
			chainHeight = height
			progress = (float64(scanHeight) / float64(chainHeight)) * 100
			message = fmt.Sprintf("Rescanning... %d/%d blocks", scanHeight, chainHeight)
		}
	}

	response := types.SyncProgressResponse{
		IsRescanning: isRescanning,
		ScanHeight:   scanHeight,
		ChainHeight:  chainHeight,
		Progress:     progress,
		Message:      message,
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(response)
}

// ListTransactionsHandler handles requests for wallet transaction history
func ListTransactionsHandler(w http.ResponseWriter, r *http.Request) {
	if rpc.WalletClient == nil {
		http.Error(w, "Wallet RPC client not initialized", http.StatusServiceUnavailable)
		return
	}

	// Parse query parameters
	query := r.URL.Query()
	count := 50 // default
	from := 0   // default

	if c := query.Get("count"); c != "" {
		if parsed, err := fmt.Sscanf(c, "%d", &count); err == nil && parsed == 1 {
			// count parsed successfully
		}
	}
	if f := query.Get("from"); f != "" {
		if parsed, err := fmt.Sscanf(f, "%d", &from); err == nil && parsed == 1 {
			// from parsed successfully
		}
	}

	// Create context with timeout
	ctx, cancel := context.WithTimeout(r.Context(), 10*time.Second)
	defer cancel()

	// Fetch transactions
	transactions, err := services.ListTransactions(ctx, count, from)
	if err != nil {
		log.Printf("Error listing transactions: %v", err)
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(transactions)
}

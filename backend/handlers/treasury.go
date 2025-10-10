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
	"time"

	"decred-pulse-backend/services"
)

// GetTreasuryInfoHandler returns current treasury status
func GetTreasuryInfoHandler(w http.ResponseWriter, r *http.Request) {
	ctx, cancel := context.WithTimeout(context.Background(), 15*time.Second)
	defer cancel()

	info, err := services.FetchTreasuryInfo(ctx)
	if err != nil {
		log.Printf("Error fetching treasury info: %v", err)
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(info)
}

// TriggerTSpendScanHandler triggers a historical blockchain scan for TSpends
func TriggerTSpendScanHandler(w http.ResponseWriter, r *http.Request) {
	// Parse request body to get startHeight (optional)
	var req struct {
		StartHeight int64 `json:"startHeight"`
	}

	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		// If no body or invalid, default to treasury activation height
		req.StartHeight = 552448
	}

	// If startHeight is 0 or invalid, use treasury activation height
	if req.StartHeight < 552448 {
		req.StartHeight = 552448
	}

	err := services.TriggerHistoricalScan(req.StartHeight)
	if err != nil {
		log.Printf("Error triggering TSpend scan: %v", err)
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(map[string]interface{}{
		"success": true,
		"message": fmt.Sprintf("Historical TSpend scan started from block %d", req.StartHeight),
	})
}

// GetTSpendScanProgressHandler returns the current scan progress
func GetTSpendScanProgressHandler(w http.ResponseWriter, r *http.Request) {
	progress, err := services.GetScanProgress()
	if err != nil {
		log.Printf("Error getting scan progress: %v", err)
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(progress)
}

// GetTSpendScanResultsHandler returns the results from the last completed scan
func GetTSpendScanResultsHandler(w http.ResponseWriter, r *http.Request) {
	results := services.GetScanResults()

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(results)
}

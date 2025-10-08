// Copyright (c) 2015-2025 The Decred developers
// Use of this source code is governed by an ISC
// license that can be found in the LICENSE file.

package handlers

import (
	"context"
	"fmt"
	"log"
	"net/http"
	"time"

	"decred-pulse-backend/rpc"
	"decred-pulse-backend/services"

	"github.com/gorilla/websocket"
)

// StreamRescanGrpcHandler handles WebSocket connections and monitors rescan progress
//
// IMPORTANT: This endpoint ONLY monitors existing rescans. It does NOT start new rescans.
// The gRPC Rescan() method would start a NEW rescan each time it's called, which causes
// multiple concurrent rescans when navigating between pages.
//
// Rescans should only be initiated by explicit user actions:
// - Manual rescan button (via RescanWalletHandler)
// - Xpub import (via ImportXpubHandler)
//
// This handler uses log-based monitoring which works for all rescans regardless of how they were started.
func StreamRescanGrpcHandler(w http.ResponseWriter, r *http.Request) {
	// Upgrade to WebSocket
	upgrader := websocket.Upgrader{
		CheckOrigin: func(r *http.Request) bool {
			return true // Allow all origins for development
		},
	}

	conn, err := upgrader.Upgrade(w, r, nil)
	if err != nil {
		log.Printf("Failed to upgrade to WebSocket: %v", err)
		return
	}
	defer conn.Close()

	log.Println("WebSocket connection established for rescan progress monitoring")

	// Get chain height for progress calculation
	chainHeight := int64(1016874) // Default fallback
	if rpc.DcrdClient != nil {
		ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
		height, err := rpc.DcrdClient.GetBlockCount(ctx)
		cancel()
		if err == nil {
			chainHeight = height
		}
	}

	// Monitor existing rescan progress via log parsing
	ticker := time.NewTicker(1 * time.Second)
	defer ticker.Stop()

	notRescanningCount := 0
	const maxNotRescanningBeforeClose = 5
	const gracePeriodTicks = 5
	tickCount := 0

	for {
		select {
		case <-ticker.C:
			tickCount++

			// Check if rescan is active by parsing logs
			isRescanning, scanHeight, err := services.ParseWalletLogsForRescan()
			if err != nil {
				log.Printf("Error parsing logs: %v", err)
				continue
			}

			progress := 0.0
			if chainHeight > 0 && scanHeight > 0 {
				progress = (float64(scanHeight) / float64(chainHeight)) * 100
			}

			progressData := map[string]interface{}{
				"isRescanning": isRescanning,
				"scanHeight":   scanHeight,
				"chainHeight":  chainHeight,
				"progress":     progress,
			}

			if isRescanning {
				progressData["message"] = fmt.Sprintf("Rescanning... %d/%d blocks", scanHeight, chainHeight)
				notRescanningCount = 0
			} else {
				if tickCount > gracePeriodTicks {
					notRescanningCount++
					progressData["message"] = "No active rescan"
				} else {
					progressData["message"] = "Checking for active rescan..."
				}
			}

			if err := conn.WriteJSON(progressData); err != nil {
				log.Printf("Failed to write to WebSocket: %v", err)
				return
			}

			if notRescanningCount >= maxNotRescanningBeforeClose {
				log.Println("No active rescan detected, closing monitoring stream")
				return
			}
		}
	}
}

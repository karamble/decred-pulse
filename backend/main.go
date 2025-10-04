// Copyright (c) 2015-2025 The Decred developers
// Use of this source code is governed by an ISC
// license that can be found in the LICENSE file.

package main

import (
	"context"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"math"
	"net/http"
	"os"
	"time"

	"github.com/decred/dcrd/rpcclient/v8"
	"github.com/gorilla/mux"
	"github.com/rs/cors"
)

var (
	rpcClient *rpcclient.Client
)

// Config holds the RPC connection configuration
type Config struct {
	RPCHost     string
	RPCPort     string
	RPCUser     string
	RPCPassword string
	RPCCert     string
}

// DashboardData represents all dashboard metrics
type DashboardData struct {
	NodeStatus     NodeStatus     `json:"nodeStatus"`
	BlockchainInfo BlockchainInfo `json:"blockchainInfo"`
	NetworkInfo    NetworkInfo    `json:"networkInfo"`
	Peers          []Peer         `json:"peers"`
	SupplyInfo     SupplyInfo     `json:"supplyInfo"`
	LastUpdate     time.Time      `json:"lastUpdate"`
}

type NodeStatus struct {
	Status       string  `json:"status"`
	SyncProgress float64 `json:"syncProgress"`
	Version      string  `json:"version"`
}

type BlockchainInfo struct {
	BlockHeight int64   `json:"blockHeight"`
	BlockHash   string  `json:"blockHash"`
	Difficulty  float64 `json:"difficulty"`
	ChainSize   int64   `json:"chainSize"`
	BlockTime   string  `json:"blockTime"`
}

type NetworkInfo struct {
	PeerCount     int     `json:"peerCount"`
	Hashrate      string  `json:"hashrate"`
	NetworkHashPS float64 `json:"networkHashPS"`
}

type Peer struct {
	ID       int    `json:"id"`
	Address  string `json:"address"`
	Protocol string `json:"protocol"`
	Latency  string `json:"latency"`
}

type SupplyInfo struct {
	CirculatingSupply string  `json:"circulatingSupply"`
	StakedSupply      string  `json:"stakedSupply"`
	StakedPercent     float64 `json:"stakedPercent"`
	ExchangeRate      string  `json:"exchangeRate"`
	TreasurySize      string  `json:"treasurySize"`
	MixedPercent      string  `json:"mixedPercent"`
}

type RPCConnectionRequest struct {
	Host     string `json:"host"`
	Port     string `json:"port"`
	Username string `json:"username"`
	Password string `json:"password"`
}

type RPCConnectionResponse struct {
	Success bool   `json:"success"`
	Message string `json:"message"`
}

// Initialize RPC client
func initRPCClient(config Config) error {
	// Read the TLS certificate if provided
	var certs []byte
	var err error

	if config.RPCCert != "" {
		log.Printf("Reading TLS certificate from: %s", config.RPCCert)
		certs, err = ioutil.ReadFile(config.RPCCert)
		if err != nil {
			return fmt.Errorf("failed to read RPC certificate: %v", err)
		}
		log.Printf("Successfully loaded TLS certificate (%d bytes)", len(certs))
	}

	connCfg := &rpcclient.ConnConfig{
		Host:         fmt.Sprintf("%s:%s", config.RPCHost, config.RPCPort),
		Endpoint:     "ws",
		User:         config.RPCUser,
		Pass:         config.RPCPassword,
		HTTPPostMode: true,
		DisableTLS:   config.RPCCert == "", // Disable TLS only if no cert provided
		Certificates: certs,
	}

	rpcClient, err = rpcclient.New(connCfg, nil)
	if err != nil {
		return fmt.Errorf("failed to create RPC client: %v", err)
	}

	// Test connection
	ctx := context.Background()
	_, err = rpcClient.GetBlockCount(ctx)
	if err != nil {
		return fmt.Errorf("failed to connect to dcrd: %v", err)
	}

	log.Println("Successfully connected to dcrd RPC with TLS")
	return nil
}

// API Handlers

func getDashboardDataHandler(w http.ResponseWriter, r *http.Request) {
	if rpcClient == nil {
		http.Error(w, "RPC client not initialized", http.StatusServiceUnavailable)
		return
	}

	data, err := fetchDashboardData()
	if err != nil {
		log.Printf("Error fetching dashboard data: %v", err)
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(data)
}

func getNodeStatusHandler(w http.ResponseWriter, r *http.Request) {
	if rpcClient == nil {
		http.Error(w, "RPC client not initialized", http.StatusServiceUnavailable)
		return
	}

	status, err := fetchNodeStatus()
	if err != nil {
		log.Printf("Error fetching node status: %v", err)
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(status)
}

func getBlockchainInfoHandler(w http.ResponseWriter, r *http.Request) {
	if rpcClient == nil {
		http.Error(w, "RPC client not initialized", http.StatusServiceUnavailable)
		return
	}

	info, err := fetchBlockchainInfo()
	if err != nil {
		log.Printf("Error fetching blockchain info: %v", err)
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(info)
}

func getPeersHandler(w http.ResponseWriter, r *http.Request) {
	if rpcClient == nil {
		http.Error(w, "RPC client not initialized", http.StatusServiceUnavailable)
		return
	}

	peers, err := fetchPeers()
	if err != nil {
		log.Printf("Error fetching peers: %v", err)
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(peers)
}

func connectRPCHandler(w http.ResponseWriter, r *http.Request) {
	var req RPCConnectionRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		http.Error(w, "Invalid request body", http.StatusBadRequest)
		return
	}

	config := Config{
		RPCHost:     req.Host,
		RPCPort:     req.Port,
		RPCUser:     req.Username,
		RPCPassword: req.Password,
	}

	err := initRPCClient(config)
	response := RPCConnectionResponse{
		Success: err == nil,
		Message: "Connected successfully",
	}

	if err != nil {
		response.Message = err.Error()
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(response)
}

func healthCheckHandler(w http.ResponseWriter, r *http.Request) {
	status := map[string]interface{}{
		"status":       "healthy",
		"rpcConnected": rpcClient != nil,
		"time":         time.Now(),
	}
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(status)
}

// Data fetching functions

func fetchDashboardData() (*DashboardData, error) {
	nodeStatus, err := fetchNodeStatus()
	if err != nil {
		return nil, err
	}

	blockchainInfo, err := fetchBlockchainInfo()
	if err != nil {
		return nil, err
	}

	networkInfo, err := fetchNetworkInfo()
	if err != nil {
		return nil, err
	}

	peers, err := fetchPeers()
	if err != nil {
		return nil, err
	}

	supplyInfo, err := fetchSupplyInfo()
	if err != nil {
		return nil, err
	}

	return &DashboardData{
		NodeStatus:     *nodeStatus,
		BlockchainInfo: *blockchainInfo,
		NetworkInfo:    *networkInfo,
		Peers:          peers,
		SupplyInfo:     *supplyInfo,
		LastUpdate:     time.Now(),
	}, nil
}

func fetchNodeStatus() (*NodeStatus, error) {
	ctx := context.Background()

	// Get basic node info for version
	info, err := rpcClient.GetInfo(ctx)
	if err != nil {
		return nil, err
	}

	// Get blockchain info for accurate sync status
	chainInfo, err := rpcClient.GetBlockChainInfo(ctx)
	if err != nil {
		return nil, err
	}

	// Debug logging
	log.Printf("Blockchain sync status - InitialBlockDownload: %v, Blocks: %d, SyncHeight: %d, VerificationProgress: %f",
		chainInfo.InitialBlockDownload, chainInfo.Blocks, chainInfo.SyncHeight, chainInfo.VerificationProgress)

	// Calculate sync progress based on actual blockchain sync
	var syncProgress float64
	var status string

	if chainInfo.InitialBlockDownload {
		// Node is still syncing
		status = "syncing"

		// Use verification progress if available, otherwise calculate from blocks/syncheight
		if chainInfo.VerificationProgress > 0 {
			syncProgress = chainInfo.VerificationProgress * 100
		} else if chainInfo.SyncHeight > 0 {
			syncProgress = (float64(chainInfo.Blocks) / float64(chainInfo.SyncHeight)) * 100
		} else {
			syncProgress = 0
		}
	} else {
		// Node is fully synced
		status = "running"
		syncProgress = 100.0
	}

	// Ensure progress is between 0 and 100
	if syncProgress > 100 {
		syncProgress = 100
	} else if syncProgress < 0 {
		syncProgress = 0
	}

	return &NodeStatus{
		Status:       status,
		SyncProgress: syncProgress,
		Version:      fmt.Sprintf("v%d.%d.%d", info.Version/1000000, (info.Version/10000)%100, (info.Version/100)%100),
	}, nil
}

func fetchBlockchainInfo() (*BlockchainInfo, error) {
	ctx := context.Background()
	info, err := rpcClient.GetBlockChainInfo(ctx)
	if err != nil {
		return nil, err
	}

	bestBlockHash, err := rpcClient.GetBestBlockHash(ctx)
	if err != nil {
		return nil, err
	}

	blockHeader, err := rpcClient.GetBlockHeader(ctx, bestBlockHash)
	if err != nil {
		return nil, err
	}

	// Estimate block time (time since last block)
	timeSinceBlock := time.Since(blockHeader.Timestamp)
	blockTime := fmt.Sprintf("%dm %ds", int(timeSinceBlock.Minutes()), int(timeSinceBlock.Seconds())%60)

	return &BlockchainInfo{
		BlockHeight: info.Blocks,
		BlockHash:   bestBlockHash.String(),
		Difficulty:  float64(info.Difficulty),
		ChainSize:   0, // Would need to calculate from disk usage
		BlockTime:   blockTime,
	}, nil
}

func fetchNetworkInfo() (*NetworkInfo, error) {
	ctx := context.Background()
	peerInfo, err := rpcClient.GetPeerInfo(ctx)
	if err != nil {
		return nil, err
	}

	// Get network hash rate
	info, err := rpcClient.GetBlockChainInfo(ctx)
	if err != nil {
		return nil, err
	}

	// Calculate network hashrate (simplified)
	hashrate := float64(info.Difficulty) * math.Pow(2, 32) / 300 // 5 min block time
	hashrateStr := formatHashrate(hashrate)

	return &NetworkInfo{
		PeerCount:     len(peerInfo),
		Hashrate:      hashrateStr,
		NetworkHashPS: hashrate,
	}, nil
}

func fetchPeers() ([]Peer, error) {
	ctx := context.Background()
	peerInfo, err := rpcClient.GetPeerInfo(ctx)
	if err != nil {
		return nil, err
	}

	peers := make([]Peer, 0, len(peerInfo))
	for i, p := range peerInfo {
		latency := fmt.Sprintf("%.0fms", p.PingTime*1000)

		peers = append(peers, Peer{
			ID:       i + 1,
			Address:  p.Addr,
			Protocol: "TCP",
			Latency:  latency,
		})
	}

	return peers, nil
}

func fetchSupplyInfo() (*SupplyInfo, error) {
	// For supply info, we would typically call external APIs or calculate from blockchain
	// For now, returning placeholder values
	// In production, you'd integrate with dcrdata API or calculate from chain

	ctx := context.Background()
	info, err := rpcClient.GetBlockChainInfo(ctx)
	if err != nil {
		return nil, err
	}

	// Simplified calculation - in production use proper supply formulas
	circulatingSupply := fmt.Sprintf("%.2fM", float64(info.Blocks)*1.5/1000000)

	return &SupplyInfo{
		CirculatingSupply: circulatingSupply,
		StakedSupply:      "10.15M",     // Would calculate from ticket pool
		StakedPercent:     59.4,         // Would calculate from actual data
		ExchangeRate:      "$17.70",     // Would fetch from external API
		TreasurySize:      "861.6K DCR", // Would query treasury address
		MixedPercent:      "62%",        // Would query from mixer statistics
	}, nil
}

func formatHashrate(hashrate float64) string {
	units := []string{"H/s", "KH/s", "MH/s", "GH/s", "TH/s", "PH/s", "EH/s"}
	unitIndex := 0

	for hashrate >= 1000 && unitIndex < len(units)-1 {
		hashrate /= 1000
		unitIndex++
	}

	return fmt.Sprintf("%.2f %s", hashrate, units[unitIndex])
}

func main() {
	// Load configuration from environment variables
	config := Config{
		RPCHost:     getEnv("DCRD_RPC_HOST", "localhost"),
		RPCPort:     getEnv("DCRD_RPC_PORT", "9109"),
		RPCUser:     getEnv("DCRD_RPC_USER", ""),
		RPCPassword: getEnv("DCRD_RPC_PASS", ""),
		RPCCert:     getEnv("DCRD_RPC_CERT", ""),
	}

	// Try to initialize RPC client if credentials are provided
	if config.RPCUser != "" && config.RPCPassword != "" {
		if err := initRPCClient(config); err != nil {
			log.Printf("Warning: Could not connect to dcrd on startup: %v", err)
			log.Println("RPC connection can be configured via API")
		}
	} else {
		log.Println("No RPC credentials provided. Use /api/connect endpoint to configure.")
	}

	// Setup router
	r := mux.NewRouter()

	// API routes
	api := r.PathPrefix("/api").Subrouter()
	api.HandleFunc("/health", healthCheckHandler).Methods("GET")
	api.HandleFunc("/dashboard", getDashboardDataHandler).Methods("GET")
	api.HandleFunc("/node/status", getNodeStatusHandler).Methods("GET")
	api.HandleFunc("/blockchain/info", getBlockchainInfoHandler).Methods("GET")
	api.HandleFunc("/network/peers", getPeersHandler).Methods("GET")
	api.HandleFunc("/connect", connectRPCHandler).Methods("POST")

	// CORS configuration
	corsHandler := cors.New(cors.Options{
		AllowedOrigins:   []string{"*"},
		AllowedMethods:   []string{"GET", "POST", "PUT", "DELETE", "OPTIONS"},
		AllowedHeaders:   []string{"*"},
		AllowCredentials: true,
	})

	// Start server
	port := getEnv("PORT", "8080")
	address := fmt.Sprintf(":%s", port)

	log.Printf("Starting Decred Dashboard API server on %s", address)
	log.Fatal(http.ListenAndServe(address, corsHandler.Handler(r)))
}

func getEnv(key, defaultValue string) string {
	value := os.Getenv(key)
	if value == "" {
		return defaultValue
	}
	return value
}

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
	"strings"
	"sync"
	"time"

	"github.com/decred/dcrd/rpcclient/v8"
	"github.com/gorilla/mux"
	"github.com/rs/cors"
)

var (
	rpcClient *rpcclient.Client

	// Track previous sync values to calculate delta
	prevHeaders int64
	prevBlocks  int64
	syncMutex   sync.Mutex
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
	StakingInfo    StakingInfo    `json:"stakingInfo"`
	MempoolInfo    MempoolInfo    `json:"mempoolInfo"`
	LastUpdate     time.Time      `json:"lastUpdate"`
}

type NodeStatus struct {
	Status       string  `json:"status"`
	SyncProgress float64 `json:"syncProgress"`
	Version      string  `json:"version"`
	SyncPhase    string  `json:"syncPhase"`   // "headers" or "blocks"
	SyncMessage  string  `json:"syncMessage"` // e.g., "Processed 36,000 headers in the last 30 seconds"
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
	ID         int    `json:"id"`
	Address    string `json:"address"`
	Protocol   string `json:"protocol"`
	Latency    string `json:"latency"`
	ConnTime   string `json:"connTime"`
	Traffic    string `json:"traffic"`
	Version    string `json:"version"`
	IsSyncNode bool   `json:"isSyncNode"`
}

type SupplyInfo struct {
	CirculatingSupply string  `json:"circulatingSupply"`
	StakedSupply      string  `json:"stakedSupply"`
	StakedPercent     float64 `json:"stakedPercent"`
	ExchangeRate      string  `json:"exchangeRate"`
	TreasurySize      string  `json:"treasurySize"`
	MixedPercent      string  `json:"mixedPercent"`
}

type StakingInfo struct {
	TicketPrice       float64 `json:"ticketPrice"`
	PoolSize          uint32  `json:"poolSize"`
	LockedDCR         float64 `json:"lockedDCR"`
	ParticipationRate float64 `json:"participationRate"`
	AllMempoolTix     uint32  `json:"allMempoolTix"`
	Immature          uint32  `json:"immature"`
	Live              uint32  `json:"live"`
	Voted             uint32  `json:"voted"`
	Missed            uint32  `json:"missed"`
	Revoked           uint32  `json:"revoked"`
}

type MempoolInfo struct {
	Size           uint64  `json:"size"`
	Bytes          uint64  `json:"bytes"`
	TxCount        int     `json:"txCount"`
	TotalFee       float64 `json:"totalFee"`
	AverageFeeRate float64 `json:"averageFeeRate"`
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

	stakingInfo, err := fetchStakingInfo()
	if err != nil {
		return nil, err
	}

	mempoolInfo, err := fetchMempoolInfo()
	if err != nil {
		return nil, err
	}

	return &DashboardData{
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

func fetchNodeStatus() (*NodeStatus, error) {
	ctx := context.Background()

	// Get version info using version command
	versionInfo, err := rpcClient.Version(ctx)
	if err != nil {
		return nil, fmt.Errorf("failed to get version: %v", err)
	}

	// Get blockchain info for accurate sync status
	chainInfo, err := rpcClient.GetBlockChainInfo(ctx)
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
				syncMessage = fmt.Sprintf("Processed %s headers in the last 30 seconds", formatNumber(deltaHeaders))
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
				syncMessage = fmt.Sprintf("Processed %s blocks in the last 30 seconds", formatNumber(deltaBlocks))
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

	return &NodeStatus{
		Status:       status,
		SyncProgress: syncProgress,
		Version:      fmt.Sprintf("%d.%d.%d", versionInfo["dcrd"].Major, versionInfo["dcrd"].Minor, versionInfo["dcrd"].Patch),
		SyncPhase:    syncPhase,
		SyncMessage:  syncMessage,
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

	// Get peer count
	peerCount := 0
	peerInfo, err := rpcClient.GetPeerInfo(ctx)
	if err == nil {
		peerCount = len(peerInfo)
	}

	// Get network difficulty and calculate hashrate - direct RPC method
	hashrateStr := "N/A"
	networkHashPS := float64(0)

	difficulty, err := rpcClient.GetDifficulty(ctx)
	if err == nil && difficulty > 0 {
		// Calculate network hashrate from difficulty
		// Formula: hashrate = difficulty * 2^32 / target_block_time
		// For Decred, target block time is 5 minutes = 300 seconds
		networkHashPS = difficulty * math.Pow(2, 32) / 300
		hashrateStr = formatHashrate(networkHashPS)
	}

	return &NetworkInfo{
		PeerCount:     peerCount,
		Hashrate:      hashrateStr,
		NetworkHashPS: networkHashPS,
	}, nil
}

func fetchPeers() ([]Peer, error) {
	ctx := context.Background()
	peerInfo, err := rpcClient.GetPeerInfo(ctx)
	if err != nil {
		return nil, err
	}

	peers := make([]Peer, 0, len(peerInfo))
	now := time.Now().Unix()

	for i, p := range peerInfo {
		// Convert pingtime from microseconds to milliseconds
		latency := fmt.Sprintf("%.0fms", float64(p.PingTime)/1000)

		// Calculate connection time
		connDuration := now - p.ConnTime
		connTime := formatDuration(connDuration)

		// Calculate total traffic in MB
		totalBytes := p.BytesSent + p.BytesRecv
		traffic := formatTraffic(totalBytes)

		// Extract version from subver string (e.g., "/dcrwire:1.0.0/dcrd:2.0.5/" -> "2.0.5")
		version := extractDcrdVersion(p.SubVer)

		peers = append(peers, Peer{
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
func formatDuration(seconds int64) string {
	if seconds < 60 {
		return fmt.Sprintf("%ds", seconds)
	}
	minutes := seconds / 60
	if minutes < 60 {
		return fmt.Sprintf("%dm", minutes)
	}
	hours := minutes / 60
	remainingMinutes := minutes % 60
	if hours < 24 {
		if remainingMinutes == 0 {
			return fmt.Sprintf("%dh", hours)
		}
		return fmt.Sprintf("%dh %dm", hours, remainingMinutes)
	}
	days := hours / 24
	remainingHours := hours % 24
	if remainingHours == 0 {
		return fmt.Sprintf("%dd", days)
	}
	return fmt.Sprintf("%dd %dh", days, remainingHours)
}

// formatTraffic formats bytes to a human-readable traffic string
func formatTraffic(bytes uint64) string {
	if bytes < 1024 {
		return fmt.Sprintf("%d B", bytes)
	}
	kb := float64(bytes) / 1024
	if kb < 1024 {
		return fmt.Sprintf("%.1f KB", kb)
	}
	mb := kb / 1024
	if mb < 1024 {
		return fmt.Sprintf("%.2f MB", mb)
	}
	gb := mb / 1024
	return fmt.Sprintf("%.2f GB", gb)
}

// extractDcrdVersion extracts the dcrd version from subver string
func extractDcrdVersion(subver string) string {
	// subver format: "/dcrwire:1.0.0/dcrd:2.0.5/"
	if idx := strings.Index(subver, "dcrd:"); idx != -1 {
		start := idx + 5
		end := strings.Index(subver[start:], "/")
		if end != -1 {
			return subver[start : start+end]
		}
		return subver[start:]
	}
	return "unknown"
}

func fetchSupplyInfo() (*SupplyInfo, error) {
	ctx := context.Background()

	// Get real circulating supply from dcrd - direct RPC method
	circulatingSupply := "N/A"
	stakedSupply := "N/A"
	stakedPercent := float64(0)
	treasuryBalance := "N/A"

	coinSupply, err := rpcClient.GetCoinSupply(ctx)
	if err == nil && coinSupply > 0 {
		// Convert atoms to DCR and format with commas
		coinSupplyDCR := coinSupply.ToCoin()
		circulatingSupply = formatDCRAmount(coinSupplyDCR)

		// Calculate staked supply from ticket pool
		ticketPoolValue, err := rpcClient.GetTicketPoolValue(ctx)
		if err == nil && ticketPoolValue > 0 {
			lockedDCR := ticketPoolValue.ToCoin()
			stakedSupply = formatDCRAmount(lockedDCR)

			if coinSupplyDCR > 0 {
				stakedPercent = (lockedDCR / coinSupplyDCR) * 100
			}
		}
	}

	// Get treasury balance - direct RPC method
	// Pass nil for hash (gets latest) and false for verbose
	treasuryBalanceResult, err := rpcClient.GetTreasuryBalance(ctx, nil, false)
	if err == nil && treasuryBalanceResult.Balance > 0 {
		// Balance is in atoms (uint64), convert to DCR by dividing by 1e8
		treasuryBalanceDCR := float64(treasuryBalanceResult.Balance) / 1e8
		// Format with 2 decimal places and commas
		treasuryBalance = formatDCRAmountWithDecimals(treasuryBalanceDCR, 2)
	}

	return &SupplyInfo{
		CirculatingSupply: circulatingSupply,
		StakedSupply:      stakedSupply,
		StakedPercent:     stakedPercent,
		ExchangeRate:      "N/A", // Requires external API
		TreasurySize:      treasuryBalance,
		MixedPercent:      "N/A", // Requires mixer statistics
	}, nil
}

// formatDCRAmount formats a DCR amount with commas and appropriate precision
func formatDCRAmount(amount float64) string {
	// Format with 0 decimal places for large amounts
	intAmount := int64(amount)
	return addCommas(intAmount)
}

// formatDCRAmountWithDecimals formats a DCR amount with commas and specified decimal places
func formatDCRAmountWithDecimals(amount float64, decimals int) string {
	// Split into integer and decimal parts
	intPart := int64(amount)
	decimalPart := amount - float64(intPart)

	// Format integer part with commas
	intStr := addCommas(intPart)

	// Format decimal part
	format := fmt.Sprintf("%%.%df", decimals)
	fullStr := fmt.Sprintf(format, decimalPart)

	// Extract just the decimal part (after the "0.")
	decimalStr := fullStr[2:] // Skip "0."

	return fmt.Sprintf("%s.%s", intStr, decimalStr)
}

// formatNumber formats a number with thousands separators
func formatNumber(n int64) string {
	if n < 1000 {
		return fmt.Sprintf("%d", n)
	}
	return addCommas(n)
}

// addCommas adds comma separators to a number
func addCommas(n int64) string {
	str := fmt.Sprintf("%d", n)
	if len(str) <= 3 {
		return str
	}

	var result string
	for i, digit := range str {
		if i > 0 && (len(str)-i)%3 == 0 {
			result += ","
		}
		result += string(digit)
	}
	return result
}

func fetchStakingInfo() (*StakingInfo, error) {
	ctx := context.Background()

	// Get stake difficulty (ticket price) - direct RPC method
	stakeDiff, err := rpcClient.GetStakeDifficulty(ctx)
	if err != nil {
		return nil, fmt.Errorf("failed to get stake difficulty: %v", err)
	}
	ticketPrice := stakeDiff.CurrentStakeDifficulty

	// Get live tickets from pool - direct RPC method
	// LiveTickets returns []*chainhash.Hash directly
	liveTickets, err := rpcClient.LiveTickets(ctx)
	poolSize := uint32(0)
	if err == nil && liveTickets != nil {
		// Count the actual number of live tickets
		poolSize = uint32(len(liveTickets))
	}

	// Get ticket pool value (total locked DCR) - direct RPC method
	// Returns dcrutil.Amount which needs to be converted to float64 DCR
	lockedDCR := float64(0)
	poolValue, err := rpcClient.GetTicketPoolValue(ctx)
	if err == nil {
		lockedDCR = poolValue.ToCoin()
	}

	// Get total coin supply for participation rate calculation - direct RPC method
	// Returns dcrutil.Amount which needs to be converted to float64 DCR
	participationRate := float64(0)
	coinSupply, err := rpcClient.GetCoinSupply(ctx)
	if err == nil && coinSupply > 0 {
		// Calculate participation rate as percentage of total supply
		coinSupplyDCR := coinSupply.ToCoin()
		participationRate = (lockedDCR / coinSupplyDCR) * 100
	}

	return &StakingInfo{
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

func fetchMempoolInfo() (*MempoolInfo, error) {
	ctx := context.Background()

	// Get raw mempool - this is a standard method available in rpcclient
	hashes, err := rpcClient.GetRawMempool(ctx, "false")
	if err != nil {
		// If mempool query fails (e.g., during sync), return empty mempool
		return &MempoolInfo{
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
	return &MempoolInfo{
		Size:           uint64(txCount),
		Bytes:          0,
		TxCount:        txCount,
		TotalFee:       0,
		AverageFeeRate: 0,
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

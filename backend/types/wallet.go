// Copyright (c) 2015-2025 The Decred developers
// Use of this source code is governed by an ISC
// license that can be found in the LICENSE file.

package types

import "time"

// WalletDashboardData represents all wallet dashboard metrics
type WalletDashboardData struct {
	WalletStatus WalletStatus  `json:"walletStatus"`
	AccountInfo  AccountInfo   `json:"accountInfo"`
	Accounts     []AccountInfo `json:"accounts"`
	LastUpdate   time.Time     `json:"lastUpdate"`
}

type WalletStatus struct {
	Status           string  `json:"status"` // "locked", "unlocked", "syncing", "synced", "no_wallet"
	SyncProgress     float64 `json:"syncProgress"`
	SyncHeight       int64   `json:"syncHeight"`
	BestBlockHash    string  `json:"bestBlockHash"`
	Version          string  `json:"version"`
	Unlocked         bool    `json:"unlocked"`
	RescanInProgress bool    `json:"rescanInProgress"`
	SyncMessage      string  `json:"syncMessage"`
}

type AccountInfo struct {
	AccountName        string  `json:"accountName"`
	TotalBalance       float64 `json:"totalBalance"`
	SpendableBalance   float64 `json:"spendableBalance"`
	ImmatureBalance    float64 `json:"immatureBalance"`
	UnconfirmedBalance float64 `json:"unconfirmedBalance"`
	AccountNumber      uint32  `json:"accountNumber"`
}

type Transaction struct {
	TxID          string    `json:"txid"`
	Amount        float64   `json:"amount"`
	Fee           float64   `json:"fee"`
	Confirmations int64     `json:"confirmations"`
	Time          time.Time `json:"time"`
	Type          string    `json:"type"` // "send", "receive", "ticket", etc.
	Comment       string    `json:"comment"`
}

type Address struct {
	Address string `json:"address"`
	Account string `json:"account"`
	Used    bool   `json:"used"`
	Path    string `json:"path"` // BIP44 path
}

type ImportXpubRequest struct {
	Xpub        string `json:"xpub"`
	AccountName string `json:"accountName"`
	Rescan      bool   `json:"rescan"`
}

type ImportXpubResponse struct {
	Success    bool   `json:"success"`
	Message    string `json:"message"`
	AccountNum uint32 `json:"accountNum,omitempty"`
}

type RescanRequest struct {
	BeginHeight int32 `json:"beginHeight"`
}

type RescanResponse struct {
	Success bool   `json:"success"`
	Message string `json:"message"`
}

// Copyright (c) 2015-2025 The Decred developers
// Use of this source code is governed by an ISC
// license that can be found in the LICENSE file.

package rpc

import (
	"context"
	"fmt"
	"io/ioutil"
	"log"

	"github.com/decred/dcrd/rpcclient/v8"
)

var (
	// DcrdClient is the RPC client for dcrd
	DcrdClient *rpcclient.Client

	// WalletClient is the RPC client for dcrwallet
	WalletClient *rpcclient.Client
)

// Config holds the RPC connection configuration
type Config struct {
	RPCHost     string
	RPCPort     string
	RPCUser     string
	RPCPassword string
	RPCCert     string
}

// InitDcrdClient initializes the dcrd RPC client
func InitDcrdClient(config Config) error {
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

	DcrdClient, err = rpcclient.New(connCfg, nil)
	if err != nil {
		return fmt.Errorf("failed to create RPC client: %v", err)
	}

	// Test connection
	ctx := context.Background()
	_, err = DcrdClient.GetBlockCount(ctx)
	if err != nil {
		return fmt.Errorf("failed to connect to dcrd: %v", err)
	}

	log.Println("Successfully connected to dcrd RPC with TLS")
	return nil
}

// InitWalletClient initializes the dcrwallet RPC client
func InitWalletClient(config Config) error {
	// Read the TLS certificate if provided
	var certs []byte
	var err error

	if config.RPCCert != "" {
		log.Printf("Reading wallet TLS certificate from: %s", config.RPCCert)
		certs, err = ioutil.ReadFile(config.RPCCert)
		if err != nil {
			return fmt.Errorf("failed to read wallet RPC certificate: %v", err)
		}
		log.Printf("Successfully loaded wallet TLS certificate (%d bytes)", len(certs))
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

	WalletClient, err = rpcclient.New(connCfg, nil)
	if err != nil {
		return fmt.Errorf("failed to create wallet RPC client: %v", err)
	}

	// Test connection with getinfo
	ctx := context.Background()
	_, err = WalletClient.GetInfo(ctx)
	if err != nil {
		// Wallet might be locked or not initialized, but connection is OK
		log.Printf("Wallet RPC connected but getinfo failed (may be locked): %v", err)
	} else {
		log.Println("Successfully connected to dcrwallet RPC with TLS")
	}

	return nil
}

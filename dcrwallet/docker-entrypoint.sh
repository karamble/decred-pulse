#!/bin/sh
# Copyright (c) 2015-2025 The Decred developers
# Use of this source code is governed by an ISC
# license that can be found in the LICENSE file.

set -e

WALLET_DIR="/home/dcrwallet/.dcrwallet"
WALLET_DB="${WALLET_DIR}/mainnet/wallet.db"
PASSPHRASE_FILE="${WALLET_DIR}/.passphrase"
CERT_DIR="/certs"
RPC_CERT="${CERT_DIR}/rpc.cert"
RPC_KEY="${CERT_DIR}/rpc.key"

# Create wallet directory if it doesn't exist
mkdir -p "${WALLET_DIR}"

# Check if wallet database exists
if [ ! -f "${WALLET_DB}" ]; then
    echo "============================================"
    echo "No wallet database found. Creating new wallet..."
    echo "============================================"
    
    # Generate a secure random passphrase for private keys
    echo "Generating secure passphrase..."
    PRIVATE_PASS=$(openssl rand -base64 32)
    echo "${PRIVATE_PASS}" > "${PASSPHRASE_FILE}"
    chmod 600 "${PASSPHRASE_FILE}"
    echo "✓ Passphrase generated and stored securely"
    
    # Create wallet using expect to automate interactive prompts
    echo "Creating wallet database..."
    
    # Create expect script for automated wallet creation
    cat > /tmp/create-wallet.exp <<'EXPECTEOF'
#!/usr/bin/expect -f
set timeout 120
set passphrase [lindex $argv 0]

log_user 1
spawn dcrwallet --create

expect {
    -re "Enter the private passphrase.*:" {
        send "$passphrase\r"
        exp_continue
    }
    -re "Confirm.*:" {
        send "$passphrase\r"
        exp_continue
    }
    -re "additional layer.*:" {
        send "no\r"
        exp_continue
    }
    -re "existing.*seed.*:" {
        send "no\r"
        exp_continue
    }
    -re "Your wallet generation seed is:" {
        expect -re "(\[a-z\]+ ){5,}"
        set seed $expect_out(0,string)
        puts "\n============================================"
        puts "WALLET SEED GENERATED (stored for backup)"
        puts "============================================\n"
        exp_continue
    }
    -re "Hex:" {
        expect -re "\[0-9a-f\]{32,}"
        exp_continue
    }
    -re "(?i)(enter|type).*(OK|ok)" {
        send "OK\r"
        exp_continue
    }
    -re "(?i)additional account.*extended public key" {
        send "no\r"
        exp_continue
    }
    -re "Creating the wallet" {
        puts "\n✓ Wallet creation successful\n"
    }
    eof {
        puts "\n✓ Wallet setup complete\n"
    }
    timeout {
        puts "\nERROR: Timeout during wallet creation\n"
        exit 1
    }
}

exit 0
EXPECTEOF
    
    chmod +x /tmp/create-wallet.exp
    
    # Run the expect script
    if /tmp/create-wallet.exp "${PRIVATE_PASS}"; then
        echo "✓ Wallet created successfully"
        echo ""
        echo "============================================"
        echo "IMPORTANT NOTES:"
        echo "1. This is a WATCH-ONLY wallet setup"
        echo "2. You need to import an xpub key via the web interface"
        echo "3. The seed above is for backup purposes only"
        echo "4. Private passphrase is stored in: ${PASSPHRASE_FILE}"
        echo "============================================"
        echo ""
        rm -f /tmp/create-wallet.exp
    else
        echo "ERROR: Failed to create wallet"
        rm -f /tmp/create-wallet.exp
        exit 1
    fi
else
    echo "✓ Wallet database found, loading existing wallet..."
    
    # Read existing passphrase
    if [ -f "${PASSPHRASE_FILE}" ]; then
        PRIVATE_PASS=$(cat "${PASSPHRASE_FILE}")
        echo "✓ Passphrase loaded from storage"
    else
        echo "ERROR: Wallet database exists but passphrase file not found!"
        echo "Please restore the passphrase file or recreate the wallet."
        exit 1
    fi
fi

# Verify RPC certificates exist (shared from dcrd)
if [ ! -f "${RPC_CERT}" ] || [ ! -f "${RPC_KEY}" ]; then
    echo "WARNING: RPC certificates not found at ${CERT_DIR}"
    echo "Waiting for dcrd to generate certificates..."
    
    # Wait up to 60 seconds for certificates
    for i in $(seq 1 60); do
        if [ -f "${RPC_CERT}" ] && [ -f "${RPC_KEY}" ]; then
            echo "✓ RPC certificates found"
            break
        fi
        sleep 1
        if [ $i -eq 60 ]; then
            echo "ERROR: Timeout waiting for RPC certificates"
            exit 1
        fi
    done
else
    echo "✓ Using RPC certificates from dcrd"
fi

# Start dcrwallet with configuration
echo "Starting dcrwallet..."
echo "Configuration:"
echo "  RPC Listen: 0.0.0.0:9110"
echo "  gRPC Listen: 0.0.0.0:9111"
echo "  dcrd Connection: ${DCRD_RPC_HOST:-dcrd}:9109"
echo ""

# Execute dcrwallet with all provided arguments
exec dcrwallet \
    --username="${DCRWALLET_RPC_USER}" \
    --password="${DCRWALLET_RPC_PASS}" \
    --rpclisten=0.0.0.0:9110 \
    --rpccert="${RPC_CERT}" \
    --rpckey="${RPC_KEY}" \
    --dcrdusername="${DCRD_RPC_USER}" \
    --dcrdpassword="${DCRD_RPC_PASS}" \
    --rpcconnect="${DCRD_RPC_HOST:-dcrd}:9109" \
    --cafile="${RPC_CERT}" \
    --grpclisten=0.0.0.0:9111 \
    --pass="${PRIVATE_PASS}" \
    "$@"


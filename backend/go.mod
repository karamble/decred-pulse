module decred-pulse-backend

go 1.21

require (
	decred.org/dcrwallet/v4 v4.1.0
	github.com/decred/dcrd/rpcclient/v8 v8.0.1
	github.com/gorilla/mux v1.8.1
	github.com/gorilla/websocket v1.5.1
	github.com/rs/cors v1.10.1
	google.golang.org/grpc v1.65.0
)

// Exclude old genproto to avoid ambiguous import
exclude google.golang.org/genproto v0.0.0-20200526211855-cb27e3aa2013

require (
	github.com/agl/ed25519 v0.0.0-20170116200512-5312a6153412 // indirect
	github.com/dchest/siphash v1.2.3 // indirect
	github.com/decred/base58 v1.0.5 // indirect
	github.com/decred/dcrd/blockchain/stake/v5 v5.0.1 // indirect
	github.com/decred/dcrd/chaincfg/chainhash v1.0.4 // indirect
	github.com/decred/dcrd/chaincfg/v3 v3.2.1 // indirect
	github.com/decred/dcrd/crypto/blake256 v1.0.1 // indirect
	github.com/decred/dcrd/crypto/ripemd160 v1.0.2 // indirect
	github.com/decred/dcrd/database/v3 v3.0.2 // indirect
	github.com/decred/dcrd/dcrec v1.0.1 // indirect
	github.com/decred/dcrd/dcrec/edwards/v2 v2.0.3 // indirect
	github.com/decred/dcrd/dcrec/secp256k1/v4 v4.3.0 // indirect
	github.com/decred/dcrd/dcrjson/v4 v4.1.0 // indirect
	github.com/decred/dcrd/dcrutil/v4 v4.0.2 // indirect
	github.com/decred/dcrd/gcs/v4 v4.1.0 // indirect
	github.com/decred/dcrd/rpc/jsonrpc/types/v4 v4.3.0 // indirect
	github.com/decred/dcrd/txscript/v4 v4.1.1 // indirect
	github.com/decred/dcrd/wire v1.7.0 // indirect
	github.com/decred/go-socks v1.1.0 // indirect
	github.com/decred/slog v1.2.0 // indirect
	github.com/klauspost/cpuid/v2 v2.2.5 // indirect
	golang.org/x/net v0.25.0 // indirect
	golang.org/x/sys v0.20.0 // indirect
	golang.org/x/text v0.15.0 // indirect
	google.golang.org/genproto/googleapis/rpc v0.0.0-20240528184218-531527333157 // indirect
	google.golang.org/protobuf v1.34.1 // indirect
	lukechampine.com/blake3 v1.3.0 // indirect
)

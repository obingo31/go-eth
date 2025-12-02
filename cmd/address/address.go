package main

import (
	"context"
	"flag"
	"fmt"
	"log"

	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/ethclient"
)

func main() {
	const defaultAddr = "0xF4D17Dd253A5a21555bF1a5B9B7285ed764AF706"
	const defaultRPC = "http://localhost:8545"

	addrFlag := flag.String("addr", "", "hex-encoded Ethereum address")
	rpcFlag := flag.String("rpc", defaultRPC, "Ethereum RPC endpoint (defaults to local dev chain)")
	flag.Parse()

	target := *addrFlag
	if target == "" {
		target = defaultAddr
	}

	if !common.IsHexAddress(target) {
		log.Fatalf("invalid address: %s", target)
	}

	addr := common.HexToAddress(target)

	fmt.Printf("Input: %s\n", target)
	fmt.Printf("Checksum: %s\n", addr.Hex())
	fmt.Printf("Bytes: %x\n", addr.Bytes())
	fmt.Printf("Hash: %s\n", common.BytesToHash(addr.Bytes()).Hex())

	client, err := ethclient.Dial(*rpcFlag)
	if err != nil {
		log.Fatalf("dial RPC: %v", err)
	}
	defer client.Close()

	ctx := context.Background()
	balanceWei, err := client.BalanceAt(ctx, addr, nil) // nil block => latest balance
	if err != nil {
		log.Fatalf("query balance: %v", err)
	}

	fmt.Printf("Balance (wei): %s\n", balanceWei.String())
}

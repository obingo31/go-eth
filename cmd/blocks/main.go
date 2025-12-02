package main

import (
	"context"
	"flag"
	"fmt"
	"log"
	"math/big"
	"time"

	"github.com/ethereum/go-ethereum/ethclient"
)

func main() {
	rpcFlag := flag.String("rpc", "http://localhost:8545", "Ethereum RPC endpoint")
	blockFlag := flag.Int64("block", 5671744, "block number to inspect (-1 for latest)")
	flag.Parse()

	client, err := ethclient.Dial(*rpcFlag)
	if err != nil {
		log.Fatalf("dial RPC: %v", err)
	}
	defer client.Close()

	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	header, err := client.HeaderByNumber(ctx, nil)
	if err != nil {
		log.Fatalf("fetch latest header: %v", err)
	}
	fmt.Printf("Latest block: %s\n", header.Number.String())

	var blockNumber *big.Int
	if *blockFlag >= 0 {
		blockNumber = big.NewInt(*blockFlag)
	}

	block, err := client.BlockByNumber(ctx, blockNumber)
	if err != nil {
		log.Fatalf("fetch block: %v", err)
	}

	fmt.Printf("Block %d details\n", block.NumberU64())
	fmt.Printf("  Timestamp : %d\n", block.Time())
	fmt.Printf("  Difficulty: %s\n", block.Difficulty().String())
	fmt.Printf("  Hash      : %s\n", block.Hash().Hex())
	fmt.Printf("  Tx count  : %d\n", len(block.Transactions()))

	count, err := client.TransactionCount(ctx, block.Hash())
	if err != nil {
		log.Fatalf("transaction count: %v", err)
	}
	fmt.Printf("TransactionCount API: %d\n", count)
}

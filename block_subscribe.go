package main

import (
	"context"
	"flag"
	"fmt"
	"log"

	"github.com/ethereum/go-ethereum/core/types"
	"github.com/ethereum/go-ethereum/ethclient"
)

func main() {
	subscribeToBlocks()
}

func subscribeToBlocks() {
	wsURL := flag.String("ws", "ws://127.0.0.1:8545", "WebSocket RPC endpoint (Ganache, Anvil, etc.)")
	flag.Parse()

	ctx := context.Background()

	client, err := ethclient.Dial(*wsURL)
	if err != nil {
		log.Fatalf("dial websocket: %v", err)
	}
	defer client.Close()

	headers := make(chan *types.Header)
	sub, err := client.SubscribeNewHead(ctx, headers)
	if err != nil {
		log.Fatalf("subscribe new heads: %v", err)
	}
	defer sub.Unsubscribe()

	for {
		select {
		case err := <-sub.Err():
			log.Fatalf("subscription error: %v", err)
		case header := <-headers:
			fmt.Printf("new head %s\n", header.Hash().Hex())

			block, err := client.BlockByHash(ctx, header.Hash())
			if err != nil {
				log.Fatalf("fetch block %s: %v", header.Hash(), err)
			}

			fmt.Printf("block hash: %s\n", block.Hash().Hex())
			fmt.Printf("number   : %d\n", block.Number().Uint64())
			fmt.Printf("time     : %d\n", block.Time())
			fmt.Printf("nonce    : %d\n", block.Nonce())
			fmt.Printf("tx count : %d\n", len(block.Transactions()))
		}
	}
}

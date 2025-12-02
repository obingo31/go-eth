package main

import (
	"context"
	"flag"
	"fmt"
	"log"
	"time"

	"github.com/ethereum/go-ethereum"
	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/core/types"
	"github.com/ethereum/go-ethereum/ethclient"
)

func main() {
	wsFlag := flag.String("ws", "ws://127.0.0.1:8545", "WebSocket RPC endpoint")
	addrFlag := flag.String("addr", "", "Contract address to filter events for")
	flag.Parse()

	if *addrFlag == "" {
		log.Fatal("--addr is required")
	}

	ctx := context.Background()
	client, err := ethclient.DialContext(ctx, *wsFlag)
	if err != nil {
		log.Fatalf("dial websocket: %v", err)
	}
	defer client.Close()

	contractAddress := common.HexToAddress(*addrFlag)
	query := ethereum.FilterQuery{Addresses: []common.Address{contractAddress}}

	logsCh := make(chan types.Log)
	sub, err := client.SubscribeFilterLogs(ctx, query, logsCh)
	if err != nil {
		log.Fatalf("subscribe logs: %v", err)
	}
	defer sub.Unsubscribe()

	log.Printf("listening for logs from %s over %s", contractAddress.Hex(), *wsFlag)

	for {
		select {
		case err := <-sub.Err():
			log.Fatalf("subscription error: %v", err)
		case event := <-logsCh:
			fmt.Printf("[%s] block=%d tx=%s topics=%d data=%d bytes\n",
				time.Now().Format(time.RFC3339),
				event.BlockNumber,
				event.TxHash.Hex(),
				len(event.Topics),
				len(event.Data))
		}
	}
}

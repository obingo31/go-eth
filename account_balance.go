package main

import (
	"context"
	"flag"
	"fmt"
	"log"
	"math"
	"math/big"

	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/ethclient"
)

func main() {
	const defaultAddr = "0xF4D17Dd253A5a21555bF1a5B9B7285ed764AF706"
	const defaultRPC = "https://mainnet.infura.io/v3/<PROJECT_ID>"

	rpcFlag := flag.String("rpc", defaultRPC, "Ethereum RPC endpoint (Infura, Alchemy, local node, etc.)")
	addrFlag := flag.String("addr", defaultAddr, "Hex-encoded Ethereum account address")
	blockFlag := flag.Int64("block", 5532993, "Historical block number to inspect")
	flag.Parse()

	if !common.IsHexAddress(*addrFlag) {
		log.Fatalf("invalid address: %s", *addrFlag)
	}

	client, err := ethclient.Dial(*rpcFlag)
	if err != nil {
		log.Fatalf("dial RPC: %v", err)
	}
	defer client.Close()

	ctx := context.Background()
	account := common.HexToAddress(*addrFlag)

	latestBalance, err := client.BalanceAt(ctx, account, nil)
	if err != nil {
		log.Fatalf("latest balance: %v", err)
	}
	fmt.Printf("Latest balance (wei): %s\n", latestBalance)

	historicalBlock := big.NewInt(*blockFlag)
	historicalBalance, err := client.BalanceAt(ctx, account, historicalBlock)
	if err != nil {
		log.Fatalf("balance at block %d: %v", *blockFlag, err)
	}
	fmt.Printf("Balance at block %d (wei): %s\n", historicalBlock, historicalBalance)

	fbalance := new(big.Float).SetInt(historicalBalance)
	denom := big.NewFloat(math.Pow10(18))
	ethValue := new(big.Float).Quo(fbalance, denom)
	fmt.Printf("Balance at block %d (ETH): %s\n", historicalBlock, ethValue.Text('f', 18))

	pendingBalance, err := client.PendingBalanceAt(ctx, account)
	if err != nil {
		log.Fatalf("pending balance: %v", err)
	}
	fmt.Printf("Pending balance (wei): %s\n", pendingBalance)
}

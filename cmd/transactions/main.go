package main

import (
	"context"
	"flag"
	"fmt"
	"log"
	"math/big"
	"time"

	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/core/types"
	"github.com/ethereum/go-ethereum/ethclient"
)

func main() {
	rpcFlag := flag.String("rpc", "http://localhost:8545", "Ethereum RPC endpoint")
	blockFlag := flag.Int64("block", 5671744, "block number to inspect")
	hashFlag := flag.String("tx", "", "specific transaction hash to fetch")
	flag.Parse()

	client, err := ethclient.Dial(*rpcFlag)
	if err != nil {
		log.Fatalf("dial RPC: %v", err)
	}
	defer client.Close()

	ctx, cancel := context.WithTimeout(context.Background(), 15*time.Second)
	defer cancel()

	blockNumber := big.NewInt(*blockFlag)
	block, err := client.BlockByNumber(ctx, blockNumber)
	if err != nil {
		log.Fatalf("fetch block %d: %v", blockNumber, err)
	}

	fmt.Printf("Inspecting block %d (%s) with %d txns\n", block.NumberU64(), block.Hash(), len(block.Transactions()))

	chainID, err := client.NetworkID(ctx)
	if err != nil {
		log.Fatalf("fetch network ID: %v", err)
	}

	signer := types.LatestSignerForChainID(chainID)

	for _, tx := range block.Transactions() {
		fmt.Printf("Tx %s\n", tx.Hash().Hex())
		fmt.Printf("  Value   : %s\n", tx.Value().String())
		fmt.Printf("  Gas     : %d\n", tx.Gas())
		fmt.Printf("  GasPrice: %s\n", tx.GasPrice().String())
		fmt.Printf("  Nonce   : %d\n", tx.Nonce())
		fmt.Printf("  DataLen : %d\n", len(tx.Data()))
		if to := tx.To(); to != nil {
			fmt.Printf("  To      : %s\n", to.Hex())
		} else {
			fmt.Printf("  To      : <contract creation>\n")
		}

		if from, err := types.Sender(signer, tx); err == nil {
			fmt.Printf("  From    : %s\n", from.Hex())
		}

		receipt, err := client.TransactionReceipt(ctx, tx.Hash())
		if err != nil {
			log.Fatalf("receipt for %s: %v", tx.Hash(), err)
		}
		fmt.Printf("  Status  : %d\n", receipt.Status)
	}

	blockHash := block.Hash()
	count, err := client.TransactionCount(ctx, blockHash)
	if err != nil {
		log.Fatalf("transaction count: %v", err)
	}
	fmt.Printf("TransactionCount(%s) = %d\n", blockHash.Hex(), count)

	for idx := uint(0); idx < count; idx++ {
		tx, err := client.TransactionInBlock(ctx, blockHash, idx)
		if err != nil {
			log.Fatalf("tx %d in block %s: %v", idx, blockHash, err)
		}
		fmt.Printf("TransactionInBlock idx %d -> %s\n", idx, tx.Hash().Hex())
	}

	if *hashFlag != "" {
		txHash := common.HexToHash(*hashFlag)
		tx, isPending, err := client.TransactionByHash(ctx, txHash)
		if err != nil {
			log.Fatalf("transaction by hash: %v", err)
		}
		fmt.Printf("TransactionByHash %s -> pending=%v\n", tx.Hash().Hex(), isPending)
	}
}

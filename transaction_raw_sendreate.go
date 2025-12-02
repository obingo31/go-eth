package main

import (
	"context"
	"encoding/hex"
	"fmt"
	"log"

	"github.com/ethereum/go-ethereum/core/types"
	"github.com/ethereum/go-ethereum/ethclient"
	"github.com/ethereum/go-ethereum/rlp"
)

func main() {
	client, err := ethclient.Dial("http://127.0.0.1:8545")
	if err != nil {
		log.Fatal(err)
	}

	rawTx := "f87104847735940082520894e236d0200aa21c7896975d7a2419150149541a95880de0b6b3a764000080860335c0e77a55a03d0970a032b6338248c04dacd2db7087b8c0d5eae970fb43a9dd21dfced2bf1da031d6306f4856441fb8b29745cad9d3d0461f9053156167d1b5c1f1c8cfd73f72"

	rawTxBytes, err := hex.DecodeString(rawTx)
	if err != nil {
		log.Fatal(err)
	}

	tx := new(types.Transaction)
	if err := rlp.DecodeBytes(rawTxBytes, &tx); err != nil {
		log.Fatal(err)
	}

	if err := client.SendTransaction(context.Background(), tx); err != nil {
		log.Fatal(err)
	}

	fmt.Printf("tx sent: %s", tx.Hash().Hex())
}

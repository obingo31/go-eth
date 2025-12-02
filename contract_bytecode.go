package main

import (
	"context"
	"encoding/hex"
	"flag"
	"fmt"
	"log"

	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/ethclient"
)

func main() {
	rpcFlag := flag.String("rpc", "http://127.0.0.1:8545", "Ethereum RPC endpoint")
	addrFlag := flag.String("addr", "", "contract address to inspect")
	flag.Parse()

	if *addrFlag == "" {
		log.Fatal("--addr is required")
	}

	client, err := ethclient.Dial(*rpcFlag)
	if err != nil {
		log.Fatalf("dial RPC: %v", err)
	}
	defer client.Close()

	address := common.HexToAddress(*addrFlag)
	code, err := client.CodeAt(context.Background(), address, nil)
	if err != nil {
		log.Fatalf("fetch bytecode: %v", err)
	}

	fmt.Println(hex.EncodeToString(code))
}

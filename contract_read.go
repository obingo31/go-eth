package main

import (
	"context"
	"flag"
	"fmt"
	"log"

	"github.com/ethereum/go-ethereum/accounts/abi/bind"
	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/ethclient"

	store "github.com/obingo31/go-eth/contracts"
)

func main() {
	rpcFlag := flag.String("rpc", "http://127.0.0.1:8545", "Ethereum RPC endpoint (Ganache, Anvil, etc.)")
	addrFlag := flag.String("addr", "", "Store contract address")
	flag.Parse()

	if *addrFlag == "" {
		log.Fatal("--addr is required (Store contract address)")
	}

	client, err := ethclient.Dial(*rpcFlag)
	if err != nil {
		log.Fatalf("dial RPC: %v", err)
	}
	defer client.Close()

	contractAddr := common.HexToAddress(*addrFlag)
	instance, err := store.NewStore(contractAddr, client)
	if err != nil {
		log.Fatalf("bind Store: %v", err)
	}

	opts := &bind.CallOpts{Context: context.Background()}
	version, err := instance.Version(opts)
	if err != nil {
		log.Fatalf("read version: %v", err)
	}

	fmt.Printf("Store version: %s\n", version)
}

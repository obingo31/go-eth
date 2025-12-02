package main

import (
	"context"
	"crypto/ecdsa"
	"fmt"
	"log"
	"math/big"

	"github.com/ethereum/go-ethereum/accounts/abi/bind"
	"github.com/ethereum/go-ethereum/crypto"
	"github.com/ethereum/go-ethereum/ethclient"

	store "github.com/obingo31/go-eth/contracts"
)

const (
	rpcURL       = "http://127.0.0.1:8545"
	deployerKey  = "4f3edf983ac636a65a842ce7c78d9aa706d3b113bce9c46f30d7d21715b23b1d" // Ganache --deterministic account[0]
	storeVersion = "local-1.0.0"
)

func main() {
	ctx := context.Background()

	client, err := ethclient.Dial(rpcURL)
	if err != nil {
		log.Fatalf("dial RPC: %v", err)
	}
	defer client.Close()

	privateKey, err := crypto.HexToECDSA(deployerKey)
	if err != nil {
		log.Fatalf("parse private key: %v", err)
	}

	publicKey, ok := privateKey.Public().(*ecdsa.PublicKey)
	if !ok {
		log.Fatal("unable to cast public key to ECDSA")
	}
	fromAddress := crypto.PubkeyToAddress(*publicKey)

	chainID, err := client.NetworkID(ctx)
	if err != nil {
		log.Fatalf("fetch chain id: %v", err)
	}

	auth, err := bind.NewKeyedTransactorWithChainID(privateKey, chainID)
	if err != nil {
		log.Fatalf("build transactor: %v", err)
	}

	nonce, err := client.PendingNonceAt(ctx, fromAddress)
	if err != nil {
		log.Fatalf("fetch nonce: %v", err)
	}
	auth.Nonce = big.NewInt(int64(nonce))

	tipCap, err := client.SuggestGasTipCap(ctx)
	supports1559 := err == nil
	if err != nil {
		log.Printf("warn: eth_maxPriorityFeePerGas unavailable, using legacy gas pricing: %v", err)
	}

	baseFee, err := client.SuggestGasPrice(ctx)
	if err != nil {
		log.Fatalf("suggest base fee: %v", err)
	}

	if supports1559 {
		auth.GasTipCap = tipCap
		auth.GasFeeCap = new(big.Int).Add(baseFee, tipCap)
		auth.GasPrice = nil
	} else {
		auth.GasTipCap = nil
		auth.GasFeeCap = nil
		auth.GasPrice = baseFee
	}
	auth.GasLimit = uint64(500000)
	auth.Value = big.NewInt(0)
	auth.Context = ctx

	address, tx, instance, err := store.DeployStore(auth, client, storeVersion)
	if err != nil {
		log.Fatalf("deploy store: %v", err)
	}

	fmt.Printf("Deployer    : %s\n", fromAddress.Hex())
	fmt.Printf("Contract    : %s\n", address.Hex())
	fmt.Printf("Tx Hash     : %s\n", tx.Hash().Hex())

	version, err := instance.Version(&bind.CallOpts{Context: ctx})
	if err != nil {
		log.Fatalf("read version: %v", err)
	}
	fmt.Printf("Version set : %s\n", version)
}

package main

import (
	"context"
	"crypto/ecdsa"
	"flag"
	"fmt"
	"log"
	"math/big"

	"github.com/ethereum/go-ethereum/accounts/abi/bind"
	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/crypto"
	"github.com/ethereum/go-ethereum/ethclient"

	store "github.com/obingo31/go-eth/contracts"
)

func main() {
	rpcFlag := flag.String("rpc", "http://127.0.0.1:8545", "Ethereum RPC endpoint")
	addrFlag := flag.String("addr", "", "Store contract address")
	privFlag := flag.String("priv", "4f3edf983ac636a65a842ce7c78d9aa706d3b113bce9c46f30d7d21715b23b1d", "hex private key (Ganache default account[0])")
	keyFlag := flag.String("key", "foo", "key to store (<=32 bytes)")
	valFlag := flag.String("val", "bar", "value to store (<=32 bytes)")
	flag.Parse()

	if *addrFlag == "" {
		log.Fatal("--addr is required")
	}

	client, err := ethclient.Dial(*rpcFlag)
	if err != nil {
		log.Fatalf("dial RPC: %v", err)
	}
	defer client.Close()

	privateKey, err := crypto.HexToECDSA(*privFlag)
	if err != nil {
		log.Fatalf("parse private key: %v", err)
	}

	instance, auth, err := prepareWriter(context.Background(), client, privateKey, *addrFlag)
	if err != nil {
		log.Fatalf("prepare writer: %v", err)
	}

	key := stringToBytes32(*keyFlag)
	value := stringToBytes32(*valFlag)

	tx, err := instance.SetItem(auth, key, value)
	if err != nil {
		log.Fatalf("set item: %v", err)
	}
	fmt.Printf("tx sent: %s\n", tx.Hash().Hex())

	// Read back to confirm
	stored, err := instance.Items(&bind.CallOpts{Context: auth.Context}, key)
	if err != nil {
		log.Fatalf("read item: %v", err)
	}
	fmt.Printf("stored value: %s\n", bytes32ToString(stored))
}

func prepareWriter(ctx context.Context, client *ethclient.Client, privateKey *ecdsa.PrivateKey, addr string) (*store.Store, *bind.TransactOpts, error) {
	chainID, err := client.NetworkID(ctx)
	if err != nil {
		return nil, nil, fmt.Errorf("fetch chain id: %w", err)
	}

	auth, err := bind.NewKeyedTransactorWithChainID(privateKey, chainID)
	if err != nil {
		return nil, nil, fmt.Errorf("build transactor: %w", err)
	}

	fromAddress := auth.From
	nonce, err := client.PendingNonceAt(ctx, fromAddress)
	if err != nil {
		return nil, nil, fmt.Errorf("fetch nonce: %w", err)
	}
	auth.Nonce = big.NewInt(int64(nonce))

	tipCap, err := client.SuggestGasTipCap(ctx)
	supports1559 := err == nil
	if err != nil {
		log.Printf("warn: eth_maxPriorityFeePerGas unavailable, using legacy gas price: %v", err)
	}

	baseFee, err := client.SuggestGasPrice(ctx)
	if err != nil {
		return nil, nil, fmt.Errorf("suggest gas price: %w", err)
	}

	if supports1559 {
		auth.GasTipCap = tipCap
		auth.GasFeeCap = new(big.Int).Add(baseFee, tipCap)
		auth.GasPrice = nil
	} else {
		auth.GasPrice = baseFee
		auth.GasTipCap = nil
		auth.GasFeeCap = nil
	}

	auth.GasLimit = 500000
	auth.Value = big.NewInt(0)
	auth.Context = ctx

	contractAddr := common.HexToAddress(addr)
	instance, err := store.NewStore(contractAddr, client)
	if err != nil {
		return nil, nil, fmt.Errorf("bind store: %w", err)
	}

	return instance, auth, nil
}

func stringToBytes32(input string) [32]byte {
	var out [32]byte
	copy(out[:], []byte(input))
	return out
}

func bytes32ToString(b [32]byte) string {
	n := 0
	for n < len(b) && b[n] != 0 {
		n++
	}
	return string(b[:n])
}

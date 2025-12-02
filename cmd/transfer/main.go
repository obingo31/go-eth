package main

import (
	"context"
	"crypto/ecdsa"
	"flag"
	"fmt"
	"log"
	"math/big"
	"strings"
	"time"

	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/core/types"
	"github.com/ethereum/go-ethereum/crypto"
	"github.com/ethereum/go-ethereum/ethclient"
)

func main() {
	rpcFlag := flag.String("rpc", "https://sepolia.infura.io/v3/b3ce18e518ab499cb2975c8952bb0a47", "Ethereum RPC endpoint")
	privFlag := flag.String("priv", "", "hex-encoded private key (without passphrase)")
	toFlag := flag.String("to", "", "destination address")
	valueFlag := flag.String("value", "1000000000000000", "amount to send in wei (default 0.001 ETH)")
	gasLimitFlag := flag.Uint64("gas-limit", 21000, "gas limit for the transfer")
	flag.Parse()

	if *privFlag == "" || *toFlag == "" {
		log.Fatal("both --priv and --to are required")
	}

	value := new(big.Int)
	if _, ok := value.SetString(*valueFlag, 10); !ok {
		log.Fatalf("invalid value %q", *valueFlag)
	}

	privateKey, err := crypto.HexToECDSA(strings.TrimPrefix(*privFlag, "0x"))
	if err != nil {
		log.Fatalf("invalid private key: %v", err)
	}

	client, err := ethclient.Dial(*rpcFlag)
	if err != nil {
		log.Fatalf("dial RPC: %v", err)
	}
	defer client.Close()

	ctx, cancel := context.WithTimeout(context.Background(), 15*time.Second)
	defer cancel()

	publicKey := privateKey.Public()
	publicKeyECDSA, ok := publicKey.(*ecdsa.PublicKey)
	if !ok {
		log.Fatal("error casting public key to ECDSA")
	}

	fromAddress := crypto.PubkeyToAddress(*publicKeyECDSA)
	nonce, err := client.PendingNonceAt(ctx, fromAddress)
	if err != nil {
		log.Fatalf("fetch nonce: %v", err)
	}

	gasPrice, err := client.SuggestGasPrice(ctx)
	if err != nil {
		log.Fatalf("suggest gas price: %v", err)
	}

	toAddress := common.HexToAddress(*toFlag)
	tx := types.NewTransaction(nonce, toAddress, value, *gasLimitFlag, gasPrice, nil)

	chainID, err := client.NetworkID(ctx)
	if err != nil {
		log.Fatalf("fetch chain ID: %v", err)
	}

	signer := types.LatestSignerForChainID(chainID)
	signedTx, err := types.SignTx(tx, signer, privateKey)
	if err != nil {
		log.Fatalf("sign tx: %v", err)
	}

	if err := client.SendTransaction(ctx, signedTx); err != nil {
		log.Fatalf("send tx: %v", err)
	}

	fmt.Printf("tx sent: %s\n", signedTx.Hash().Hex())
}

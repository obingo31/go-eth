package main

import (
	"context"
	"crypto/ecdsa"
	"flag"
	"fmt"
	"log"
	"math/big"
	"strings"

	"github.com/ethereum/go-ethereum"
	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/common/hexutil"
	"github.com/ethereum/go-ethereum/core/types"
	"github.com/ethereum/go-ethereum/crypto"
	"github.com/ethereum/go-ethereum/ethclient"
	"golang.org/x/crypto/sha3"
)

func main() {
	rpcFlag := flag.String("rpc", "https://sepolia.infura.io/v3/b3ce18e518ab499cb2975c8952bb0a47", "Ethereum RPC endpoint")
	privFlag := flag.String("priv", "", "hex-encoded private key (no passphrase)")
	contractFlag := flag.String("contract", "0xF4D17Dd253A5a21555bF1a5B9B7285ed764AF706", "ERC-20 token contract address")
	toFlag := flag.String("to", "0x5bb34D0bf5DC32df87Ae454DEb17001F808b986b", "recipient address")
	amountFlag := flag.String("amount", "1000000000000000000000", "token amount in the smallest unit (wei-style)")
	flag.Parse()

	if *privFlag == "" {
		log.Fatal("--priv is required")
	}
	if !common.IsHexAddress(*contractFlag) || !common.IsHexAddress(*toFlag) {
		log.Fatal("invalid --contract or --to address")
	}

	privateKey, err := crypto.HexToECDSA(strings.TrimPrefix(*privFlag, "0x"))
	if err != nil {
		log.Fatalf("parse private key: %v", err)
	}

	amount := new(big.Int)
	if _, ok := amount.SetString(*amountFlag, 10); !ok {
		log.Fatalf("invalid amount: %s", *amountFlag)
	}

	tokenAddress := common.HexToAddress(*contractFlag)
	recipient := common.HexToAddress(*toFlag)

	transferFnSignature := []byte("transfer(address,uint256)")
	hash := sha3.NewLegacyKeccak256()
	hash.Write(transferFnSignature)
	methodID := hash.Sum(nil)[:4]
	fmt.Printf("Method ID      : %s\n", hexutil.Encode(methodID))

	paddedAddress := common.LeftPadBytes(recipient.Bytes(), 32)
	fmt.Printf("Padded address : %s\n", hexutil.Encode(paddedAddress))

	paddedAmount := common.LeftPadBytes(amount.Bytes(), 32)
	fmt.Printf("Padded amount  : %s\n", hexutil.Encode(paddedAmount))

	var data []byte
	data = append(data, methodID...)
	data = append(data, paddedAddress...)
	data = append(data, paddedAmount...)
	fmt.Printf("Calldata       : %s\n", hexutil.Encode(data))

	client, err := ethclient.Dial(*rpcFlag)
	if err != nil {
		log.Fatalf("dial RPC: %v", err)
	}
	defer client.Close()

	publicKey, ok := privateKey.Public().(*ecdsa.PublicKey)
	if !ok {
		log.Fatal("unable to cast public key to ECDSA")
	}

	fromAddress := crypto.PubkeyToAddress(*publicKey)
	nonce, err := client.PendingNonceAt(context.Background(), fromAddress)
	if err != nil {
		log.Fatalf("fetch nonce: %v", err)
	}

	value := big.NewInt(0)
	gasPrice, err := client.SuggestGasPrice(context.Background())
	if err != nil {
		log.Fatalf("suggest gas price: %v", err)
	}

	gasLimit, err := client.EstimateGas(context.Background(), ethereum.CallMsg{
		From: fromAddress,
		To:   &tokenAddress,
		Data: data,
	})
	if err != nil {
		log.Fatalf("estimate gas: %v", err)
	}
	fmt.Printf("Gas limit      : %d\n", gasLimit)

	tx := types.NewTransaction(nonce, tokenAddress, value, gasLimit, gasPrice, data)

	chainID, err := client.NetworkID(context.Background())
	if err != nil {
		log.Fatalf("fetch chain ID: %v", err)
	}

	signedTx, err := types.SignTx(tx, types.NewEIP155Signer(chainID), privateKey)
	if err != nil {
		log.Fatalf("sign tx: %v", err)
	}

	if err := client.SendTransaction(context.Background(), signedTx); err != nil {
		log.Fatalf("send tx: %v", err)
	}

	fmt.Printf("tx sent: %s\n", signedTx.Hash().Hex())
}

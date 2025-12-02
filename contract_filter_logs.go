package main

import (
	"context"
	"flag"
	"fmt"
	"log"
	"math/big"
	"strings"

	"github.com/ethereum/go-ethereum"
	"github.com/ethereum/go-ethereum/accounts/abi"
	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/crypto"
	"github.com/ethereum/go-ethereum/ethclient"
	token "github.com/obingo31/go-eth/token"
)

type LogTransfer struct {
	From   common.Address
	To     common.Address
	Tokens *big.Int `abi:"value"`
}

type LogApproval struct {
	TokenOwner common.Address
	Spender    common.Address
	Tokens     *big.Int `abi:"value"`
}

func main() {
	rpcFlag := flag.String("rpc", "https://mainnet.infura.io/v3/b3ce18e518ab499cb2975c8952bb0a47", "Ethereum RPC endpoint")
	fromFlag := flag.String("from", "6383820", "start block number")
	toFlag := flag.String("to", "6383840", "end block number (inclusive)")
	addrFlag := flag.String("addr", "0xe41d2489571d322189246dafa5ebde1f4699f498", "ERC-20 contract address")
	flag.Parse()

	fromBlock, ok := new(big.Int).SetString(*fromFlag, 10)
	if !ok {
		log.Fatalf("invalid --from block: %s", *fromFlag)
	}
	toBlock, ok := new(big.Int).SetString(*toFlag, 10)
	if !ok {
		log.Fatalf("invalid --to block: %s", *toFlag)
	}

	client, err := ethclient.Dial(*rpcFlag)
	if err != nil {
		log.Fatalf("dial RPC: %v", err)
	}
	defer client.Close()

	contractAddress := common.HexToAddress(*addrFlag)
	query := ethereum.FilterQuery{
		FromBlock: fromBlock,
		ToBlock:   toBlock,
		Addresses: []common.Address{contractAddress},
	}

	logs, err := client.FilterLogs(context.Background(), query)
	if err != nil {
		log.Fatalf("filter logs: %v", err)
	}

	contractABI, err := abi.JSON(strings.NewReader(token.TokenABI))
	if err != nil {
		log.Fatalf("parse token ABI: %v", err)
	}

	transferSig := crypto.Keccak256Hash([]byte("Transfer(address,address,uint256)"))
	approvalSig := crypto.Keccak256Hash([]byte("Approval(address,address,uint256)"))

	for _, vLog := range logs {
		fmt.Printf("Log Block Number: %d\n", vLog.BlockNumber)
		fmt.Printf("Log Index: %d\n", vLog.Index)

		switch vLog.Topics[0].Hex() {
		case transferSig.Hex():
			fmt.Println("Log Name: Transfer")
			var transferEvent LogTransfer
			if err := contractABI.UnpackIntoInterface(&transferEvent, "Transfer", vLog.Data); err != nil {
				log.Fatalf("unpack transfer: %v", err)
			}
			transferEvent.From = common.HexToAddress(vLog.Topics[1].Hex())
			transferEvent.To = common.HexToAddress(vLog.Topics[2].Hex())
			fmt.Printf("From: %s\n", transferEvent.From.Hex())
			fmt.Printf("To: %s\n", transferEvent.To.Hex())
			fmt.Printf("Tokens: %s\n", transferEvent.Tokens.String())
		case approvalSig.Hex():
			fmt.Println("Log Name: Approval")
			var approvalEvent LogApproval
			if err := contractABI.UnpackIntoInterface(&approvalEvent, "Approval", vLog.Data); err != nil {
				log.Fatalf("unpack approval: %v", err)
			}
			approvalEvent.TokenOwner = common.HexToAddress(vLog.Topics[1].Hex())
			approvalEvent.Spender = common.HexToAddress(vLog.Topics[2].Hex())
			fmt.Printf("Token Owner: %s\n", approvalEvent.TokenOwner.Hex())
			fmt.Printf("Spender: %s\n", approvalEvent.Spender.Hex())
			fmt.Printf("Tokens: %s\n", approvalEvent.Tokens.String())
		default:
			fmt.Println("Log Name: Unknown event signature")
		}
		fmt.Println()
	}
}

package main

import (
	"flag"
	"fmt"
	"log"

	"github.com/ethereum/go-ethereum/common/hexutil"
	"github.com/ethereum/go-ethereum/crypto"
)

func main() {
	msgFlag := flag.String("msg", "hello", "message to hash and sign")
	privFlag := flag.String("priv", "4f3edf983ac636a65a842ce7c78d9aa706d3b113bce9c46f30d7d21715b23b1d", "hex private key for signing")
	flag.Parse()

	privateKey, err := crypto.HexToECDSA(*privFlag)
	if err != nil {
		log.Fatalf("parse private key: %v", err)
	}

	hash := crypto.Keccak256Hash([]byte(*msgFlag))
	fmt.Printf("hash: %s\n", hash.Hex())

	signature, err := crypto.Sign(hash.Bytes(), privateKey)
	if err != nil {
		log.Fatalf("sign message: %v", err)
	}

	fmt.Printf("signature: %s\n", hexutil.Encode(signature))
}

package main

import (
	"bytes"
	"crypto/ecdsa"
	"flag"
	"fmt"
	"log"

	"github.com/ethereum/go-ethereum/common/hexutil"
	"github.com/ethereum/go-ethereum/crypto"
)

func main() {
	msgFlag := flag.String("msg", "hello", "message to hash and sign")
	privFlag := flag.String("priv", "4f3edf983ac636a65a842ce7c78d9aa706d3b113bce9c46f30d7d21715b23b1d", "hex private key used for signing")
	flag.Parse()

	privateKey, err := crypto.HexToECDSA(*privFlag)
	if err != nil {
		log.Fatalf("parse private key: %v", err)
	}

	publicKey, ok := privateKey.Public().(*ecdsa.PublicKey)
	if !ok {
		log.Fatal("public key type assertion failed")
	}
	publicKeyBytes := crypto.FromECDSAPub(publicKey)

	hash := crypto.Keccak256Hash([]byte(*msgFlag))
	fmt.Printf("hash: %s\n", hash.Hex())

	signature, err := crypto.Sign(hash.Bytes(), privateKey)
	if err != nil {
		log.Fatalf("sign message: %v", err)
	}
	fmt.Printf("signature: %s\n", hexutil.Encode(signature))

	sigPublicKey, err := crypto.Ecrecover(hash.Bytes(), signature)
	if err != nil {
		log.Fatalf("ecrecover: %v", err)
	}
	fmt.Printf("Ecrecover matches: %t\n", bytes.Equal(sigPublicKey, publicKeyBytes))

	sigPublicKeyECDSA, err := crypto.SigToPub(hash.Bytes(), signature)
	if err != nil {
		log.Fatalf("SigToPub: %v", err)
	}
	recoveredBytes := crypto.FromECDSAPub(sigPublicKeyECDSA)
	fmt.Printf("SigToPub matches: %t\n", bytes.Equal(recoveredBytes, publicKeyBytes))

	signatureNoRecoverID := signature[:len(signature)-1]
	verified := crypto.VerifySignature(publicKeyBytes, hash.Bytes(), signatureNoRecoverID)
	fmt.Printf("VerifySignature: %t\n", verified)
}

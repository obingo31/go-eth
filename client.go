package main

import (
	"fmt"
	"log"

	"github.com/ethereum/go-ethereum/ethclient"
)

func main() {
	const infuraURL = "https://mainnet.infura.io/v3/<PROJECT_ID>" // replace with your Infura project endpoint

	client, err := ethclient.Dial(infuraURL)
	if err != nil {
		log.Fatal(err)
	}

	fmt.Println("we have a connection")
	_ = client // placeholder for later use
}

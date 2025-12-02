package main

import (
	"context"
	"fmt"
	"log"
	"os"
	"runtime"

	"github.com/ethereum/go-ethereum/common/hexutil"
	"github.com/ethereum/go-ethereum/whisper/shhclient"
	"github.com/ethereum/go-ethereum/whisper/whisperv6"
)

func main() {
	client, err := shhclient.Dial("ws://127.0.0.1:8546")
	if err != nil {
		log.Fatal(err)
	}

	keyID, err := client.NewKeyPair(context.Background())
	if err != nil {
		log.Fatal(err)
	}
	fmt.Println(keyID)

	messages := make(chan *whisperv6.Message)
	criteria := whisperv6.Criteria{
		PrivateKeyID: keyID,
	}
	sub, err := client.SubscribeMessages(context.Background(), criteria, messages)
	if err != nil {
		log.Fatal(err)
	}

	go func() {
		for {
			select {
			case err := <-sub.Err():
				log.Fatal(err)
			case message := <-messages:
				fmt.Printf(string(message.Payload))
				os.Exit(0)
			}
		}
	}()

	publicKey, err := client.PublicKey(context.Background(), keyID)
	if err != nil {
		log.Print(err)
	}
	fmt.Println(hexutil.Encode(publicKey))

	message := whisperv6.NewMessage{
		Payload:   []byte("Hello"),
		PublicKey: publicKey,
		TTL:       60,
		PowTime:   2,
		PowTarget: 2.5,
	}

	messageHash, err := client.Post(context.Background(), message)
	if err != nil {
		log.Fatal(err)
	}
	fmt.Println(messageHash)

	runtime.Goexit()
}

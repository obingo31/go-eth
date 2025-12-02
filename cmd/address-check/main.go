package main

import (
	"context"
	"flag"
	"fmt"
	"log"
	"regexp"
	"time"

	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/ethclient"
)

type preset struct {
	contract string
	account  string
}

var presets = map[string]preset{
	"mainnet": {
		contract: "0xe41d2489571d322189246dafa5ebde1f4699f498", // 0x Protocol Token (ZRX)
		account:  "0x407d73d8a49eeb85d32cf465507dd71d507100c1",
	},
	"ganache": {
		contract: "0xF4D17Dd253A5a21555bF1a5B9B7285ed764AF706", // Ganache account 0
		account:  "0x5bb34D0bf5DC32df87Ae454DEb17001F808b986b", // Ganache account 1
	},
}

func main() {
	rpcFlag := flag.String("rpc", "http://localhost:8545", "Ethereum RPC endpoint to inspect")
	contractFlag := flag.String("contract", "", "address expected to be a contract (defaults depend on --preset)")
	accountFlag := flag.String("account", "", "address expected to be an EOA (defaults depend on --preset)")
	presetFlag := flag.String("preset", "mainnet", "address preset to use: mainnet or ganache")
	flag.Parse()

	defaults, ok := presets[*presetFlag]
	if !ok {
		log.Fatalf("unknown preset %q (expected %v)", *presetFlag, presetNames())
	}

	if *contractFlag == "" {
		*contractFlag = defaults.contract
	}

	if *accountFlag == "" {
		*accountFlag = defaults.account
	}

	re := regexp.MustCompile(`^0x[0-9a-fA-F]{40}$`)

	fmt.Printf("is valid (%s): %v\n", *contractFlag, re.MatchString(*contractFlag))
	fmt.Printf("is valid (%s): %v\n", *accountFlag, re.MatchString(*accountFlag))

	client, err := ethclient.Dial(*rpcFlag)
	if err != nil {
		log.Fatalf("dial RPC: %v", err)
	}
	defer client.Close()

	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	contractAddr := common.HexToAddress(*contractFlag)
	isContract, err := hasCode(ctx, client, contractAddr)
	if err != nil {
		log.Fatalf("fetch contract bytecode: %v", err)
	}
	fmt.Printf("is contract (%s): %v\n", contractAddr.Hex(), isContract)

	accountAddr := common.HexToAddress(*accountFlag)
	isAccountContract, err := hasCode(ctx, client, accountAddr)
	if err != nil {
		log.Fatalf("fetch account bytecode: %v", err)
	}
	fmt.Printf("is contract (%s): %v\n", accountAddr.Hex(), isAccountContract)
}

func hasCode(ctx context.Context, client *ethclient.Client, addr common.Address) (bool, error) {
	bytecode, err := client.CodeAt(ctx, addr, nil)
	if err != nil {
		return false, err
	}
	return len(bytecode) > 0, nil
}

func presetNames() []string {
	names := make([]string, 0, len(presets))
	for name := range presets {
		names = append(names, name)
	}
	return names
}

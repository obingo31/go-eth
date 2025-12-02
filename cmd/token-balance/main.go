package main

import (
	"context"
	"flag"
	"fmt"
	"log"
	"math"
	"math/big"
	"strings"

	"github.com/ethereum/go-ethereum/accounts/abi/bind"
	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/ethclient"

	"github.com/obingo31/go-eth/token"
)

func main() {
	rpcFlag := flag.String("rpc", "https://mainnet.infura.io/v3/YOUR_KEY", "Ethereum RPC endpoint")
	contractFlag := flag.String("contract", "", "ERC-20 token contract address")
	accountFlag := flag.String("account", "", "Address whose balance should be fetched")
	flag.Parse()

	if !common.IsHexAddress(*contractFlag) {
		log.Fatal("--contract must be a valid hex address")
	}
	if !common.IsHexAddress(*accountFlag) {
		log.Fatal("--account must be a valid hex address")
	}

	rpcURL := strings.TrimSpace(*rpcFlag)
	if rpcURL == "" {
		log.Fatal("--rpc must not be empty")
	}

	client, err := ethclient.Dial(rpcURL)
	if err != nil {
		log.Fatalf("dial RPC: %v", err)
	}
	defer client.Close()

	tokenAddress := common.HexToAddress(*contractFlag)
	instance, err := token.NewToken(tokenAddress, client)
	if err != nil {
		log.Fatalf("instantiate token binding: %v", err)
	}

	callOpts := &bind.CallOpts{Context: context.Background()}

	accountAddress := common.HexToAddress(*accountFlag)
	balance, err := instance.BalanceOf(callOpts, accountAddress)
	if err != nil {
		log.Fatalf("fetch balance: %v", err)
	}

	name, err := instance.Name(callOpts)
	if err != nil {
		log.Fatalf("fetch name: %v", err)
	}

	symbol, err := instance.Symbol(callOpts)
	if err != nil {
		log.Fatalf("fetch symbol: %v", err)
	}

	decimals, err := instance.Decimals(callOpts)
	if err != nil {
		log.Fatalf("fetch decimals: %v", err)
	}

	fmt.Printf("name: %s\n", name)
	fmt.Printf("symbol: %s\n", symbol)
	fmt.Printf("decimals: %d\n", decimals)
	fmt.Printf("wei balance: %s\n", balance)

	balanceFloat := toDecimal(balance, decimals)
	fmt.Printf("balance: %f\n", balanceFloat)
}

func toDecimal(value *big.Int, decimals uint8) *big.Float {
	if value == nil {
		return big.NewFloat(0)
	}
	fVal := new(big.Float).SetInt(value)
	denom := new(big.Float).SetFloat64(math.Pow10(int(decimals)))
	if denom.Cmp(big.NewFloat(0)) == 0 {
		return fVal
	}
	return new(big.Float).Quo(fVal, denom)
}

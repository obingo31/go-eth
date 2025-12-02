package main

import (
	"flag"
	"fmt"
	"log"

	hdwallet "github.com/miguelmota/go-ethereum-hdwallet"
)

func main() {
	mnemonicFlag := flag.String("mnemonic", "room city turn fan foam nuclear wrestle general brother cool know lab", "BIP39 mnemonic to derive from")
	pathFlag := flag.String("path", "m/44'/60'/0'/0", "base derivation path (without trailing index)")
	countFlag := flag.Int("n", 5, "number of addresses to derive starting from index 0")
	showPriv := flag.Bool("show-priv", false, "display private keys (only for testing!)")
	flag.Parse()

	wallet, err := hdwallet.NewFromMnemonic(*mnemonicFlag)
	if err != nil {
		log.Fatalf("create wallet from mnemonic: %v", err)
	}

	for i := 0; i < *countFlag; i++ {
		path := fmt.Sprintf("%s/%d", *pathFlag, i)
		derivationPath, err := hdwallet.ParseDerivationPath(path)
		if err != nil {
			log.Fatalf("parse derivation path %s: %v", path, err)
		}

		account, err := wallet.Derive(derivationPath, false)
		if err != nil {
			log.Fatalf("derive account %s: %v", path, err)
		}

		fmt.Printf("%s -> %s\n", path, account.Address.Hex())
		if *showPriv {
			privHex, err := wallet.PrivateKeyHex(account)
			if err != nil {
				log.Fatalf("export privkey for %s: %v", account.Address.Hex(), err)
			}
			fmt.Printf("    priv: %s\n", privHex)
		}
	}
}

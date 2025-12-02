package main

import (
    "flag"
    "fmt"
    "log"
    "os"
    "path/filepath"

    "github.com/ethereum/go-ethereum/accounts/keystore"
)

func main() {
    dirFlag := flag.String("dir", "./keystore", "directory to save keystore files")
    passwordFlag := flag.String("password", "", "password used to encrypt the keystore (required)")
    countFlag := flag.Int("n", 1, "number of accounts to generate")
    lightFlag := flag.Bool("light", false, "use light scrypt parameters (faster but weaker)")
    flag.Parse()

    if *passwordFlag == "" {
        log.Fatal("password is required")
    }

    if *countFlag < 1 {
        log.Fatal("n must be >= 1")
    }

    if err := os.MkdirAll(*dirFlag, 0o700); err != nil {
        log.Fatalf("failed to create keystore directory: %v", err)
    }

    scryptN := keystore.StandardScryptN
    scryptP := keystore.StandardScryptP
    if *lightFlag {
        scryptN = keystore.LightScryptN
        scryptP = keystore.LightScryptP
    }

    ks := keystore.NewKeyStore(*dirFlag, scryptN, scryptP)

    for i := 0; i < *countFlag; i++ {
        account, err := ks.NewAccount(*passwordFlag)
        if err != nil {
            log.Fatalf("failed to create account %d: %v", i+1, err)
        }

        filename := filepath.Base(account.URL.Path)
        fmt.Printf("Account %d: %s\n", i+1, account.Address.Hex())
        fmt.Printf("Keystore file: %s\n", filepath.Join(*dirFlag, filename))
    }
}

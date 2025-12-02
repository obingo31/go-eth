package main

import (
	"flag"
	"fmt"
	"log"
	"os"

	"github.com/ethereum/go-ethereum/accounts/keystore"
)

func main() {
	mode := flag.String("mode", "create", "operation to perform: create or import")
	dir := flag.String("dir", "./keystore", "target directory for keystore files")
	password := flag.String("password", "", "password used to encrypt/decrypt the keystore (required)")
	importPath := flag.String("import-file", "", "path to an existing keystore JSON (required for import mode)")
	light := flag.Bool("light", false, "use light scrypt parameters (faster but weaker)")
	flag.Parse()

	if *password == "" {
		log.Fatal("password is required")
	}

	if err := os.MkdirAll(*dir, 0o700); err != nil {
		log.Fatalf("failed to create keystore directory: %v", err)
	}

	scryptN := keystore.StandardScryptN
	scryptP := keystore.StandardScryptP
	if *light {
		scryptN = keystore.LightScryptN
		scryptP = keystore.LightScryptP
	}

	ks := keystore.NewKeyStore(*dir, scryptN, scryptP)

	switch *mode {
	case "create":
		createAccount(ks, *password)
	case "import":
		if *importPath == "" {
			log.Fatal("import mode requires --import-file pointing to a keystore JSON")
		}
		importAccount(ks, *importPath, *password)
	default:
		log.Fatalf("unknown mode %q (expected 'create' or 'import')", *mode)
	}
}

func createAccount(ks *keystore.KeyStore, password string) {
	account, err := ks.NewAccount(password)
	if err != nil {
		log.Fatalf("failed to create account: %v", err)
	}

	fmt.Printf("Created account: %s\n", account.Address.Hex())
	fmt.Printf("Keystore file: %s\n", account.URL.Path)
}

func importAccount(ks *keystore.KeyStore, filePath, password string) {
	jsonBytes, err := os.ReadFile(filePath)
	if err != nil {
		log.Fatalf("failed to read keystore file %s: %v", filePath, err)
	}

	account, err := ks.Import(jsonBytes, password, password)
	if err != nil {
		log.Fatalf("failed to import account: %v", err)
	}

	fmt.Printf("Imported account: %s\n", account.Address.Hex())
	fmt.Printf("Stored in: %s\n", account.URL.Path)
}

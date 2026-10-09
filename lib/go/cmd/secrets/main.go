// Command secrets seeds a Backend's secret store from dotenv files — the
// command behind `just inject <svc>`:
//
//	SECRETS_PROVIDER=aws-secretsmanager SECRETS_ID=trieoh/identityx \
//	  go run lib/cmd/secrets inject .env api/identityx/.env
//
// Later files win over earlier ones. The store is chosen exactly as the
// Backend chooses it at boot (lib/secrets/provider), so what is injected is
// what is read.
package main

import (
	"context"
	"errors"
	"fmt"
	"maps"
	"os"
	"time"

	"lib/secrets"
	"lib/secrets/provider"
)

func main() {
	err := run(os.Args[1:])
	if err != nil {
		fmt.Fprintln(os.Stderr, "secrets:", err)
		os.Exit(1)
	}
}

func run(args []string) error {
	if len(args) < 2 || args[0] != "inject" {
		return errors.New("usage: secrets inject <dotenv file>")
	}
	values := make(map[string]string)
	for _, path := range args[1:] {
		file, err := secrets.ParseDotenv(path)
		if err != nil {
			return fmt.Errorf("read %s: %w", path, err)
		}
		maps.Copy(values, file)
	}

	ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
	defer cancel()
	store, err := provider.StoreFromEnv(ctx)
	if err != nil {
		return err
	}
	err = store.Put(ctx, values)
	if err != nil {
		return err
	}
	fmt.Printf("injected %d values into %s %s\n", len(values), os.Getenv("SECRETS_PROVIDER"), os.Getenv("SECRETS_ID"))
	return nil
}

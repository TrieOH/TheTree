// Package provider picks the lib/secrets adapter from the environment. It is
// the only place that knows every adapter, so lib/secrets itself stays free
// of vendor SDKs.
//
//	SECRETS_PROVIDER  env (default) | aws-secretsmanager
//	SECRETS_ID        the secret's name/ARN (required for aws-secretsmanager)
package provider

import (
	"context"
	"fmt"
	"log"
	"os"

	"lib/secrets"
	"lib/secrets/awssm"
)

const (
	Env               = "env"
	AWSSecretsManager = "aws-secretsmanager"
)

// FromEnv returns the configured Source.
func FromEnv(ctx context.Context) (secrets.Source, error) { //nolint:ireturn
	switch name := os.Getenv("SECRETS_PROVIDER"); name {
	case "", Env:
		return secrets.Env{}, nil
	case AWSSecretsManager:
		return awssm.New(ctx, os.Getenv("SECRETS_ID"))
	default:
		return nil, fmt.Errorf("unknown SECRETS_PROVIDER %q (want %q or %q)", name, Env, AWSSecretsManager)
	}
}

// StoreFromEnv returns the configured Source as a writable Store, for the
// inject command. The env provider has nothing to write to.
func StoreFromEnv(ctx context.Context) (secrets.Store, error) { //nolint:ireturn
	src, err := FromEnv(ctx)
	if err != nil {
		return nil, err
	}
	store, ok := src.(secrets.Store)
	if !ok {
		return nil, fmt.Errorf("SECRETS_PROVIDER %q is not writable; set it to %q to inject", os.Getenv("SECRETS_PROVIDER"), AWSSecretsManager)
	}
	return store, nil
}

// Load applies the configured Source to the process environment. Call it
// before parsing config.
func Load(ctx context.Context) error {
	src, err := FromEnv(ctx)
	if err != nil {
		return err
	}
	applied, err := secrets.Apply(ctx, src)
	if err != nil {
		return err
	}
	if len(applied) > 0 {
		log.Printf("secrets: exported %d values from the secret store", len(applied))
	}
	return nil
}

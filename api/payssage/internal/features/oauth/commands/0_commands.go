package commands

import (
	"crypto/rand"
	"encoding/hex"
	"payssage/ports"
)

type Commands struct {
	wallets    ports.WalletRepo
	orgs       ports.OrganizationRepo
	oauth      ports.OAuthStateRepo
	collectors ports.CollectorRepo
	sellers    ports.SellerRepo
}

func NewCommands(
	wallets ports.WalletRepo,
	orgs ports.OrganizationRepo,
	oauth ports.OAuthStateRepo,
	collectors ports.CollectorRepo,
	sellers ports.SellerRepo,
) *Commands {
	return &Commands{
		wallets:    wallets,
		orgs:       orgs,
		oauth:      oauth,
		collectors: collectors,
		sellers:    sellers,
	}
}

func generateState() (string, error) {
	b := make([]byte, 32)
	_, err := rand.Read(b)
	if err != nil {
		return "", err
	}
	return hex.EncodeToString(b), nil
}

package oauth_handler

import (
	"TriePayments/internal/core/application/oauth/commands"
)

type Handler struct {
	commands *commands.CommandService
}

func NewOAuthHandler(
	commands *commands.CommandService,
) *Handler {
	return &Handler{
		commands: commands,
	}
}

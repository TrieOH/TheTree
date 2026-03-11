package oauth_handler

import (
	"TriePayments/internal/core/application/oauth/commands"
	"TriePayments/internal/core/application/oauth/queries"
)

type Handler struct {
	commands *commands.CommandService
	queries  *queries.QueryService
}

func NewOAuthHandler(
	commands *commands.CommandService,
	queries *queries.QueryService,
) *Handler {
	return &Handler{
		commands: commands,
		queries:  queries,
	}
}

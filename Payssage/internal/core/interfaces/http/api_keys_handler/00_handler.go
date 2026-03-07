package workspaces_handler

import (
	"TriePayments/internal/core/application/api_keys/commands"
	"TriePayments/internal/core/application/api_keys/queries"
)

type Handler struct {
	commands *commands.CommandService
	queries  *queries.QueryService
}

func NewApiKeysHandler(
	commands *commands.CommandService,
	queries *queries.QueryService,
) *Handler {
	return &Handler{
		commands: commands,
		queries:  queries,
	}
}

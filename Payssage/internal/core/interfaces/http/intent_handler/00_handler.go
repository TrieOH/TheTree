package workspaces_handler

import (
	"TriePayments/internal/core/application/intents/commands"
	"TriePayments/internal/core/application/intents/queries"
)

type Handler struct {
	commands *commands.CommandService
	queries  *queries.QueryService
}

func NewIntentsHandler(
	commands *commands.CommandService,
	queries *queries.QueryService,
) *Handler {
	return &Handler{
		commands: commands,
		queries:  queries,
	}
}

package handlers

import (
	"net/http"

	"github.com/MintzyG/fun"
)

func (handler *Handlers) List(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()
	out, err := handler.queries.List(ctx)
	if err != nil {
		fun.Error(err).Send(w)
		return
	}

	fun.OK().WithData(out).Send(w)
}

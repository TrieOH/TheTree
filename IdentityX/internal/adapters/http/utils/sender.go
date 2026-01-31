package utils

import (
	"GoAuth/internal/apierr"
	"net/http"

	resp "github.com/MintzyG/FastUtilitiesNet/response"
	"github.com/MintzyG/fail"
)

func Sender(e *fail.Error, module string, w http.ResponseWriter) (*resp.Response, bool) {
	trrs, err := apierr.HTTPResponseTranslator().Translate(e)
	if err != nil {
		resp.InternalServerError().WithData(err).WithModule(module).Send(w)
		return nil, false
	}
	if rs, ok := trrs.(*resp.Response); ok {
		return rs.WithModule(module), true
	} else {
		resp.InternalServerError("couldn't cast to response").WithModule(module).Send(w)
		return nil, false
	}
}

package events_test

import (
	"os"
	"testing"

	"github.com/MintzyG/fun"
)

func TestMain(m *testing.M) {
	fun.SetConfig(fun.Config{
		DefaultModule:        "test",
		DefaultContentType:   "application/json",
		EnableSizeValidation: true,
	})
	os.Exit(m.Run())
}

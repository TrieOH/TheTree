package database

import (
	"errors"
	"strconv"
	"testing"

	"github.com/jackc/pgx/v5"
)

var errThingNotFound = errors.New("thing not found")

func TestOne_MapsTheRow(t *testing.T) {
	got, err := One(7, nil, NewErrorHandler(errThingNotFound), strconv.Itoa)
	if err != nil || got == nil || *got != "7" {
		t.Fatalf("One = %v, %v; want \"7\", nil", got, err)
	}
}

func TestOne_NamesTheError(t *testing.T) {
	got, err := One(0, pgx.ErrNoRows, NewErrorHandler(errThingNotFound), strconv.Itoa)
	if got != nil || !errors.Is(err, errThingNotFound) || !errors.Is(err, ErrNotFound) {
		t.Fatalf("One = %v, %v; want nil, the named not-found error", got, err)
	}
}

func TestMany_MapsTheRows(t *testing.T) {
	got, err := Many([]int{1, 2}, nil, NewErrorHandler(errThingNotFound), strconv.Itoa)
	if err != nil || len(got) != 2 || got[0] != "1" || got[1] != "2" {
		t.Fatalf("Many = %v, %v; want [1 2], nil", got, err)
	}
}

func TestMany_PassesInternalErrorsThrough(t *testing.T) {
	boom := errors.New("boom")
	got, err := Many([]int(nil), boom, NewErrorHandler(errThingNotFound), strconv.Itoa)
	if got != nil || !errors.Is(err, boom) {
		t.Fatalf("Many = %v, %v; want nil, boom", got, err)
	}
}

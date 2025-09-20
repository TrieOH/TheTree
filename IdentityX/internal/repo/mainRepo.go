package handler

import (
	resp "github.com/MintzyG/GoResponse/response"
	"net/http"
)

type GoAuthRepo struct {
	Name string
}

func NewGoAuthRepo(name string) *GoAuthRepo {
	return &GoAuthRepo{Name: name}
}

func (h *GoAuthRepo) Hi() {
}

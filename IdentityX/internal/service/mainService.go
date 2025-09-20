package service

import (
	"strings"
)

type GoAuthService struct {
	Name string
}

func NewGoAuthService(name string) *GoAuthService {
	return &GoAuthService{Name: name}
}

func (h *GoAuthService) Hi(name string) string {
	name = strings.ToLower(name)
	words := strings.Fields(name)
	for i, w := range words {
		if len(w) > 0 {
			words[i] = strings.ToUpper(string(w[0])) + w[1:]
		}
	}
	return strings.Join(words, " ")
}

package models

import "github.com/google/uuid"

type Signature struct {
	ID        uuid.UUID `json:"id"`
	EditionID uuid.UUID `json:"edition_id"`
	Title     string    `json:"title"`
	URL       string    `json:"url"`
}

type AddSignatureRequest struct {
	Title string `json:"title"`
	URL   string `json:"url"`
}

func (r AddSignatureRequest) ToInput(editionID uuid.UUID) AddSignatureInput {
	return AddSignatureInput{
		Title:     r.Title,
		URL:       r.URL,
		EditionID: editionID,
	}
}

type AddSignatureInput struct {
	EditionID uuid.UUID `json:"edition_id"`
	Title     string    `json:"title"`
	URL       string    `json:"url"`
}

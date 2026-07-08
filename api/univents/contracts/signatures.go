package contracts

import "github.com/google/uuid"

type Signature struct {
	ID        uuid.UUID `json:"id"`
	EditionID uuid.UUID `json:"edition_id"`
	Title     string    `json:"title"`
	URL       string    `json:"url"`
	PosX      int       `json:"pos_x"`
	PosY      int       `json:"pos_y"`
}

type AddSignatureRequest struct {
	Title string `json:"title"`
	URL   string `json:"url"`
	PosX  int    `json:"pos_x"`
	PosY  int    `json:"pos_y"`
}

func (r AddSignatureRequest) ToInput(editionID uuid.UUID) AddSignatureInput {
	return AddSignatureInput{
		Title:     r.Title,
		URL:       r.URL,
		EditionID: editionID,
		PosX:      r.PosX,
		PosY:      r.PosY,
	}
}

type AddSignatureInput struct {
	EditionID uuid.UUID `json:"edition_id"`
	Title     string    `json:"title"`
	URL       string    `json:"url"`
	PosX      int       `json:"pos_x"`
	PosY      int       `json:"pos_y"`
}

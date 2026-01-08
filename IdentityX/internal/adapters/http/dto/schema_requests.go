package dto

type DraftSchemaRequest struct {
	SchemaType string `json:"schema_type" validate:"required,oneof=core context sub-context"`
	Title      string `json:"title" validate:"required,max=255"`
	FlowID     string `json:"flow_id" validate:"required,max=63"`
}

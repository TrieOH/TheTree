package handlers

import (
	"GoAuth/internal/adapters/http/dto"
	"GoAuth/internal/adapters/http/validation"
	"GoAuth/internal/domain/field"
	"GoAuth/internal/ports/inbounds"
	"net/http"
	"strconv"

	resp "github.com/MintzyG/FastUtilitiesNet/response"
	"github.com/go-chi/chi/v5"
	"github.com/google/uuid"
)

type SchemaFieldsHandler struct {
	fields inbounds.SchemaFieldsService
}

func NewSchemaFieldsHandler(uc inbounds.SchemaFieldsService) *SchemaFieldsHandler {
	return &SchemaFieldsHandler{fields: uc}
}

// Create godoc
// @Summary Create fields for a schema version
// @Description Creates fields for a specific version of a schema.
// @Tags schema-fields
// @Accept json
// @Produce json
// @Param Cookie header string true "Cookie: access_token=xxx; refresh_token=yyy"
// @Param project_id path string true "Project ID"
// @Param schema_id path string true "Schema ID"
// @Param version path int true "Schema Version Number"
// @Param fieldInfo body dto.CreateFieldRequest true "Field creation information"
// @Success 201 {array} dto.FieldResponse
// @Failure 400 {object} ErrorResponse
// @Failure 401 {object} ErrorResponse
// @Failure 500 {object} ErrorResponse
// @Router /projects/{project_id}/schemas/{schema_id}/v{version} [post]
func (handler *SchemaFieldsHandler) Create(w http.ResponseWriter, r *http.Request) {
	projectID, rs := getUUID(r, "project_id")
	if rs != nil {
		rs.Send(w)
		return
	}

	schemaID, rs := getUUID(r, "schema_id")
	if rs != nil {
		rs.Send(w)
		return
	}

	version := chi.URLParam(r, "version")
	if version == "" {
		resp.BadRequest("missing version parameter").Send(w)
		return
	}

	versionNumber, err := strconv.Atoi(version)
	if err != nil {
		resp.BadRequest("invalid version parameter").AddTrace(err).Send(w)
		return
	}

	if versionNumber <= 0 {
		resp.BadRequest("version must be >= 1").Send(w)
		return
	}

	var req dto.CreateFieldRequest
	if err := validation.ValidateInto(r, &req); err != nil {
		resp.FromError(err).Send(w)
		return
	}

	in := inbounds.SchemaFieldInput{
		ProjectID:     projectID,
		SchemaID:      schemaID,
		VersionNumber: versionNumber,
		Fields:        dto.FieldParamSliceToInputFieldSlice(req.Fields),
	}

	ctx := r.Context()
	res, err := handler.fields.Create(ctx, in)
	if err != nil {
		resp.FromError(err).Send(w)
		return
	}

	if res.Warnings != nil {
		response := resp.Created("created fields with warnings").
			WithData(dto.OutputFieldSliceToFieldResponseSlice(res.Fields))
		for _, warning := range res.Warnings {
			response.AddTrace(warning)
		}
		response.Send(w)
		return
	}

	resp.Created("created fields").
		WithData(dto.OutputFieldSliceToFieldResponseSlice(res.Fields)).
		Send(w)
}

// EditField godoc
// @Summary Edit a field in a schema version draft
// @Description Updates a field's properties (only provided fields are updated)
// @Tags schema-fields
// @Accept json
// @Produce json
// @Param Cookie header string true "Cookie: access_token=xxx; refresh_token=yyy"
// @Param project_id path string true "Project ID"
// @Param schema_id path string true "Schema ID"
// @Param version path int true "Schema Version Number"
// @Param field_id path string true "Field Object ID"
// @Param fieldInfo body dto.EditFieldRequest true "Field update information"
// @Success 200 {object} dto.FieldResponse
// @Failure 400 {object} ErrorResponse
// @Failure 401 {object} ErrorResponse
// @Failure 409 {object} ErrorResponse
// @Failure 500 {object} ErrorResponse
// @Router /projects/{project_id}/schemas/{schema_id}/v{version}/fields/{field_id} [patch]
func (handler *SchemaFieldsHandler) EditField(w http.ResponseWriter, r *http.Request) {
	projectID, rs := getUUID(r, "project_id")
	if rs != nil {
		rs.Send(w)
		return
	}

	schemaID, rs := getUUID(r, "schema_id")
	if rs != nil {
		rs.Send(w)
		return
	}

	fieldID, rs := getUUID(r, "field_id")
	if rs != nil {
		rs.Send(w)
		return
	}

	version := chi.URLParam(r, "version")
	if version == "" {
		resp.BadRequest("missing version parameter").Send(w)
		return
	}

	versionNumber, err := strconv.Atoi(version)
	if err != nil {
		resp.BadRequest("invalid version parameter").AddTrace(err).Send(w)
		return
	}

	if versionNumber <= 0 {
		resp.BadRequest("version must be >= 1").Send(w)
		return
	}

	var req dto.EditFieldRequest
	if err := validation.ValidateInto(r, &req); err != nil {
		resp.FromError(err).Send(w)
		return
	}

	in := inbounds.EditFieldInput{
		ProjectID:     projectID,
		SchemaID:      schemaID,
		VersionNumber: versionNumber,
		FieldObjectID: fieldID,
		Key:           req.Key,
		Title:         req.Title,
		Description:   req.Description,
		Placeholder:   req.Placeholder,
		Type:          req.Type,
		Required:      req.Required,
		Mutable:       req.Mutable,
		DefaultValue:  req.DefaultValue,
		Position:      req.Position,
	}

	ctx := r.Context()
	updatedField, err := handler.fields.EditField(ctx, in)
	if err != nil {
		resp.FromError(err).Send(w)
		return
	}

	resp.OK().WithData(dto.OutputFieldToFieldResponse(inbounds.FieldToOutputField(updatedField))).Send(w)
}

// DeleteField godoc
// @Summary Delete a field from a schema version draft
// @Description Deletes a field and all its options and rules
// @Tags schema-fields
// @Accept json
// @Produce json
// @Param Cookie header string true "Cookie: access_token=xxx; refresh_token=yyy"
// @Param project_id path string true "Project ID"
// @Param schema_id path string true "Schema ID"
// @Param version path int true "Schema Version Number"
// @Param field_id path string true "Field Object ID"
// @Success 204 "No Content"
// @Failure 400 {object} ErrorResponse
// @Failure 401 {object} ErrorResponse
// @Failure 409 {object} ErrorResponse
// @Failure 500 {object} ErrorResponse
// @Router /projects/{project_id}/schemas/{schema_id}/v{version}/fields/{field_id} [delete]
func (handler *SchemaFieldsHandler) DeleteField(w http.ResponseWriter, r *http.Request) {
	projectID, rs := getUUID(r, "project_id")
	if rs != nil {
		rs.Send(w)
		return
	}

	schemaID, rs := getUUID(r, "schema_id")
	if rs != nil {
		rs.Send(w)
		return
	}

	fieldID, rs := getUUID(r, "field_id")
	if rs != nil {
		rs.Send(w)
		return
	}

	version := chi.URLParam(r, "version")
	if version == "" {
		resp.BadRequest("missing version parameter").Send(w)
		return
	}

	versionNumber, err := strconv.Atoi(version)
	if err != nil {
		resp.BadRequest("invalid version parameter").AddTrace(err).Send(w)
		return
	}

	if versionNumber <= 0 {
		resp.BadRequest("version must be >= 1").Send(w)
		return
	}

	in := inbounds.DeleteFieldInput{
		ProjectID:     projectID,
		SchemaID:      schemaID,
		VersionNumber: versionNumber,
		FieldObjectID: fieldID,
	}

	ctx := r.Context()
	if err := handler.fields.DeleteField(ctx, in); err != nil {
		resp.FromError(err).Send(w)
		return
	}

	resp.NoContent().Send(w)
}

// SetFieldOptions godoc
// @Summary Set options for a field (replaces all existing options)
// @Description Replaces all options for a field. Only allowed on draft versions.
// @Tags schema-fields
// @Accept json
// @Produce json
// @Param Cookie header string true "Cookie: access_token=xxx; refresh_token=yyy"
// @Param project_id path string true "Project ID"
// @Param schema_id path string true "Schema ID"
// @Param version path int true "Schema Version Number"
// @Param field_id path string true "Field Object ID"
// @Param options body dto.SetFieldOptionsRequest true "Options to set"
// @Success 200 {array} dto.OptionResponse
// @Failure 400 {object} ErrorResponse
// @Failure 401 {object} ErrorResponse
// @Failure 409 {object} ErrorResponse
// @Failure 500 {object} ErrorResponse
// @Router /projects/{project_id}/schemas/{schema_id}/v{version}/fields/{field_id}/options [put]
func (handler *SchemaFieldsHandler) SetFieldOptions(w http.ResponseWriter, r *http.Request) {
	projectID, rs := getUUID(r, "project_id")
	if rs != nil {
		rs.Send(w)
		return
	}

	schemaID, rs := getUUID(r, "schema_id")
	if rs != nil {
		rs.Send(w)
		return
	}

	fieldID, rs := getUUID(r, "field_id")
	if rs != nil {
		rs.Send(w)
		return
	}

	version := chi.URLParam(r, "version")
	if version == "" {
		resp.BadRequest("missing version parameter").Send(w)
		return
	}

	versionNumber, err := strconv.Atoi(version)
	if err != nil {
		resp.BadRequest("invalid version parameter").AddTrace(err).Send(w)
		return
	}

	if versionNumber <= 0 {
		resp.BadRequest("version must be >= 1").Send(w)
		return
	}

	var req dto.SetFieldOptionsRequest
	if err := validation.ValidateInto(r, &req); err != nil {
		resp.FromError(err).Send(w)
		return
	}

	options := make([]inbounds.InputOption, len(req.Options))
	for i, opt := range req.Options {
		options[i] = inbounds.InputOption{
			Value:    opt.Value,
			Label:    opt.Label,
			Position: opt.Position,
		}
	}

	in := inbounds.SetFieldOptionsInput{
		ProjectID:     projectID,
		SchemaID:      schemaID,
		VersionNumber: versionNumber,
		FieldObjectID: fieldID,
		Options:       options,
	}

	ctx := r.Context()
	result, err := handler.fields.SetFieldOptions(ctx, in)
	if err != nil {
		resp.FromError(err).Send(w)
		return
	}

	resp.OK().WithData(dto.OptionSliceToResponse(result)).Send(w)
}

// DeleteFieldOption godoc
// @Summary Delete a specific option from a field
// @Description Deletes a single option by ID. Only allowed on draft versions.
// @Tags schema-fields
// @Accept json
// @Produce json
// @Param Cookie header string true "Cookie: access_token=xxx; refresh_token=yyy"
// @Param project_id path string true "Project ID"
// @Param schema_id path string true "Schema ID"
// @Param version path int true "Schema Version Number"
// @Param field_id path string true "Field Object ID"
// @Param option_id path string true "Option ID"
// @Success 204 "No Content"
// @Failure 400 {object} ErrorResponse
// @Failure 401 {object} ErrorResponse
// @Failure 409 {object} ErrorResponse
// @Failure 500 {object} ErrorResponse
// @Router /projects/{project_id}/schemas/{schema_id}/v{version}/fields/{field_id}/options/{option_id} [delete]
func (handler *SchemaFieldsHandler) DeleteFieldOption(w http.ResponseWriter, r *http.Request) {
	projectID, rs := getUUID(r, "project_id")
	if rs != nil {
		rs.Send(w)
		return
	}

	schemaID, rs := getUUID(r, "schema_id")
	if rs != nil {
		rs.Send(w)
		return
	}

	fieldID, rs := getUUID(r, "field_id")
	if rs != nil {
		rs.Send(w)
		return
	}

	optionID, rs := getUUID(r, "option_id")
	if rs != nil {
		rs.Send(w)
		return
	}

	version := chi.URLParam(r, "version")
	if version == "" {
		resp.BadRequest("missing version parameter").Send(w)
		return
	}

	versionNumber, err := strconv.Atoi(version)
	if err != nil {
		resp.BadRequest("invalid version parameter").AddTrace(err).Send(w)
		return
	}

	if versionNumber <= 0 {
		resp.BadRequest("version must be >= 1").Send(w)
		return
	}

	in := inbounds.DeleteFieldOptionInput{
		ProjectID:     projectID,
		SchemaID:      schemaID,
		VersionNumber: versionNumber,
		FieldObjectID: fieldID,
		OptionID:      optionID,
	}

	ctx := r.Context()
	if err := handler.fields.DeleteFieldOption(ctx, in); err != nil {
		resp.FromError(err).Send(w)
		return
	}

	resp.NoContent().Send(w)
}

// SetVisibilityRules godoc
// @Summary Set visibility rules for a field (replaces all existing rules)
// @Description Replaces all visibility rules for a field. Only allowed on draft versions.
// @Tags schema-fields
// @Accept json
// @Produce json
// @Param Cookie header string true "Cookie: access_token=xxx; refresh_token=yyy"
// @Param project_id path string true "Project ID"
// @Param schema_id path string true "Schema ID"
// @Param version path int true "Schema Version Number"
// @Param field_id path string true "Field Object ID"
// @Param rules body dto.SetVisibilityRulesRequest true "Visibility rules to set"
// @Success 200 {array} dto.VisibilityRuleResponse
// @Failure 400 {object} ErrorResponse
// @Failure 401 {object} ErrorResponse
// @Failure 409 {object} ErrorResponse
// @Failure 500 {object} ErrorResponse
// @Router /projects/{project_id}/schemas/{schema_id}/v{version}/fields/{field_id}/visibility-rules [put]
func (handler *SchemaFieldsHandler) SetVisibilityRules(w http.ResponseWriter, r *http.Request) {
	projectID, rs := getUUID(r, "project_id")
	if rs != nil {
		rs.Send(w)
		return
	}

	schemaID, rs := getUUID(r, "schema_id")
	if rs != nil {
		rs.Send(w)
		return
	}

	fieldID, rs := getUUID(r, "field_id")
	if rs != nil {
		rs.Send(w)
		return
	}

	version := chi.URLParam(r, "version")
	if version == "" {
		resp.BadRequest("missing version parameter").Send(w)
		return
	}

	versionNumber, err := strconv.Atoi(version)
	if err != nil {
		resp.BadRequest("invalid version parameter").AddTrace(err).Send(w)
		return
	}

	if versionNumber <= 0 {
		resp.BadRequest("version must be >= 1").Send(w)
		return
	}

	var req dto.SetVisibilityRulesRequest
	if err := validation.ValidateInto(r, &req); err != nil {
		resp.FromError(err).Send(w)
		return
	}

	rules := make([]inbounds.InputVisibilityRule, len(req.VisibilityRules))
	for i, rule := range req.VisibilityRules {
		rules[i] = inbounds.InputVisibilityRule{
			DependsOnFieldKey: rule.DependsOnFieldKey,
			Operator:          rule.Operator,
			Value:             rule.Value,
		}
	}

	in := inbounds.SetVisibilityRulesInput{
		ProjectID:       projectID,
		SchemaID:        schemaID,
		VersionNumber:   versionNumber,
		FieldObjectID:   fieldID,
		VisibilityRules: rules,
	}

	ctx := r.Context()
	result, err := handler.fields.SetVisibilityRules(ctx, in)
	if err != nil {
		resp.FromError(err).Send(w)
		return
	}

	resp.OK().WithData(dto.VisibilityRuleSliceToResponse(result)).Send(w)
}

// EditVisibilityRule godoc
// @Summary Edit a visibility rule
// @Description Updates a visibility rule's properties (only provided fields are updated). Only allowed on draft versions.
// @Tags schema-fields
// @Accept json
// @Produce json
// @Param Cookie header string true "Cookie: access_token=xxx; refresh_token=yyy"
// @Param project_id path string true "Project ID"
// @Param schema_id path string true "Schema ID"
// @Param version path int true "Schema Version Number"
// @Param field_id path string true "Field Object ID"
// @Param rule_id path string true "Rule ID"
// @Param rule body dto.EditVisibilityRuleRequest true "Visibility rule update information"
// @Success 200 {object} dto.VisibilityRuleResponse
// @Failure 400 {object} ErrorResponse
// @Failure 401 {object} ErrorResponse
// @Failure 409 {object} ErrorResponse
// @Failure 500 {object} ErrorResponse
// @Router /projects/{project_id}/schemas/{schema_id}/v{version}/fields/{field_id}/visibility-rules/{rule_id} [patch]
func (handler *SchemaFieldsHandler) EditVisibilityRule(w http.ResponseWriter, r *http.Request) {
	projectID, rs := getUUID(r, "project_id")
	if rs != nil {
		rs.Send(w)
		return
	}

	schemaID, rs := getUUID(r, "schema_id")
	if rs != nil {
		rs.Send(w)
		return
	}

	fieldID, rs := getUUID(r, "field_id")
	if rs != nil {
		rs.Send(w)
		return
	}

	ruleID, rs := getUUID(r, "rule_id")
	if rs != nil {
		rs.Send(w)
		return
	}

	version := chi.URLParam(r, "version")
	if version == "" {
		resp.BadRequest("missing version parameter").Send(w)
		return
	}

	versionNumber, err := strconv.Atoi(version)
	if err != nil {
		resp.BadRequest("invalid version parameter").AddTrace(err).Send(w)
		return
	}

	if versionNumber <= 0 {
		resp.BadRequest("version must be >= 1").Send(w)
		return
	}

	var req dto.EditVisibilityRuleRequest
	if err := validation.ValidateInto(r, &req); err != nil {
		resp.FromError(err).Send(w)
		return
	}

	var dependsOnFieldID *uuid.UUID
	if req.DependsOnFieldID != nil {
		id, err := uuid.Parse(*req.DependsOnFieldID)
		if err != nil {
			resp.BadRequest("invalid depends_on_field_id").Send(w)
			return
		}
		dependsOnFieldID = &id
	}

	in := inbounds.EditVisibilityRuleInput{
		ProjectID:        projectID,
		SchemaID:         schemaID,
		VersionNumber:    versionNumber,
		FieldObjectID:    fieldID,
		RuleID:           ruleID,
		DependsOnFieldID: dependsOnFieldID,
		Operator:         req.Operator,
		Value:            req.Value,
	}

	ctx := r.Context()
	updatedRule, err := handler.fields.EditVisibilityRule(ctx, in)
	if err != nil {
		resp.FromError(err).Send(w)
		return
	}

	resp.OK().WithData(dto.VisibilityRuleSliceToResponse([]field.VisibilityRule{*updatedRule})).Send(w)
}

// DeleteVisibilityRule godoc
// @Summary Delete a visibility rule
// @Description Deletes a single visibility rule by ID. Only allowed on draft versions.
// @Tags schema-fields
// @Accept json
// @Produce json
// @Param Cookie header string true "Cookie: access_token=xxx; refresh_token=yyy"
// @Param project_id path string true "Project ID"
// @Param schema_id path string true "Schema ID"
// @Param version path int true "Schema Version Number"
// @Param field_id path string true "Field Object ID"
// @Param rule_id path string true "Rule ID"
// @Success 204 "No Content"
// @Failure 400 {object} ErrorResponse
// @Failure 401 {object} ErrorResponse
// @Failure 409 {object} ErrorResponse
// @Failure 500 {object} ErrorResponse
// @Router /projects/{project_id}/schemas/{schema_id}/v{version}/fields/{field_id}/visibility-rules/{rule_id} [delete]
func (handler *SchemaFieldsHandler) DeleteVisibilityRule(w http.ResponseWriter, r *http.Request) {
	projectID, rs := getUUID(r, "project_id")
	if rs != nil {
		rs.Send(w)
		return
	}

	schemaID, rs := getUUID(r, "schema_id")
	if rs != nil {
		rs.Send(w)
		return
	}

	fieldID, rs := getUUID(r, "field_id")
	if rs != nil {
		rs.Send(w)
		return
	}

	ruleID, rs := getUUID(r, "rule_id")
	if rs != nil {
		rs.Send(w)
		return
	}

	version := chi.URLParam(r, "version")
	if version == "" {
		resp.BadRequest("missing version parameter").Send(w)
		return
	}

	versionNumber, err := strconv.Atoi(version)
	if err != nil {
		resp.BadRequest("invalid version parameter").AddTrace(err).Send(w)
		return
	}

	if versionNumber <= 0 {
		resp.BadRequest("version must be >= 1").Send(w)
		return
	}

	in := inbounds.DeleteVisibilityRuleInput{
		ProjectID:     projectID,
		SchemaID:      schemaID,
		VersionNumber: versionNumber,
		FieldObjectID: fieldID,
		RuleID:        ruleID,
	}

	ctx := r.Context()
	if err := handler.fields.DeleteVisibilityRule(ctx, in); err != nil {
		resp.FromError(err).Send(w)
		return
	}

	resp.NoContent().Send(w)
}

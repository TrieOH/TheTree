import type { FieldDefinition } from "#/shared/model/form-types";

/**
 * Behaviour options for select-type fields.
 */
const SELECT_BEHAVIOUR_OPTIONS = [
    { label: "Checkbox", value: "checkbox" },
    { label: "Radio", value: "radio" },
    { label: "Dropdown Checkbox", value: "dropdown-checkbox" },
    { label: "Dropdown Radio", value: "dropdown-radio" },
] as const;

/**
 * Value type options for select-type fields.
 */
const SELECT_VALUE_TYPE_OPTIONS = [
    { label: "String", value: "string" },
    { label: "Email", value: "email" },
    { label: "Integer", value: "int" },
    { label: "Float", value: "float" },
    { label: "Date", value: "date" },
    { label: "Time", value: "time" },
    { label: "Datetime", value: "datetime" },
    { label: "Phone", value: "phone" },
    { label: "URL", value: "url" },
] as const;

/**
 * Field type options for the type selector.
 */
const FIELD_TYPE_OPTIONS = [
    { label: "String", value: "string" },
    { label: "Email", value: "email" },
    { label: "Integer", value: "int" },
    { label: "Float", value: "float" },
    { label: "Boolean", value: "bool" },
    { label: "Date", value: "date" },
    { label: "Time", value: "time" },
    { label: "Datetime", value: "datetime" },
    { label: "Select", value: "select" },
    { label: "File", value: "file" },
    { label: "Phone", value: "phone" },
    { label: "URL", value: "url" },
] as const;

/**
 * Boolean selector options (required field).
 */
const BOOL_OPTIONS = [
    { label: "Yes", value: "true" },
    { label: "No", value: "false" },
] as const;

/** Text-based field types that support placeholder and text default_value. */
const TEXT_BASED_TYPES = [
    "string", "email", "int", "float", "phone", "url",
    "date", "time", "datetime",
];

/**
 * Base field definitions shared between create and edit field forms.
 */
const BASE_FIELD_DEFS: FieldDefinition<Record<string, unknown>>[] = [
    {
        name: "title",
        label: "Field Label",
        type: "text",
        placeholder: "e.g. Full Name",
    },
    {
        name: "key",
        label: "Field Key",
        type: "text",
        placeholder: "e.g. full_name",
    },
    {
        name: "description",
        label: "Description",
        type: "textarea",
        rows: 3,
        placeholder: "e.g. Enter your full name as it appears on your ID.",
        required: false,
    },
    {
        name: "type",
        label: "Field Type",
        type: "select",
        placeholder: "Select a type…",
        options: [...FIELD_TYPE_OPTIONS],
    },
    {
        name: "required",
        label: "Required",
        type: "boolean",
        placeholder: "User must select an option to submit",
    },
    {
        name: "placeholder",
        label: "Placeholder",
        type: "text",
        placeholder: "e.g. Type your answer here…",
        required: false,
        dependsOn: { field: "type", value: TEXT_BASED_TYPES },
    },
    {
        name: "default_value",
        label: "Default Value",
        type: "text",
        placeholder: "e.g. John Doe",
        required: false,
        dependsOn: { field: "type", value: [...TEXT_BASED_TYPES, "select"] },
    },
];

/**
 * Select config fields that appear when type === "select".
 */
const SELECT_CONFIG_FIELD_DEFS: FieldDefinition<Record<string, unknown>>[] = [
    {
        name: "select_config.options",
        label: "Options",
        type: "textarea",
        rows: 5,
        placeholder: "One option per line\ne.g.\nOption 1\nOption 2\nOption 3",
        required: false,
        dependsOn: { field: "type", value: "select" },
    },
    {
        name: "select_config.behaviour",
        label: "Select Behaviour",
        type: "select",
        placeholder: "Select behaviour…",
        options: [...SELECT_BEHAVIOUR_OPTIONS],
        dependsOn: { field: "type", value: "select" },
    },
    {
        name: "select_config.value_type",
        label: "Value Type",
        type: "select",
        placeholder: "Select value type…",
        options: [...SELECT_VALUE_TYPE_OPTIONS],
        dependsOn: { field: "type", value: "select" },
    },
];

/**
 * Returns the full list of field definitions for a field form (create or edit),
 * appending select_config fields that show conditionally when type === "select".
 *
 * @param extraFields - Additional fields to append at the end (e.g. position_hint for create).
 */
export function getFieldFormDefs<T = Record<string, unknown>>(
    ...extraFields: FieldDefinition<T>[]
): FieldDefinition<T>[] {
    return [
        ...(BASE_FIELD_DEFS as FieldDefinition<T>[]),
        ...(SELECT_CONFIG_FIELD_DEFS as FieldDefinition<T>[]),
        ...extraFields,
    ];
}

export { FIELD_TYPE_OPTIONS, BOOL_OPTIONS, SELECT_BEHAVIOUR_OPTIONS, SELECT_VALUE_TYPE_OPTIONS };

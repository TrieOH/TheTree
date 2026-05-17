-- +goose Up
CREATE TABLE forms (
    id UUID PRIMARY KEY DEFAULT uuidv7(),
    owner_id UUID NOT NULL,
    created_by UUID NOT NULL,
    namespace_id UUID REFERENCES namespaces(id)
        ON DELETE CASCADE,

    name VARCHAR(255) NOT NULL,
    CONSTRAINT uniq_form_name_per_namespace UNIQUE (namespace_id, name),
    status VARCHAR(8) NOT NULL DEFAULT 'draft',
    CONSTRAINT chk_forms_valid_status CHECK (status IN ('draft', 'open', 'closed', 'archived')),
    is_public BOOLEAN DEFAULT FALSE,
    is_multi_response BOOLEAN DEFAULT FALSE,

    opened_at TIMESTAMPTZ,
    closed_at TIMESTAMPTZ,
    archived_at TIMESTAMPTZ,
    CONSTRAINT chk_forms_valid_status_state CHECK (
        (status = 'open' AND opened_at IS NOT NULL)
            OR (status = 'closed' AND closed_at IS NOT NULL)
            OR (status = 'archived' AND archived_at IS NOT NULL)
            OR (status = 'draft')
        ),

    created_at TIMESTAMPTZ NOT NULL DEFAULT now(),
    updated_at TIMESTAMPTZ NOT NULL DEFAULT now()
);

CREATE INDEX idx_forms_owner_id ON forms (owner_id);

CREATE TABLE form_members (
    form_id UUID NOT NULL REFERENCES forms(id)
        ON DELETE CASCADE,
    user_id UUID NOT NULL,

    role VARCHAR(32) NOT NULL,
    CONSTRAINT chk_valid_namespace_member_role CHECK (
        role in ('viewer', 'editor', 'admin')
    ),

    added_at TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    added_by UUID NOT NULL,

    PRIMARY KEY (user_id, form_id)
);

CREATE INDEX idx_form_members_role ON namespace_members (role);

-- +goose Down
DROP INDEX IF EXISTS idx_form_members_role;
DROP TABLE IF EXISTS form_members;
DROP INDEX IF EXISTS idx_forms_owner_id;
DROP TABLE IF EXISTS forms;

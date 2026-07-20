-- +goose Up
CREATE TABLE programs (
    id               UUID PRIMARY KEY DEFAULT uuidv7(),
    edition_id       UUID NOT NULL REFERENCES editions(id) ON DELETE CASCADE,
    kind             TEXT NOT NULL DEFAULT 'activity',
    CONSTRAINT chk_programs_kind_valid CHECK (kind IN ('activity', 'checkpoint')),
    name             VARCHAR(256) NOT NULL,
    description      TEXT,
    min_access_level INT,
    staff_only       BOOLEAN NOT NULL DEFAULT FALSE,
    price            INT NOT NULL DEFAULT 0,
    token_cost       INT NOT NULL DEFAULT 0,
    created_at       TIMESTAMPTZ NOT NULL DEFAULT now(),
    updated_at       TIMESTAMPTZ NOT NULL DEFAULT now(),
    deleted_at       TIMESTAMPTZ
);

CREATE TABLE program_occurrences (
    id           UUID PRIMARY KEY DEFAULT uuidv7(),
    program_id   UUID NOT NULL REFERENCES programs(id) ON DELETE CASCADE,
    starts_at    TIMESTAMPTZ NOT NULL,
    ends_at      TIMESTAMPTZ NOT NULL,
    CONSTRAINT chk_program_occurrences_dates_valid CHECK (ends_at > starts_at),
    max_capacity INT,
    created_at   TIMESTAMPTZ NOT NULL DEFAULT now(),
    updated_at   TIMESTAMPTZ NOT NULL DEFAULT now(),
    deleted_at   TIMESTAMPTZ
);

CREATE TABLE program_participations (
    id              UUID PRIMARY KEY DEFAULT uuidv7(),
    occurrence_id   UUID NOT NULL REFERENCES program_occurrences(id) ON DELETE CASCADE,
    registration_id UUID NOT NULL REFERENCES registrations(id) ON DELETE CASCADE,
    status          TEXT NOT NULL DEFAULT 'registered',
    CONSTRAINT chk_program_participations_status_valid CHECK (
    status IN ('registered', 'attended', 'no_show', 'cancelled')
    ),
    created_at      TIMESTAMPTZ NOT NULL DEFAULT now(),
    updated_at      TIMESTAMPTZ NOT NULL DEFAULT now()
);
-- +goose Down
DROP TABLE IF EXISTS program_participations;
DROP TABLE IF EXISTS program_occurrences;
DROP TABLE IF EXISTS programs;
-- +goose Up
-- Store (issue-61 split 3): purchases, purchase_items, ws_tokens.

CREATE TABLE purchases (
    id                 UUID PRIMARY KEY DEFAULT uuidv7(),
    edition_id         UUID NOT NULL REFERENCES editions(id) ON DELETE CASCADE,
    purchaser_id       UUID NOT NULL,
    status             TEXT NOT NULL DEFAULT 'pending',
    CONSTRAINT chk_purchases_status_valid CHECK (
        status IN ('pending', 'approved', 'expired', 'cancelled')
    ),
    status_reason      TEXT,
    total_cents        BIGINT NOT NULL DEFAULT 0,
    currency           TEXT NOT NULL DEFAULT 'BRL',
    payment_method     TEXT,
    payssage_seller_id UUID,
    payssage_intent_id UUID, -- the correlation key (D2): never a provider-specific id
    qr_code            TEXT,
    qr_code_base64     TEXT,
    expires_at         TIMESTAMPTZ NOT NULL DEFAULT now() + interval '10 minutes 1 second',
    river_job_id       BIGINT,
    created_at         TIMESTAMPTZ NOT NULL DEFAULT now(),
    updated_at         TIMESTAMPTZ,
    deleted_at         TIMESTAMPTZ
);

-- One active (pending) purchase per purchaser per edition; the checkout
-- (split 7) 409s on a second one.
CREATE UNIQUE INDEX uniq_purchases_pending_per_purchaser_edition
    ON purchases (purchaser_id, edition_id)
    WHERE status = 'pending' AND deleted_at IS NULL;

CREATE INDEX idx_purchases_edition ON purchases(edition_id) WHERE deleted_at IS NULL;
CREATE INDEX idx_purchases_purchaser ON purchases(purchaser_id) WHERE deleted_at IS NULL;

CREATE TABLE purchase_items (
    id                  UUID PRIMARY KEY DEFAULT uuidv7(),
    purchase_id         UUID NOT NULL REFERENCES purchases(id) ON DELETE CASCADE,
    item_type           TEXT NOT NULL,
    CONSTRAINT chk_purchase_items_type_valid CHECK (
        item_type IN ('ticket', 'product', 'program_occurrence')
    ),
    item_id             UUID NOT NULL,
    quantity            INT NOT NULL DEFAULT 1,
    CONSTRAINT chk_purchase_items_quantity_positive CHECK (quantity > 0),
    unit_price_cents    BIGINT NOT NULL DEFAULT 0,
    -- materialization links (D4): the pending rows the purchase owns
    registration_id      UUID REFERENCES registrations(id),
    product_purchase_id  UUID REFERENCES product_purchases(id),
    participation_id     UUID REFERENCES program_participations(id)
);

-- Line uniqueness is per item type (checkout, split 7):
--   * tickets are one row PER PERSON (gifting) — one line per attendee,
--     each with its own registration_id, so the same ticket type may
--     appear on multiple lines;
--   * products are one line per item, quantity > 1 allowed;
--   * program occurrences are one line per occurrence (quantity = 1),
--     attached to the ticket's registration.
CREATE UNIQUE INDEX uniq_purchase_items_ticket_unit
    ON purchase_items (purchase_id, item_id, registration_id)
    WHERE item_type = 'ticket';
CREATE UNIQUE INDEX uniq_purchase_items_product_line
    ON purchase_items (purchase_id, item_id)
    WHERE item_type = 'product';
CREATE UNIQUE INDEX uniq_purchase_items_program_line
    ON purchase_items (purchase_id, item_id)
    WHERE item_type = 'program_occurrence';

CREATE INDEX idx_purchase_items_item ON purchase_items(item_type, item_id);

CREATE TABLE ws_tokens (
    id          UUID PRIMARY KEY DEFAULT uuidv7(),
    purchase_id UUID NOT NULL REFERENCES purchases(id) ON DELETE CASCADE,
    user_id     UUID NOT NULL,
    token_hash  TEXT NOT NULL,
    expires_at  TIMESTAMPTZ NOT NULL,
    used_at     TIMESTAMPTZ,
    CONSTRAINT uniq_ws_tokens_token_hash UNIQUE (token_hash)
);

CREATE INDEX idx_ws_tokens_purchase ON ws_tokens(purchase_id);

-- +goose Down
DROP TABLE IF EXISTS ws_tokens;
DROP TABLE IF EXISTS purchase_items;
DROP TABLE IF EXISTS purchases;

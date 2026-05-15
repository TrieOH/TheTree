-- +goose Up

CREATE TYPE permission_type AS ENUM (
    'activity',
    'product',
    'checkpoint'
);

-- Ticket permissions (what each ticket can access)
CREATE TABLE ticket_permissions (
    id UUID PRIMARY KEY DEFAULT uuidv7(),
    ticket_id UUID NOT NULL REFERENCES tickets(id) ON DELETE CASCADE,

    permission_type permission_type NOT NULL, -- 'activity', 'product', 'checkpoint'
    activity_id UUID NULL REFERENCES activities(id)
        ON DELETE CASCADE
        ON UPDATE CASCADE,
    product_id UUID NULL REFERENCES products(id)
        ON DELETE CASCADE
        ON UPDATE CASCADE,
    checkpoint_id UUID NULL REFERENCES checkpoints(id)
        ON DELETE CASCADE
        ON UPDATE CASCADE,

    created_at TIMESTAMPTZ NOT NULL DEFAULT now(),

    CONSTRAINT chk_type_matches_target CHECK (
        (permission_type = 'activity'   AND activity_id   IS NOT NULL AND product_id  IS NULL AND checkpoint_id IS NULL) OR
        (permission_type = 'product'    AND product_id    IS NOT NULL AND activity_id IS NULL AND checkpoint_id IS NULL) OR
        (permission_type = 'checkpoint' AND checkpoint_id IS NOT NULL AND activity_id IS NULL AND product_id    IS NULL)
    )
);

-- one permission per ticket per activity
CREATE UNIQUE INDEX uq_ticket_permission_activity
    ON ticket_permissions(ticket_id, activity_id)
    WHERE activity_id IS NOT NULL;

-- one permission per ticket per product
CREATE UNIQUE INDEX uq_ticket_permission_product
    ON ticket_permissions(ticket_id, product_id)
    WHERE product_id IS NOT NULL;

-- one permission per ticket per checkpoint
CREATE UNIQUE INDEX uq_ticket_permission_checkpoint
    ON ticket_permissions(ticket_id, checkpoint_id)
    WHERE checkpoint_id IS NOT NULL;

-- +goose Down
DROP INDEX IF EXISTS uq_ticket_permission_checkpoint;
DROP INDEX IF EXISTS uq_ticket_permission_product;
DROP INDEX IF EXISTS uq_ticket_permission_activity;
DROP TABLE IF EXISTS ticket_permissions;
DROP TYPE IF EXISTS permission_type

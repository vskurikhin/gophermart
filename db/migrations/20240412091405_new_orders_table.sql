-- +goose Up
-- +goose StatementBegin
CREATE TABLE IF NOT EXISTS "orders"
(
    "login"       varchar,
    "number"      varchar,
    "status_id"   integer   NOT NULL,
    "accrual"     numeric,
    "uploaded_at" timestamp,
    "created_at"  timestamp NOT NULL,
    "update_at"   timestamp,
    PRIMARY KEY ("login", "number"),
    UNIQUE ("login", "number"),
    UNIQUE ("number")
);

ALTER TABLE "orders"
    ADD FOREIGN KEY ("login") REFERENCES "users" ("login");
ALTER TABLE "orders"
    ADD FOREIGN KEY ("status_id") REFERENCES "status" ("id");

DROP TRIGGER IF EXISTS set_created_at_in_order ON "orders";
CREATE TRIGGER set_created_at_in_order
    BEFORE UPDATE
    ON "orders"
    FOR EACH ROW
EXECUTE FUNCTION set_created_at();

DROP TRIGGER IF EXISTS set_update_at_in_order ON "orders";
CREATE TRIGGER set_update_at_in_order
    BEFORE UPDATE
    ON "orders"
    FOR EACH ROW
EXECUTE FUNCTION set_update_at();
-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
DROP TABLE IF EXISTS "orders";
-- +goose StatementEnd

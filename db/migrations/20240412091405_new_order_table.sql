-- +goose Up
-- +goose StatementBegin
CREATE TABLE IF NOT EXISTS "order"
(
    "login"       varchar,
    "number"      varchar,
    "status_id"   integer   NOT NULL,
    "uploaded_at" timestamp,
    "created_at"  timestamp NOT NULL,
    "update_at"   timestamp,
    PRIMARY KEY ("login", "number"),
    UNIQUE ("login", "number"),
    UNIQUE ("number")
);

ALTER TABLE "order"
    ADD FOREIGN KEY ("login") REFERENCES "user" ("login");
ALTER TABLE "order"
    ADD FOREIGN KEY ("status_id") REFERENCES "status" ("id");

DROP TRIGGER IF EXISTS set_created_at_in_order ON "order";
CREATE TRIGGER set_created_at_in_order
    BEFORE UPDATE
    ON "order"
    FOR EACH ROW
EXECUTE FUNCTION set_created_at();

DROP TRIGGER IF EXISTS set_update_at_in_order ON "order";
CREATE TRIGGER set_update_at_in_order
    BEFORE UPDATE
    ON "order"
    FOR EACH ROW
EXECUTE FUNCTION set_update_at();
-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
DROP TABLE IF EXISTS "order";
-- +goose StatementEnd

-- +goose Up
-- +goose StatementBegin
CREATE TABLE IF NOT EXISTS "balance"
(
    "login"      varchar PRIMARY KEY,
    "current"    numeric   NOT NULL DEFAULT 0,
    "withdrawn"  numeric   NOT NULL DEFAULT 0,
    "created_at" timestamp NOT NULL,
    "update_at"  timestamp
);

ALTER TABLE "balance"
    ADD FOREIGN KEY ("login") REFERENCES "users" ("login");

DROP TRIGGER IF EXISTS set_created_at_in_balance ON "balance";
CREATE TRIGGER set_created_at_in_balance
    BEFORE UPDATE
    ON "balance"
    FOR EACH ROW
EXECUTE FUNCTION set_created_at();

DROP TRIGGER IF EXISTS set_update_at_in_balance ON "balance";
CREATE TRIGGER set_update_at_in_balance
    BEFORE UPDATE
    ON "balance"
    FOR EACH ROW
EXECUTE FUNCTION set_update_at();
-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
DROP TABLE IF EXISTS "balance";
-- +goose StatementEnd

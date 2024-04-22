-- +goose Up
-- +goose StatementBegin
CREATE TABLE IF NOT EXISTS "users"
(
    "login"      varchar PRIMARY KEY,
    "password"   varchar,
    "created_at" timestamp NOT NULL,
    "update_at"  timestamp
);

DROP TRIGGER IF EXISTS set_created_at_in_user ON "users";
CREATE TRIGGER set_created_at_in_status
    BEFORE UPDATE
    ON "users"
    FOR EACH ROW
EXECUTE FUNCTION set_created_at();

DROP TRIGGER IF EXISTS set_update_at_in_user ON "users";
CREATE TRIGGER set_update_at_in_status
    BEFORE UPDATE
    ON "users"
    FOR EACH ROW
EXECUTE FUNCTION set_update_at();
-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
DROP TABLE IF EXISTS "users";
-- +goose StatementEnd
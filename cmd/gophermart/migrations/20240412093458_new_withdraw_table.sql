-- +goose Up
-- +goose StatementBegin
CREATE TABLE "withdraw"
(
    "login"        varchar,
    "number"       varchar,
    "sum"          numeric   NOT NULL DEFAULT 0,
    "status_id"    integer   NOT NULL,
    "processed_at" timestamp,
    "created_at"   timestamp NOT NULL,
    "update_at"    timestamp,
    PRIMARY KEY ("login", "number"),
    UNIQUE ("number")
);

ALTER TABLE "withdraw"
    ADD FOREIGN KEY ("status_id") REFERENCES "status" ("id");
ALTER TABLE "withdraw"
    ADD FOREIGN KEY ("login") REFERENCES "users" ("login");

DROP TRIGGER IF EXISTS set_created_at_in_withdraw ON "withdraw";
CREATE TRIGGER set_created_at_in_withdraw
    BEFORE UPDATE
    ON "withdraw"
    FOR EACH ROW
EXECUTE FUNCTION set_created_at();

DROP TRIGGER IF EXISTS set_update_at_in_withdraw ON "withdraw";
CREATE TRIGGER set_update_at_in_withdraw
    BEFORE UPDATE
    ON "withdraw"
    FOR EACH ROW
EXECUTE FUNCTION set_update_at();
-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
DROP TABLE IF EXISTS "withdraw";
-- +goose StatementEnd

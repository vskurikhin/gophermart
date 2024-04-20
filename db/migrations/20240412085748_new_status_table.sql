-- +goose Up
-- +goose StatementBegin
CREATE TABLE IF NOT EXISTS "status"
(
    "id"         SERIAL PRIMARY KEY,
    "status"     varchar UNIQUE,
    "created_at" timestamp NOT NULL,
    "update_at"  timestamp
);

DROP TRIGGER IF EXISTS set_created_at_in_status ON "status";
CREATE TRIGGER set_created_at_in_status
    BEFORE UPDATE
    ON "status"
    FOR EACH ROW
EXECUTE FUNCTION set_created_at();

DROP TRIGGER IF EXISTS set_update_at_in_status ON "status";
CREATE TRIGGER set_update_at_in_status
    BEFORE UPDATE
    ON "status"
    FOR EACH ROW
EXECUTE FUNCTION set_update_at();

INSERT INTO "status" ("id", "status", "created_at")
VALUES (0, 'NEW', now()),
       (1, 'REGISTERED', now()),
       (2, 'PROCESSING', now()),
       (3, 'INVALID', now()),
       (4, 'PROCESSED', now());
-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
DELETE
FROM "status"
WHERE "status" IN ('NEW', 'PROCESSING', 'INVALID', 'PROCESSED');

DROP TABLE IF EXISTS "status";
-- +goose StatementEnd

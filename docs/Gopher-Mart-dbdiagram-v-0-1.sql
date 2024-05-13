CREATE TABLE "users"
(
    "login"      varchar PRIMARY KEY,
    "password"   varchar,
    "role"       varchar,
    "created_at" timestamp,
    "update_at"  timestamp
);

CREATE TABLE "balance"
(
    "login"      varchar PRIMARY KEY,
    "current"    numeric,
    "withdrawn"  numeric,
    "created_at" timestamp,
    "update_at"  timestamp
);

CREATE TABLE "orders"
(
    "login"       varchar,
    "number"      varchar,
    "status_id"   integer,
    "accrual"     numeric,
    "uploaded_at" timestamp,
    "created_at"  timestamp,
    "update_at"   timestamp,
    PRIMARY KEY ("login", "number")
);

CREATE TABLE "status"
(
    "id"         integer PRIMARY KEY,
    "status"     varchar,
    "created_at" timestamp,
    "update_at"  timestamp
);

CREATE TABLE "withdraw"
(
    "login"        varchar,
    "number"       varchar,
    "sum"          numeric,
    "status_id"    integer,
    "processed_at" timestamp,
    "created_at"   timestamp,
    "update_at"    timestamp,
    PRIMARY KEY ("login", "number")
);

ALTER TABLE "balance"
    ADD FOREIGN KEY ("login") REFERENCES "users" ("login");

ALTER TABLE "orders"
    ADD FOREIGN KEY ("login") REFERENCES "users" ("login");

ALTER TABLE "orders"
    ADD FOREIGN KEY ("status_id") REFERENCES "status" ("id");

ALTER TABLE "withdraw"
    ADD FOREIGN KEY ("login") REFERENCES "users" ("login");

ALTER TABLE "withdraw"
    ADD FOREIGN KEY ("status_id") REFERENCES "status" ("id");

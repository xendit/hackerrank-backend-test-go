BEGIN;

CREATE TABLE IF NOT EXISTS "user" (
    "id" SERIAL PRIMARY KEY,
    "first_name" VARCHAR(255) NOT NULL,
    "last_name" VARCHAR(255) NOT NULL,
    "address" VARCHAR(255) NOT NULL,
    "is_active" bool,
    "created_time" TIMESTAMP NOT NULL DEFAULT (now() AT TIME ZONE 'utc'),
    "updated_time" TIMESTAMP NOT NULL DEFAULT (now() AT TIME ZONE 'utc')
);

COMMIT;
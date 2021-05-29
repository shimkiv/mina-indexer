-- +goose Up
CREATE EXTENSION IF NOT EXISTS "uuid-ossp";

CREATE TYPE OWNER_TYPE AS ENUM ('validator', 'delegator', 'unknown');

-- +goose Down
DROP EXTENSION IF EXISTS "uuid-ossp";

DROP DOMAIN IF EXISTS OWNER_TYPE;
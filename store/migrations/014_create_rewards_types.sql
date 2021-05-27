-- +goose Up
CREATE EXTENSION IF NOT EXISTS "uuid-ossp";

CREATE TYPE REWARD_OWNER_TYPE AS ENUM ('validator', 'delegator', 'unknown');

-- +goose Down
DROP EXTENSION IF EXISTS "uuid-ossp";

DROP DOMAIN IF EXISTS REWARD_OWNER_TYPE;
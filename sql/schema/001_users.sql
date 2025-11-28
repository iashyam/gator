-- +goose up
CREATE TABLE users(
    id UUID PRIMARY KEY UNIQUE, 
    created_At TIMESTAMP NOT NULL,
    updated_At TIMESTAMP NOT NULL,
    name TEXT NOT NULL UNIQUE
);

-- +goose down
DROP TABLE users;
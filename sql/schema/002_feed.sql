-- +goose up
CREATE TABLE feeds(
    id UUID PRIMARY KEY UNIQUE, 
    created_At TIMESTAMP NOT NULL,
    updated_At TIMESTAMP NOT NULL,
    name TEXT NOT NULL UNIQUE,
    url TEXT NOT NULL, 
    user_id UUID NOT NULL REFERENCES users(id) ON DELETE CASCADE
);

-- +goose down
DROP TABLE feeds;

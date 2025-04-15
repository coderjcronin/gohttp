-- +goose Up
CREATE TABLE chirps (
    id UUID PRIMARY KEY,
    created_at TIMESTAMP NOT NULL,
    updated_at TIMESTAMP NOT NULL,
    body TEXT UNIQUE NOT NULL,
    user_id UUID NOT NULL,
    CONSTRAINT fk_user_id 
    FOREIGN KEY (user_id) 
    REFERENCES users(id)
);

-- +goose Down
DROP TABLE chirps;
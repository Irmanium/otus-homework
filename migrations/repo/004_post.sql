-- +goose Up
CREATE TABLE post (
    id uuid,
    user_id uuid,
    text text,
    updated_at timestamp with time zone DEFAULT now(),
    PRIMARY KEY (id)
);

CREATE INDEX post_idx_1
    ON post(user_id);

-- +goose Down
DROP TABLE post;
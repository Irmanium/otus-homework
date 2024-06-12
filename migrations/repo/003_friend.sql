-- +goose Up
CREATE TABLE friend (
    user_id uuid,
    friend_id uuid,
    PRIMARY KEY(user_id, friend_id)
);

-- +goose Down
DROP TABLE friend;
-- +goose Up
CREATE TABLE user_profile (
    id uuid,
    first_name text,
    second_name text,
    birthdate date,
    biography text,
    city text,
    password_hash text,
    PRIMARY KEY(id)
);

-- +goose Down
DROP TABLE user_profile;
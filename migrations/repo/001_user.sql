-- +goose Up
CREATE TABLE user_profile (
    id uuid,
    first_name text NOT NULL,
    second_name text NOT NULL,
    birthdate date NOT NULL,
    biography text NOT NULL,
    city text NOT NULL,
    password_hash text NOT NULL,
    PRIMARY KEY(id)
);

-- +goose Down
DROP TABLE user_profile;
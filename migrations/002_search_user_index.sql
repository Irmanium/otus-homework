-- +goose Up
CREATE INDEX user_profile_idx_1
    ON user_profile(first_name text_pattern_ops, second_name text_pattern_ops);

-- +goose Down
DROP INDEX user_profile_idx_1;
-- +goose Up
CREATE TABLE messages (
    id uuid,
    dialog_id text,
    from_user_id uuid,
    to_user_id uuid,
    text text,
    PRIMARY KEY (dialog_id, id)
);

SELECT create_distributed_table('messages', 'dialog_id');

-- +goose Down
DROP TABLE messages;
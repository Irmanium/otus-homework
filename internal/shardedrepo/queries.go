package shardedrepo

const (
	sendMessageQuery = `
		INSERT INTO messages (id, dialog_id, from_user_id, to_user_id, text)
			VALUES ($1, $2, $3, $4, $5);`

	getDialogQuery = `
		SELECT from_user_id, to_user_id, text
			FROM messages
			WHERE dialog_id = $1;`
)

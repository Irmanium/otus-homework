package repo

const (
	getUserQuery = `
		SELECT first_name, second_name, birthdate, biography, city
			FROM user_profile
			WHERE id = $1;`

	registerQuery = `
		INSERT INTO user_profile (id, first_name, second_name, birthdate, biography, city, password_hash)
			VALUES ($1, $2, $3, $4, $5, $6, $7);`

	getPasswordQuery = `
		SELECT password_hash
			FROM user_profile
			WHERE id = $1;`

	searchUserQuery = `
		SELECT id, first_name, second_name, birthdate, biography, city
			FROM user_profile
			WHERE first_name LIKE $1 || '%' AND second_name LIKE $2 || '%'
			ORDER BY id;`

	setFriendQuery = `
		INSERT INTO friend (user_id, friend_id)
			VALUES ($1, $2)
			ON CONFLICT DO NOTHING;`

	deleteFriendQuery = `
		DELETE FROM friend
			WHERE (user_id = $1 AND friend_id = $2) OR (friend_id = $1 AND user_id = $2);`

	createPostQuery = `
		INSERT INTO post (id, user_id, text)
			VALUES ($1, $2, $3);`

	updatePostQuery = `
		UPDATE post
			SET text = $2, updated_at = now()
			WHERE id = $1;`

	deletePostQuery = `
		DELETE FROM post
			WHERE id = $1;`

	getPostQuery = `
		SELECT user_id, text, updated_at
			FROM post
			WHERE id = $1;`

	getFeedQuery = `
		SELECT id, text, user_id
			FROM post
			WHERE user_id IN (
    			SELECT friend_id
    				FROM friend
    				WHERE user_id = $1
				UNION
				SELECT user_id AS friend_id
    				FROM friend
					WHERE friend_id = $1
			)
			ORDER BY updated_at
			LIMIT $2;`

	getFriendsQuery = `
		SELECT friend_id
			FROM friend
			WHERE user_id = $1
		UNION
		SELECT user_id AS friend_id
			FROM friend
			WHERE friend_id = $1;`
)

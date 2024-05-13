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
)

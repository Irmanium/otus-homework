package repository

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
)

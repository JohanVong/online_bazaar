package stmts

const (
	get_user = `
	SELECT 
		user_uid, 
		username, 
		pw_hash, 
		email, 
		COALESCE (phone, '') AS phone,
		country_uid, 
		country,
		history_uid,
		created_at, 
		COALESCE (updated_at, '0001-01-01') AS updated_at, 
		COALESCE (deleted_at, '0001-01-01') AS deleted_at
	FROM users 
	JOIN countries USING (country_uid) 
	JOIN histories USING (history_uid)`
	GET_USER_BY_NAME = get_user + " WHERE username = $1;"
	GET_USER_BY_PK   = get_user + " WHERE user_uid = $1;"

	INSERT_USER = "INSERT INTO users (user_uid, username, pw_hash, email, phone, country_uid, history_uid) VALUES ($1, $2, $3, $4, $5, $6, $7);"
)

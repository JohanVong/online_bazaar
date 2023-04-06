package stmts

const (
	INSERT_HISTORY = "INSERT INTO histories (history_uid, created_at, updated_at, deleted_at) VALUES ($1, $2, $3, $4);"
	UPDATE_HISTORY = "UPDATE histories SET updated_at = $1 WHERE history_uid = $2;"
	DELETE_HISTORY = "UPDATE histories SET deleted_at = $1 WHERE history_uid = $2;"
)

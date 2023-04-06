package stmts

const (
	GET_COUNTRY_PK = "SELECT country_uid FROM countries WHERE country = $1;"
	GET_COUNTRIES  = "SELECT country_uid, country FROM countries;"
)

package db

import (
	"database/sql"
	"errors"

	"github.com/JohanVong/online_bazaar/internal/db/stmts"
	"github.com/JohanVong/online_bazaar/pkg/models"
)

// CountryModel - модель сущности countries
type CountryModel struct {
	DB *sql.DB
}

// GetList() - метод, который достает список всех стран из БД
func (c *CountryModel) GetList() ([]*models.CountryOutput, error) {
	var countries []*models.CountryOutput

	rows, err := c.DB.Query(stmts.GET_COUNTRIES)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, errors.New("No records found")
		}
		return nil, err
	}
	defer rows.Close()

	for rows.Next() {
		c := &models.CountryOutput{}
		err = rows.Scan(&c.UUID, &c.Name)
		if err != nil {
			return nil, err
		}
		countries = append(countries, c)
	}
	if err = rows.Err(); err != nil {
		return nil, err
	}

	return countries, nil
}

// GetByName() - метод, который достает ключ по названию страны
func (c *CountryModel) GetByName(name string) (string, error) {
	var cuid string

	row := c.DB.QueryRow(stmts.GET_COUNTRY_PK, name)
	err := row.Scan(&cuid)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return "", errors.New("Provided country does not exist")
		}
		return "", err
	}

	return cuid, nil
}

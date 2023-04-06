package db

import (
	"database/sql"
	"errors"
	"time"

	"github.com/gofrs/uuid"

	"github.com/JohanVong/online_bazaar/internal/db/stmts"
	"github.com/JohanVong/online_bazaar/pkg/models"
)

// UserModel - модель сущности users
type UserModel struct {
	DB *sql.DB
}

// Insert() - метод для создания новой записи о пользователе
func (u *UserModel) Insert(input *models.UserSignupInput) error {
	var (
		cuid string
		err  error
	)

	huid, _ := uuid.NewV6()
	uid, _ := uuid.NewV6()

	tx, err := u.DB.Begin()
	if err != nil {
		return err
	}

	row := tx.QueryRow(stmts.GET_COUNTRY_PK, input.Country)
	err = row.Scan(&cuid)
	if err != nil {
		tx.Rollback()
		if errors.Is(err, sql.ErrNoRows) {
			return errors.New("Provided country does not exist")
		}
		return err
	}

	_, err = tx.Exec(stmts.INSERT_HISTORY, huid.String(), time.Now(), nil, nil)
	if err != nil {
		tx.Rollback()
		return err
	}

	_, err = tx.Exec(stmts.INSERT_USER, uid.String(), input.Username, input.Password, input.Email, input.Phone, cuid, huid.String())
	if err != nil {
		tx.Rollback()
		return err
	}

	tx.Commit()
	return nil
}

// Get() - метод для получения данных о пользователе по ключу или юзернейму
func (u *UserModel) Get(key string, byPK bool) (*models.UserOutput, error) {
	var (
		stmt string
		uodb models.UserOutput
		err  error
	)

	if byPK {
		stmt = stmts.GET_USER_BY_PK
	} else {
		stmt = stmts.GET_USER_BY_NAME
	}

	row := u.DB.QueryRow(stmt, key)
	err = row.Scan(
		&uodb.UserUID,
		&uodb.Username,
		&uodb.Hash,
		&uodb.Email,
		&uodb.Phone,
		&uodb.CountryUID,
		&uodb.Country,
		&uodb.HistoryUID,
		&uodb.CreatedAt,
		&uodb.UpdatedAt,
		&uodb.DeletedAt,
	)
	if err != nil {
		return nil, err
	}

	return &uodb, nil
}

// Update() - метод для обновления некоторых данных в записи пользователя
func (u *UserModel) Update(uid string, input *models.UserUpdateInput) error {
	var (
		huid    string
		counter int
	)

	tx, err := u.DB.Begin()
	if err != nil {
		return err
	}

	if input.Email != "" {
		counter++
		_, err := tx.Exec("UPDATE users SET email = $1 WHERE user_uid = $2", input.Email, uid)
		if err != nil {
			tx.Rollback()
			return err
		}
	}

	if input.Phone != "" {
		counter++
		_, err := tx.Exec("UPDATE users SET phone = $1 WHERE user_uid = $2", input.Phone, uid)
		if err != nil {
			tx.Rollback()
			return err
		}
	}

	if input.Country != "" {
		counter++
		_, err := tx.Exec("UPDATE users SET country_uid = $1 WHERE user_uid = $2", input.Country, uid)
		if err != nil {
			tx.Rollback()
			return err
		}
	}

	if counter <= 0 {
		return errors.New("Nothing to update")
	}

	row := tx.QueryRow("SELECT history_uid FROM users WHERE user_uid = $1", uid)
	err = row.Scan(&huid)
	if err != nil {
		tx.Rollback()
		return err
	}

	_, err = tx.Exec(stmts.UPDATE_HISTORY, time.Now(), huid)
	if err != nil {
		tx.Rollback()
		return err
	}

	tx.Commit()
	return nil
}

// UpdatePassword() - метод для обновления пароля у записи пользователя
func (u *UserModel) UpdatePassword(uid string, input *models.UpdateUserPasswordInput) error {
	var huid string

	tx, err := u.DB.Begin()
	if err != nil {
		return err
	}

	_, err = tx.Exec("UPDATE users SET pw_hash = $1 WHERE user_uid = $2", input.NewPassword, uid)
	if err != nil {
		tx.Rollback()
		return err
	}

	row := tx.QueryRow("SELECT history_uid FROM users WHERE user_uid = $1", uid)
	err = row.Scan(&huid)
	if err != nil {
		tx.Rollback()
		return err
	}

	_, err = tx.Exec(stmts.UPDATE_HISTORY, time.Now(), huid)
	if err != nil {
		tx.Rollback()
		return err
	}

	tx.Commit()
	return nil
}

/*
Delete() - метод для фейкового удаления пользователя.
На самом деле он заполняет поле deleted_at в истории пользователя.
При следующих запросах миддлвер видит это и блокирует доступ к пользователю,
но запись с его данными остается в БД.
*/
func (u *UserModel) Delete(uid string) error {
	var huid string

	tx, err := u.DB.Begin()
	if err != nil {
		return err
	}

	row := tx.QueryRow("SELECT history_uid FROM users WHERE user_uid = $1", uid)
	err = row.Scan(&huid)
	if err != nil {
		tx.Rollback()
		return err
	}

	_, err = tx.Exec(stmts.DELETE_HISTORY, time.Now(), huid)
	if err != nil {
		tx.Rollback()
		return err
	}

	tx.Commit()
	return nil
}

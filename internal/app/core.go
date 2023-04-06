package app

import (
	"database/sql"
	"fmt"
	"io/ioutil"
	"log"
	"os"

	"github.com/go-playground/validator"
	"github.com/labstack/echo/v4"

	"github.com/JohanVong/online_bazaar/pkg/models"
	"github.com/JohanVong/online_bazaar/pkg/models/db"
	"github.com/JohanVong/online_bazaar/pkg/models/mock"
	"github.com/JohanVong/online_bazaar/tools"
)

// core - ядро приложения
type core struct {
	sign     string
	echo     *echo.Echo
	infoLog  *log.Logger
	errorLog *log.Logger
	users    interface {
		Insert(*models.UserSignupInput) error
		Get(string, bool) (*models.UserOutput, error)
		Update(string, *models.UserUpdateInput) error
		UpdatePassword(string, *models.UpdateUserPasswordInput) error
		Delete(string) error
	}
	countries interface {
		GetList() ([]*models.CountryOutput, error)
		GetByName(string) (string, error)
	}
}

// getConnDB() - функция, устанавливающая соединение с постгрес
func getConnDB() (*sql.DB, error) {
	psqlInfo := fmt.Sprintf("host=%s port=%v user=%s password=%s dbname=%s sslmode=disable",
		os.Getenv("PSQL_HOST"), os.Getenv("PSQL_PORT"), os.Getenv("PSQL_USER"), os.Getenv("PSQL_PASS"), os.Getenv("PSQL_NAME"))

	db, err := sql.Open("postgres", psqlInfo)
	if err != nil {
		return nil, err
	}

	err = db.Ping()
	if err != nil {
		return nil, err
	}

	return db, nil
}

// AssembleAndGo() - собирает ядро и запускает приложение
func AssembleAndGo() {
	conn, err := getConnDB()
	if err != nil {
		panic(err)
	}
	defer conn.Close()

	appCore := &core{
		sign:      os.Getenv("SIGN"),
		echo:      echo.New(),
		infoLog:   log.New(os.Stdout, "INFO:\t", log.Ldate|log.Ltime),
		errorLog:  log.New(os.Stdout, "ERROR:\t", log.Ldate|log.Ltime|log.Lshortfile),
		users:     &db.UserModel{DB: conn},
		countries: &db.CountryModel{DB: conn},
	}
	appCore.echo.Validator = &tools.CustomValidator{Validator: validator.New()}
	appCore.configureRouting()

	appCore.errorLog.Fatal(appCore.echo.Start(":8080"))
}

// assembleTestCore() - собирает тестовое ядро
func assembleTestCore() *core {
	testCore := &core{
		sign:      os.Getenv("TEST_SIGN"),
		echo:      echo.New(),
		infoLog:   log.New(ioutil.Discard, "", 0),
		errorLog:  log.New(ioutil.Discard, "", 0),
		users:     &mock.UserModel{},
		countries: &mock.CountryModel{},
	}
	testCore.echo.Validator = &tools.CustomValidator{Validator: validator.New()}
	testCore.configureRouting()

	return testCore
}

package app

import (
	"fmt"
	"net/http"
	"runtime/debug"

	"github.com/golang-jwt/jwt"
)

// apiResponse - структура ответа приложения
type apiResponse struct {
	Data  interface{} `json:"Data,omitempty"`
	Error string      `json:"Error,omitempty"`
}

// respondOK() - метод приложения для успешного ответа
func (ac *core) respondOK(data interface{}) (int, interface{}) {
	resp := apiResponse{
		Data: "OK",
	}

	if data != nil {
		resp.Data = data
	}
	ac.infoLog.Println("Request successful")

	return http.StatusOK, resp
}

// validationError() - отдает клиенту ошибку валидации
func (ac *core) validationError(err error) (int, interface{}) {
	resp := apiResponse{
		Error: "Data validation failed",
	}
	ac.errorLog.Println(err.Error())

	return http.StatusBadRequest, resp
}

// bindError() - отдает клиенту ошибку о неправильном формате входных данных
func (ac *core) bindError(err error) (int, interface{}) {
	resp := apiResponse{
		Error: "Wrong data format",
	}
	ac.errorLog.Println(err.Error())

	return http.StatusBadRequest, resp
}

// unauthorized() - метод приложения для ответа со статусом Unauthorized
func (ac *core) unauthorized(text string) (int, interface{}) {
	resp := apiResponse{
		Error: text,
	}
	ac.errorLog.Println(text)

	return http.StatusUnauthorized, resp
}

// serverError() - метод приложения для ответа и обработки внутренней ошибки сервера
func (ac *core) serverError(err error) (int, interface{}) {
	resp := apiResponse{
		Error: err.Error(),
	}
	ac.errorLog.Println(err.Error())

	return http.StatusInternalServerError, resp
}

// appPanic() - метод для отдачи наружу паники
func (ac *core) appPanic(err error) (int, interface{}) {
	resp := apiResponse{
		Error: err.Error(),
	}

	trace := fmt.Sprintf("%s\n%s", err.Error(), debug.Stack())
	ac.errorLog.Println(trace)

	return http.StatusInternalServerError, resp
}

// generateToken() - метод для генерации токена
func (ac *core) generateToken(claims jwt.MapClaims, isStandard bool) (string, error) {
	var token *jwt.Token

	if isStandard {
		token = jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	} else {
		token = jwt.NewWithClaims(jwt.SigningMethodHS384, claims)
	}

	t, err := token.SignedString([]byte(ac.sign))
	if err != nil {
		return "", err
	}

	return t, nil
}

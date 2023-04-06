package app

import (
	"errors"
	"fmt"
	"net/http"

	"github.com/golang-jwt/jwt"
	"github.com/labstack/echo/v4"
)

// authorize() - авторизационный миддлвер для пользователя
func (ac *core) authorize(next echo.HandlerFunc) echo.HandlerFunc {
	return func(c echo.Context) error {
		sign := []byte(ac.sign)

		keyFunc := func(t *jwt.Token) (interface{}, error) {
			if t.Method.Alg() != "HS256" {
				return nil, errors.New("Wrong signing method")
			}
			return sign, nil
		}

		auth := c.Request().Header.Get("Authorization")
		auth = auth[7:]

		token, err := jwt.Parse(auth, keyFunc)
		if err != nil {
			c.JSON(http.StatusUnauthorized, map[string]string{"Error": err.Error()})
			return err
		}

		claims := token.Claims.(jwt.MapClaims)
		uid := claims["UID"].(string)

		uo, err := ac.users.Get(uid, true)
		if err != nil {
			c.JSON(http.StatusUnauthorized, map[string]string{"Error": err.Error()})
			return err
		}

		if !uo.DeletedAt.IsZero() {
			c.JSON(http.StatusUnauthorized, map[string]string{"Error": "User was deleted"})
			return err
		}

		c.Set("uid", uid)
		return next(c)
	}
}

// recoverPanic() - миддлвер для обработки паник
func (ac *core) recoverPanic(next echo.HandlerFunc) echo.HandlerFunc {
	return func(c echo.Context) error {
		defer func() {
			if err := recover(); err != nil {
				c.JSON(ac.appPanic(fmt.Errorf("%s", err)))
			}
		}()

		return next(c)
	}
}

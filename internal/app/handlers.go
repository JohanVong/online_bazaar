package app

import (
	"crypto/sha512"
	"encoding/base64"
	"time"

	"github.com/golang-jwt/jwt"
	"github.com/labstack/echo/v4"

	"github.com/JohanVong/online_bazaar/pkg/models"
)

// testAlive() - проверка типа пинг-понг
func (ac *core) testAlive(c echo.Context) error {
	return c.JSON(ac.respondOK("We are ok!"))
}

// signupUser() - хэндлер для регистрации пользователя
func (ac *core) signupUser(c echo.Context) error {
	var (
		user models.UserSignupInput
		err  error
	)

	if err = c.Bind(&user); err != nil {
		return c.JSON(ac.bindError(err))
	}

	if err = c.Validate(&user); err != nil {
		return c.JSON(ac.validationError(err))
	}

	hash64 := sha512.Sum512([]byte(user.Password))
	hash := base64.StdEncoding.EncodeToString(hash64[:])
	user.Password = hash

	err = ac.users.Insert(&user)
	if err != nil {
		return c.JSON(ac.serverError(err))
	}

	return c.JSON(ac.respondOK("OK"))
}

// loginUser() - хэндлер для аутентификации пользователя
func (ac *core) loginUser(c echo.Context) error {
	var (
		uli  models.UserLoginInput
		ulo  models.UserLoginOutput
		uodb *models.UserOutput
		err  error
	)

	if err = c.Bind(&uli); err != nil {
		return c.JSON(ac.bindError(err))
	}

	uodb, err = ac.users.Get(uli.Username, false)
	if err != nil {
		return c.JSON(ac.unauthorized("Wrong credentials provided"))
	}

	hash64 := sha512.Sum512([]byte(uli.Password))
	hash := base64.StdEncoding.EncodeToString(hash64[:])
	if hash != uodb.Hash {
		return c.JSON(ac.unauthorized("Wrong credentials provided"))
	}

	if !uodb.DeletedAt.IsZero() {
		return c.JSON(ac.unauthorized("User was deleted"))
	}

	claims := jwt.MapClaims{}
	claims["UID"] = uodb.UserUID
	claims["exp"] = time.Now().Add(time.Minute * 3).Unix()

	// Create token with claims & generate encoded one
	token, err := ac.generateToken(claims, true)
	if err != nil {
		return c.JSON(ac.serverError(err))
	}

	ulo.Token = token
	return c.JSON(ac.respondOK(ulo))
}

// updateUser() - хэндлер для обновления данных пользователя
func (ac *core) updateUser(c echo.Context) error {
	var (
		uus models.UserUpdateInput
		err error
	)

	uid := c.Get("uid").(string)

	if err = c.Bind(&uus); err != nil {
		return c.JSON(ac.bindError(err))
	}

	if uus.Email != "" {
		err = c.Validate(&uus)
		if err != nil {
			return c.JSON(ac.validationError(err))
		}
	}

	if uus.Country != "" {
		cuid, err := ac.countries.GetByName(uus.Country)
		if err != nil {
			return c.JSON(ac.serverError(err))
		}

		uus.Country = cuid
	}

	err = ac.users.Update(uid, &uus)
	if err != nil {
		return c.JSON(ac.serverError(err))
	}

	return c.JSON(ac.respondOK("OK"))
}

// updateUserPassword() - хэндлер для обновления пароля пользователя
func (ac *core) updateUserPassword(c echo.Context) error {
	var (
		upi models.UpdateUserPasswordInput
		err error
	)

	uid := c.Get("uid").(string)

	if err = c.Bind(&upi); err != nil {
		return c.JSON(ac.bindError(err))
	}

	if err = c.Validate(&upi); err != nil {
		return c.JSON(ac.validationError(err))
	}

	hash64 := sha512.Sum512([]byte(upi.NewPassword))
	hash := base64.StdEncoding.EncodeToString(hash64[:])
	upi.NewPassword = hash

	err = ac.users.UpdatePassword(uid, &upi)
	if err != nil {
		return c.JSON(ac.serverError(err))
	}

	return c.JSON(ac.respondOK("OK"))
}

// deleteUser() - хэндлер для "удаления" пользователя.
// Важно: Пользователь не будет удален, но будет деактивирован
func (ac *core) deleteUser(c echo.Context) error {
	uid := c.Get("uid").(string)

	err := ac.users.Delete(uid)
	if err != nil {
		return c.JSON(ac.serverError(err))
	}

	return c.JSON(ac.respondOK("OK"))
}

// getCountries - хэндлер для получения списка стран из БД
func (ac *core) getCountries(c echo.Context) error {
	co, err := ac.countries.GetList()
	if err != nil {
		return c.JSON(ac.serverError(err))
	}

	return c.JSON(ac.respondOK(co))
}

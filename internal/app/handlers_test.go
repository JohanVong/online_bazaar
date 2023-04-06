package app

import (
	"fmt"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestPing(t *testing.T) {
	want := `{"Data":"We are ok!"}`
	req := httptest.NewRequest("", "/", nil)
	rec := httptest.NewRecorder()

	testCore := assembleTestCore()
	c := testCore.echo.NewContext(req, rec)

	if assert.NoError(t, testCore.testAlive(c)) {
		assert.Equal(t, http.StatusOK, rec.Code)
		assert.Equal(t, want, strings.TrimSpace(rec.Body.String()))
	}
}

func TestSignupUser(t *testing.T) {
	testCore := assembleTestCore()

	tests := []struct {
		input    string
		wantCode int
		wantBody interface{}
	}{
		{ // good request
			`{
				"Username":"TestUser",
				"Password":"TestPassword",
				"Email":"testuser@mail.test",
				"Phone":"87776665544",
				"Country":"TestCountry"
			}`,
			200,
			`{"Data":"OK"}`,
		},
		{ // wrong json
			`{
				"Username":"TestUser",
				"Password":"TestPassword",
				"Email":"testuser@mail.test",
				"Phone":"87776665544",
				"Country":"TestCountry",
			}`,
			400,
			`{"Error":"Wrong data format"}`,
		},
		{ // username validation error
			`{
				"Username":"",
				"Password":"TestPassword",
				"Email":"testuser@mail.test",
				"Phone":"87776665544",
				"Country":"TestCountry"
			}`,
			400,
			`{"Error":"Data validation failed"}`,
		},
		{ // password validation error
			`{
				"Username":"TestUser",
				"Password":"123",
				"Email":"testuser@mail.test",
				"Phone":"87776665544",
				"Country":"TestCountry"
			}`,
			400,
			`{"Error":"Data validation failed"}`,
		},
		{ // email validation error
			`{
				"Username":"TestUser",
				"Password":"TestPassword",
				"Email":"testmail.test",
				"Phone":"87776665544",
				"Country":"TestCountry"
			}`,
			400,
			`{"Error":"Data validation failed"}`,
		},
		{ // unique constraint violation
			`{
				"Username":"Exists",
				"Password":"TestPassword",
				"Email":"testuser@mail.test",
				"Phone":"87776665544",
				"Country":"TestCountry"
			}`,
			500,
			`{"Error":"duplicate key value violates unique constraint"}`,
		},
		{ // non-existing country
			`{
				"Username":"Test",
				"Password":"12345678",
				"Email":"testuser@mail.com",
				"Phone":"87776665544",
				"Country":"No existing country"
			}`,
			500,
			`{"Error":"Provided country does not exist"}`,
		},
	}

	for _, tt := range tests {
		req := httptest.NewRequest("", "/", strings.NewReader(tt.input))
		req.Header.Set("Content-Type", "application/json")
		rec := httptest.NewRecorder()
		c := testCore.echo.NewContext(req, rec)

		if assert.NoError(t, testCore.signupUser(c)) {
			assert.Equal(t, tt.wantCode, rec.Code)
			assert.Equal(t, tt.wantBody, strings.TrimSpace(rec.Body.String()))
		}
	}

}

func TestLoginUser(t *testing.T) {
	testCore := assembleTestCore()

	tests := []struct {
		input    string
		wantCode int
		wantBody interface{}
	}{
		{ // good request
			`{
				"Username":"TestUser",
				"Password": "TestPassword"
			}`,
			200,
			`{"Data":`,
		},
		{ // wrong json
			`{
				"Username":"TestUser",
				"Password": "TestPassword",
			}`,
			400,
			`{"Error":"Wrong data format"}`,
		},
		{ // wrong username
			`{
				"Username":"Test",
				"Password": "TestPassword"
			}`,
			401,
			`{"Error":"Wrong credentials provided"}`,
		},
		{ // wrong password
			`{
				"Username":"TestUser",
				"Password": "Test"
			}`,
			401,
			`{"Error":"Wrong credentials provided"}`,
		},
		{ // deleted user
			`{
				"Username":"DeletedUser",
				"Password": "TestPassword"
			}`,
			401,
			`{"Error":"User was deleted"}`,
		},
	}

	for _, tt := range tests {
		req := httptest.NewRequest("", "/", strings.NewReader(tt.input))
		req.Header.Set("Content-Type", "application/json")
		rec := httptest.NewRecorder()
		c := testCore.echo.NewContext(req, rec)

		if assert.NoError(t, testCore.loginUser(c)) {
			assert.Equal(t, tt.wantCode, rec.Code)
			if rec.Code == 200 {
				assert.Equal(t, tt.wantBody, strings.TrimSpace(rec.Body.String()[:8]))
			} else {
				assert.Equal(t, tt.wantBody, strings.TrimSpace(rec.Body.String()))
			}
		}
	}
}

func TestUpdateUser(t *testing.T) {
	testCore := assembleTestCore()

	tests := []struct {
		input    string
		wantCode int
		wantBody interface{}
	}{
		{ // good request
			`{
				"Email": "",
				"Phone": "",
				"Country": "TestCountry2"
			}`,
			200,
			`{"Data":"OK"}`,
		},
		{ // wrong json
			`{
				"Email": "",
				"Phone": "",
				"Country": "TestCountry2",
			}`,
			400,
			`{"Error":"Wrong data format"}`,
		},
		{ // validation error, wrong format email
			`{
				"Email": "dfsgahf",
				"Phone": "",
				"Country": ""
			}`,
			400,
			`{"Error":"Data validation failed"}`,
		},
		{ // non-existing country
			`{
				"Email": "",
				"Phone": "",
				"Country": "Non-Existing Country"
			}`,
			500,
			`{"Error":"Provided country does not exist"}`,
		},
		{ // nothing to update
			`{
				"Email": "",
				"Phone": "",
				"Country": ""
			}`,
			500,
			`{"Error":"Nothing to update"}`,
		},
	}

	for _, tt := range tests {
		req := httptest.NewRequest("", "/", strings.NewReader(tt.input))
		req.Header.Set("Content-Type", "application/json")
		rec := httptest.NewRecorder()
		c := testCore.echo.NewContext(req, rec)
		c.Set("uid", "uuid.v6[1]")

		if assert.NoError(t, testCore.updateUser(c)) {
			assert.Equal(t, tt.wantCode, rec.Code)
			assert.Equal(t, tt.wantBody, strings.TrimSpace(rec.Body.String()))
		}
	}
}

func TestUpdateUserPassword(t *testing.T) {
	testCore := assembleTestCore()

	tests := []struct {
		input    string
		uid      string
		wantCode int
		wantBody interface{}
	}{
		{ // good request
			`{
				"NewPassword": "12345678"
			}`,
			"uuid.v6[1]",
			200,
			`{"Data":"OK"}`,
		},
		{ // wrong json
			`{
				"NewPassword": "12345678",
			}`,
			"uuid.v6[1]",
			400,
			`{"Error":"Wrong data format"}`,
		},
		{ // validation error (required)
			`{
				"NewPassword": ""
			}`,
			"uuid.v6[1]",
			400,
			`{"Error":"Data validation failed"}`,
		},
		{ // validation error (min len)
			`{
				"NewPassword": ""
			}`,
			"uuid.v6[1]",
			400,
			`{"Error":"Data validation failed"}`,
		},
		{ // non-existing user
			`{
				"NewPassword": "12345678"
			}`,
			"uuid.v6[93]",
			500,
			`{"Error":"No record found"}`,
		},
	}

	for _, tt := range tests {
		req := httptest.NewRequest("", "/", strings.NewReader(tt.input))
		req.Header.Set("Content-Type", "application/json")
		rec := httptest.NewRecorder()
		c := testCore.echo.NewContext(req, rec)
		c.Set("uid", tt.uid)

		if assert.NoError(t, testCore.updateUserPassword(c)) {
			assert.Equal(t, tt.wantCode, rec.Code)
			assert.Equal(t, tt.wantBody, strings.TrimSpace(rec.Body.String()))
		}
	}
}

func TestDeleteUser(t *testing.T) {
	testCore := assembleTestCore()

	tests := []struct {
		uid      string
		wantCode int
		wantBody interface{}
	}{
		{ // good request
			"uuid.v6[1]",
			200,
			`{"Data":"OK"}`,
		},
		{ // no such user
			"uuid.v6[2]",
			500,
			`{"Error":"No record found"}`,
		},
	}

	i := 1
	for _, tt := range tests {
		req := httptest.NewRequest("", "/", nil)
		rec := httptest.NewRecorder()
		c := testCore.echo.NewContext(req, rec)
		c.Set("uid", fmt.Sprintf("uuid.v6[%v]", i))

		if assert.NoError(t, testCore.deleteUser(c)) {
			assert.Equal(t, tt.wantCode, rec.Code)
			assert.Equal(t, tt.wantBody, strings.TrimSpace(rec.Body.String()))
		}
		i++
	}
}

func TestGetCountries(t *testing.T) {
	wantBody := `{"Data":[{"Name":"TestCountry"},{"Name":"TestCountry2"},{"Name":"TestCountry3"}]}`

	testCore := assembleTestCore()
	req := httptest.NewRequest("", "/", nil)
	rec := httptest.NewRecorder()
	c := testCore.echo.NewContext(req, rec)

	if assert.NoError(t, testCore.getCountries(c)) {
		assert.Equal(t, http.StatusOK, rec.Code)
		assert.Equal(t, strings.TrimSpace(wantBody), strings.TrimSpace(rec.Body.String()))
	}
}

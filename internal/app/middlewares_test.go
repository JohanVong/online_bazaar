package app

import (
	"io/ioutil"
	"net/http"
	"strings"
	"testing"
	"time"

	"github.com/golang-jwt/jwt"
	"github.com/stretchr/testify/assert"
)

func TestAuthorize(t *testing.T) {
	var (
		token string
		err   error
	)

	tests := []struct {
		alg      int
		uid      string
		wantCode int
		wantBody interface{}
	}{
		{ // pass middleware
			2,
			"uuid.v6[1]",
			http.StatusOK,
			`{"Data":"We are ok!"}`,
		},
		{ // wrongly signed token
			3,
			"uuid.v6[1]",
			http.StatusUnauthorized,
			`{"Error":"Wrong signing method"}`,
		},
		{ // bad token, parsing error
			0,
			"uuid.v6[1]",
			http.StatusUnauthorized,
			`{"Error":"token contains an invalid number of segments"}`,
		},
		{ // unexisting user tries to authorize
			2,
			"uuid.v6[93]",
			http.StatusUnauthorized,
			`{"Error":"No record found"}`,
		},
		{ // deleted user tries to authorize
			2,
			"uuid.v6[4]",
			http.StatusUnauthorized,
			`{"Error":"User was deleted"}`,
		},
		{ // panic
			2,
			"panic",
			http.StatusInternalServerError,
			`{"Error":"test panic!"}`,
		},
	}

	testCore := assembleTestCore()
	go func() {
		testCore.errorLog.Fatal(testCore.echo.Start(":63246"))
	}()
	time.Sleep(1 * time.Second)

	for _, tt := range tests {
		claims := jwt.MapClaims{}
		claims["UID"] = tt.uid
		claims["exp"] = time.Now().Add(time.Minute * 3).Unix()

		switch tt.alg {
		case 0:
			token = "bad"
		case 2:
			token, err = testCore.generateToken(claims, true)
		default:
			token, err = testCore.generateToken(claims, false)
		}
		if err != nil {
			t.Fatal(err)
		}

		req, err := http.NewRequest(http.MethodGet, "http://127.0.0.1:63246/test/auth", nil)
		if err != nil {
			t.Fatal(err)
		}
		req.Header.Set("Authorization", "Bearer "+token)

		res, err := http.DefaultClient.Do(req)
		if res != nil {
			defer res.Body.Close()
		}
		if err != nil {
			t.Fatal(err)
		}
		body, err := ioutil.ReadAll(res.Body)

		assert.Equal(t, tt.wantCode, res.StatusCode)
		assert.Equal(t, tt.wantBody, strings.TrimSpace(string(body)))
	}
}

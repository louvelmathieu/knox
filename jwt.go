package knox

import (
	"encoding/json"
	"errors"
	"fmt"
	"github.com/dgrijalva/jwt-go"
	"net/http"
	"os"
	"strings"
)

func VerifyJWTToken(tokenString string) (map[string]interface{}, error) {
	if strings.Index(tokenString, "Bearer ") == 0 {
		tokenString = tokenString[len("Bearer "):]
	}

	token, err := jwt.Parse(tokenString, func(token *jwt.Token) (interface{}, error) {

		// Don't forget to validate the alg is what you expect:
		if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, fmt.Errorf("Unexpected signing method: %v.", token.Header["alg"])
		}

		return []byte(os.Getenv("JWT_SECRET")), nil
	})

	if err == nil {
		if info, ok := token.Claims.(jwt.MapClaims); ok && token.Valid {
			return info, nil
		} else {
			return nil, errors.New("Invalid token")
		}
	}

	return nil, err
}

func jwtMiddleware(next http.Handler, requireAuthentification bool) http.Handler {

	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		// Skip Auth for OPTIONS
		if r.Method == "OPTIONS" {
			next.ServeHTTP(w, r)
			return
		}

		// authentification and check Authorization Header
		var tokenString = r.Header.Get("Authorization")
		if tokenString != "" {

			data, err := VerifyJWTToken(tokenString)
			if err != nil {
				error := Error{
					http.StatusBadRequest,
					"Invalid token",
				}
				error.Send(w, http.StatusUnauthorized)
				return
			}

			jdata, _ := json.Marshal(data)
			r.Header.Set("JWT", string(jdata))

			next.ServeHTTP(w, r)
			return
		}

		// Some path can have optinal authentification
		if requireAuthentification == true {
			error := Error{
				http.StatusUnauthorized,
				"Empty token",
			}
			error.Send(w, http.StatusUnauthorized)
			return
		} else {
			next.ServeHTTP(w, r)
			return
		}
	})
}

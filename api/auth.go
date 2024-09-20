package api

import (
	"net/http"
	"os"

	"github.com/golang-jwt/jwt"
)

// Middleware для проверки прав пользователя
func Auth(next http.HandlerFunc) http.HandlerFunc {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		pass := os.Getenv("TODO_PASSWORD")
		if len(pass) > 0 {
			var requestJwt string

			cookie, err := r.Cookie("token")
			if err == nil {
				requestJwt = cookie.Value
			}

			// Проверка токена
			jwtToken, err := jwt.Parse(requestJwt, func(t *jwt.Token) (secret interface{}, err error) {
				return []byte(pass), nil
			})
			if err != nil || !jwtToken.Valid {
				returnError(w, "ошибка авторизации", 401)
				return
			}
		}
		next(w, r)
	})
}

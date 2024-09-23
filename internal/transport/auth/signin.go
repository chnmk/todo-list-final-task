package auth

import (
	"bytes"
	"encoding/json"
	"log"
	"net/http"

	"github.com/chnmk/todo-list-final-task/internal/transport"
	"github.com/golang-jwt/jwt"
)

// Обрабатывает запросы к /api/signin. При успешном запросе возвращает JSON с JWT-токеном.
func AuthHandler(w http.ResponseWriter, r *http.Request) {
	if transport.EnvPassword == "" {
		transport.ReturnError(w, "авторизация не предусмотрена", 400)
		return
	}
	if r.Method != http.MethodPost {
		transport.ReturnError(w, "неожиданный метод запроса, ожидался POST", 400)
		return
	}

	// Получение данных
	var passwordStruct transport.PasswordStruct
	var buf bytes.Buffer

	_, err := buf.ReadFrom(r.Body)
	if err != nil {
		transport.ReturnError(w, err.Error(), 500)
		return
	}

	if err = json.Unmarshal(buf.Bytes(), &passwordStruct); err != nil {
		transport.ReturnError(w, err.Error(), 500)
		return
	}

	// Проверка пароля и создание JWT токена
	if passwordStruct.Password != transport.EnvPassword {
		transport.ReturnError(w, "неправильный пароль", 401)
		return
	}

	jwtToken := jwt.New(jwt.SigningMethodHS256)

	signedToken, err := jwtToken.SignedString([]byte(transport.EnvPassword))
	if err != nil {
		transport.ReturnError(w, err.Error(), 500)
		return
	}

	// Вывод токена в консоль для тестов
	log.Println("JWT ТОКЕН: " + signedToken)

	var tokenStruct transport.TokenStruct
	tokenStruct.Token = signedToken

	resp, err := json.Marshal(tokenStruct)
	if err != nil {
		transport.ReturnError(w, err.Error(), 500)
		return
	}

	w.WriteHeader(http.StatusOK)
	w.Write([]byte(resp))
}

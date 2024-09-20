package api

import (
	"bytes"
	"encoding/json"
	"log"
	"net/http"

	"github.com/golang-jwt/jwt"
)

// Обрабатывает запросы к /api/signin. При успешном запросе возвращает JSON с JWT-токеном.
func AuthHandler(w http.ResponseWriter, r *http.Request) {
	if EnvPassword == "" {
		returnError(w, "авторизация не предусмотрена", 400)
		return
	}
	if r.Method != http.MethodPost {
		returnError(w, "неожиданный метод запроса, ожидался POST", 400)
		return
	}

	// Получение данных
	var passwordStruct PasswordStruct
	var buf bytes.Buffer

	_, err := buf.ReadFrom(r.Body)
	if err != nil {
		returnError(w, err.Error(), 500)
		return
	}

	if err = json.Unmarshal(buf.Bytes(), &passwordStruct); err != nil {
		returnError(w, err.Error(), 500)
		return
	}

	// Проверка пароля и создание JWT токена
	if passwordStruct.Password != EnvPassword {
		returnError(w, "неправильный пароль", 401)
		return
	}

	jwtToken := jwt.New(jwt.SigningMethodHS256)

	signedToken, err := jwtToken.SignedString([]byte(EnvPassword))
	if err != nil {
		returnError(w, err.Error(), 500)
		return
	}

	// Вывод токена в консоль для тестов
	log.Println("JWT ТОКЕН: " + signedToken)

	var tokenStruct TokenStruct
	tokenStruct.Token = signedToken

	resp, err := json.Marshal(tokenStruct)
	if err != nil {
		returnError(w, err.Error(), 500)
		return
	}

	w.WriteHeader(http.StatusOK)
	w.Write([]byte(resp))
}

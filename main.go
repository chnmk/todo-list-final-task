package main

import (
	"fmt"
	"log"
	"net/http"
	"os"

	"github.com/chnmk/todo-list-final-task/tests"
	"github.com/joho/godotenv"
)

var webDir = "./web/"

func main() {
	port := getEnv()
	fmt.Printf("Port: %s\n", port)

	http.Handle("/", http.FileServer(http.Dir(webDir)))
	err := http.ListenAndServe(port, nil)
	if err != nil {
		panic(err)
	}
}

func getEnv() (port string) {
	// Попытка найти .env файл
	err := godotenv.Load()
	if err != nil {
		log.Print(".env file not found")
	}

	// Загружает TODO_PORT из окружения
	// Если он отсутствует, пишет стандартный из settings.go
	port, exists := os.LookupEnv("TODO_PORT")
	if exists {
		port = fmt.Sprintf(":%s", port)
	} else {
		port = fmt.Sprintf(":%d", tests.Port)
	}

	return
}

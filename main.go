package main

import (
	"fmt"
	"log"
	"net/http"
	"os"
	"strings"

	"github.com/chnmk/todo-list-final-task/internal/database"
	"github.com/chnmk/todo-list-final-task/internal/transport"
	"github.com/chnmk/todo-list-final-task/internal/transport/auth"
	"github.com/chnmk/todo-list-final-task/internal/transport/middleware"
	"github.com/chnmk/todo-list-final-task/tests"
	"github.com/joho/godotenv"
)

var webDir = "./web/"

func main() {
	log.Println("Initialization...")

	// Переменные окружения
	port, databaseDir, password, authRequired := getEnv()
	log.Println("Port: " + port)

	// Подключение к БД
	db := database.SetupDB(databaseDir)
	defer db.Close()

	transport.DatabaseDir = databaseDir
	transport.DatabaseFile = db

	// Маршрутизация
	http.HandleFunc("/api/signin", auth.AuthHandler)
	http.HandleFunc("/api/nextdate", transport.NextDate)

	if !authRequired {
		http.HandleFunc("/api/task", transport.TaskRequest)
		http.HandleFunc("/api/tasks", transport.TasksRequest)
		http.HandleFunc("/api/task/done", transport.TaskDone)
	} else {
		transport.EnvPassword = password
		http.HandleFunc("/api/task", middleware.Auth(transport.TaskRequest))
		http.HandleFunc("/api/tasks", middleware.Auth(transport.TasksRequest))
		http.HandleFunc("/api/task/done", middleware.Auth(transport.TaskDone))
	}

	http.Handle("/", http.FileServer(http.Dir(webDir)))

	log.Println("Starting server...")

	// Запуск сервера
	err := http.ListenAndServe(port, nil)
	if err != nil {
		panic(err)
	}

	log.Println("Shutting down...")
}

func getEnv() (port string, dbpath string, password string, authRequired bool) {
	// Попытка найти .env файл
	err := godotenv.Load()
	if err != nil {
		log.Print(".env file not found")
	}

	log.Print("Reading from .env file...")

	// Загружает TODO_PORT из окружения
	// Если он отсутствует, пишет стандартный из settings.go
	port, exists := os.LookupEnv("TODO_PORT")
	if exists {
		log.Print("TODO_PORT: " + port)
		port = fmt.Sprintf(":%s", port)
	} else {
		log.Print("TODO_PORT not found in env variables, using default")
		port = fmt.Sprintf(":%d", tests.Port)
	}

	// Загружает путь к базе данных
	dbpath, exists = os.LookupEnv("TODO_DBFILE")
	if exists {
		log.Print("TODO_DBFILE: " + dbpath)
		dbpath = strings.ReplaceAll(dbpath, "../", "")
	} else {
		log.Print("TODO_DBFILE not found in env variables, using default")
		dbpath = strings.ReplaceAll(tests.DBFile, "../", "")
	}

	// Загружает пароль
	password, exists = os.LookupEnv("TODO_PASSWORD")
	if exists {
		log.Print("TODO_PASSWORD: " + password)
		authRequired = true
	} else {
		log.Print("TODO_PASSWORD not found in env variables, using default")
	}

	return
}

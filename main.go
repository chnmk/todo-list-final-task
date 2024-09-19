package main

import (
	"fmt"
	"log"
	"net/http"
	"os"
	"strings"

	"github.com/chnmk/todo-list-final-task/api"
	"github.com/chnmk/todo-list-final-task/db"
	"github.com/chnmk/todo-list-final-task/tests"
	"github.com/joho/godotenv"
)

var webDir = "./web/"

func main() {
	port, databaseDir, password, authRequired := getEnv()
	fmt.Printf("Port: %s\n", port)

	db.SetupDB(databaseDir)
	api.DatabaseDir = databaseDir

	http.HandleFunc("/api/signin", api.AuthHandler)
	http.HandleFunc("/api/nextdate", api.NextDate)

	if !authRequired {
		http.HandleFunc("/api/task", api.TaskRequest)
		http.HandleFunc("/api/tasks", api.TasksRequest)
		http.HandleFunc("/api/task/done", api.TaskDone)
	} else {
		api.EnvPassword = password
		http.HandleFunc("/api/task", api.Auth(api.TaskRequest))
		http.HandleFunc("/api/tasks", api.Auth(api.TasksRequest))
		http.HandleFunc("/api/task/done", api.Auth(api.TaskDone))
	}

	http.Handle("/", http.FileServer(http.Dir(webDir)))

	err := http.ListenAndServe(port, nil)
	if err != nil {
		panic(err)
	}
}

func getEnv() (port string, dbpath string, password string, authRequired bool) {
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

	// Загружает путь к базе данных
	dbpath, exists = os.LookupEnv("TODO_DBFILE")
	if exists {
		dbpath = strings.ReplaceAll(dbpath, "../", "")
	} else {
		dbpath = strings.ReplaceAll(tests.DBFile, "../", "")
	}

	// Загружает пароль
	password, exists = os.LookupEnv("TODO_PASSWORD")
	if exists {
		authRequired = true
	}

	return
}

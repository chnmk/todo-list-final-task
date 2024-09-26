package transport

import (
	"bytes"
	"database/sql"
	"encoding/json"
	"errors"
	"log"
	"net/http"
	"time"

	"github.com/chnmk/todo-list-final-task/internal/services"
	"github.com/chnmk/todo-list-final-task/tests"
)

type TasksArray struct {
	Tasks []services.Task `json:"tasks"`
}

type ResponseValid struct {
	Id string `json:"id"`
}

type ResponseInvalid struct {
	Error string `json:"error"`
}

type PasswordStruct struct {
	Password string `json:"password"`
}

type TokenStruct struct {
	Token string `json:"token"`
}

var EnvPassword string
var DatabaseFile *sql.DB
var DatabaseDir = tests.DBFile

// Выбирает хендлер в зависимости от метода запроса к /api/task, либо возвращает ошибку
func TaskRequest(w http.ResponseWriter, r *http.Request) {
	log.Println("New request to /api/task, method: " + r.Method)

	switch r.Method {
	case http.MethodPost:
		taskPOST(w, r)
	case http.MethodGet:
		taskGET(w, r)
	case http.MethodPut:
		taskPUT(w, r)
	case http.MethodDelete:
		taskDELETE(w, r)
	default:
		log.Println("Error: unexpected request method")
		ReturnError(w, "неожиданный метод запроса", 500)
		return
	}
}

// Записывает текст ошибки в JSON-файл и возвращает ответ с указанным кодом ошибки.
func ReturnError(w http.ResponseWriter, msg string, status int) {
	var e ResponseInvalid
	e.Error = msg

	resp, err := json.Marshal(e)
	if err != nil {
		e.Error = "ошибка при записи ответа"
		ReturnError(w, msg, 500)
	}

	if status == 400 {
		w.WriteHeader(http.StatusBadRequest)
	} else if status == 401 {
		w.WriteHeader(http.StatusUnauthorized)
	} else {
		w.WriteHeader(http.StatusInternalServerError)
	}

	log.Println("Error: " + e.Error)
	w.Write(resp)
}

// Читает тело запроса и проверяет полученные данные на корректность.
func checkRequest(w http.ResponseWriter, r *http.Request) (services.Task, error) {
	now := time.Now().Format("20060102")
	nowParsed, err := time.Parse("20060102", now)
	if err != nil {
		ReturnError(w, err.Error(), 500)
		return services.Task{}, err
	}

	var task services.Task
	var buf bytes.Buffer

	_, err = buf.ReadFrom(r.Body)
	if err != nil {
		ReturnError(w, err.Error(), 500)
		return services.Task{}, err
	}

	if err = json.Unmarshal(buf.Bytes(), &task); err != nil {
		ReturnError(w, err.Error(), 500)
		return services.Task{}, err
	}

	// Проверка на наличие title
	if task.Title == "" {
		ReturnError(w, "отсутствует название задачи", 400)
		return services.Task{}, errors.New("отсутствует название задачи")
	}

	// Если дата пустая или не указана, берётся сегодняшнее число
	if task.Date == "" {
		task.Date = now
	}

	// Проверяет формат даты 20060102
	dateParsed, err := time.Parse("20060102", task.Date)
	if err != nil {
		ReturnError(w, "некорректный формат времени", 400)
		return services.Task{}, err
	}

	// Если дата меньше сегодняшнего числа...
	if dateParsed.Before(nowParsed) {
		if task.Repeat != "" {
			// При указанном правиле повторения нужно вычислить и записать в таблицу дату выполнения
			task.Date, err = services.NextDate(time.Now(), task.Date, task.Repeat)
			if err != nil {
				ReturnError(w, err.Error(), 500)
				return services.Task{}, err
			}
		} else {
			// Если правило повторения не указано или равно пустой строке, подставляется сегодняшнее число
			task.Date = time.Now().Format("20060102")
		}

	}

	return task, nil
}

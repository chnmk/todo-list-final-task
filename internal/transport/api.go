package transport

import (
	"bytes"
	"encoding/json"
	"errors"
	"net/http"
	"time"

	"github.com/chnmk/todo-list-final-task/internal/services"
	"github.com/chnmk/todo-list-final-task/tests"
)

type Task struct {
	Id      string `json:"id"`
	Date    string `json:"date"`
	Title   string `json:"title"`
	Comment string `json:"comment"`
	Repeat  string `json:"repeat"`
}

type TasksArray struct {
	Tasks []Task `json:"tasks"`
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
var DatabaseDir = tests.DBFile

// Выбирает хендлер в зависимости от метода запроса к /api/task, либо возвращает ошибку
func TaskRequest(w http.ResponseWriter, r *http.Request) {
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

	w.Write(resp)
}

// Читает тело запроса и проверяет полученные данные на корректность.
func checkRequest(w http.ResponseWriter, r *http.Request) (Task, error) {
	now := time.Now().Format("20060102")
	nowParsed, err := time.Parse("20060102", now)
	if err != nil {
		ReturnError(w, err.Error(), 500)
		return Task{}, err
	}

	var task Task
	var buf bytes.Buffer

	_, err = buf.ReadFrom(r.Body)
	if err != nil {
		ReturnError(w, err.Error(), 500)
		return Task{}, err
	}

	if err = json.Unmarshal(buf.Bytes(), &task); err != nil {
		ReturnError(w, err.Error(), 500)
		return Task{}, err
	}

	// Проверка на наличие title
	if task.Title == "" {
		ReturnError(w, "отсутствует название задачи", 400)
		return Task{}, errors.New("отсутствует название задачи")
	}

	// Если дата пустая или не указана, берётся сегодняшнее число
	if task.Date == "" {
		task.Date = now
	}

	// Проверяет формат даты 20060102
	dateParsed, err := time.Parse("20060102", task.Date)
	if err != nil {
		ReturnError(w, "некорректный формат времени", 400)
		return Task{}, err
	}

	// Если дата меньше сегодняшнего числа...
	if dateParsed.Before(nowParsed) {
		if task.Repeat != "" {
			// При указанном правиле повторения нужно вычислить и записать в таблицу дату выполнения
			task.Date, err = services.NextDate(time.Now(), task.Date, task.Repeat)
			if err != nil {
				ReturnError(w, err.Error(), 500)
				return Task{}, err
			}
		} else {
			// Если правило повторения не указано или равно пустой строке, подставляется сегодняшнее число
			task.Date = time.Now().Format("20060102")
		}

	}

	return task, nil
}

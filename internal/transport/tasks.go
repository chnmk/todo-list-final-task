package transport

import (
	"database/sql"
	"encoding/json"
	"log"
	"net/http"
	"time"

	"github.com/chnmk/todo-list-final-task/internal/services"
)

// Обрабатывает запросы к /api/tasks, возвращает ближайшие 10 задач.
func TasksRequest(w http.ResponseWriter, r *http.Request) {
	log.Println("New request to /api/tasks")
	w.Header().Set("Content-Type", "application/json; charset=UTF-8")

	// Обработка параметров поиска
	var searchDate bool
	search := r.FormValue("search")

	searchParsed, ok := time.Parse("02.01.2006", search)
	if ok == nil {
		searchDate = true
		search = searchParsed.Format("20060102")
	} else {
		search = "%" + search + "%" // Для подстановки в запрос с LIKE в SQL
	}

	// Составление строки запроса
	var tasks []services.Task

	query := "SELECT * FROM scheduler ORDER BY date LIMIT 10"
	if search != "" && !searchDate {
		query = "SELECT * FROM scheduler WHERE title LIKE :search OR comment LIKE :search ORDER BY date LIMIT 10"
	} else if search != "" && searchDate {
		query = "SELECT * FROM scheduler WHERE date = :search ORDER BY date LIMIT 10"
	}

	// Выполнение нужного запроса
	rows, err := DatabaseFile.Query(query, sql.Named("search", search))
	if err != nil {
		ReturnError(w, err.Error(), 500)
		return
	}

	defer rows.Close()

	// Чтение данных из базы
	for rows.Next() {
		task := services.Task{}

		err := rows.Scan(&task.Id, &task.Date, &task.Title, &task.Comment, &task.Repeat)
		if err != nil {
			ReturnError(w, err.Error(), 500)
			return
		}

		tasks = append(tasks, task)
	}

	// Запись ответа
	var result TasksArray

	if tasks != nil {
		result.Tasks = tasks
	} else {
		result.Tasks = []services.Task{}
	}

	resp, err := json.Marshal(result)
	if err != nil {
		ReturnError(w, err.Error(), 500)
		return
	}

	log.Println("Success!")
	w.WriteHeader(http.StatusOK)
	w.Write(resp)
}

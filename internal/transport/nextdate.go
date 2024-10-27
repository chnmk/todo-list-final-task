package transport

import (
	"net/http"
	"time"

	"github.com/chnmk/todo-list-final-task/internal/services"
)

// Возвращает следующую дату в соответствии с данными запроса.
func NextDate(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodGet {
		ReturnError(w, "неожиданный метод запроса, ожидался GET", 400)
		return
	}

	now := r.FormValue("now")
	date := r.FormValue("date")
	repeat := r.FormValue("repeat")

	nowTime, err := time.Parse("20060102", now)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	response, err := services.NextDate(nowTime, date, repeat)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.Write([]byte(response))
}

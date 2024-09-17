package api

import (
	"net/http"
	"time"

	"github.com/chnmk/todo-list-final-task/services"
)

func NextDate(w http.ResponseWriter, r *http.Request) {
	var responseInvalid ResponseInvalid

	if r.Method != http.MethodGet {
		responseInvalid.Error = "неожиданный метод запроса, ожидался GET"
		returnInvalid(w, responseInvalid, 500)
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
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	w.Write([]byte(response))
}

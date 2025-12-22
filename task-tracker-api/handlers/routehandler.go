package handlers

import "net/http"

func TasksHandler(w http.ResponseWriter, r *http.Request) {
	switch r.Method {
	case http.MethodPost:
		CreateTaskHandler(w, r)
	case http.MethodPut:
		UpdateTaskHandler(w, r)
	case http.MethodDelete:
		DeleteTaskHandler(w, r)
	case http.MethodGet:
		GetUserTasksHandler(w, r)
	default:
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
	}
}

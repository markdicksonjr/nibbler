package nibbler

import "net/http"

func Write404Json(w http.ResponseWriter) {
	w.WriteHeader(http.StatusNotFound)
	w.Header().Set("Content-Type", "application/json")
	w.Write([]byte(`{"result": "not found"}`))
}

func Write500Json(w http.ResponseWriter, message string) {
	w.WriteHeader(http.StatusInternalServerError)
	w.Header().Set("Content-Type", "application/json")
	w.Write([]byte(`{"result": "` + message + `"}`))
}
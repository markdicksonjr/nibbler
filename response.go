package nibbler

import "net/http"

// WriteJson is some syntactic sugar to allow for a quick way to write JSON responses with a status code
func WriteJson(w http.ResponseWriter, content string, code int) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(code)
	w.Write([]byte(content))
}

// Write200Json is some syntactic sugar to allow for a quick way to write JSON responses with an OK code
func Write200Json(w http.ResponseWriter, content string) {
	WriteJson(w, content, http.StatusOK)
}

// Write401Json
func Write401Json(w http.ResponseWriter) {
	WriteJson(w, `{"result": "not authorized"}`, http.StatusUnauthorized)
}

// Write404Json is some syntactic sugar to allow for a quick way to write JSON responses with a StatusNotFound code
func Write404Json(w http.ResponseWriter) {
	WriteJson(w, `{"result": "not found"}`, http.StatusNotFound)
}

// Write500Json is some syntactic sugar to allow for a quick way to write JSON responses with an InternalServererror code
func Write500Json(w http.ResponseWriter, message string) {
	WriteJson(w, `{"result": "` + message + `"}`, http.StatusInternalServerError)
}

package nibbler

import (
	"encoding/json"
	"net/http"
)

// WriteJson is some syntactic sugar to allow for a quick way to write JSON responses with a status code
func WriteJson(w http.ResponseWriter, content string, code int) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(code)
	w.Write([]byte(content))
}

// WriteJson is some syntactic sugar to allow for a quick way to write JSON responses from structs with a status code
func WriteStructToJson(w http.ResponseWriter, content interface{}, code int) {
	r, err := json.Marshal(content)
	if err != nil {
		Write500Json(w, err.Error())
		return
	}
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(code)
	w.Write(r)
}

// Write200Json is some syntactic sugar to allow for a quick way to write JSON responses with an OK code
func Write200Json(w http.ResponseWriter, content string) {
	WriteJson(w, content, http.StatusOK)
}

// Write201Json is some syntactic sugar to allow for a quick way to write JSON responses with an OK code
func Write201Json(w http.ResponseWriter, content string) {
	WriteJson(w, content, http.StatusCreated)
}

// Write204 is some syntactic sugar to allow for a quick way to write JSON responses with a no-content code
func Write204(w http.ResponseWriter, content string) {
	w.WriteHeader(http.StatusNoContent)
	w.Write(nil)
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

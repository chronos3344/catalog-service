package httph

import (
	"net/http"
)

// Error структура для ошибок API
type Error struct {
	Message string `json:"error"` // Добавьте json тэг со значением "error"
}

// ErrorApply отправка ошибки с указанием статуса и сообщения
func ErrorApply(w http.ResponseWriter, code int, message string) {
	// Записываем заголовок
	//  наш EncodeJSON
	w.Header().Set("Content-Type", "application/json; charset=utf-8")
	w.WriteHeader(code)

	_ = EncodeJSON(w, &Error{message})
}

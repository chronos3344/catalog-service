package httph

import (
	"net/http"
)

// SendRaw отправляет сырые байты с указанным MIME-типом и статус-кодом
func SendRaw(w http.ResponseWriter, statusCode int, mimeType string, data []byte) {
	// 1. Если mimeType не пустой, установите его в заголовок Content-Type
	if mimeType != "" {
		w.Header().Set("Content-Type", MIMEApplicationJSONCharsetUTF8)
	}

	// 2. Запишите статус код через w.WriteHeader(statusCode)
	w.WriteHeader(statusCode)

	// 3. Если data не пустой, запишите его в тело ответа через w.Write(data)
	//    Ошибку можно проигнорировать
	if len(data) > 0 {
		_, _ = w.Write(data)
	}
}

// SendEmpty отправляет ответ без тела (только статус-код)
func SendEmpty(w http.ResponseWriter, statusCode int) {
	// Вызываем SendRaw, передаем только статус-код, остальные значения - дефолтные.
	SendRaw(w, statusCode, "", nil)
}

// SendEncodedWithMIME кодирует объект и отправляет с указанным MIME-типом
func SendEncodedWithMIME(w http.ResponseWriter, r *http.Request, statusCode int, mimeType string, obj any) {
	// Вызываем SendRaw и передаем только статус код и MIME.
	SendRaw(w, statusCode, MIMEApplicationJSONCharsetUTF8, nil)

	// Вызываем наш EncodeJSON и в случае ошибки вызываем ErrorApply.
	err := EncodeJSON(w, obj)
	if err != nil {
		ErrorApply(w, http.StatusBadRequest, "Неверный формат")
		return
	}
}

// SendEncoded отправляет объект в формате JSON с указанным статус-кодом
func SendEncoded(w http.ResponseWriter, r *http.Request, statusCode int, obj any) {
	// Вызываем SendEncodedWithMIME и передаем нашу констату MIME с Application/JSON в UTF-8 (Такой формат используется в 90% случаев.)
	SendEncodedWithMIME(w, r, statusCode, MIMEApplicationJSONCharsetUTF8, obj)
	// А если нам пригодится что-то особенное, то просто вызовем SendEncodedWithMIME прямо в обработчике.
}

package binding

import (
	"net/http"
)

func ScanAndValidateJSON(r *http.Request, to any) error {
	return scanAndValidate(r, to, bJSON)
}

func ScanAndValidateQuery(r *http.Request, to any) error {
	return scanAndValidate(r, to, bQuery)
}

func scanAndValidate(r *http.Request, to any, b Binding) error {
	// Вызовите b.Bind(r, to) и верните результат.
	// Bind() уже выполняет и парсинг, и валидацию.

	return b.Bind(r, to)
}

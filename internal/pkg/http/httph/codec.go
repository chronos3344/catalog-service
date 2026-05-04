package httph

import (
	"encoding/json"
	"net/http"
)

func EncodeJSON(w http.ResponseWriter, data interface{}) error {
	// Используйте json.NewEncoder(w).Encode(data)
	// Верните ошибку если она возникла
	return json.NewEncoder(w).Encode(data)
}

func DecodeJSON(r *http.Request, v interface{}) error {
	// Используйте json.NewDecoder(r.Body).Decode(v)
	// Верните ошибку если она возникла
	return json.NewDecoder(r.Body).Decode(v)
}

package httph

import (
	"encoding/json"
	"log"
	"net/http"
)

func EncodeJSON(w http.ResponseWriter, data interface{}) error {
	// Используйте json.NewEncoder(w).Encode(data)
	// Верните ошибку если она возникла
	err := json.NewEncoder(w).Encode(data)
	if err != nil {
		log.Printf("failed to encode response: %v", err)
	}
	return err
}

func DecodeJSON(r *http.Request, v interface{}) error {
	// Используйте json.NewDecoder(r.Body).Decode(v)
	// Верните ошибку если она возникла
	err := json.NewDecoder(r.Body).Decode(v)
	if err != nil {
		log.Printf("failed to decode response: %v", err)
	}
	return err
}

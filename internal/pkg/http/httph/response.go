package httph

import (
	"encoding/json"
	"errors"
	"net/http"
)

type httpCoder interface {
	error
	HTTPStatus() int
}

type HTTPError struct {
	status int
	msg    string
}

func (e *HTTPError) Error() string {
	return e.msg
}

func (e *HTTPError) HTTPStatus() int {
	return e.status
}

func NewHTTPError(status int, msg string) *HTTPError {
	return &HTTPError{
		status: status,
		msg:    msg,
	}
}

func HandleError(w http.ResponseWriter, r *http.Request, err error) {
	ErrorApply(r, err)

	var hc httpCoder
	if errors.As(err, &hc) {
		status := hc.HTTPStatus()
		ErrorApplyStatusCode(r, status)
		sendError(w, status, hc)
		return
	}
	ErrorApplyStatusCode(r, http.StatusInternalServerError)
	sendError(w, http.StatusInternalServerError, err)
}

func sendError(w http.ResponseWriter, status int, err error) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(status)

	response := map[string]interface{}{
		"error": map[string]interface{}{
			"status":  status,
			"message": err.Error(),
		},
	}
	_ = json.NewEncoder(w).Encode(response)
}

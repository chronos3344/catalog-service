package binding

import (
	"errors"
	"net/http"

	"github.com/chronos3344/catalog-service/internal/pkg/http/httph"
)

type jsonBinding struct{}

func (jsonBinding) Name() string {
	return "JSON"
}

func (jsonBinding) Bind(req *http.Request, obj any) error {
	// 1. Проверьте req и req.Body на nil, в случае nil верните ошибку "invalid request"
	if req == nil || req.Body == nil {
		return errors.New("invalid request")
	}

	// 2. Вызовите httph.DecodeJSON(req, obj), обработайте ошибку
	err := httph.DecodeJSON(req, obj)
	if err != nil {
		return err
	}

	return validate(obj)
}

package req

import (
	"net/http"

	"shorty/pkg/res"
)

// HandleBody декодирует тело запроса в структуру T и проверяет её валидность.
func HandleBody[T any](w *http.ResponseWriter, r *http.Request) (*T, error) {
	body, err := Decode[T](r.Body)
	if err != nil {
		res.JSON(*w, err.Error(), http.StatusBadRequest)
		return nil, err
	}
	err = IsValidate(body)
	if err != nil {
		res.JSON(*w, err.Error(), http.StatusBadRequest)
		return nil, err
	}
	return &body, nil
}

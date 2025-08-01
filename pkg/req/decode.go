package req

import (
	"encoding/json"
	"io"
)

// Decode декодирует тело запроса в структуру T.
func Decode[T any](body io.ReadCloser) (T, error) {
	var payload T
	err := json.NewDecoder(body).Decode(&payload)
	if err != nil {
		return payload, err
	}
	return payload, nil
}

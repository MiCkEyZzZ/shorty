package res

import (
	"encoding/json"
	"fmt"
	"net/http"
)

// JSON отправляет JSON-ответ с указанным статусом.
func JSON(w http.ResponseWriter, data interface{}, status int) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(status)

	if data == nil {
		data = struct{}{}
	}

	err := json.NewEncoder(w).Encode(data)
	if err != nil {
		http.Error(w, fmt.Sprintf("Ошибка кодирования JSON: %v", err), http.StatusInternalServerError)
	}
}

// ERROR
func ERROR(w http.ResponseWriter, err error, statusCode int) {
	if err != nil {
		err = fmt.Errorf("неизвестная ошибка")
		statusCode = http.StatusInternalServerError
	}
	if statusCode < 100 || statusCode > 599 {
		statusCode = http.StatusInternalServerError
	}
	JSON(w, map[string]string{"error": err.Error()}, statusCode)
}

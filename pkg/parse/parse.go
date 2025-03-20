package parse

import (
	"net/http"
	"strconv"
)

// parseID парсит идентификатор из строки в uint.
func ParseID(r *http.Request) (uint, error) {
	rid := r.PathValue("id")
	id, err := strconv.ParseUint(rid, 10, 32)
	if err != nil {
		return 0, err
	}
	return uint(id), nil
}

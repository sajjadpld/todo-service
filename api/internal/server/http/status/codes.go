package status

import "net/http"

var MappedStatuses = map[HttpMappedStatus]int{
	Success:      http.StatusOK,
	Created:      http.StatusCreated,
	Updated:      http.StatusNoContent,
	Validate:     http.StatusUnprocessableEntity,
	Failed:       http.StatusBadRequest,
	NotFound:     http.StatusNotFound,
	Unauthorized: http.StatusUnauthorized,
	Conflict:     http.StatusConflict,
	ItemExist:    http.StatusConflict,
}

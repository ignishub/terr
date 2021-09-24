package apierror

import (
	"encoding/json"
	"net/http"
)

var includeDebugInfo bool = false

// IncludeDebugInfo determines whether to add debug info into http response.
func IncludeDebugInfo(value bool) {
	includeDebugInfo = value
}

// ServeError serves error over http. When err is not Error it converts to InternalServerError.
func ServeError(w http.ResponseWriter, err error) {
	var e *Error
	var ok bool
	e, ok = err.(*Error)
	if !ok {
		e = InternalServerError("UNKNOWN_ERROR", err.Error())
	}

	if !includeDebugInfo {
		e.Debug = nil
	}

	w.WriteHeader(e.HTTPStatusCode)
	w.Header().Add("Content-Type", "application/json")
	json.NewEncoder(w).Encode(e)
}

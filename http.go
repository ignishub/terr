package terr

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
	e := From(err)
	if !includeDebugInfo {
		e.Debug = nil
	}

	w.WriteHeader(e.HTTPStatusCode)
	w.Header().Add("Content-Type", "application/json")
	json.NewEncoder(w).Encode(e)
}

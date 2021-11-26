package httperror

import (
	"encoding/json"
	"net/http"

	"github.com/ignishub/terr"
)

var (
	WithoutDebugInfo = true
	WithoutDetails   = false
)

// ServeError serves error over http. When err is not Error it converts to InternalServerError.
func ServeError(w http.ResponseWriter, err error) {
	e := terr.From(err)

	if WithoutDebugInfo {
		e.Debug = nil
	}

	if WithoutDetails {
		e.Details = nil
	}

	w.WriteHeader(e.HTTPStatusCode)
	w.Header().Add("Content-Type", "application/json")
	err = json.NewEncoder(w).Encode(e)
	if err != nil {
		panic(err)
	}
}

package apierror

import "encoding/json"

// Error represents API error.
type Error struct {
	// HTTPStatusCode http status code that will be returned with error.
	HTTPStatusCode int `json:"-"`
	// Code is constant string message that represents error.
	Code string `json:"code"`
	// Message is message which describes error for user.
	Message string `json:"message"`
	// Details contains additional info.
	Details map[string]interface{} `json:"details,omitempty"`
	// Debug contains debug information that can be skipped in production.
	Debug *map[string]interface{} `json:"debug,omitempty"`
}

// Equal returns true if errors have same code.
func Equal(err1, err2 Error) bool {
	return err1.Code == err2.Code
}

func (e *Error) MarshalJSON() ([]byte, error) {
	return json.Marshal(e)
}

func (e Error) Error() string {
	return e.Message
}

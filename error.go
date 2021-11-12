package terr

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
	Debug map[string]interface{} `json:"debug,omitempty"`
}

// From returns Error structure from all kind of error types.
func From(err error) *Error {
	e, ok := err.(*Error)
	if !ok {
		e = InternalServerError("UNKNOWN_ERROR", err.Error())
	}
	return e
}

// Equal returns true if errors have same code.
func Equal(e1, e2 error) bool {
	if e1 == nil && e2 == nil {
		return true
	}
	if (e1 == nil && e2 != nil) || (e1 != nil && e2 == nil) {
		return false
	}

	err1, ok := e1.(*Error)
	if !ok {
		return e1.Error() == e2.Error()
	}

	err2, ok := e2.(*Error)
	if !ok {
		return e1.Error() == e2.Error()
	}

	return err1.Code == err2.Code
}

// Error реализация интерфейса ошибки.
func (e *Error) Error() string {
	return e.Message
}

// WithDebug добавляет отладочную информацию в ошибку.
func (e *Error) WithDebug(field string, value interface{}) *Error {
	if e.Debug == nil {
		e.Debug = make(map[string]interface{})
	}
	e.Debug[field] = value
	return e
}

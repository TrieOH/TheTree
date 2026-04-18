package idx

type SdkError struct {
	Message string `json:"message"`
	Cause   error  `json:"cause"`
}

func (e SdkError) Error() string {
	if e.Cause != nil {
		if e.Message != "" {
			return e.Message + ": " + e.Cause.Error()
		}
		return e.Cause.Error()
	}
	return e.Message
}

type ApiError struct {
	ErrorID string   `json:"error_id"`
	Message string   `json:"message"`
	Trace   []string `json:"trace"`
	Code    int      `json:"code"`
}

func (e ApiError) Error() string {
	if e.ErrorID != "" {
		return "[" + e.ErrorID + "] " + e.Message
	}
	return e.Message
}

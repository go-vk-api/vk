package vk

// MethodError represents a VK API method call error.
type MethodError struct {
	Code          int64  `json:"error_code"`
	Message       string `json:"error_msg"`
	RequestParams []struct {
		Key   string `json:"key"`
		Value string `json:"value"`
	} `json:"request_params"`
}

func (err *MethodError) Error() string {
	return "vk: " + err.Message
}

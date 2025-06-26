package modeldto

type ResponseStatusDto struct {
	Success        bool   `json:"success"`
	ResponseTimeMs int64  `json:"response_time_ms"`
	Latency        int64  `json:"latency"`
	ErrorCode      string `json:"error_code"`
	ErrorMessage   string `json:"error_message"`
}

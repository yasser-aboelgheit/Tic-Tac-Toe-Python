package response

// ErrorResponse default/base error responses on all servers.
type ErrorResponse struct {
	Error   error  `json:"error"    xml:"error"    yaml:"error"`
	TraceID string `json:"trace_id" xml:"trace_id" yaml:"trace_id"`
}

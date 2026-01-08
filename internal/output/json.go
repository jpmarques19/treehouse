package output

import (
	"encoding/json"
	"fmt"
	"io"
	"os"
)

// writer is the output destination (stdout by default, overridable for testing)
var writer io.Writer = os.Stdout

// Response is the standard JSON response wrapper for all CLI output
type Response struct {
	Success bool        `json:"success"`
	Data    interface{} `json:"data,omitempty"`
	Error   *ErrorInfo  `json:"error,omitempty"`
}

// ErrorInfo contains error details for failed operations
type ErrorInfo struct {
	Code    string `json:"code"`
	Message string `json:"message"`
}

// SetWriter sets the output writer (used for testing)
func SetWriter(w io.Writer) {
	writer = w
}

// PrintSuccess outputs a success response with optional data
func PrintSuccess(data interface{}) {
	resp := Response{Success: true, Data: data}
	printJSON(resp)
}

// PrintError outputs an error response with code and message
func PrintError(code, message string) {
	resp := Response{
		Success: false,
		Error:   &ErrorInfo{Code: code, Message: message},
	}
	printJSON(resp)
}

// printJSON marshals and prints the response as JSON
func printJSON(v interface{}) {
	bytes, err := json.Marshal(v)
	if err != nil {
		// Fallback to a valid error JSON if marshal fails
		bytes = []byte(`{"success":false,"error":{"code":"MARSHAL_ERROR","message":"failed to marshal response"}}`)
	}
	fmt.Fprintln(writer, string(bytes))
}

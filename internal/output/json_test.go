package output

import (
	"bytes"
	"encoding/json"
	"testing"
)

func TestResponse_Success(t *testing.T) {
	tests := []struct {
		name     string
		data     interface{}
		wantJSON string
	}{
		{
			name:     "nil data",
			data:     nil,
			wantJSON: `{"success":true}`,
		},
		{
			name:     "string data",
			data:     "hello",
			wantJSON: `{"success":true,"data":"hello"}`,
		},
		{
			name:     "map data",
			data:     map[string]string{"version": "0.4.0"},
			wantJSON: `{"success":true,"data":{"version":"0.4.0"}}`,
		},
		{
			name:     "struct data",
			data:     struct{ Name string }{"test"},
			wantJSON: `{"success":true,"data":{"Name":"test"}}`,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			var buf bytes.Buffer
			oldWriter := writer
			writer = &buf
			t.Cleanup(func() { writer = oldWriter })

			PrintSuccess(tt.data)

			got := buf.String()
			// Remove trailing newline for comparison
			got = got[:len(got)-1]

			if got != tt.wantJSON {
				t.Errorf("PrintSuccess() = %v, want %v", got, tt.wantJSON)
			}
		})
	}
}

func TestResponse_Error(t *testing.T) {
	tests := []struct {
		name     string
		code     string
		message  string
		wantJSON string
	}{
		{
			name:     "no command error",
			code:     "NO_COMMAND",
			message:  "No command specified",
			wantJSON: `{"success":false,"error":{"code":"NO_COMMAND","message":"No command specified"}}`,
		},
		{
			name:     "not found error",
			code:     "NOT_FOUND",
			message:  "Resource not found",
			wantJSON: `{"success":false,"error":{"code":"NOT_FOUND","message":"Resource not found"}}`,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			var buf bytes.Buffer
			oldWriter := writer
			writer = &buf
			t.Cleanup(func() { writer = oldWriter })

			PrintError(tt.code, tt.message)

			got := buf.String()
			// Remove trailing newline for comparison
			got = got[:len(got)-1]

			if got != tt.wantJSON {
				t.Errorf("PrintError() = %v, want %v", got, tt.wantJSON)
			}
		})
	}
}

func TestResponse_JSONValidity(t *testing.T) {
	tests := []struct {
		name string
		fn   func()
	}{
		{
			name: "success response",
			fn: func() {
				PrintSuccess(map[string]string{"key": "value"})
			},
		},
		{
			name: "error response",
			fn: func() {
				PrintError("TEST_ERROR", "Test message")
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			var buf bytes.Buffer
			oldWriter := writer
			writer = &buf
			t.Cleanup(func() { writer = oldWriter })

			tt.fn()

			var resp Response
			if err := json.Unmarshal(buf.Bytes(), &resp); err != nil {
				t.Errorf("Output is not valid JSON: %v", err)
			}
		})
	}
}

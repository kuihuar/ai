package generic

import (
	"bytes"
	"io"
	"os"
	"testing"
)

func TestPrint(t *testing.T) {
	tests := []struct {
		name string // description of this test case
		// Named input parameters for target function.
		input any
		// Expected output from target function.
		want string
		// v T
	}{
		{"int", 123, "123\n"},
		{"string", "hello", "hello\n"},
		{"bool", true, "true\n"},
		{"struct", struct{ A int }{100}, "{100}\n"},
		{"slice", []int{1, 2}, "[1 2]\n"},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			out := captureStdout(func() {
				Print(tt.input)
			})
			if out != tt.want {
				t.Errorf("Print() output = %q, want %q", out, tt.want)
			}
		})
	}
}

func captureStdout(f func()) string {
	old := os.Stdout
	r, w, _ := os.Pipe()
	os.Stdout = w

	f()
	_ = w.Close()
	os.Stdout = old

	var buf bytes.Buffer
	_, _ = io.Copy(&buf, r)
	return buf.String()
}

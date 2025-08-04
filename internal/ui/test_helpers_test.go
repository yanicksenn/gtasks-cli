package ui

import (
	"bytes"
	"io"
	"os"
	"testing"
)

// CaptureOutput is a helper to capture stdout from a function.
func CaptureOutput(t *testing.T, f func()) string {
	t.Helper()
	oldStdout := os.Stdout
	r, w, _ := os.Pipe()
	os.Stdout = w

	f()

	w.Close()
	os.Stdout = oldStdout
	var buf bytes.Buffer
	io.Copy(&buf, r)
	return buf.String()
}

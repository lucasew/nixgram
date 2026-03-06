package errreporter

import (
	"bytes"
	"errors"
	"log"
	"strings"
	"testing"
)

func TestReportError(t *testing.T) {
	var buf bytes.Buffer
	log.SetOutput(&buf)
	defer log.SetOutput(nil)

	// Test nil error
	ReportError(nil, nil)
	if buf.Len() > 0 {
		t.Errorf("Expected nothing logged for nil error, got: %s", buf.String())
	}

	// Test non-nil error
	err := errors.New("test error")
	context := map[string]interface{}{"key": "value"}
	ReportError(err, context)

	output := buf.String()
	if !strings.Contains(output, "ERROR: test error") {
		t.Errorf("Expected output to contain 'ERROR: test error', got: %s", output)
	}
	if !strings.Contains(output, "Context: map[key:value]") {
		t.Errorf("Expected output to contain 'Context: map[key:value]', got: %s", output)
	}
}

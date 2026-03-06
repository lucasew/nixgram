package errreporter

import (
	"log"
)

// ReportError centralizes error reporting for the application.
// All code paths that handle unexpected errors MUST funnel through this function.
func ReportError(err error, context map[string]interface{}) {
	if err == nil {
		return
	}
	log.Printf("ERROR: %v | Context: %v", err, context)
}

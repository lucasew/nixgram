package errreporter

import (
	"log"
)

// ReportError acts as the centralized error-reporting function for the application.
// All unexpected errors must be funneled through this function as per project guidelines.
// This prevents silent failures and provides a single sink for Sentry/logging integration.
func ReportError(err error, metadata map[string]interface{}) {
	if err == nil {
		return
	}

	// If metadata is provided, log it alongside the error
	if len(metadata) > 0 {
		log.Printf("ERROR: %v | Metadata: %v\n", err, metadata)
	} else {
		log.Printf("ERROR: %v\n", err)
	}
}

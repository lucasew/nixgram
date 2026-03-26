package errreporter

import (
	"log"
)

// ReportError centralizes error reporting for the application.
// It logs the error with its contextual message. In the future,
// this can be wired to Sentry or other external observability tools.
func ReportError(contextMsg string, err error) {
	if err != nil {
		log.Printf("[ERROR] %s: %v\n", contextMsg, err)
	}
}

package errreporter

import (
	"log"
	"runtime/debug"
)

// ReportError is the centralized error reporting function.
// All code paths that handle unexpected errors MUST funnel through this function.
func ReportError(err error, contextMsg string) {
	if err == nil {
		return
	}
	log.Printf("ERROR: %s: %v\nStack trace:\n%s", contextMsg, err, debug.Stack())
}

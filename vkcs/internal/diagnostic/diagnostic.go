package diagnostic

import (
	"fmt"

	"github.com/hashicorp/terraform-plugin-framework/diag"
)

type DiagnosticWithRequestID struct {
	detail    string
	summary   string
	RequestID string
}

func NewErrorWithRequestID(summary, detail, requestID string) DiagnosticWithRequestID {
	return DiagnosticWithRequestID{
		summary:   summary,
		detail:    detail,
		RequestID: requestID,
	}
}

func (d DiagnosticWithRequestID) Severity() diag.Severity {
	return diag.SeverityError
}

func (d DiagnosticWithRequestID) Summary() string {
	return d.summary
}

func (d DiagnosticWithRequestID) Detail() string {
	return fmt.Sprintf("%s\nRequest ID: %s", d.detail, d.RequestID)
}

func (d DiagnosticWithRequestID) Equal(o diag.Diagnostic) bool {
	if d.Detail() != o.Detail() {
		return false
	}

	if d.Summary() != o.Summary() {
		return false
	}
	return true
}

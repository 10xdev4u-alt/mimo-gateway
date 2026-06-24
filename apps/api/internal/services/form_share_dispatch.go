package services

import (
	"fmt"

	"gorm.io/gorm"
)

// SharedResourceSubmission is the result of a public form submission —
// the created record's ID and human label, both safe to return to
// anonymous visitors.
type SharedResourceSubmission struct {
	ID    string
	Label string
}

// SubmitSharedForm dispatches a public form submission to the right
// resource service based on the FormShare's ResourceName. fields is
// a free-form map (validated by the resource service's own binding
// rules), since public submissions don't carry the operator's typed
// struct context.
//
// Adding a new resource? grit generate resource appends a case to
// the switch below at the auto-dispatch marker. Each case re-marshals
// fields into the typed model via json.Marshal(fields) — that's why
// the parameter is named "fields" rather than "body".
func SubmitSharedForm(db *gorm.DB, resourceName string, fields map[string]interface{}) (*SharedResourceSubmission, error) {
	switch resourceName {
	// grit:form-share:dispatch
	default:
		return nil, fmt.Errorf("public submission disabled for %q (no dispatch case registered)", resourceName)
	}
}

// Package common provides shared types used across PRD, MRD, and TRD documents.
package common

import (
	"time"

	core "github.com/grokify/prism-core"
)

// Person is an alias for core.Person.
// Represents an individual contributor.
type Person = core.Person

// Approver represents a person with approval authority.
// Note: Uses *time.Time for ApprovedAt to support omitempty in JSON.
type Approver struct {
	Person
	ApprovedAt *time.Time `json:"approvedAt,omitempty"`
	Approved   bool       `json:"approved"`
	Comments   string     `json:"comments,omitempty"`
}

// FormatPersonMarkdown formats a Person for markdown display.
// Delegates to core.FormatPersonMarkdown.
func FormatPersonMarkdown(p Person) string {
	return core.FormatPersonMarkdown(p)
}

// FormatPeopleMarkdown formats a slice of Person for markdown display.
// Delegates to core.FormatPeopleMarkdown.
func FormatPeopleMarkdown(people []Person) string {
	return core.FormatPeopleMarkdown(people)
}

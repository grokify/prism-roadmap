// Package goals provides a framework-agnostic wrapper for goal-setting systems.
//
// Organizations use OKR (Objectives and Key Results), V2MOM (Vision, Values,
// Methods, Obstacles, Measures), or VSEM (Vision, Strategy, Execution,
// Metrics) for strategic planning. This package provides a unified interface
// that supports each framework through a discriminated union pattern.
//
// Use Goals in PRDs and roadmaps to support organizations using either framework,
// with common rendering and reporting regardless of the underlying goal system.
package goals

import (
	"github.com/grokify/prism-roadmap/goals/okr"
	"github.com/grokify/prism-roadmap/goals/v2mom"
	"github.com/grokify/prism-roadmap/goals/vsem"
)

// Framework identifies the goal-setting framework in use.
type Framework string

const (
	// FrameworkOKR indicates OKR (Objectives and Key Results) framework.
	FrameworkOKR Framework = "okr"
	// FrameworkV2MOM indicates V2MOM (Vision, Values, Methods, Obstacles, Measures) framework.
	FrameworkV2MOM Framework = "v2mom"
	// FrameworkVSEM indicates VSEM (Vision, Strategy, Execution, Metrics) framework.
	FrameworkVSEM Framework = "vsem"
)

// Goals is a framework-agnostic container for organizational goals.
// It supports OKR, V2MOM, and VSEM through a discriminated union pattern.
// Exactly one of OKR, V2MOM, or VSEM should be set based on the Framework
// field.
type Goals struct {
	// Framework identifies which goal system is in use ("okr", "v2mom", or
	// "vsem").
	Framework Framework `json:"framework"`

	// OKR contains OKR data when Framework is "okr".
	OKR *okr.OKRSet `json:"okr,omitempty"`

	// V2MOM contains V2MOM data when Framework is "v2mom".
	V2MOM *v2mom.V2MOM `json:"v2mom,omitempty"`

	// VSEM contains VSEM data when Framework is "vsem".
	VSEM *vsem.VSEM `json:"vsem,omitempty"`
}

// NewOKR creates a Goals wrapper containing OKR data.
func NewOKR(okrs *okr.OKRSet) *Goals {
	return &Goals{
		Framework: FrameworkOKR,
		OKR:       okrs,
	}
}

// NewV2MOM creates a Goals wrapper containing V2MOM data.
func NewV2MOM(v *v2mom.V2MOM) *Goals {
	return &Goals{
		Framework: FrameworkV2MOM,
		V2MOM:     v,
	}
}

// NewVSEM creates a Goals wrapper containing VSEM data.
func NewVSEM(v *vsem.VSEM) *Goals {
	return &Goals{
		Framework: FrameworkVSEM,
		VSEM:      v,
	}
}

// IsOKR returns true if this Goals uses the OKR framework.
func (g *Goals) IsOKR() bool {
	return g != nil && g.Framework == FrameworkOKR
}

// IsV2MOM returns true if this Goals uses the V2MOM framework.
func (g *Goals) IsV2MOM() bool {
	return g != nil && g.Framework == FrameworkV2MOM
}

// IsVSEM returns true if this Goals uses the VSEM framework.
func (g *Goals) IsVSEM() bool {
	return g != nil && g.Framework == FrameworkVSEM
}

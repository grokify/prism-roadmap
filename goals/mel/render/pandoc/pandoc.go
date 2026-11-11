// Package pandoc renders a MEL plan as Pandoc Markdown: the Markdown body with
// YAML front matter set up for a PDF through LuaLaTeX, with sans-serif fonts,
// a table of contents, and numbered sections. The output is deterministic; the
// date is the plan's last update, when it has one.
//
// To convert the output to PDF:
//
//	pandoc plan.md -o plan.pdf --pdf-engine=lualatex
package pandoc

import (
	"fmt"
	"strings"

	"github.com/grokify/prism-roadmap/goals/mel"
	"github.com/grokify/prism-roadmap/goals/mel/render"
	"github.com/grokify/prism-roadmap/goals/mel/render/markdown"
)

// Renderer implements render.Renderer for Pandoc Markdown.
type Renderer struct{}

// New returns a Pandoc Markdown renderer.
func New() *Renderer {
	return &Renderer{}
}

// Format returns "pandoc".
func (r *Renderer) Format() string {
	return "pandoc"
}

// FileExtension returns ".md".
func (r *Renderer) FileExtension() string {
	return ".md"
}

// Render converts the plan to Pandoc Markdown with YAML front matter.
func (r *Renderer) Render(p *mel.Plan, opts *render.Options) ([]byte, error) {
	if p == nil {
		return nil, fmt.Errorf("pandoc: nil MEL plan")
	}
	var b strings.Builder
	b.WriteString("---\n")
	fmt.Fprintf(&b, "title: %s\n", quote(markdown.Title(p)))
	if m := p.Metadata; m != nil {
		var sub []string
		for _, s := range []string{m.Team, m.Period} {
			if s != "" {
				sub = append(sub, s)
			}
		}
		if len(sub) > 0 {
			fmt.Fprintf(&b, "subtitle: %s\n", quote(strings.Join(sub, " · ")))
		}
		if m.Author != "" {
			fmt.Fprintf(&b, "author: %s\n", quote(m.Author))
		}
		if !m.UpdatedAt.IsZero() {
			fmt.Fprintf(&b, "date: %s\n", quote(m.UpdatedAt.Format("2006-01-02")))
		}
	}
	b.WriteString(`documentclass: article
papersize: a4
geometry:
  - margin=2cm
fontsize: 11pt
mainfont: "DejaVu Sans"
sansfont: "DejaVu Sans"
monofont: "DejaVu Sans Mono"
header-includes:
  - \usepackage{booktabs}
  - \usepackage{longtable}
  - \usepackage{array}
colorlinks: true
linkcolor: blue
urlcolor: blue
toc: true
toc-depth: 2
numbersections: true
---

`)
	b.Write(markdown.Body(p, opts))
	return []byte(strings.TrimRight(b.String(), "\n") + "\n"), nil
}

// quote writes s as a YAML single-quoted scalar, which needs no escapes but a
// doubled quote.
func quote(s string) string {
	return "'" + strings.ReplaceAll(s, "'", "''") + "'"
}

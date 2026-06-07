//go:build ignore

// This program generates rendered outputs for canvas examples.
// Run with: go run generate.go
package main

import (
	"encoding/json"
	"fmt"
	"os"
	"path/filepath"

	"github.com/grokify/prism-roadmap/canvas"
	"github.com/grokify/prism-roadmap/canvas/render"
	"github.com/grokify/prism-roadmap/canvas/render/d2"
	"github.com/grokify/prism-roadmap/canvas/render/lit"
	"github.com/grokify/prism-roadmap/canvas/render/mermaid"
	"github.com/grokify/prism-roadmap/canvas/render/svg"
)

func main() {
	// Register renderers
	registry := render.NewRegistry()
	registry.Register(d2.NewD2Renderer())
	registry.Register(mermaid.NewMermaidRenderer())
	registry.Register(lit.NewLitRenderer())
	registry.Register(svg.NewSVGRenderer())

	// Process each example
	examples := []struct {
		name     string
		jsonDir  string
		jsonFile string // If different from name
		opts     *render.Options
	}{
		{"bmc", "bmc", "", render.BMCOptions()},
		{"ost", "ost", "", render.OSTOptions()},
		{"opportunity_flow", "opportunity", "opportunity", render.OpportunityOptions()},
		{"opportunity_grid", "opportunity", "opportunity", render.OpportunityGridOptions()},
		{"feature", "feature", "", render.FeatureOptions()},
		{"leanux", "leanux", "", render.LeanUXOptions()},
	}

	for _, ex := range examples {
		fmt.Printf("Processing %s...\n", ex.name)

		// Determine JSON file name
		jsonName := ex.name
		if ex.jsonFile != "" {
			jsonName = ex.jsonFile
		}

		// Read JSON
		jsonPath := filepath.Join(ex.jsonDir, jsonName+"_example.json")
		data, err := os.ReadFile(jsonPath)
		if err != nil {
			fmt.Printf("  Error reading %s: %v\n", jsonPath, err)
			continue
		}

		// Parse canvas
		var c canvas.Canvas
		if err := json.Unmarshal(data, &c); err != nil {
			fmt.Printf("  Error parsing %s: %v\n", jsonPath, err)
			continue
		}

		// Generate D2
		d2Output, err := registry.Render(&c, render.FormatD2, ex.opts)
		if err != nil {
			fmt.Printf("  Error rendering D2: %v\n", err)
		} else {
			outPath := filepath.Join(ex.jsonDir, ex.name+"_example.d2")
			if err := os.WriteFile(outPath, d2Output, 0644); err != nil {
				fmt.Printf("  Error writing %s: %v\n", outPath, err)
			} else {
				fmt.Printf("  Created %s\n", outPath)
			}
		}

		// Generate Mermaid
		mmdOutput, err := registry.Render(&c, render.FormatMermaid, ex.opts)
		if err != nil {
			fmt.Printf("  Error rendering Mermaid: %v\n", err)
		} else {
			outPath := filepath.Join(ex.jsonDir, ex.name+"_example.mmd")
			if err := os.WriteFile(outPath, mmdOutput, 0644); err != nil {
				fmt.Printf("  Error writing %s: %v\n", outPath, err)
			} else {
				fmt.Printf("  Created %s\n", outPath)
			}
		}

		// Generate Lit JSON
		litOutput, err := registry.Render(&c, render.FormatLit, ex.opts)
		if err != nil {
			fmt.Printf("  Error rendering Lit: %v\n", err)
		} else {
			outPath := filepath.Join(ex.jsonDir, ex.name+"_example.lit.json")
			if err := os.WriteFile(outPath, litOutput, 0644); err != nil {
				fmt.Printf("  Error writing %s: %v\n", outPath, err)
			} else {
				fmt.Printf("  Created %s\n", outPath)
			}
		}

		// Generate native SVG for grid layouts (Opportunity Canvas and Lean UX Canvas)
		supportsNativeSVG := ex.opts.GridLayout && (c.Type == canvas.CanvasTypeOpportunity || c.Type == canvas.CanvasTypeLeanUX)
		if supportsNativeSVG {
			svgOutput, err := registry.Render(&c, render.FormatSVG, ex.opts)
			if err != nil {
				fmt.Printf("  Error rendering native SVG: %v\n", err)
			} else {
				outPath := filepath.Join(ex.jsonDir, ex.name+"_example.svg")
				if err := os.WriteFile(outPath, svgOutput, 0644); err != nil {
					fmt.Printf("  Error writing %s: %v\n", outPath, err)
				} else {
					fmt.Printf("  Created %s (native)\n", outPath)
				}
			}
		}

		// Generate HTML viewer (with native SVG for grid layouts)
		var svgContent string
		if supportsNativeSVG {
			svgOutput, err := registry.Render(&c, render.FormatSVG, ex.opts)
			if err == nil {
				svgContent = string(svgOutput)
			}
		}
		htmlContent := generateHTML(ex.name, string(litOutput), string(mmdOutput), svgContent)
		outPath := filepath.Join(ex.jsonDir, ex.name+"_example.html")
		if err := os.WriteFile(outPath, []byte(htmlContent), 0644); err != nil {
			fmt.Printf("  Error writing %s: %v\n", outPath, err)
		} else {
			fmt.Printf("  Created %s\n", outPath)
		}
	}

	fmt.Println("\nDone! To generate SVGs, run: d2 <file>.d2 <file>.svg")
}

func generateHTML(name, litJSON, mermaidCode, nativeSVG string) string {
	// Determine diagram content - use native SVG if available, else Mermaid
	diagramContent := fmt.Sprintf(`<div class="mermaid">
%s
                </div>`, mermaidCode)
	if nativeSVG != "" {
		diagramContent = fmt.Sprintf(`<div class="svg-container" style="text-align:center;">
%s
                </div>`, nativeSVG)
	}

	// Determine if we need Mermaid script
	mermaidScript := `<script src="https://cdn.jsdelivr.net/npm/mermaid/dist/mermaid.min.js"></script>`
	mermaidInit := `mermaid.initialize({ startOnLoad: true, theme: 'default' });`
	if nativeSVG != "" {
		mermaidScript = ""
		mermaidInit = ""
	}

	return fmt.Sprintf(`<!DOCTYPE html>
<html lang="en">
<head>
    <meta charset="UTF-8">
    <meta name="viewport" content="width=device-width, initial-scale=1.0">
    <title>%s Canvas Example</title>
    %s
    <style>
        :root {
            --primary: #2563eb;
            --bg: #f8fafc;
            --card: #ffffff;
            --text: #1e293b;
            --border: #e2e8f0;
        }
        * { box-sizing: border-box; }
        body {
            font-family: -apple-system, BlinkMacSystemFont, 'Segoe UI', Roboto, sans-serif;
            background: var(--bg);
            color: var(--text);
            margin: 0;
            padding: 20px;
        }
        .container { max-width: 1400px; margin: 0 auto; }
        h1 { color: var(--primary); margin-bottom: 10px; }
        .subtitle { color: #64748b; margin-bottom: 30px; }
        .tabs {
            display: flex;
            gap: 10px;
            margin-bottom: 20px;
            border-bottom: 2px solid var(--border);
            padding-bottom: 10px;
        }
        .tab {
            padding: 10px 20px;
            background: var(--card);
            border: 1px solid var(--border);
            border-radius: 8px 8px 0 0;
            cursor: pointer;
            font-weight: 500;
        }
        .tab.active { background: var(--primary); color: white; border-color: var(--primary); }
        .panel { display: none; }
        .panel.active { display: block; }
        .card {
            background: var(--card);
            border: 1px solid var(--border);
            border-radius: 12px;
            padding: 24px;
            margin-bottom: 20px;
            box-shadow: 0 1px 3px rgba(0,0,0,0.1);
        }
        .mermaid { text-align: center; }
        pre {
            background: #1e293b;
            color: #e2e8f0;
            padding: 20px;
            border-radius: 8px;
            overflow-x: auto;
            font-size: 13px;
            line-height: 1.5;
        }
        .grid { display: grid; grid-template-columns: repeat(auto-fit, minmax(300px, 1fr)); gap: 20px; }
        .data-card h3 { margin-top: 0; color: var(--primary); }
        .data-card ul { margin: 0; padding-left: 20px; }
        .data-card li { margin-bottom: 8px; }
        .badge {
            display: inline-block;
            padding: 2px 8px;
            border-radius: 4px;
            font-size: 12px;
            font-weight: 500;
        }
        .badge-validated { background: #dcfce7; color: #166534; }
        .badge-running { background: #fef3c7; color: #92400e; }
        .badge-must { background: #fee2e2; color: #991b1b; }
        .badge-should { background: #dbeafe; color: #1e40af; }
    </style>
</head>
<body>
    <div class="container">
        <h1>%s Canvas</h1>
        <p class="subtitle">Interactive canvas visualization with multiple rendering formats</p>

        <div class="tabs">
            <div class="tab active" onclick="showTab('diagram')">Diagram</div>
            <div class="tab" onclick="showTab('data')">Data View</div>
            <div class="tab" onclick="showTab('mermaid')">Mermaid Code</div>
            <div class="tab" onclick="showTab('json')">JSON</div>
        </div>

        <div id="diagram" class="panel active">
            <div class="card">
                %s
            </div>
        </div>

        <div id="data" class="panel">
            <div class="card">
                <div id="data-view"></div>
            </div>
        </div>

        <div id="mermaid" class="panel">
            <div class="card">
                <pre><code>%s</code></pre>
            </div>
        </div>

        <div id="json" class="panel">
            <div class="card">
                <pre><code id="json-view"></code></pre>
            </div>
        </div>
    </div>

    <script>
        %s

        const canvasData = %s;

        function showTab(name) {
            document.querySelectorAll('.tab').forEach(t => t.classList.remove('active'));
            document.querySelectorAll('.panel').forEach(p => p.classList.remove('active'));
            event.target.classList.add('active');
            document.getElementById(name).classList.add('active');
        }

        // Render JSON
        document.getElementById('json-view').textContent = JSON.stringify(canvasData, null, 2);

        // Render data view
        function renderDataView() {
            const container = document.getElementById('data-view');
            const data = canvasData.data;
            if (!data) return;

            let html = '<div class="grid">';

            // Generic renderer based on canvas type
            const type = canvasData.canvasType;

            if (type === 'bmc' && data.metadata) {
                html += renderSection('Customer Segments', data.customerSegments?.map(s => s.name + ': ' + (s.description || '')) || []);
                html += renderSection('Value Propositions', data.valuePropositions?.map(v => v.description) || []);
                html += renderSection('Key Resources', data.keyResources?.map(r => r.name + ' (' + r.type + ')') || []);
                html += renderSection('Key Partners', data.keyPartnerships?.map(p => p.partner) || []);
                html += renderSection('Revenue Streams', data.revenueStreams?.map(r => r.description) || []);
                html += renderSection('Cost Structure', data.costStructure?.map(c => c.description) || []);
            } else if (type === 'ost' && data.outcome) {
                html += '<div class="data-card"><h3>Outcome</h3><p>' + data.outcome.description + '</p>';
                if (data.outcome.metric) html += '<p><strong>Metric:</strong> ' + data.outcome.metric + '</p>';
                if (data.outcome.target) html += '<p><strong>Target:</strong> ' + data.outcome.target + '</p>';
                html += '</div>';
                for (const opp of (data.outcome.opportunities || [])) {
                    html += '<div class="data-card"><h3>Opportunity: ' + opp.description.slice(0, 50) + '</h3>';
                    html += '<ul>';
                    for (const sol of (opp.solutions || [])) {
                        html += '<li>' + sol.description + ' <span class="badge badge-' + sol.status + '">' + sol.status + '</span></li>';
                    }
                    html += '</ul></div>';
                }
            } else if (type === 'opportunity') {
                html += renderSection('Problems', data.problems?.map(p => p.description) || []);
                html += renderSection('Users', data.users?.map(u => u.name + ': ' + (u.description || '')) || []);
                html += '<div class="data-card"><h3>Value Proposition</h3><p>' + (data.valueProposition?.statement || '') + '</p></div>';
                html += renderSection('User Value', data.userValue || []);
                html += renderSection('Business Value', data.businessValue || []);
            } else if (type === 'feature') {
                html += '<div class="data-card" style="grid-column: 1 / -1"><h3>Idea Statement</h3><p>' + (data.ideaStatement || '') + '</p></div>';
                html += renderSection('Situations', data.situations?.map(s => s.description) || []);
                html += renderSection('Problems', data.problems?.map(p => p.description) || []);
                html += renderSection('Value', data.value?.map(v => v.description) || []);
                html += renderBadgedSection('Capabilities', data.capabilities?.map(c => ({text: c.description, badge: c.priority, class: 'badge-' + c.priority})) || []);
                html += renderSection('Restrictions', data.restrictions || []);
                html += renderSection('Limitations', data.limitations || []);
            } else if (type === 'leanux') {
                html += '<div class="data-card" style="grid-column: 1 / -1"><h3>Business Problem</h3><p>' + (data.businessProblem || '') + '</p></div>';
                html += renderSection('Business Outcomes', data.businessOutcomes?.map(o => o.description + ' (' + o.target + ')') || []);
                html += renderSection('Users', data.users?.map(u => u.name) || []);
                html += renderSection('User Outcomes', data.userOutcomes?.map(o => o.description) || []);
                html += renderSection('Solutions', data.solutions?.map(s => s.description) || []);
                html += renderBadgedSection('Hypotheses', data.hypotheses?.map(h => ({
                    text: h.weBelieve,
                    badge: h.validated === true ? 'validated' : (h.validated === false ? 'invalidated' : 'untested'),
                    class: 'badge-' + (h.validated === true ? 'validated' : 'running')
                })) || []);
            }

            html += '</div>';
            container.innerHTML = html;
        }

        function renderSection(title, items) {
            if (!items || items.length === 0) return '';
            return '<div class="data-card"><h3>' + title + '</h3><ul>' +
                items.map(i => '<li>' + i + '</li>').join('') + '</ul></div>';
        }

        function renderBadgedSection(title, items) {
            if (!items || items.length === 0) return '';
            return '<div class="data-card"><h3>' + title + '</h3><ul>' +
                items.map(i => '<li>' + i.text + ' <span class="badge ' + i.class + '">' + i.badge + '</span></li>').join('') + '</ul></div>';
        }

        renderDataView();
    </script>
</body>
</html>
`, name, mermaidScript, name, diagramContent, mermaidCode, mermaidInit, litJSON)
}

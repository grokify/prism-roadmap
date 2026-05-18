# Integration Examples

Examples of integrating prism-execution with other tools and workflows.

## CI/CD Integration

### GitHub Actions

Validate PRDs in your CI pipeline:

```yaml
name: PRD Validation

on:
  pull_request:
    paths:
      - '**/*.prd.json'

jobs:
  validate:
    runs-on: ubuntu-latest
    steps:
      - uses: actions/checkout@v4

      - name: Setup Go
        uses: actions/setup-go@v5
        with:
          go-version: '1.21'

      - name: Install validator
        run: go install github.com/grokify/prism-execution/cmd/srequirements@latest

      - name: Validate PRDs
        run: |
          for file in $(find . -name "*.prd.json"); do
            echo "Validating $file..."
            srequirements validate "$file"
          done
```

### Scoring Gate

Block PRs with low-scoring PRDs:

```yaml
- name: Score PRD
  run: |
    SCORE=$(srequirements score docs/feature.prd.json --output=score)
    if [ "$SCORE" -lt 70 ]; then
      echo "PRD score ($SCORE%) below threshold (70%)"
      exit 1
    fi
```

## Slack Notifications

Send PRD summaries to Slack:

```go
package main

import (
    "bytes"
    "encoding/json"
    "net/http"

    "github.com/grokify/prism-execution/requirements/prd"
)

func notifySlack(doc *prd.Document, webhookURL string) error {
    scores := prd.Score(doc)
    execView := prd.GenerateExecView(doc, scores)

    message := map[string]interface{}{
        "blocks": []map[string]interface{}{
            {
                "type": "header",
                "text": map[string]string{
                    "type": "plain_text",
                    "text": "PRD Review: " + doc.Metadata.Title,
                },
            },
            {
                "type": "section",
                "fields": []map[string]string{
                    {"type": "mrkdwn", "text": "*Decision:*\n" + execView.Decision},
                    {"type": "mrkdwn", "text": "*Score:*\n" + execView.ScoreSummary},
                },
            },
        },
    }

    body, _ := json.Marshal(message)
    _, err := http.Post(webhookURL, "application/json", bytes.NewReader(body))
    return err
}
```

## Notion Integration

Export PRD to Notion:

```go
package main

import (
    "github.com/grokify/prism-execution/requirements/prd"
    // hypothetical notion client
    "github.com/example/notion-go"
)

func exportToNotion(doc *prd.Document, pageID string) error {
    client := notion.NewClient(os.Getenv("NOTION_TOKEN"))

    // Generate PM view markdown
    pmView := prd.GeneratePMView(doc)
    markdown := prd.RenderPMMarkdown(pmView)

    // Create Notion page
    return client.CreatePage(notion.Page{
        Parent:   notion.PageID(pageID),
        Title:    doc.Metadata.Title,
        Content:  markdown,
        Properties: map[string]interface{}{
            "Status":  doc.Metadata.Status,
            "Version": doc.Metadata.Version,
            "Author":  doc.Metadata.Authors[0].Name,
        },
    })
}
```

## Jira Integration

Create Jira epics from PRD requirements:

```go
func createJiraEpics(doc *prd.Document) error {
    client := jira.NewClient(os.Getenv("JIRA_URL"), os.Getenv("JIRA_TOKEN"))

    for _, req := range doc.Requirements.Functional {
        if req.MoSCoW == prd.MoSCoWMust {
            epic := jira.Epic{
                Summary:     req.Title,
                Description: req.Description,
                Labels:      []string{"prd:" + doc.Metadata.ID},
                Priority:    mapPriority(req.Priority),
            }

            if err := client.CreateEpic(epic); err != nil {
                return err
            }
        }
    }
    return nil
}
```

## Confluence Export

Publish 6-pager to Confluence:

```go
func publishToConfluence(doc *prd.Document, spaceKey, parentID string) error {
    client := confluence.NewClient(
        os.Getenv("CONFLUENCE_URL"),
        os.Getenv("CONFLUENCE_TOKEN"),
    )

    // Generate 6-pager
    sixPager := prd.GenerateSixPagerView(doc)
    content := prd.RenderSixPagerMarkdown(sixPager)

    return client.CreatePage(confluence.Page{
        SpaceKey: spaceKey,
        ParentID: parentID,
        Title:    doc.Metadata.Title + " - 6-Pager",
        Body:     content,
    })
}
```

## Goals Sync

Sync PRD metrics with OKR tracking:

```go
func syncWithOKRTool(doc *prd.Document) error {
    if doc.Goals == nil || doc.Goals.OKR == nil {
        return nil
    }

    client := okrtool.NewClient(os.Getenv("OKR_API_KEY"))

    for _, obj := range doc.Goals.OKR.Objectives {
        for _, kr := range obj.KeyResults {
            // Update KR progress in external tool
            err := client.UpdateKeyResult(kr.ID, okrtool.Update{
                Score:   kr.Score,
                Current: kr.Current,
            })
            if err != nil {
                return err
            }
        }
    }
    return nil
}
```

## Automated Reviews

Weekly PRD review automation:

```go
func weeklyReview() {
    prds, _ := filepath.Glob("prds/*.prd.json")

    var report strings.Builder
    report.WriteString("# Weekly PRD Review\n\n")

    for _, path := range prds {
        doc, _ := prd.Load(path)
        scores := prd.Score(doc)

        report.WriteString(fmt.Sprintf("## %s\n", doc.Metadata.Title))
        report.WriteString(fmt.Sprintf("- Score: %.0f%%\n", scores.OverallScore*100))
        report.WriteString(fmt.Sprintf("- Decision: %s\n", scores.Decision))

        if len(scores.Triggers) > 0 {
            report.WriteString("- Issues:\n")
            for _, t := range scores.Triggers {
                report.WriteString(fmt.Sprintf("  - %s\n", t.Issue))
            }
        }
        report.WriteString("\n")
    }

    // Send report
    sendEmail("team@example.com", "Weekly PRD Review", report.String())
}
```

## Template Generation

Generate PRD from templates:

```go
func generateFromTemplate(templateName string, data map[string]string) *prd.Document {
    templates := map[string]func(map[string]string) *prd.Document{
        "feature":  createFeaturePRD,
        "platform": createPlatformPRD,
        "api":      createAPIPRD,
    }

    generator := templates[templateName]
    if generator == nil {
        generator = createFeaturePRD
    }

    return generator(data)
}

func createFeaturePRD(data map[string]string) *prd.Document {
    doc := prd.New(prd.GenerateID(), data["title"],
        prd.Person{Name: data["author"]})

    doc.ExecutiveSummary.ProblemStatement = data["problem"]
    doc.ExecutiveSummary.ProposedSolution = data["solution"]

    return doc
}
```

## Watch Mode

Auto-validate on file changes:

```go
func watchPRDs(dir string) {
    watcher, _ := fsnotify.NewWatcher()
    defer watcher.Close()

    watcher.Add(dir)

    for {
        select {
        case event := <-watcher.Events:
            if strings.HasSuffix(event.Name, ".prd.json") {
                if event.Op&fsnotify.Write == fsnotify.Write {
                    validateAndNotify(event.Name)
                }
            }
        case err := <-watcher.Errors:
            log.Println("error:", err)
        }
    }
}

func validateAndNotify(path string) {
    doc, err := prd.Load(path)
    if err != nil {
        notify("Validation Error", err.Error())
        return
    }

    result := prd.Validate(doc)
    if !result.Valid {
        notify("PRD Invalid", formatErrors(result.Errors))
    }
}
```

## Next Steps

- [PRD Examples](prd-examples.md)
- [Quick Start](../getting-started/quickstart.md)
- [API Reference](../api/reference.md)

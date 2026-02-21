package prd

import (
	"encoding/json"
	"fmt"
	"os"
	"path/filepath"
	"strconv"
	"strings"
	"time"
)

// DefaultFilename is the standard PRD filename.
const DefaultFilename = "PRD.json"

// Load reads a Document from a JSON file.
func Load(path string) (*Document, error) {
	data, err := os.ReadFile(filepath.Clean(path))
	if err != nil {
		return nil, fmt.Errorf("reading PRD file: %w", err)
	}

	var doc Document
	if err := json.Unmarshal(data, &doc); err != nil {
		return nil, fmt.Errorf("parsing PRD JSON: %w", err)
	}

	return &doc, nil
}

// Save writes a Document to a JSON file.
func Save(doc *Document, path string) error {
	doc.Metadata.UpdatedAt = time.Now()

	data, err := json.MarshalIndent(doc, "", "  ")
	if err != nil {
		return fmt.Errorf("marshaling PRD: %w", err)
	}

	if err := os.MkdirAll(filepath.Dir(path), 0755); err != nil {
		return fmt.Errorf("creating directory: %w", err)
	}

	if err := os.WriteFile(path, data, 0600); err != nil {
		return fmt.Errorf("writing PRD file: %w", err)
	}

	return nil
}

// New creates a new Document with required fields initialized.
func New(id, title string, authors ...Person) *Document {
	now := time.Now()

	doc := &Document{
		Metadata: Metadata{
			ID:        id,
			Title:     title,
			Version:   "1.0.0",
			Status:    StatusDraft,
			CreatedAt: now,
			UpdatedAt: now,
			Authors:   authors,
		},
		ExecutiveSummary: ExecutiveSummary{},
		Objectives: Objectives{
			OKRs: []OKR{},
		},
		Personas:     []Persona{},
		UserStories:  []UserStory{},
		Requirements: Requirements{},
		Roadmap: Roadmap{
			Phases: []Phase{},
		},
		RevisionHistory: []RevisionRecord{
			{
				Version: "1.0.0",
				Changes: []string{"Initial PRD creation"},
				Trigger: TriggerInitial,
				Date:    now,
			},
		},
	}

	return doc
}

// GenerateID generates a PRD ID based on the current date.
// Format: PRD-YYYY-DDD where DDD is the day of year.
func GenerateID() string {
	now := time.Now()
	return fmt.Sprintf("PRD-%d-%03d", now.Year(), now.YearDay())
}

// GenerateIDWithPrefix generates an ID with a custom prefix.
// Format: PREFIX-YYYY-DDD where DDD is the day of year.
func GenerateIDWithPrefix(prefix string) string {
	now := time.Now()
	return fmt.Sprintf("%s-%d-%03d", prefix, now.Year(), now.YearDay())
}

// AddRevision records a revision in the Document history.
func (doc *Document) AddRevision(changes []string, trigger RevisionTriggerType, author string) {
	// Increment version
	doc.Metadata.Version = incrementVersion(doc.Metadata.Version)
	doc.Metadata.UpdatedAt = time.Now()

	doc.RevisionHistory = append(doc.RevisionHistory, RevisionRecord{
		Version: doc.Metadata.Version,
		Changes: changes,
		Trigger: trigger,
		Date:    time.Now(),
		Author:  author,
	})
}

// incrementVersion increments the patch version number.
// Example: "1.0.0" -> "1.0.1", "2.3.4" -> "2.3.5"
func incrementVersion(version string) string {
	parts := strings.Split(version, ".")
	if len(parts) != 3 {
		return version + ".1"
	}

	patch, err := strconv.Atoi(parts[2])
	if err != nil {
		return version
	}

	return fmt.Sprintf("%s.%s.%d", parts[0], parts[1], patch+1)
}

// ToJSON converts the Document to JSON string.
func (doc *Document) ToJSON() (string, error) {
	data, err := json.MarshalIndent(doc, "", "  ")
	if err != nil {
		return "", fmt.Errorf("marshaling document: %w", err)
	}
	return string(data), nil
}

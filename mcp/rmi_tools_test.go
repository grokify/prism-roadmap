package mcp

import (
	"context"
	"os"
	"path/filepath"
	"testing"

	"github.com/grokify/prism-roadmap/prioritization"
	"github.com/grokify/prism-roadmap/rmi"
)

func setupTestServer(t *testing.T) (*Server, string, func()) {
	t.Helper()

	tmpDir, err := os.MkdirTemp("", "mcp-test-*")
	if err != nil {
		t.Fatalf("creating temp dir: %v", err)
	}

	cleanup := func() {
		_ = os.RemoveAll(tmpDir)
	}

	return NewServer(), tmpDir, cleanup
}

func createTestRoadmap(t *testing.T, dir string) string {
	t.Helper()

	filePath := filepath.Join(dir, "test-roadmap.json")
	svc := rmi.NewService()

	_, err := svc.Create(rmi.CreateInput{
		ID:          "item-1",
		Name:        "Test Item 1",
		Description: "First test item",
		MoSCoW:      prioritization.MoSCoWMustHave,
		Quarter:     "Q1 2025",
		Owner:       "Alice",
		Tags:        []string{"api", "core"},
	})
	if err != nil {
		t.Fatalf("creating test item 1: %v", err)
	}

	_, err = svc.Create(rmi.CreateInput{
		ID:          "item-2",
		Name:        "Test Item 2",
		Description: "Second test item",
		MoSCoW:      prioritization.MoSCoWShouldHave,
		Quarter:     "Q2 2025",
		Owner:       "Bob",
		Tags:        []string{"ui"},
	})
	if err != nil {
		t.Fatalf("creating test item 2: %v", err)
	}

	_, err = svc.Create(rmi.CreateInput{
		ID:          "item-3",
		Name:        "Test Item 3",
		Description: "Third test item",
		MoSCoW:      prioritization.MoSCoWCouldHave,
		Quarter:     "Q1 2025",
		Owner:       "Alice",
		Tags:        []string{"api", "experimental"},
	})
	if err != nil {
		t.Fatalf("creating test item 3: %v", err)
	}

	if err := svc.SaveAs(filePath); err != nil {
		t.Fatalf("saving test roadmap: %v", err)
	}

	return filePath
}

func TestServer_listRMIs(t *testing.T) {
	srv, tmpDir, cleanup := setupTestServer(t)
	defer cleanup()

	filePath := createTestRoadmap(t, tmpDir)

	t.Run("list all items", func(t *testing.T) {
		_, output, err := srv.listRMIs(context.Background(), nil, ListRMIsInput{
			File: filePath,
		})
		if err != nil {
			t.Fatalf("listRMIs() error = %v", err)
		}

		if output.Error != "" {
			t.Errorf("unexpected error: %s", output.Error)
		}
		if output.Total != 3 {
			t.Errorf("Total = %d, want 3", output.Total)
		}
		if len(output.Items) != 3 {
			t.Errorf("len(Items) = %d, want 3", len(output.Items))
		}
	})

	t.Run("filter by moscow", func(t *testing.T) {
		_, output, err := srv.listRMIs(context.Background(), nil, ListRMIsInput{
			File:   filePath,
			MoSCoW: "must_have",
		})
		if err != nil {
			t.Fatalf("listRMIs() error = %v", err)
		}

		if output.Total != 1 {
			t.Errorf("Total = %d, want 1", output.Total)
		}
	})

	t.Run("filter by quarter", func(t *testing.T) {
		_, output, err := srv.listRMIs(context.Background(), nil, ListRMIsInput{
			File:    filePath,
			Quarter: "Q1 2025",
		})
		if err != nil {
			t.Fatalf("listRMIs() error = %v", err)
		}

		if output.Total != 2 {
			t.Errorf("Total = %d, want 2", output.Total)
		}
	})

	t.Run("with limit", func(t *testing.T) {
		_, output, err := srv.listRMIs(context.Background(), nil, ListRMIsInput{
			File:  filePath,
			Limit: 2,
		})
		if err != nil {
			t.Fatalf("listRMIs() error = %v", err)
		}

		if len(output.Items) != 2 {
			t.Errorf("len(Items) = %d, want 2", len(output.Items))
		}
	})

	t.Run("nonexistent file", func(t *testing.T) {
		_, output, err := srv.listRMIs(context.Background(), nil, ListRMIsInput{
			File: filepath.Join(tmpDir, "nonexistent.json"),
		})
		if err != nil {
			t.Fatalf("listRMIs() error = %v", err)
		}

		if output.Error == "" {
			t.Error("expected error for nonexistent file")
		}
	})
}

func TestServer_getRMI(t *testing.T) {
	srv, tmpDir, cleanup := setupTestServer(t)
	defer cleanup()

	filePath := createTestRoadmap(t, tmpDir)

	t.Run("get existing item", func(t *testing.T) {
		_, output, err := srv.getRMI(context.Background(), nil, GetRMIInput{
			File: filePath,
			ID:   "item-1",
		})
		if err != nil {
			t.Fatalf("getRMI() error = %v", err)
		}

		if output.Error != "" {
			t.Errorf("unexpected error: %s", output.Error)
		}
		if output.Item == nil {
			t.Fatal("Item is nil")
		}
		if output.Item.ID != "item-1" {
			t.Errorf("ID = %s, want item-1", output.Item.ID)
		}
		if output.Item.Name != "Test Item 1" {
			t.Errorf("Name = %s, want Test Item 1", output.Item.Name)
		}
	})

	t.Run("get nonexistent item", func(t *testing.T) {
		_, output, err := srv.getRMI(context.Background(), nil, GetRMIInput{
			File: filePath,
			ID:   "nonexistent",
		})
		if err != nil {
			t.Fatalf("getRMI() error = %v", err)
		}

		if output.Error == "" {
			t.Error("expected error for nonexistent ID")
		}
	})
}

func TestServer_createRMI(t *testing.T) {
	srv, tmpDir, cleanup := setupTestServer(t)
	defer cleanup()

	filePath := filepath.Join(tmpDir, "new-roadmap.json")

	t.Run("create new item", func(t *testing.T) {
		_, output, err := srv.createRMI(context.Background(), nil, CreateRMIInput{
			File:        filePath,
			ID:          "new-item",
			Name:        "New Item",
			Description: "A new item",
			MoSCoW:      "must_have",
			Quarter:     "Q3 2025",
			Owner:       "Charlie",
			Tags:        []string{"new"},
		})
		if err != nil {
			t.Fatalf("createRMI() error = %v", err)
		}

		if output.Error != "" {
			t.Errorf("unexpected error: %s", output.Error)
		}
		if !output.Created {
			t.Error("Created = false, want true")
		}
		if output.Item == nil {
			t.Fatal("Item is nil")
		}
		if output.Item.ID != "new-item" {
			t.Errorf("ID = %s, want new-item", output.Item.ID)
		}
	})

	t.Run("invalid moscow priority", func(t *testing.T) {
		_, output, err := srv.createRMI(context.Background(), nil, CreateRMIInput{
			File:   filePath,
			ID:     "invalid-item",
			Name:   "Invalid Item",
			MoSCoW: "invalid_priority",
		})
		if err != nil {
			t.Fatalf("createRMI() error = %v", err)
		}

		if output.Error == "" {
			t.Error("expected error for invalid moscow priority")
		}
	})
}

func TestServer_updateRMI(t *testing.T) {
	srv, tmpDir, cleanup := setupTestServer(t)
	defer cleanup()

	filePath := createTestRoadmap(t, tmpDir)

	t.Run("update existing item", func(t *testing.T) {
		_, output, err := srv.updateRMI(context.Background(), nil, UpdateRMIInput{
			File:     filePath,
			ID:       "item-1",
			Name:     "Updated Item 1",
			Status:   "in_progress",
			Progress: 50,
		})
		if err != nil {
			t.Fatalf("updateRMI() error = %v", err)
		}

		if output.Error != "" {
			t.Errorf("unexpected error: %s", output.Error)
		}
		if !output.Updated {
			t.Error("Updated = false, want true")
		}
		if output.Item == nil {
			t.Fatal("Item is nil")
		}
		if output.Item.Name != "Updated Item 1" {
			t.Errorf("Name = %s, want Updated Item 1", output.Item.Name)
		}
		if output.Item.Status != "in_progress" {
			t.Errorf("Status = %s, want in_progress", output.Item.Status)
		}
	})

	t.Run("update nonexistent item", func(t *testing.T) {
		_, output, err := srv.updateRMI(context.Background(), nil, UpdateRMIInput{
			File: filePath,
			ID:   "nonexistent",
			Name: "New Name",
		})
		if err != nil {
			t.Fatalf("updateRMI() error = %v", err)
		}

		if output.Error == "" {
			t.Error("expected error for nonexistent ID")
		}
	})
}

func TestServer_deleteRMI(t *testing.T) {
	srv, tmpDir, cleanup := setupTestServer(t)
	defer cleanup()

	filePath := createTestRoadmap(t, tmpDir)

	t.Run("delete existing item", func(t *testing.T) {
		_, output, err := srv.deleteRMI(context.Background(), nil, DeleteRMIInput{
			File: filePath,
			ID:   "item-1",
		})
		if err != nil {
			t.Fatalf("deleteRMI() error = %v", err)
		}

		if output.Error != "" {
			t.Errorf("unexpected error: %s", output.Error)
		}
		if !output.Deleted {
			t.Error("Deleted = false, want true")
		}

		// Verify deletion
		_, getOutput, _ := srv.getRMI(context.Background(), nil, GetRMIInput{
			File: filePath,
			ID:   "item-1",
		})
		if getOutput.Error == "" {
			t.Error("item should be deleted")
		}
	})

	t.Run("delete nonexistent item", func(t *testing.T) {
		_, output, err := srv.deleteRMI(context.Background(), nil, DeleteRMIInput{
			File: filePath,
			ID:   "nonexistent",
		})
		if err != nil {
			t.Fatalf("deleteRMI() error = %v", err)
		}

		if output.Error == "" {
			t.Error("expected error for nonexistent ID")
		}
	})
}

func TestServer_rmiSummary(t *testing.T) {
	srv, tmpDir, cleanup := setupTestServer(t)
	defer cleanup()

	filePath := createTestRoadmap(t, tmpDir)

	_, output, err := srv.rmiSummary(context.Background(), nil, RMISummaryInput{
		File: filePath,
	})
	if err != nil {
		t.Fatalf("rmiSummary() error = %v", err)
	}

	if output.Error != "" {
		t.Errorf("unexpected error: %s", output.Error)
	}
	if output.Summary == nil {
		t.Fatal("Summary is nil")
	}
	if output.Summary.TotalItems != 3 {
		t.Errorf("TotalItems = %d, want 3", output.Summary.TotalItems)
	}
}

func TestServer_topRMIs(t *testing.T) {
	srv, tmpDir, cleanup := setupTestServer(t)
	defer cleanup()

	filePath := createTestRoadmap(t, tmpDir)

	t.Run("top by priority", func(t *testing.T) {
		_, output, err := srv.topRMIs(context.Background(), nil, TopRMIsInput{
			File:   filePath,
			SortBy: "priority",
			Limit:  2,
		})
		if err != nil {
			t.Fatalf("topRMIs() error = %v", err)
		}

		if output.Error != "" {
			t.Errorf("unexpected error: %s", output.Error)
		}
		if output.Total != 2 {
			t.Errorf("Total = %d, want 2", output.Total)
		}
		if output.SortBy != "priority" {
			t.Errorf("SortBy = %s, want priority", output.SortBy)
		}
	})

	t.Run("top by rice", func(t *testing.T) {
		_, output, err := srv.topRMIs(context.Background(), nil, TopRMIsInput{
			File:   filePath,
			SortBy: "rice",
			Limit:  3,
		})
		if err != nil {
			t.Fatalf("topRMIs() error = %v", err)
		}

		if output.Error != "" {
			t.Errorf("unexpected error: %s", output.Error)
		}
		if output.SortBy != "rice" {
			t.Errorf("SortBy = %s, want rice", output.SortBy)
		}
	})

	t.Run("default values", func(t *testing.T) {
		_, output, err := srv.topRMIs(context.Background(), nil, TopRMIsInput{
			File: filePath,
		})
		if err != nil {
			t.Fatalf("topRMIs() error = %v", err)
		}

		if output.SortBy != "priority" {
			t.Errorf("SortBy = %s, want priority (default)", output.SortBy)
		}
	})
}

func TestServer_ServiceCaching(t *testing.T) {
	srv, tmpDir, cleanup := setupTestServer(t)
	defer cleanup()

	filePath := createTestRoadmap(t, tmpDir)

	// First call should load the service
	_, output1, _ := srv.listRMIs(context.Background(), nil, ListRMIsInput{
		File: filePath,
	})
	if output1.Error != "" {
		t.Fatalf("first listRMIs error: %s", output1.Error)
	}

	// Second call should use cached service
	_, output2, _ := srv.listRMIs(context.Background(), nil, ListRMIsInput{
		File: filePath,
	})
	if output2.Error != "" {
		t.Fatalf("second listRMIs error: %s", output2.Error)
	}

	// Both should return same results
	if output1.Total != output2.Total {
		t.Errorf("cached results differ: %d vs %d", output1.Total, output2.Total)
	}
}

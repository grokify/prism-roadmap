package rmi

import (
	"os"
	"path/filepath"
	"testing"

	"github.com/grokify/prism-roadmap/prioritization"
)

func setupTestService(t *testing.T) (*Service, func()) {
	t.Helper()

	tmpDir, err := os.MkdirTemp("", "rmi-test-*")
	if err != nil {
		t.Fatalf("creating temp dir: %v", err)
	}

	cleanup := func() {
		os.RemoveAll(tmpDir)
	}

	return NewService(), cleanup
}

func mustCreate(t *testing.T, svc *Service, input CreateInput) {
	t.Helper()
	if _, err := svc.Create(input); err != nil {
		t.Fatalf("Create() error = %v", err)
	}
}

func TestService_Create(t *testing.T) {
	svc, cleanup := setupTestService(t)
	defer cleanup()

	input := CreateInput{
		ID:          "rmi-1",
		Name:        "SSO Integration",
		Description: "Add single sign-on support",
		MoSCoW:      prioritization.MoSCoWMustHave,
		Quarter:     "Q3 2026",
		Owner:       "alice",
		Tags:        []string{"security", "enterprise"},
	}

	item, err := svc.Create(input)
	if err != nil {
		t.Fatalf("Create() error = %v", err)
	}

	if item.ID != input.ID {
		t.Errorf("ID = %s, want %s", item.ID, input.ID)
	}
	if item.Name != input.Name {
		t.Errorf("Name = %s, want %s", item.Name, input.Name)
	}
	if item.MoSCoW != input.MoSCoW {
		t.Errorf("MoSCoW = %s, want %s", item.MoSCoW, input.MoSCoW)
	}
	if item.Status != RMIStatusPlanned {
		t.Errorf("Status = %s, want %s", item.Status, RMIStatusPlanned)
	}
	if len(item.Tags) != 2 {
		t.Errorf("len(Tags) = %d, want 2", len(item.Tags))
	}
}

func TestService_Create_Duplicate(t *testing.T) {
	svc, cleanup := setupTestService(t)
	defer cleanup()

	input := CreateInput{
		ID:     "rmi-1",
		Name:   "Test Item",
		MoSCoW: prioritization.MoSCoWMustHave,
	}

	_, err := svc.Create(input)
	if err != nil {
		t.Fatalf("First Create() error = %v", err)
	}

	_, err = svc.Create(input)
	if err == nil {
		t.Error("Expected error for duplicate ID")
	}
}

func TestService_Get(t *testing.T) {
	svc, cleanup := setupTestService(t)
	defer cleanup()

	input := CreateInput{
		ID:     "rmi-1",
		Name:   "Test Item",
		MoSCoW: prioritization.MoSCoWMustHave,
	}

	created, _ := svc.Create(input)

	item, err := svc.Get("rmi-1")
	if err != nil {
		t.Fatalf("Get() error = %v", err)
	}

	if item.ID != created.ID {
		t.Errorf("ID = %s, want %s", item.ID, created.ID)
	}
}

func TestService_Get_NotFound(t *testing.T) {
	svc, cleanup := setupTestService(t)
	defer cleanup()

	_, err := svc.Get("nonexistent")
	if err == nil {
		t.Error("Expected error for nonexistent ID")
	}
}

func TestService_List(t *testing.T) {
	svc, cleanup := setupTestService(t)
	defer cleanup()

	// Create multiple items
	items := []CreateInput{
		{ID: "rmi-1", Name: "Item 1", MoSCoW: prioritization.MoSCoWMustHave, Quarter: "Q1 2026"},
		{ID: "rmi-2", Name: "Item 2", MoSCoW: prioritization.MoSCoWShouldHave, Quarter: "Q1 2026"},
		{ID: "rmi-3", Name: "Item 3", MoSCoW: prioritization.MoSCoWCouldHave, Quarter: "Q2 2026"},
		{ID: "rmi-4", Name: "Item 4", MoSCoW: prioritization.MoSCoWMustHave, Quarter: "Q2 2026"},
		{ID: "rmi-5", Name: "Item 5", MoSCoW: prioritization.MoSCoWWontHave, Quarter: "Q3 2026"},
	}
	for _, input := range items {
		if _, err := svc.Create(input); err != nil {
			t.Fatalf("Create() error = %v", err)
		}
	}

	// Test no filter
	result := svc.List(ListFilter{})
	if len(result) != 5 {
		t.Errorf("len(result) = %d, want 5", len(result))
	}

	// Test MoSCoW filter
	result = svc.List(ListFilter{MoSCoW: prioritization.MoSCoWMustHave})
	if len(result) != 2 {
		t.Errorf("len(result) = %d, want 2", len(result))
	}

	// Test quarter filter
	result = svc.List(ListFilter{Quarter: "Q1 2026"})
	if len(result) != 2 {
		t.Errorf("len(result) = %d, want 2", len(result))
	}

	// Test limit
	result = svc.List(ListFilter{Limit: 3})
	if len(result) != 3 {
		t.Errorf("len(result) = %d, want 3", len(result))
	}
}

func TestService_Update(t *testing.T) {
	svc, cleanup := setupTestService(t)
	defer cleanup()

	if _, err := svc.Create(CreateInput{
		ID:     "rmi-1",
		Name:   "Test Item",
		MoSCoW: prioritization.MoSCoWShouldHave,
	}); err != nil {
		t.Fatalf("Create() error = %v", err)
	}

	newName := "Updated Item"
	newMoSCoW := prioritization.MoSCoWMustHave
	newStatus := RMIStatusInProgress

	item, updated, err := svc.Update("rmi-1", UpdateInput{
		Name:   &newName,
		MoSCoW: &newMoSCoW,
		Status: &newStatus,
	})
	if err != nil {
		t.Fatalf("Update() error = %v", err)
	}

	if !updated {
		t.Error("Expected updated = true")
	}
	if item.Name != newName {
		t.Errorf("Name = %s, want %s", item.Name, newName)
	}
	if item.MoSCoW != newMoSCoW {
		t.Errorf("MoSCoW = %s, want %s", item.MoSCoW, newMoSCoW)
	}
	if item.Status != newStatus {
		t.Errorf("Status = %s, want %s", item.Status, newStatus)
	}
}

func TestService_Update_NoChange(t *testing.T) {
	svc, cleanup := setupTestService(t)
	defer cleanup()

	if _, err := svc.Create(CreateInput{
		ID:     "rmi-1",
		Name:   "Test Item",
		MoSCoW: prioritization.MoSCoWMustHave,
	}); err != nil {
		t.Fatalf("Create() error = %v", err)
	}

	sameName := "Test Item"
	_, updated, err := svc.Update("rmi-1", UpdateInput{
		Name: &sameName,
	})
	if err != nil {
		t.Fatalf("Update() error = %v", err)
	}

	if updated {
		t.Error("Expected updated = false when value unchanged")
	}
}

func TestService_Delete(t *testing.T) {
	svc, cleanup := setupTestService(t)
	defer cleanup()

	if _, err := svc.Create(CreateInput{
		ID:     "rmi-1",
		Name:   "Test Item",
		MoSCoW: prioritization.MoSCoWMustHave,
	}); err != nil {
		t.Fatalf("Create() error = %v", err)
	}

	err := svc.Delete("rmi-1")
	if err != nil {
		t.Fatalf("Delete() error = %v", err)
	}

	_, err = svc.Get("rmi-1")
	if err == nil {
		t.Error("Expected error after delete")
	}
}

func TestService_Summary(t *testing.T) {
	svc, cleanup := setupTestService(t)
	defer cleanup()

	// Create items with different statuses and priorities
	mustCreate(t, svc, CreateInput{ID: "rmi-1", Name: "Item 1", MoSCoW: prioritization.MoSCoWMustHave, Quarter: "Q1 2026"})
	mustCreate(t, svc, CreateInput{ID: "rmi-2", Name: "Item 2", MoSCoW: prioritization.MoSCoWMustHave, Quarter: "Q1 2026"})
	mustCreate(t, svc, CreateInput{ID: "rmi-3", Name: "Item 3", MoSCoW: prioritization.MoSCoWShouldHave, Quarter: "Q2 2026"})

	summary := svc.Summary()

	if summary.TotalItems != 3 {
		t.Errorf("TotalItems = %d, want 3", summary.TotalItems)
	}
	if summary.MoSCoWCounts[prioritization.MoSCoWMustHave] != 2 {
		t.Errorf("MoSCoWCounts[MustHave] = %d, want 2", summary.MoSCoWCounts[prioritization.MoSCoWMustHave])
	}
	if summary.ActionableCount != 3 {
		t.Errorf("ActionableCount = %d, want 3", summary.ActionableCount)
	}
	if summary.QuarterCounts["Q1 2026"] != 2 {
		t.Errorf("QuarterCounts[Q1 2026] = %d, want 2", summary.QuarterCounts["Q1 2026"])
	}
}

func TestService_TopByPriority(t *testing.T) {
	svc, cleanup := setupTestService(t)
	defer cleanup()

	mustCreate(t, svc, CreateInput{ID: "rmi-1", Name: "Item 1", MoSCoW: prioritization.MoSCoWCouldHave})
	mustCreate(t, svc, CreateInput{ID: "rmi-2", Name: "Item 2", MoSCoW: prioritization.MoSCoWMustHave})
	mustCreate(t, svc, CreateInput{ID: "rmi-3", Name: "Item 3", MoSCoW: prioritization.MoSCoWShouldHave})

	top := svc.TopByPriority(2)

	if len(top) != 2 {
		t.Fatalf("len(top) = %d, want 2", len(top))
	}
	if top[0].MoSCoW != prioritization.MoSCoWMustHave {
		t.Errorf("top[0].MoSCoW = %s, want must_have", top[0].MoSCoW)
	}
	if top[1].MoSCoW != prioritization.MoSCoWShouldHave {
		t.Errorf("top[1].MoSCoW = %s, want should_have", top[1].MoSCoW)
	}
}

func TestService_FileOperations(t *testing.T) {
	tmpDir, err := os.MkdirTemp("", "rmi-file-test-*")
	if err != nil {
		t.Fatalf("creating temp dir: %v", err)
	}
	defer os.RemoveAll(tmpDir)

	filePath := filepath.Join(tmpDir, "roadmap.json")

	// Create and save
	svc := NewService()
	mustCreate(t, svc, CreateInput{ID: "rmi-1", Name: "Item 1", MoSCoW: prioritization.MoSCoWMustHave})
	mustCreate(t, svc, CreateInput{ID: "rmi-2", Name: "Item 2", MoSCoW: prioritization.MoSCoWShouldHave})

	if err := svc.SaveAs(filePath); err != nil {
		t.Fatalf("SaveAs() error = %v", err)
	}

	// Load and verify
	loadedSvc, err := NewServiceFromFile(filePath)
	if err != nil {
		t.Fatalf("NewServiceFromFile() error = %v", err)
	}

	items := loadedSvc.List(ListFilter{})
	if len(items) != 2 {
		t.Errorf("loaded len(items) = %d, want 2", len(items))
	}

	// Verify Save() works after loading
	mustCreate(t, svc, CreateInput{ID: "rmi-3", Name: "Item 3", MoSCoW: prioritization.MoSCoWCouldHave})
	if err := svc.Save(); err != nil {
		t.Fatalf("Save() error = %v", err)
	}
}

func TestRoadmapItemSet_WriteFile_ReadFile(t *testing.T) {
	tmpDir, err := os.MkdirTemp("", "rmi-io-test-*")
	if err != nil {
		t.Fatalf("creating temp dir: %v", err)
	}
	defer os.RemoveAll(tmpDir)

	filePath := filepath.Join(tmpDir, "items.json")

	// Create set and save
	set := NewRoadmapItemSet()
	set.Description = "Test roadmap"
	set.Quarter = "Q3 2026"
	set.Add(*NewRoadmapItem("rmi-1", "Test Item", prioritization.MoSCoWMustHave))

	if err := set.WriteFile(filePath); err != nil {
		t.Fatalf("WriteFile() error = %v", err)
	}

	// Read and verify
	loadedSet, err := ReadFile(filePath)
	if err != nil {
		t.Fatalf("ReadFile() error = %v", err)
	}

	if loadedSet.Description != "Test roadmap" {
		t.Errorf("Description = %s, want Test roadmap", loadedSet.Description)
	}
	if loadedSet.Quarter != "Q3 2026" {
		t.Errorf("Quarter = %s, want Q3 2026", loadedSet.Quarter)
	}
	if len(loadedSet.Items) != 1 {
		t.Errorf("len(Items) = %d, want 1", len(loadedSet.Items))
	}
	if loadedSet.Items[0].ID != "rmi-1" {
		t.Errorf("Items[0].ID = %s, want rmi-1", loadedSet.Items[0].ID)
	}
}

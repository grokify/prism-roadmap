package mcp

import (
	"context"
	"path/filepath"

	"github.com/modelcontextprotocol/go-sdk/mcp"

	"github.com/grokify/prism-roadmap/prioritization"
	"github.com/grokify/prism-roadmap/rmi"
)

// RMI Tool Input/Output types

// ListRMIsInput is the input for list_rmis.
type ListRMIsInput struct {
	File    string `json:"file" jsonschema:"description=Path to roadmap JSON file"`
	MoSCoW  string `json:"moscow,omitempty" jsonschema:"description=Filter by MoSCoW priority (must_have, should_have, could_have, wont_have)"`
	Status  string `json:"status,omitempty" jsonschema:"description=Filter by status (planned, in_progress, completed, blocked, cancelled, deferred)"`
	Quarter string `json:"quarter,omitempty" jsonschema:"description=Filter by quarter (e.g., Q3 2026)"`
	Limit   int    `json:"limit,omitempty" jsonschema:"description=Maximum results (default 50)"`
	SortBy  string `json:"sort_by,omitempty" jsonschema:"description=Sort by (priority, rice, market_signal)"`
}

// ListRMIsOutput is the output for list_rmis.
type ListRMIsOutput struct {
	File    string            `json:"file"`
	Items   []rmi.RoadmapItem `json:"items"`
	Total   int               `json:"total"`
	Summary *rmi.Summary      `json:"summary,omitempty"`
	Error   string            `json:"error,omitempty"`
}

// GetRMIInput is the input for get_rmi.
type GetRMIInput struct {
	File string `json:"file" jsonschema:"description=Path to roadmap JSON file"`
	ID   string `json:"id" jsonschema:"description=Roadmap item ID"`
}

// GetRMIOutput is the output for get_rmi.
type GetRMIOutput struct {
	File  string           `json:"file"`
	Item  *rmi.RoadmapItem `json:"item,omitempty"`
	Error string           `json:"error,omitempty"`
}

// CreateRMIInput is the input for create_rmi.
type CreateRMIInput struct {
	File        string   `json:"file" jsonschema:"description=Path to roadmap JSON file"`
	ID          string   `json:"id" jsonschema:"description=Unique item ID"`
	Name        string   `json:"name" jsonschema:"description=Item name"`
	Description string   `json:"description,omitempty" jsonschema:"description=Item description"`
	MoSCoW      string   `json:"moscow,omitempty" jsonschema:"description=Optional MoSCoW priority (must_have, should_have, could_have, wont_have); empty means not yet prioritized"`
	Quarter     string   `json:"quarter,omitempty" jsonschema:"description=Target quarter (e.g., Q3 2026)"`
	Owner       string   `json:"owner,omitempty" jsonschema:"description=Item owner"`
	Tags        []string `json:"tags,omitempty" jsonschema:"description=Tags for categorization"`
}

// CreateRMIOutput is the output for create_rmi.
type CreateRMIOutput struct {
	File    string           `json:"file"`
	Item    *rmi.RoadmapItem `json:"item,omitempty"`
	Created bool             `json:"created"`
	Error   string           `json:"error,omitempty"`
}

// UpdateRMIInput is the input for update_rmi.
type UpdateRMIInput struct {
	File     string   `json:"file" jsonschema:"description=Path to roadmap JSON file"`
	ID       string   `json:"id" jsonschema:"description=Item ID to update"`
	Name     string   `json:"name,omitempty" jsonschema:"description=New name"`
	MoSCoW   string   `json:"moscow,omitempty" jsonschema:"description=New MoSCoW priority"`
	Status   string   `json:"status,omitempty" jsonschema:"description=New status"`
	Quarter  string   `json:"quarter,omitempty" jsonschema:"description=New quarter"`
	Owner    string   `json:"owner,omitempty" jsonschema:"description=New owner"`
	Progress int      `json:"progress,omitempty" jsonschema:"description=Progress percentage (0-100)"`
	Tags     []string `json:"tags,omitempty" jsonschema:"description=New tags (replaces existing)"`
}

// UpdateRMIOutput is the output for update_rmi.
type UpdateRMIOutput struct {
	File    string           `json:"file"`
	Item    *rmi.RoadmapItem `json:"item,omitempty"`
	Updated bool             `json:"updated"`
	Error   string           `json:"error,omitempty"`
}

// DeleteRMIInput is the input for delete_rmi.
type DeleteRMIInput struct {
	File string `json:"file" jsonschema:"description=Path to roadmap JSON file"`
	ID   string `json:"id" jsonschema:"description=Item ID to delete"`
}

// DeleteRMIOutput is the output for delete_rmi.
type DeleteRMIOutput struct {
	File    string `json:"file"`
	ID      string `json:"id"`
	Deleted bool   `json:"deleted"`
	Error   string `json:"error,omitempty"`
}

// RMISummaryInput is the input for rmi_summary.
type RMISummaryInput struct {
	File string `json:"file" jsonschema:"description=Path to roadmap JSON file"`
}

// RMISummaryOutput is the output for rmi_summary.
type RMISummaryOutput struct {
	File    string       `json:"file"`
	Summary *rmi.Summary `json:"summary,omitempty"`
	Error   string       `json:"error,omitempty"`
}

// TopRMIsInput is the input for top_rmis.
type TopRMIsInput struct {
	File   string `json:"file" jsonschema:"description=Path to roadmap JSON file"`
	SortBy string `json:"sort_by,omitempty" jsonschema:"description=Sort by (priority, rice, market_signal). Default: priority"`
	Limit  int    `json:"limit,omitempty" jsonschema:"description=Number of items to return (default 10)"`
}

// TopRMIsOutput is the output for top_rmis.
type TopRMIsOutput struct {
	File   string            `json:"file"`
	Items  []rmi.RoadmapItem `json:"items"`
	Total  int               `json:"total"`
	SortBy string            `json:"sort_by"`
	Error  string            `json:"error,omitempty"`
}

// RegisterRMITools registers RMI-related MCP tools.
func (s *Server) RegisterRMITools(server *mcp.Server) {
	mcp.AddTool(server, &mcp.Tool{
		Name:        "list_rmis",
		Description: "List roadmap items with optional filtering by MoSCoW priority, status, or quarter.",
	}, s.listRMIs)

	mcp.AddTool(server, &mcp.Tool{
		Name:        "get_rmi",
		Description: "Get details for a specific roadmap item including priority scores and effort estimates.",
	}, s.getRMI)

	mcp.AddTool(server, &mcp.Tool{
		Name:        "create_rmi",
		Description: "Create a new roadmap item with MoSCoW prioritization.",
	}, s.createRMI)

	mcp.AddTool(server, &mcp.Tool{
		Name:        "update_rmi",
		Description: "Update an existing roadmap item's status, priority, or other fields.",
	}, s.updateRMI)

	mcp.AddTool(server, &mcp.Tool{
		Name:        "delete_rmi",
		Description: "Delete a roadmap item by ID.",
	}, s.deleteRMI)

	mcp.AddTool(server, &mcp.Tool{
		Name:        "rmi_summary",
		Description: "Get aggregated statistics about all roadmap items including counts by MoSCoW and status.",
	}, s.rmiSummary)

	mcp.AddTool(server, &mcp.Tool{
		Name:        "top_rmis",
		Description: "Get top roadmap items sorted by priority score, RICE score, or market signal.",
	}, s.topRMIs)
}

// Tool implementations

func (s *Server) listRMIs(ctx context.Context, req *mcp.CallToolRequest, input ListRMIsInput) (*mcp.CallToolResult, ListRMIsOutput, error) {
	svc, err := s.getService(input.File)
	if err != nil {
		return nil, ListRMIsOutput{
			File:  filepath.Base(input.File),
			Error: err.Error(),
		}, nil
	}

	if input.Limit == 0 {
		input.Limit = 50
	}

	filter := rmi.ListFilter{
		MoSCoW:  prioritization.MoSCoWPriority(input.MoSCoW),
		Status:  rmi.RMIStatus(input.Status),
		Quarter: input.Quarter,
		Limit:   input.Limit,
		SortBy:  input.SortBy,
	}

	items := svc.List(filter)
	summary := svc.Summary()

	return nil, ListRMIsOutput{
		File:    filepath.Base(input.File),
		Items:   items,
		Total:   len(items),
		Summary: summary,
	}, nil
}

func (s *Server) getRMI(ctx context.Context, req *mcp.CallToolRequest, input GetRMIInput) (*mcp.CallToolResult, GetRMIOutput, error) {
	svc, err := s.getService(input.File)
	if err != nil {
		return nil, GetRMIOutput{
			File:  filepath.Base(input.File),
			Error: err.Error(),
		}, nil
	}

	item, err := svc.Get(input.ID)
	if err != nil {
		return nil, GetRMIOutput{
			File:  filepath.Base(input.File),
			Error: err.Error(),
		}, nil
	}

	return nil, GetRMIOutput{
		File: filepath.Base(input.File),
		Item: item,
	}, nil
}

func (s *Server) createRMI(ctx context.Context, req *mcp.CallToolRequest, input CreateRMIInput) (*mcp.CallToolResult, CreateRMIOutput, error) {
	svc := s.getOrCreateService(input.File)

	// MoSCoW is optional — empty means not yet prioritized. Validate only
	// when supplied (mirroring updateRMI).
	moscow := prioritization.MoSCoWPriority(input.MoSCoW)
	if moscow != prioritization.MoSCoWUnspecified && !prioritization.IsValidMoSCoWPriority(moscow) {
		return nil, CreateRMIOutput{
			File:  filepath.Base(input.File),
			Error: "invalid moscow priority: " + input.MoSCoW,
		}, nil
	}

	createInput := rmi.CreateInput{
		ID:          input.ID,
		Name:        input.Name,
		Description: input.Description,
		MoSCoW:      moscow,
		Quarter:     input.Quarter,
		Owner:       input.Owner,
		Tags:        input.Tags,
	}

	item, err := svc.Create(createInput)
	if err != nil {
		return nil, CreateRMIOutput{
			File:  filepath.Base(input.File),
			Error: err.Error(),
		}, nil
	}

	if err := svc.SaveAs(input.File); err != nil {
		return nil, CreateRMIOutput{
			File:  filepath.Base(input.File),
			Error: err.Error(),
		}, nil
	}

	return nil, CreateRMIOutput{
		File:    filepath.Base(input.File),
		Item:    item,
		Created: true,
	}, nil
}

func (s *Server) updateRMI(ctx context.Context, req *mcp.CallToolRequest, input UpdateRMIInput) (*mcp.CallToolResult, UpdateRMIOutput, error) {
	svc, err := s.getService(input.File)
	if err != nil {
		return nil, UpdateRMIOutput{
			File:  filepath.Base(input.File),
			Error: err.Error(),
		}, nil
	}

	updateInput := rmi.UpdateInput{}

	if input.Name != "" {
		updateInput.Name = &input.Name
	}
	if input.MoSCoW != "" {
		moscow := prioritization.MoSCoWPriority(input.MoSCoW)
		updateInput.MoSCoW = &moscow
	}
	if input.Status != "" {
		status := rmi.RMIStatus(input.Status)
		updateInput.Status = &status
	}
	if input.Quarter != "" {
		updateInput.Quarter = &input.Quarter
	}
	if input.Owner != "" {
		updateInput.Owner = &input.Owner
	}
	if input.Progress > 0 {
		updateInput.Progress = &input.Progress
	}
	if input.Tags != nil {
		updateInput.Tags = input.Tags
	}

	item, updated, err := svc.Update(input.ID, updateInput)
	if err != nil {
		return nil, UpdateRMIOutput{
			File:  filepath.Base(input.File),
			Error: err.Error(),
		}, nil
	}

	if updated {
		if err := svc.Save(); err != nil {
			return nil, UpdateRMIOutput{
				File:  filepath.Base(input.File),
				Error: err.Error(),
			}, nil
		}
	}

	return nil, UpdateRMIOutput{
		File:    filepath.Base(input.File),
		Item:    item,
		Updated: updated,
	}, nil
}

func (s *Server) deleteRMI(ctx context.Context, req *mcp.CallToolRequest, input DeleteRMIInput) (*mcp.CallToolResult, DeleteRMIOutput, error) {
	svc, err := s.getService(input.File)
	if err != nil {
		return nil, DeleteRMIOutput{
			File:  filepath.Base(input.File),
			ID:    input.ID,
			Error: err.Error(),
		}, nil
	}

	if err := svc.Delete(input.ID); err != nil {
		return nil, DeleteRMIOutput{
			File:  filepath.Base(input.File),
			ID:    input.ID,
			Error: err.Error(),
		}, nil
	}

	if err := svc.Save(); err != nil {
		return nil, DeleteRMIOutput{
			File:  filepath.Base(input.File),
			ID:    input.ID,
			Error: err.Error(),
		}, nil
	}

	return nil, DeleteRMIOutput{
		File:    filepath.Base(input.File),
		ID:      input.ID,
		Deleted: true,
	}, nil
}

func (s *Server) rmiSummary(ctx context.Context, req *mcp.CallToolRequest, input RMISummaryInput) (*mcp.CallToolResult, RMISummaryOutput, error) {
	svc, err := s.getService(input.File)
	if err != nil {
		return nil, RMISummaryOutput{
			File:  filepath.Base(input.File),
			Error: err.Error(),
		}, nil
	}

	summary := svc.Summary()

	return nil, RMISummaryOutput{
		File:    filepath.Base(input.File),
		Summary: summary,
	}, nil
}

func (s *Server) topRMIs(ctx context.Context, req *mcp.CallToolRequest, input TopRMIsInput) (*mcp.CallToolResult, TopRMIsOutput, error) {
	svc, err := s.getService(input.File)
	if err != nil {
		return nil, TopRMIsOutput{
			File:  filepath.Base(input.File),
			Error: err.Error(),
		}, nil
	}

	if input.Limit == 0 {
		input.Limit = 10
	}
	if input.SortBy == "" {
		input.SortBy = "priority"
	}

	var items []rmi.RoadmapItem
	switch input.SortBy {
	case "rice":
		items = svc.TopByRICE(input.Limit)
	case "market_signal":
		items = svc.TopByMarketSignal(input.Limit)
	default:
		items = svc.TopByPriority(input.Limit)
	}

	return nil, TopRMIsOutput{
		File:   filepath.Base(input.File),
		Items:  items,
		Total:  len(items),
		SortBy: input.SortBy,
	}, nil
}

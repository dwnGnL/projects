package models

type ProjectResp struct {
	ProjectID       int64            `json:"project_id"`
	Title           string           `json:"title"`
	Description     string           `json:"description"`
	MediaID         int64            `json:"media_id"`
	Type            string           `json:"type"`
	BusinessOwner   string           `json:"business_owner"`
	LegacyEntity    string           `json:"legacy_entity"`
	Cluster         string           `json:"cluster"`
	Stage           string           `json:"stage"`
	Phase           string           `json:"phase"`
	OwnerID         string           `json:"owner_id"`
	OwnerPhoto      string           `json:"owner_photo"`
	Hidden          int64            `json:"hidden"`
	Created         int64            `json:"created"`
	Category        string           `json:"category"`
	Template        *ProjectTemplate `json:"template"`
	Region          string           `json:"region"`
	Status          string           `json:"status"`
	Priority        int              `json:"priority"`
	PipelineManager string           `json:"pipeline_manager"`
	ProjectManager  string           `json:"project_manager"`
}

type ProjectTemplate struct {
	Stage []Stage `json:"stage"`
}

type Stage struct {
	StageID      int64  `json:"stage_id"`
	Order        int    `json:"order"`
	Title        string `json:"title"`
	DateStart    string `json:"date_start"`
	Description  string `json:"description"`
	DateEnd      string `json:"date_end"`
	ProjectID    int64  `json:"project_id"`
	WorkspaceID  int64  `json:"workspace_id"`
	ActionPlanID int64  `json:"action_plan_id"`
	Hidden       bool   `json:"hidden"`

	Milestone []Milestone `json:"milestone"`
}

type Milestones struct {
	Milestones []Milestone `json:"milestones"`
}

type Milestone struct {
	MilestoneID int64          `json:"milestone_id"`
	ProjectID   int64          `json:"project_id"`
	StageID     int64          `json:"stage_id"`
	Order       int            `json:"order"`
	Status      string         `json:"status"`
	DateStart   string         `json:"date_start"`
	Description string         `json:"description"`
	DateEnd     string         `json:"date_end"`
	Title       string         `json:"title"`
	AssignID    string         `json:"assign_id"`
	Epic        []EpicResponse `json:"epic,omitempty"`
	Task        []Task         `json:"task,omitempty"`
}

type ActionPlanResp struct {
	ActionPlanID int64   `json:"action_plan_id"`
	Stage        []Stage `json:"stage,omitempty"`
	ProjectID    int64   `json:"project_id"`
	Title        string  `json:"title"`
	Created      int64   `json:"created"`
	Status       string  `json:"status"`
}

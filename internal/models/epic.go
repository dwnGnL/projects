package models

type EpicRequest struct {
	ID          int64  `json:"id,omitempty"`
	WorkspaceID int64  `json:"workspace_id,omitempty"`
	ProjectID   int64  `json:"project_id,omitempty"`
	StageID     int64  `json:"stage_id,omitempty"`
	MilestoneID int64  `json:"milestone_id,omitempty"`
	Title       string `json:"title,omitempty"`
	Description string `json:"description,omitempty"`
}

type EpicResponse struct {
	ID          int64  `json:"id,omitempty"`
	WorkspaceID int64  `json:"workspace_id,omitempty"`
	ProjectID   int64  `json:"project_id,omitempty"`
	StageID     int64  `json:"stage_id,omitempty"`
	MilestoneID int64  `json:"milestone_id,omitempty"`
	Title       string `json:"title,omitempty"`
	Description string `json:"description,omitempty"`
	Task        []Task `json:"task"`
}

package models

type ProjectReq struct {
	Title           *string `json:"title"`
	Description     *string `json:"description"`
	MediaID         *int64  `json:"media_id"`
	Type            *string `json:"type"`
	BusinessOwner   *string `json:"business_owner"`
	Stage           *string `json:"stage"`
	Phase           *string `json:"phase"`
	LegacyEntity    *string `json:"legacy_entity"`
	Cluster         *string `json:"cluster"`
	OwnerID         *string `json:"owner_id"`
	OwnerPhoto      *string `json:"owner_photo"`
	Region          *string `json:"region"`
	Status          *string `json:"status"`
	Priority        *int    `json:"priority"`
	PipelineManager *string `json:"pipeline_manager"`
	ProjectManager  *string `json:"project_manager"`
}

type MilestoneFilter struct {
	MilestoneID int   `json:"milestone_id"`
	ProjectID   int64 `json:"project_id"`
}

type ProjectFilter struct {
	Cluster *string `json:"cluster"`
	Type    *string `json:"type"`
	Stage   *string `json:"stage"`
}

type ActionPlan struct {
	Title       string `json:"title"`
	WorkspaceID int64  `json:"workspace_id"`
	ProjectID   int64  `json:"project_id"`
	Status      string `json:"status"`
}

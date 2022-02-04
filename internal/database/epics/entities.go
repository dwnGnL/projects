package epics

type EpicEntity struct {
	ID           int64  `gorm:"column:id"`
	WorkspaceID  int64  `gorm:"column:workspace_id"`
	ActionPlanID int64  `gorm:"column:action_plan_id"`
	ProjectID    int64  `gorm:"column:project_id"`
	StageID      int64  `gorm:"stage_id"`
	MilestoneID  int64  `gorm:"milestone_id"`
	Title        string `gorm:"title"`
	Description  string `gorm:"description"`
}

package tasks

import "projects/internal/database/epics"

type TaskEntity struct {
	ID           int64            `gorm:"column:id"`
	MilestoneID  int64            `gorm:"column:milestone_id"`
	EpicID       int64            `gorm:"column:epic_id"`
	Epic         epics.EpicEntity `gorm:"-"`
	ActionPlanID int64            `gorm:"column:action_plan_id"`
}

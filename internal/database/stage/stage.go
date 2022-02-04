package stage

import (
	"gorm.io/gorm"

	"projects/internal/database/actionPlan"
)

type StageInter interface {
	CreateMany(stages []StageEntity) ([]StageEntity, error)
	Update(stages StageEntity) (StageEntity, error)
	GetByProjectID(projectID int64) []StageEntity
	GetByActionPlan(actionPlanID int64) []StageEntity
	DeleteStage(stageID int64) error
}

type StageEntity struct {
	StageID      int64  `gorm:"column:stage_id;primary_key;autoIncrement"`
	ProjectID    int64  `gorm:"column:project_id"`
	WorkspaceID  int64  `gorm:"column:workspace_id"`
	ActionPlanID int64  `gorm:"column:action_plan_id"`
	Description  string `gorm:"column:description"`
	DateStart    string `gorm:"column:date_start"`
	DateStop     string `gorm:"column:date_stop"`
	Hidden       bool   `gorm:"column:hidden;default:false"`
	Order        int    `gorm:"column:order"`
	Title        string `gorm:"column:title"`
}

func (StageEntity) TableName() string {
	return "stage"
}

type stage struct {
	db *gorm.DB
}

func New(dbr *gorm.DB) StageInter {

	return &stage{db: dbr}
}

func (s stage) CreateMany(shedules []StageEntity) ([]StageEntity, error) {
	var newStages []StageEntity
	for _, sh := range shedules {
		var ac actionPlan.ActionPlanEntity
		s.db.Select("workspace_id, project_id").Where("action_plan_id", sh.ActionPlanID).First(&ac)
		sh.WorkspaceID = ac.WorkspaceID
		sh.ProjectID = ac.ProjectID

		newStages = append(newStages, sh)
	}

	if err := s.db.Create(&newStages).Error; err != nil {
		return []StageEntity{}, err
	}
	return newStages, nil
}

func (s stage) Update(shedules StageEntity) (StageEntity, error) {
	s.db.Updates(&shedules)
	return shedules, nil
}
func (s stage) GetByProjectID(projectID int64) []StageEntity {
	var schedules []StageEntity
	s.db.Where("project_id = ? and hidden = false", projectID).Order(`"order"`).Find(&schedules)
	return schedules
}

func (s stage) GetByActionPlan(actionPlanID int64) []StageEntity {
	var stages []StageEntity
	s.db.Where("action_plan_id = ? and hidden = false", actionPlanID).Order(`"order"`).Find(&stages)
	return stages
}

func (s stage) DeleteStage(stageID int64) error {
	return s.db.Where("stage_id = ?", stageID).Delete(StageEntity{}).Error
}

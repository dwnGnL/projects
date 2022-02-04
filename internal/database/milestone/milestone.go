package milestone

import (
	"database/sql/driver"

	"gorm.io/gorm"

	"projects/internal/database/stage"
)

type Status string

const (
	NewStatus  Status = "new"
	Hold       Status = "hold"
	Cancelled  Status = "cancelled"
	InProgress Status = "in progress"
	Completed  Status = "completed"
)

func (s *Status) Scan(value interface{}) error {
	*s = Status(value.(string))
	return nil
}

func (s Status) Value() (driver.Value, error) {
	return string(s), nil
}

func (s *Status) String() string {
	return string(*s)
}

type milestoneInter interface {
	CreateMany(Stages []MilestoneEntity) ([]MilestoneEntity, error)
	Update(shedules MilestoneEntity) (MilestoneEntity, error)
	GetByProjectID(projectID int64) []MilestoneEntity
	GetByStageID(stageID int64) []MilestoneEntity
	GetMilestoneByID(id int64) MilestoneEntity
	GetByActionPlan(actionPlanID int64) []MilestoneEntity
	DeleteByID(milestoneID int64) error
}

type MilestoneEntity struct {
	MilestoneID  int64  `gorm:"column:milestone_id;primary_key;autoIncrement"`
	StageID      int64  `gorm:"column:stage_id"`
	WorkspaceID  int64  `gorm:"column:workspace_id"`
	ActionPlanID int64  `gorm:"column:action_plan_id"`
	ProjectID    int64  `gorm:"column:project_id"`
	Order        int    `gorm:"column:order"`
	Status       Status `gorm:"column:status;type:enum_status;default:'new'"`
	Description  string `gorm:"column:description"`
	DateStart    string `gorm:"column:date_start"`
	DateStop     string `gorm:"column:date_stop"`
	Hidden       bool   `gorm:"column:hidden;default:false"`
	Title        string `gorm:"column:title"`
	AssignID     string `gorm:"column:assign_id"`
}

func (MilestoneEntity) TableName() string {
	return "milestone"
}

type milestone struct {
	db *gorm.DB
}

func New(dbr *gorm.DB) milestoneInter {

	return &milestone{db: dbr}
}

func (s milestone) CreateMany(milestones []MilestoneEntity) ([]MilestoneEntity, error) {
	var newMilestones []MilestoneEntity
	for _, m := range milestones {
		var st stage.StageEntity
		s.db.Select("action_plan_id, workspace_id, project_id").Where("stage_id", m.StageID).First(&st)
		m.ActionPlanID = st.ActionPlanID
		m.WorkspaceID = st.WorkspaceID
		m.ProjectID = st.ProjectID

		newMilestones = append(newMilestones, m)
	}

	if err := s.db.Create(&newMilestones).Error; err != nil {
		return []MilestoneEntity{}, err
	}
	return newMilestones, nil
}

func (s milestone) Update(milestone MilestoneEntity) (MilestoneEntity, error) {
	s.db.Model(&milestone).Updates(&milestone)
	return milestone, nil
}
func (s milestone) GetByProjectID(projectID int64) []MilestoneEntity {
	var milestones []MilestoneEntity
	s.db.Where("project_id = ? and hidden = false", projectID).Order(`"order"`).Find(&milestones)
	return milestones
}

func (s milestone) GetByStageID(stageID int64) []MilestoneEntity {
	var milestones []MilestoneEntity
	s.db.Where("stage_id = ? and hidden = false", stageID).Order(`"order"`).Find(&milestones)
	return milestones
}
func (s milestone) GetByActionPlan(actionPlanID int64) []MilestoneEntity {
	var miles []MilestoneEntity
	s.db.Where("action_plan_id = ? and hidden = false", actionPlanID).Order(`"order"`).Find(&miles)
	return miles
}

func (s milestone) GetMilestoneByID(id int64) MilestoneEntity {
	var ms MilestoneEntity
	s.db.Where("milestone_id = ?", id).First(&ms)

	return ms
}

func (s milestone) DeleteByID(milestoneID int64) error {
	if err := s.db.Model(MilestoneEntity{}).Where("milestone_id = ?", milestoneID).Update("hidden", true).Error; err != nil {
		return err
	}

	return nil
}

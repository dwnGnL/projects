package actionPlan

import (
	"time"

	"database/sql/driver"

	"gorm.io/gorm"

	"projects/internal/database/projects"
)

type ActionStatus string

const (
	Active   ActionStatus = "active"
	Archived ActionStatus = "archived"
	Draft    ActionStatus = "draft"
)

func (as *ActionStatus) Scan(value interface{}) error {
	*as = ActionStatus(value.(string))
	return nil
}

func (as ActionStatus) Value() (driver.Value, error) {
	return string(as), nil
}

func (as *ActionStatus) String() string {
	return string(*as)
}

type ActionPlanInter interface {
	Create(aPE *ActionPlanEntity) error
	Delete(id int64) error
	Update(id int64, title string, status string) (ActionPlanEntity, error)
	Get(id int64) (ActionPlanEntity, error)
	GetByProjectID(projectID int64) ([]ActionPlanEntity, error)
}

type ActionPlanEntity struct {
	ActionPlanID int64        `gorm:"column:action_plan_id;primary_key;autoIncrement"`
	WorkspaceID  int64        `gorm:"column:workspace_id"`
	ProjectID    int64        `gorm:"column:project_id"`
	Title        string       `gorm:"column:title"`
	Created      int64        `gorm:"column:created"`
	Status       ActionStatus `gorm:"column:status;type:enum_ac_status;default:'active'"`

	Phase   projects.PhaseEntity
	PhaseID int64 `gorm:"column:phase_id, omitempty"`
}

func (ActionPlanEntity) TableName() string {
	return "action_plan"
}

func (a *ActionPlanEntity) BeforeCreate(_ *gorm.DB) (err error) {
	a.Created = time.Now().Unix()
	return
}

type actionPlan struct {
	db *gorm.DB
}

func New(dbr *gorm.DB) ActionPlanInter {

	return &actionPlan{db: dbr}
}

func (aP actionPlan) Create(aPE *ActionPlanEntity) error {
	if err := aP.db.Create(aPE).Error; err != nil {
		return err
	}
	return nil
}

func (aP actionPlan) Get(id int64) (ActionPlanEntity, error) {
	var acPlan ActionPlanEntity

	if err := aP.db.Find(&acPlan, id).Error; err != nil {
		return ActionPlanEntity{}, err
	}
	return acPlan, nil
}

func (aP actionPlan) GetByProjectID(projectID int64) ([]ActionPlanEntity, error) {
	var acPlan []ActionPlanEntity

	if err := aP.db.Where("project_id = ?", projectID).Find(&acPlan).Error; err != nil {
		return []ActionPlanEntity{}, err
	}
	return acPlan, nil
}

func (aP actionPlan) Delete(id int64) error {
	if err := aP.db.Delete(ActionPlanEntity{}, id).Error; err != nil {
		return err
	}
	return nil
}

func (aP actionPlan) Update(id int64, title string, status string) (ActionPlanEntity, error) {
	updates := make(map[string]interface{})
	if title != "" {
		updates["title"] = title
	}
	if status != "" {
		updates["status"] = status
	}
	if err := aP.db.Model(&ActionPlanEntity{}).Where("action_plan_id = ?", id).Updates(updates).Error; err != nil {
		return ActionPlanEntity{}, err
	}
	var acPlan ActionPlanEntity
	if err := aP.db.Find(&acPlan, id).Error; err != nil {
		return ActionPlanEntity{}, err
	}
	return acPlan, nil
}

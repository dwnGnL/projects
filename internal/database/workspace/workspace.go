package workspace

import (
	"database/sql/driver"

	"gorm.io/gorm"
)

type WSType string

const (
	Development WSType = "development"
	Marketing   WSType = "marketing"
	Legal       WSType = "legal"
	Support     WSType = "support"
	BackOffice  WSType = "back office"
)

func (c *WSType) Scan(value interface{}) error {
	*c = WSType(value.(string))
	return nil
}

func (c WSType) Value() (driver.Value, error) {
	return string(c), nil
}

type WorkspaceInter interface {
	Create(wE *WorkspaceEntity) error
}

type WorkspaceEntity struct {
	WorkspaceID int64  `gorm:"column:workspace_id;primary_key;autoIncrement"`
	ProjectID   int64  `gorm:"column:project_id"`
	Type        WSType `gorm:"type:ws_type"`
	Title       string `gorm:"column:title"`
}

func (WorkspaceEntity) TableName() string {
	return "workspace"
}

type workspace struct {
	db *gorm.DB
}

func New(dbr *gorm.DB) WorkspaceInter {

	return &workspace{db: dbr}
}

func (w workspace) Create(wE *WorkspaceEntity) error {
	if err := w.db.Create(wE).Error; err != nil {
		return err
	}
	return nil
}

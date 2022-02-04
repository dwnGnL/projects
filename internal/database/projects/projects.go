package projects

import (
	"time"

	"database/sql/driver"

	"gorm.io/gorm"
)

type Phase string

func (p *Phase) Scan(value interface{}) error {
	*p = Phase(value.(string))
	return nil
}

func (p Phase) Value() (driver.Value, error) {
	return string(p), nil
}

func (p *Phase) String() string {
	return string(*p)
}

const (
	BuildingLaunch Phase = "building&launch"
	ScaleGrowth    Phase = "scale&growth"
)

type Stage string

const (
	// B&L Phase
	Ideation    Stage = "ideation"
	Concept     Stage = "concept"
	Business    Stage = "business case"
	Development Stage = "development"
	Pilot       Stage = "pilot"
	CL          Stage = "commercial launch"

	// S&L Phase
	EG        Stage = "early growths"
	Expansion Stage = "expansion"
	ShakeOut  Stage = "shake-out"
	Mature    Stage = "mature"
	Ecosystem Stage = "ecosystem"
)

func (s *Stage) Scan(value interface{}) error {
	*s = Stage(value.(string))
	return nil
}

func (s Stage) Value() (driver.Value, error) {
	return string(s), nil
}

func (s *Stage) String() string {
	return string(*s)
}

type Type string

const (
	VentureBuilding Type = "product"
	ServiceDevt     Type = "super project"
	ServiceVB       Type = "venture building"
	Social          Type = "business"
	Internal        Type = "platform"
)

func (t *Type) Scan(value interface{}) error {
	*t = Type(value.(string))
	return nil
}

func (t Type) Value() (driver.Value, error) {
	return string(t), nil
}

func (t *Type) String() string {
	return string(*t)
}

type Category string

const (
	Project    Category = "project"
	Product    Category = "product"
	Platform   Category = "platform"
	BackOffice Category = "back office"
)

func (c *Category) Scan(value interface{}) error {
	if v, ok := value.([]uint8); ok {
		*c = Category(v)
	} else if v, ok := value.(string); ok {
		*c = Category(v)
	}
	return nil
}

func (c Category) Value() (driver.Value, error) {
	return string(c), nil
}

func (c *Category) String() string {
	return string(*c)
}

type ProjectEntity struct {
	ProjectID       int64    `gorm:"column:project_id;primary_key;autoIncrement"`
	Title           string   `gorm:"column:title"`
	Description     string   `gorm:"column:description"`
	MediaID         int64    `gorm:"column:media_id"`
	Type            Type     `gorm:"type:enum_type;default:'product'"`
	BusinessOwner   string   `gorm:"column:business_owner"`
	LegacyEntity    string   `gorm:"column:legacy_entity"`
	Cluster         string   `gorm:"column:cluster"`
	Stage           Stage    `gorm:"type:enum_stages;default:'ideation'"`
	Phase           Phase    `gorm:"type:enum_phases;default:'building&launch'"`
	OwnerID         string   `gorm:"column:owner_id"`
	Hidden          int64    `gorm:"column:hidden"`
	Category        Category `gorm:"type:category;default:'project'"`
	Created         int64    `gorm:"column:created"`
	Region          string   `gorm:"column:region"`
	Status          string   `gorm:"column:status"`
	Priority        int      `gorm:"column:priority"`
	PipelineManager string   `gorm:"column:pipeline_manager"`
	ProjectManager  string   `gorm:"column:project_manager"`
}

func (ProjectEntity) TableName() string {
	return "projects"
}

func (p *ProjectEntity) BeforeCreate(tx *gorm.DB) (err error) {
	p.Created = time.Now().Unix()
	return
}

type ProjectsInter interface {
	GetAll(filter ProjectFilter) ([]ProjectEntity, error)
	Get(id int64) (ProjectEntity, error)
	Create(ProjectEntity) (ProjectEntity, error)
	Update(id int64, updateColumns map[string]interface{}) (ProjectEntity, error)
	Delete(id int64) error
}

type projects struct {
	db *gorm.DB
}

type ProjectFilter struct {
	Cluster *string
	Type    *string
	Stage   *string
}

func prepareQuery(filter ProjectFilter, dbr *gorm.DB) *gorm.DB {
	query := dbr.Model(ProjectEntity{})
	if filter.Cluster != nil {
		query = query.Where("cluster = ?", *filter.Cluster)
	}
	if filter.Type != nil {
		query = query.Where("type = ?", *filter.Type)
	}
	if filter.Stage != nil {
		query = query.Where("stage = ?", *filter.Stage)
	}
	return query
}

func New(dbr *gorm.DB) ProjectsInter {

	return &projects{db: dbr}
}

func (p projects) GetAll(filter ProjectFilter) ([]ProjectEntity, error) {
	var projectsEntity []ProjectEntity
	if err := prepareQuery(filter, p.db).Where("hidden = ?", 0).Order("priority").Find(&projectsEntity).Error; err != nil {
		return nil, err
	}
	return projectsEntity, nil
}

func (p projects) Get(id int64) (ProjectEntity, error) {
	var project ProjectEntity
	if err := p.db.Where("hidden = ?", 0).Find(&project, id).Error; err != nil {
		return ProjectEntity{}, err
	}
	return project, nil
}

func (p projects) Create(pe ProjectEntity) (ProjectEntity, error) {
	if err := p.db.Create(&pe).Error; err != nil {
		return ProjectEntity{}, err
	}
	return pe, nil
}

func (p projects) Update(id int64, updateColumns map[string]interface{}) (ProjectEntity, error) {

	if err := p.db.Model(ProjectEntity{}).Where("project_id = ?", id).Updates(updateColumns).Error; err != nil {
		return ProjectEntity{}, err
	}
	var project ProjectEntity
	if err := p.db.Find(&project, id).Error; err != nil {
		return ProjectEntity{}, err
	}
	return project, nil

}

func (p projects) Delete(id int64) error {
	if err := p.db.Model(ProjectEntity{}).Where("project_id = ?", id).Update("hidden", 1).Error; err != nil {
		return err
	}

	return nil
}

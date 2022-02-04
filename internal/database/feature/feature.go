package feature

import "gorm.io/gorm"

type featureInter interface {
}

type FeatureEntity struct {
	FeatureID   int64  `gorm:"column:feature_id;primary_key"`
	MilestoneID int64  `gorm:"column:milestone_id"`
	Title       string `gorm:"title"`
	Created     int64  `gorm:"column:created"`
}

func (FeatureEntity) TableName() string {
	return "feature"
}

type feature struct {
	db *gorm.DB
}

func New(dbr *gorm.DB) featureInter {

	return &feature{db: dbr}
}

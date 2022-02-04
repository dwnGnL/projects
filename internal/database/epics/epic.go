package epics

import (
	"gorm.io/gorm"
)

type Epic interface {
	CreateEpic(epicEntity EpicEntity) (EpicEntity, error)
	GetEpic(epicEntity EpicEntity) ([]EpicEntity, error)
	UpdateEpic(entity EpicEntity) error
	DeleteEpic(id int64) error
}

func NewEpic(db *gorm.DB) Epic {
	return &epic{db: db}
}

type epic struct {
	db *gorm.DB
}

func (e *epic) CreateEpic(epicEntity EpicEntity) (EpicEntity, error) {
	if err := e.db.Create(&epicEntity).Error; err != nil {
		return EpicEntity{}, err
	}
	return epicEntity, nil
}

func (e *epic) GetEpic(epicEntity EpicEntity) ([]EpicEntity, error) {
	var epicEntities []EpicEntity
	if err := e.db.Find(&epicEntities, epicEntity).Error; err != nil {
		return nil, err
	}
	return epicEntities, nil
}

func (e *epic) UpdateEpic(entity EpicEntity) error {
	return e.db.Updates(entity).Error
}

func (e *epic) DeleteEpic(id int64) error {
	return e.db.Delete(EpicEntity{}, id).Error
}

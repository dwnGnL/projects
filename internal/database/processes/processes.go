package processes

import (
	"gorm.io/gorm"
)

type ProcessEntity struct {
	ProcessID int64  `gorm:"column:process_id;primary_key;autoIncrement"`
	Name      string `gorm:"column:name"`
}

type ProcessInter interface {
	Create(proc *ProcessEntity) error
	GetAll() ([]ProcessEntity, error)
	Update(proc *ProcessEntity) error
	Delete(id int64) error
}

type processes struct {
	db *gorm.DB
}

func New(db *gorm.DB) ProcessInter {
	return &processes{db: db}
}

func (p *processes) Create(proc *ProcessEntity) error {
	return p.db.Create(proc).Error
}

func (p *processes) GetAll() ([]ProcessEntity, error) {
	var proc []ProcessEntity
	if err := p.db.Find(&proc).Error; err != nil {
		return nil, err
	}

	return proc, nil
}

func (p *processes) Update(proc *ProcessEntity) error {
	return p.db.Updates(proc).Error
}

func (p *processes) Delete(id int64) error {
	tx := p.db.Begin()
	defer func() {
		if r := recover(); r != nil {
			tx.Rollback()
		}
	}()
	if tx.Error != nil {
		return tx.Error
	}

	if err := tx.Table("milestone").
		Where("process_id = ?", id).
		Updates(map[string]interface{}{"process_id": nil}).Error; err != nil {
		tx.Rollback()
		return err
	}
	if err := tx.Where("process_id = ?", id).Delete(ProcessEntity{}).Error; err != nil {
		tx.Rollback()
		return err
	}

	return tx.Commit().Error
}

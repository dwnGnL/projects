package tasks

import "gorm.io/gorm"

type Task interface {
	CreateTask(taskEntity TaskEntity) (*TaskEntity, error)
	GetTaskByEpicID(epicID int64) ([]TaskEntity, error)
	GetTaskByMilestoneID(milestoneID int64) ([]TaskEntity, error)
	GetTaskByID(id int64) (TaskEntity, error)
	GetTaskByActionPlanID(id int64) ([]TaskEntity, error)
	UpdateTask(entity *TaskEntity) error
	DeleteTask(id int64) error
}

func New(db *gorm.DB) Task {
	return &task{db: db}
}

type task struct {
	db *gorm.DB
}

func (t *task) CreateTask(taskEntity TaskEntity) (*TaskEntity, error) {
	if err := t.db.Create(&taskEntity).Error; err != nil {
		return nil, err
	}

	return &taskEntity, nil
}

func (t *task) GetTaskByActionPlanID(id int64) ([]TaskEntity, error) {
	var taskEntity []TaskEntity
	if err := t.db.Where("task_entities.action_plan_id = ?", id).Find(&taskEntity).Error; err != nil {
		return nil, err
	}

	return taskEntity, nil
}

func (t *task) GetTaskByEpicID(epicID int64) ([]TaskEntity, error) {
	var taskEntity []TaskEntity
	if err := t.db.Where("epic_id = ? AND hidden = ?", epicID, 0).Find(&taskEntity).Error; err != nil {
		return nil, err
	}

	return taskEntity, nil
}

func (t *task) GetTaskByMilestoneID(milestoneID int64) ([]TaskEntity, error) {
	var taskEntity []TaskEntity
	if err := t.db.Where("milestone_id = ? AND hidden = ?", milestoneID, 0).Find(&taskEntity).Error; err != nil {
		return nil, err
	}

	return taskEntity, nil
}

func (t *task) GetTaskByID(id int64) (TaskEntity, error) {
	var taskEntity TaskEntity
	if err := t.db.Where("id = ?", id).First(&taskEntity).Error; err != nil {
		return taskEntity, err
	}

	return taskEntity, nil
}

func (t *task) UpdateTask(entity *TaskEntity) error {
	query := make(map[string]interface{})
	if entity.EpicID > 0 {
		query["epic_id"] = entity.EpicID
	}
	if entity.MilestoneID > 0 {
		query["milestone_id"] = entity.MilestoneID
	}
	if entity.ActionPlanID > 0 {
		query["action_plan_id"] = entity.ActionPlanID
	}

	return t.db.Model(TaskEntity{}).Updates(query).Error
}

func (t *task) DeleteTask(id int64) error {
	return t.db.Where("id = ?", id).Delete(TaskEntity{}).Error
}

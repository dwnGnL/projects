package task

import (
	"bytes"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"io/ioutil"
	"net/http"
	"projects/internal/database/milestone"
	"projects/internal/database/tasks"
	"projects/internal/models"
	"projects/pkg/config"
	"projects/pkg/db"
	"strconv"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/sirupsen/logrus"
	"go.uber.org/fx"
)

var Module = fx.Provide(NewTaskHandler)

type TaskHandler interface {
	GetTaskByMilestone(c *gin.Context)
	GetTaskByID(c *gin.Context)
	GetTaskByEpic(c *gin.Context)
	DeleteTask(c *gin.Context)
	UpdateTask(c *gin.Context)
	CreateTask(c *gin.Context)
	GetTask(c *gin.Context)
	GetTasksBatch(c *gin.Context)
}

type Params struct {
	fx.In
	db.DbInter
	*config.Tuner
	*logrus.Logger
}

type taskHandler struct {
	db   db.DbInter
	log  *logrus.Logger
	conf *config.Tuner
}

func NewTaskHandler(params Params) TaskHandler {
	return &taskHandler{db: params.DbInter, log: params.Logger, conf: params.Tuner}
}

func (p taskHandler) GetTasksBatch(c *gin.Context) {
	var resp []models.Task

	headers := map[string]string{"Authorization": c.GetHeader("Authorization")}

	if err := p.sendRequest(http.MethodPost, "/tasks/batch", nil, &resp, &headers); err != nil {
		p.log.Warn("sendRequest err", err)
		c.JSON(http.StatusBadGateway, err.Error())
		return
	}
	c.JSON(http.StatusOK, nil)
}

func (p taskHandler) GetTask(c *gin.Context) {
	var resp models.Task

	headers := map[string]string{"Authorization": c.GetHeader("Authorization")}

	if err := p.sendRequest(http.MethodGet, "/tasks/"+c.Param("task_id"), nil, &resp, &headers); err != nil {
		p.log.Warn("sendRequest err", err)
		c.JSON(http.StatusBadGateway, err.Error())
		return
	}
	c.JSON(http.StatusOK, resp)

}

func (p taskHandler) CreateTask(c *gin.Context) {
	var task models.TaskReq
	if err := c.ShouldBindJSON(&task); err != nil {
		p.log.Warn("bind err", err)
		c.JSON(http.StatusBadGateway, err.Error())
		return
	}
	headers := map[string]string{"Authorization": c.GetHeader("Authorization")}
	jsonB, err := json.Marshal(task)
	if err != nil {
		p.log.Warn("marshal err", err)
		c.JSON(http.StatusBadGateway, err.Error())
		return
	}
	var createdTask models.Task
	if err := p.sendRequest(http.MethodPost, "/tasks", bytes.NewBuffer(jsonB), &createdTask, &headers); err != nil {
		p.log.Warn("sendRequest err", err)
		c.JSON(http.StatusBadGateway, err.Error())
		return
	}
	if createdTask.ID == 0 {
		p.log.Warn("create Task err")
		c.JSON(http.StatusBadGateway, "create Task err")
		return
	}

	repo := tasks.New(p.db.GetDB())
	ml := milestone.New(p.db.GetDB())
	milestone := ml.GetMilestoneByID(task.MilestoneID)
	taskEntity, err := repo.CreateTask(tasks.TaskEntity{
		ID:           createdTask.ID,
		MilestoneID:  task.MilestoneID,
		EpicID:       task.EpicID,
		ActionPlanID: milestone.ActionPlanID,
	})
	createdTask.MilestoneId = task.MilestoneID
	createdTask.EpicID = task.EpicID
	if err != nil || taskEntity.ID <= 0 {
		p.log.Warn("can't create task", err)
		c.JSON(http.StatusBadGateway, "can't create task")
		return
	}

	c.JSON(http.StatusOK, createdTask)

}

func (p taskHandler) UpdateTask(c *gin.Context) {
	iDString := c.Param("id")
	ID, err := strconv.ParseInt(iDString, 10, 64)
	if err != nil {
		p.log.Warn("Parse ID err", err)
		c.JSON(http.StatusBadGateway, err.Error())
		return
	}
	var task models.TaskReq
	if err := c.ShouldBindJSON(&task); err != nil {
		p.log.Warn("bind err", err)
		c.JSON(http.StatusBadGateway, err.Error())
		return
	}
	headers := map[string]string{"Authorization": c.GetHeader("Authorization")}
	jsonB, err := json.Marshal(task)
	if err != nil {
		p.log.Warn("marshal err", err)
		c.JSON(http.StatusBadGateway, err.Error())
		return
	}
	var updatedTask models.Task
	if err := p.sendRequest(http.MethodPut, "/tasks/"+iDString, bytes.NewBuffer(jsonB), &updatedTask, &headers); err != nil {
		p.log.Warn("sendRequest err", err)
		c.JSON(http.StatusBadGateway, err.Error())
		return
	}

	repo := tasks.New(p.db.GetDB())
	if err := repo.UpdateTask(&tasks.TaskEntity{
		ID:           ID,
		MilestoneID:  updatedTask.MilestoneId,
		EpicID:       updatedTask.EpicID,
		ActionPlanID: updatedTask.ActionPlanID,
	}); err != nil {
		p.log.Warn("update task err", err)
		c.JSON(http.StatusBadGateway, err.Error())
		return
	}

	c.JSON(http.StatusOK, updatedTask)
}

func (p taskHandler) DeleteTask(c *gin.Context) {
	iDString := c.Param("id")
	id, err := strconv.Atoi(iDString)
	if err != nil {
		p.log.Warn("wrong id param", err)
		c.JSON(http.StatusBadRequest, "wrong id")
		return
	}
	if id <= 0 {
		c.JSON(http.StatusBadRequest, "wrong id")
		return
	}

	headers := map[string]string{"Authorization": c.GetHeader("Authorization")}
	if err := p.sendRequest(http.MethodDelete, "/tasks/"+fmt.Sprint(id), nil, nil, &headers); err != nil {
		p.log.Warn("sendRequest err", err)
		c.JSON(http.StatusBadGateway, err.Error())
		return
	}
	repo := tasks.New(p.db.GetDB())
	if err := repo.DeleteTask(int64(id)); err != nil {
		p.log.Warn("delete task err", err)
		c.JSON(http.StatusBadGateway, err.Error())
		return
	}

	c.JSON(http.StatusOK, gin.H{"success": "ok"})
}

func (p taskHandler) GetTaskByEpic(c *gin.Context) {
	epicIDString := c.Param("epic_id")
	epicID, err := strconv.Atoi(epicIDString)
	if err != nil {
		p.log.Warn("wrong epic id param", err)
		c.JSON(http.StatusBadRequest, "wrong epic id")
		return
	}
	if epicID <= 0 {
		c.JSON(http.StatusBadRequest, "wrong epic id")
		return
	}

	repo := tasks.New(p.db.GetDB())
	taskEntities, err := repo.GetTaskByEpicID(int64(epicID))
	if err != nil {
		p.log.Warn("can't get task by epic id", err)
		c.JSON(http.StatusBadGateway, "can' get task")
		return
	}
	var response []models.Task
	var tasksId []int64
	for _, v := range taskEntities {
		tasksId = append(tasksId, v.ID)
	}
	jsonB, err := json.Marshal(tasksId)
	if err != nil {
		p.log.Warn("marshal err", err)
		c.JSON(http.StatusBadGateway, err.Error())
		return
	}
	headers := map[string]string{"Authorization": c.GetHeader("Authorization")}
	if err := p.sendRequest(http.MethodPost, "/tasks/batch", bytes.NewBuffer([]byte(jsonB)), &response, &headers); err != nil {
		p.log.Warn("sendRequest err", err)
		c.JSON(http.StatusBadGateway, err.Error())
		return
	}
	var resp []models.Task
	for _, v := range response {
		for _, v2 := range taskEntities {
			if v.ID == v2.ID {
				resp = append(resp, models.Task{
					MilestoneId:        v2.MilestoneID,
					EpicID:             int64(epicID),
					AssigneeId:         v.AssigneeId,
					BrandId:            v.BrandId,
					BucketID:           v.BucketID,
					CreatorID:          v.CreatorID,
					EndTime:            v.EndTime,
					ID:                 v.ID,
					NumberOfAttachment: v.NumberOfAttachment,
					NumberOfComments:   v.NumberOfComments,
					Priority:           v.Priority,
					ResolvedTime:       v.ResolvedTime,
					StartTime:          v.StartTime,
					Status:             v.Status,
					StatusID:           v.StatusID,
					Title:              v.Title,
				})
			}

		}
	}
	c.JSON(http.StatusOK, resp)
}

func (p taskHandler) GetTaskByID(c *gin.Context) {
	iDString := c.Param("id")
	id, err := strconv.Atoi(iDString)
	if err != nil {
		p.log.Warn("wrong id param", err)
		c.JSON(http.StatusBadRequest, "wrong id")
		return
	}
	if id <= 0 {
		c.JSON(http.StatusBadRequest, "wrong id")
		return
	}

	var response []models.Task
	jsonB, _ := json.Marshal([]int{id})
	headers := map[string]string{"Authorization": c.GetHeader("Authorization")}
	if err := p.sendRequest(http.MethodPost, "/tasks/batch", bytes.NewBuffer([]byte(jsonB)), &response, &headers); err != nil {
		p.log.Warn("sendRequest err", err)
		c.JSON(http.StatusBadGateway, err.Error())
		return
	}

	taskEntity, err := tasks.New(p.db.GetDB()).GetTaskByID(int64(id))
	if err != nil {
		p.log.Warn("can't get task by id", err)
		c.JSON(http.StatusBadGateway, "can' get task")
		return
	}
	if taskEntity.ID > 0 && len(response) > 0 {
		response[0].EpicID = taskEntity.EpicID
		response[0].MilestoneId = taskEntity.MilestoneID
	}

	c.JSON(http.StatusOK, response)
}

func (p taskHandler) GetTaskByMilestone(c *gin.Context) {
	milestoneIDString := c.Param("milestone_id")
	milestoneID, err := strconv.Atoi(milestoneIDString)
	if err != nil {
		p.log.Warn("wrong milestone id param", err)
		c.JSON(http.StatusBadRequest, "wrong milestone id")
		return
	}
	if milestoneID <= 0 {
		c.JSON(http.StatusBadRequest, "wrong milestone id")
		return
	}

	repo := tasks.New(p.db.GetDB())
	taskEntities, err := repo.GetTaskByMilestoneID(int64(milestoneID))
	if err != nil {
		p.log.Warn("can't get task by milestone id", err)
		c.JSON(http.StatusBadGateway, "can' get task")
		return
	}
	var response []models.Task
	var tasksId []int64
	for _, v := range taskEntities {
		tasksId = append(tasksId, v.ID)
	}
	jsonB, err := json.Marshal(tasksId)
	if err != nil {
		p.log.Warn("marshal err", err)
		c.JSON(http.StatusBadGateway, err.Error())
		return
	}
	headers := map[string]string{"Authorization": c.GetHeader("Authorization")}
	if err := p.sendRequest(http.MethodPost, "/tasks/batch", bytes.NewBuffer([]byte(jsonB)), &response, &headers); err != nil {
		p.log.Warn("sendRequest err", err)
		c.JSON(http.StatusBadGateway, err.Error())
		return
	}
	var resp []models.Task
	for _, v := range response {
		for _, v2 := range taskEntities {
			if v.ID == v2.ID {
				resp = append(resp, models.Task{
					MilestoneId:        int64(milestoneID),
					EpicID:             v2.EpicID,
					AssigneeId:         v.AssigneeId,
					BrandId:            v.BrandId,
					BucketID:           v.BucketID,
					CreatorID:          v.CreatorID,
					EndTime:            v.EndTime,
					ID:                 v.ID,
					NumberOfAttachment: v.NumberOfAttachment,
					NumberOfComments:   v.NumberOfComments,
					Priority:           v.Priority,
					ResolvedTime:       v.ResolvedTime,
					StartTime:          v.StartTime,
					Status:             v.Status,
					StatusID:           v.StatusID,
					Title:              v.Title,
				})
			}

		}
	}
	c.JSON(http.StatusOK, resp)
}

func (p taskHandler) sendRequest(method, uri string, reader io.Reader, respStruct interface{}, headers *map[string]string) error {
	req, err := http.NewRequest(method, p.conf.Task.Addr+uri, reader)
	if err != nil {
		return err
	}
	client := http.Client{
		Timeout: 15 * time.Second,
	}

	// req.Header.Set("Content-Type")
	if headers != nil {
		for s, v := range *headers {
			req.Header.Set(s, v)
		}
	}
	resp, err := client.Do(req)
	if err != nil {
		return err
	}
	defer resp.Body.Close()
	if resp.StatusCode == http.StatusOK {
		body, _ := ioutil.ReadAll(resp.Body)
		if respStruct != nil {
			err = json.Unmarshal(body, &respStruct)

			if err != nil {
				return err
			}
		}

	} else {
		body, _ := ioutil.ReadAll(resp.Body)
		return errors.New(string(body))
	}

	return nil
}

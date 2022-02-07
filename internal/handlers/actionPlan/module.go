package actionPlan

import (
	"bytes"
	"encoding/json"
	"errors"
	"io"
	"io/ioutil"
	"log"
	"net/http"
	"projects/internal/database/actionPlan"
	"projects/internal/database/epics"
	"projects/internal/database/milestone"
	"projects/internal/database/stage"
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

var Module = fx.Provide(NewActionPlanHandler)

type ActionPlanHandler interface {
	DeleteActionPlan(c *gin.Context)
	UpdateActionPlan(c *gin.Context)
	DownloadActionPlan(c *gin.Context)
	GetAcPlans(c *gin.Context)
	CreateActionPlan(c *gin.Context)
}

type Params struct {
	fx.In
	db.DbInter
	*config.Tuner
	*logrus.Logger
}

type actionPlanHandler struct {
	db   db.DbInter
	log  *logrus.Logger
	conf *config.Tuner
}

func NewActionPlanHandler(params Params) ActionPlanHandler {
	return &actionPlanHandler{db: params.DbInter, log: params.Logger, conf: params.Tuner}
}

func (p actionPlanHandler) DeleteActionPlan(c *gin.Context) {
	aP := actionPlan.New(p.db.GetDB())
	id, err := strconv.ParseInt(c.Param("id"), 10, 64)
	if err != nil {
		p.log.Warnln("Param err: ", err.Error())
		c.JSON(http.StatusBadGateway, gin.H{"error": err.Error()})
		return
	}
	if id <= 0 {
		c.JSON(http.StatusBadRequest, gin.H{"error": "wrong id"})
		return
	}

	if err := aP.Delete(id); err != nil {
		p.log.Warnln("Can't delete action plan with err: ", err.Error())
		c.JSON(http.StatusBadGateway, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"result": "success"})
}

func (p actionPlanHandler) UpdateActionPlan(c *gin.Context) {
	aP := actionPlan.New(p.db.GetDB())
	id, err := strconv.ParseInt(c.Param("id"), 10, 64)
	if err != nil {
		p.log.Warnln("Param err: ", err.Error())
		c.JSON(http.StatusBadGateway, gin.H{"error": err.Error()})
		return
	}
	if id <= 0 {
		c.JSON(http.StatusBadRequest, gin.H{"error": "wrong id"})
		return
	}
	var actionPlanReq models.ActionPlan
	if err := c.ShouldBindJSON(&actionPlanReq); err != nil {
		p.log.Warnln("bind error")
		c.JSON(http.StatusBadGateway, gin.H{"error": "bind error"})
		return
	}

	if len(actionPlanReq.Title) < 2 {
		p.log.Warnln("can`t update title with less then 2 simbols")
		c.JSON(http.StatusBadGateway, gin.H{"error": "can`t update title with less then 2 simbols"})
		return
	}
	acPlanEntity, err := aP.Update(id, actionPlanReq.Title, actionPlanReq.Status)
	if err != nil {
		p.log.Warnln("Can't update action plan with err: ", err.Error())
		c.JSON(http.StatusBadGateway, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, acPlanEntity)
}

func (p actionPlanHandler) DownloadActionPlan(c *gin.Context) {
	aP := actionPlan.New(p.db.GetDB())
	st := stage.New(p.db.GetDB())
	ml := milestone.New(p.db.GetDB())
	id, err := strconv.ParseInt(c.Param("id"), 10, 64)
	if err != nil {
		p.log.Warnln("Param err: ", err.Error())
		c.JSON(http.StatusBadGateway, gin.H{"error": err.Error()})
		return
	}
	if id <= 0 {
		c.JSON(http.StatusBadRequest, gin.H{"error": "wrong id"})
		return
	}
	acPlan, err := aP.Get(id)
	if err != nil {
		p.log.Warnln("Can't update action plan with err: ", err.Error())
		c.JSON(http.StatusBadGateway, gin.H{"error": err.Error()})
		return
	}

	stages := st.GetByActionPlan(acPlan.ActionPlanID)
	miles := ml.GetByActionPlan(acPlan.ActionPlanID)
	epic, err := epics.NewEpic(p.db.GetDB()).GetEpic(epics.EpicEntity{ActionPlanID: acPlan.ActionPlanID})
	if err != nil {
		p.log.Warnln("Can't get action plan with err: ", err.Error())
		c.JSON(http.StatusBadGateway, gin.H{"error": err.Error()})
		return
	}
	epc := make(map[int64]models.EpicResponse)
	for _, e := range epic {
		epicResp := models.EpicResponse{
			ID:          e.ID,
			WorkspaceID: e.WorkspaceID,
			StageID:     e.StageID,
			MilestoneID: e.MilestoneID,
			Title:       e.Title,
			Description: e.Description,
		}
		epc[e.ID] = epicResp
	}
	task, err := tasks.New(p.db.GetDB()).GetTaskByActionPlanID(acPlan.ActionPlanID)
	if err != nil {
		p.log.Warnln("Can't get tasks with err: ", err.Error())
		c.JSON(http.StatusBadGateway, gin.H{"error": err.Error()})
		return
	}
	var tasksID []int64
	for _, t := range task {
		tasksID = append(tasksID, t.ID)
	}
	var resp []models.Task
	if len(tasksID) != 0 {

		headers := map[string]string{"Authorization": c.GetHeader("Authorization")}
		log.Println(headers)
		arrByte, err := json.Marshal(tasksID)
		if err != nil {
			p.log.Warnln("Can't get tasks with err: ", err.Error())
			c.JSON(http.StatusBadGateway, gin.H{"error": err.Error()})
			return
		}
		if err := p.sendRequest(http.MethodPost, "/tasks/batch", bytes.NewBuffer(arrByte), &resp, &headers); err != nil {
			p.log.Warn("sendRequest err", err)
			c.JSON(http.StatusBadGateway, err.Error())
			return
		}
	}
	taskMile := make(map[int64][]models.Task)
	for _, t := range task {
		if len(resp) != 0 {
			for _, v := range resp {
				if v.ID == t.ID {
					singleEpic, ok := epc[t.EpicID]
					singleTask := v
					singleTask.MilestoneId = t.MilestoneID
					singleTask.EpicID = t.EpicID
					singleTask.ActionPlanID = t.ActionPlanID

					if singleTask.EpicID == singleEpic.ID && ok {
						singleEpic.Task = append(singleEpic.Task, singleTask)
						epc[t.EpicID] = singleEpic
						continue
					}
					taskMile[t.MilestoneID] = append(taskMile[t.MilestoneID], singleTask)
				}
			}
		}
	}
	mileEpic := make(map[int64][]models.EpicResponse)
	for _, e := range epc {
		mileEpic[e.MilestoneID] = append(mileEpic[e.MilestoneID], e)
	}

	actionPlanResp := models.ActionPlanResp{
		ActionPlanID: acPlan.ActionPlanID,
		PhaseID:      acPlan.PhaseID,
		ProjectID:    acPlan.ProjectID,
		Created:      acPlan.Created,
		Title:        acPlan.Title,
		Status:       acPlan.Status.String(),
		Stage: func() []models.Stage {
			var stagesResp []models.Stage
			for _, v := range stages {
				stagesResp = append(stagesResp,
					models.Stage{StageID: v.StageID,
						Order:       v.Order,
						Title:       v.Title,
						Description: v.Description,
						DateStart:   v.DateStart,
						DateEnd:     v.DateStop,
						Milestone: func() []models.Milestone {
							var milesResp []models.Milestone
							for _, mile := range miles {
								if mile.StageID == v.StageID {
									milesResp = append(milesResp,
										models.Milestone{
											MilestoneID: mile.MilestoneID,
											Order:       mile.Order,
											StageID:     mile.StageID,
											Status:      mile.Status.String(),
											Title:       mile.Title,
											Description: mile.Description,
											DateStart:   mile.DateStart,
											DateEnd:     mile.DateStop,
											Epic:        mileEpic[mile.MilestoneID],
											Task:        taskMile[mile.MilestoneID],
										})
								}
							}
							return milesResp
						}(),
					})
			}
			return stagesResp
		}(),
	}
	if actionPlanResp.ActionPlanID == 0 {
		c.JSON(http.StatusBadGateway, gin.H{"error": "can`t find action plan with some id"})
		return
	}

	c.JSON(http.StatusOK, actionPlanResp)
}
func (p actionPlanHandler) GetAcPlans(c *gin.Context) {
	aP := actionPlan.New(p.db.GetDB())
	id, err := strconv.ParseInt(c.Param("id"), 10, 64)
	if err != nil {
		p.log.Warnln("Param err: ", err.Error())
		c.JSON(http.StatusBadGateway, gin.H{"error": err.Error()})
		return
	}
	acPlanEntities, err := aP.GetByProjectID(id)
	if err != nil {
		p.log.Warnln("GetByProjectID err: ", err.Error())
		c.JSON(http.StatusBadGateway, gin.H{"error": err.Error()})
		return
	}
	var acPlanResp []models.ActionPlanResp
	for _, v := range acPlanEntities {
		acPlanResp = append(acPlanResp, models.ActionPlanResp{
			ActionPlanID: v.ActionPlanID,
			Title:        v.Title,
			Status:       v.Status.String(),
			Created:      v.Created,
			PhaseID:      v.PhaseID,
		})
	}

	c.JSON(http.StatusOK, acPlanResp)

}

func (p actionPlanHandler) CreateActionPlan(c *gin.Context) {
	var actionPlanReq models.ActionPlan
	if err := c.ShouldBindJSON(&actionPlanReq); err != nil {
		p.log.Warnln("bind error")
		c.JSON(http.StatusBadGateway, gin.H{"error": "bind error"})
		return
	}

	acEntity := actionPlan.ActionPlanEntity{
		WorkspaceID: actionPlanReq.WorkspaceID,
		ProjectID:   actionPlanReq.ProjectID,
		Title:       actionPlanReq.Title,
		Status:      actionPlan.ActionStatus(actionPlanReq.Status),
		PhaseID:     actionPlanReq.PhaseID,
	}
	repo := actionPlan.New(p.db.GetDB())
	if err := repo.Create(&acEntity); err != nil {
		p.log.Warnln("can't create action plan entity", err)
		c.JSON(http.StatusBadGateway, gin.H{"error": "can't create action plan"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"created_id": acEntity.ActionPlanID})
}

func (p actionPlanHandler) sendRequest(method, uri string, reader io.Reader, respStruct interface{}, headers *map[string]string) error {
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

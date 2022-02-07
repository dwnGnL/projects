package milestone

import (
	"net/http"
	"projects/internal/database/milestone"
	"projects/internal/models"
	"projects/pkg/config"
	"projects/pkg/db"
	"strconv"

	"github.com/gin-gonic/gin"
	"github.com/sirupsen/logrus"
	"go.uber.org/fx"
)

var Module = fx.Provide(NewMilestoneHandler)

type MilestoneHandler interface {
	GetMilestoneByID(c *gin.Context)
	GetMilestones(c *gin.Context)
	CreateMilestone(c *gin.Context)
	EditMilestone(c *gin.Context)
	DeleteMilestone(c *gin.Context)
}

type Params struct {
	fx.In
	db.DbInter
	*config.Tuner
	*logrus.Logger
}

type milestoneHandler struct {
	db   db.DbInter
	log  *logrus.Logger
	conf *config.Tuner
}

func NewMilestoneHandler(params Params) MilestoneHandler {
	return &milestoneHandler{db: params.DbInter, log: params.Logger, conf: params.Tuner}
}

func (p milestoneHandler) GetMilestoneByID(c *gin.Context) {
	id, err := strconv.ParseInt(c.Param("id"), 10, 64)
	if err != nil {
		p.log.Warnln("Param err: ", err.Error())
		c.JSON(http.StatusBadGateway, gin.H{"error": err.Error()})
		return
	}

	repo := milestone.New(p.db.GetDB())
	milestoneEnt := repo.GetMilestoneByID(id)

	c.JSON(http.StatusOK, models.Milestone{
		MilestoneID: milestoneEnt.MilestoneID,
		ProjectID:   milestoneEnt.ProjectID,
		StageID:     milestoneEnt.StageID,
		Order:       milestoneEnt.Order,
		Status:      milestoneEnt.Status.String(),
		DateStart:   milestoneEnt.DateStart,
		DateEnd:     milestoneEnt.DateStop,
		Title:       milestoneEnt.Title,
		AssignID:    milestoneEnt.AssignID,
	})
}

func (p milestoneHandler) GetMilestones(c *gin.Context) {
	values := c.Request.URL.Query()
	sIdString := values.Get("stage_id")

	stageID, err := strconv.Atoi(sIdString)
	if err != nil {
		p.log.Warnln("Param err: ", err.Error())
		c.JSON(http.StatusBadRequest, gin.H{"error": "wrong id"})
		return
	}

	mRepo := milestone.New(p.db.GetDB())
	milestones := mRepo.GetByStageID(int64(stageID))
	if len(milestones) == 0 {
		c.JSON(http.StatusOK, gin.H{})
		return
	}

	var ms []models.Milestone
	for _, m := range milestones {
		ms = append(ms, models.Milestone{
			MilestoneID: m.MilestoneID,
			ProjectID:   m.ProjectID,
			StageID:     m.StageID,
			Order:       m.Order,
			DateStart:   m.DateStart,
			DateEnd:     m.DateStop,
			Description: m.Description,
			Title:       m.Title,
			Status:      m.Status.String(),
			AssignID:    m.AssignID,
		})
	}

	c.JSON(http.StatusOK, ms)
}

func (p milestoneHandler) CreateMilestone(c *gin.Context) {
	var milestoneReq models.Milestones
	if err := c.ShouldBindJSON(&milestoneReq); err != nil {
		p.log.Warnln("bind error")
		c.JSON(http.StatusBadGateway, gin.H{"error": "bind error"})
		return
	}

	mRepo := milestone.New(p.db.GetDB())
	var milestones []milestone.MilestoneEntity
	for _, m := range milestoneReq.Milestones {
		milestones = append(milestones, milestone.MilestoneEntity{
			Title:       m.Title,
			ProjectID:   m.ProjectID,
			StageID:     m.StageID,
			Description: m.Description,
			Order:       m.Order,
			Status:      milestone.NewStatus,
			DateStart:   m.DateStart,
			DateStop:    m.DateEnd,
			AssignID:    m.AssignID,
			ProcessID:   m.ProcessID,
		})
	}

	miles, err := mRepo.CreateMany(milestones)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	var ms []models.Milestone
	for _, m := range miles {
		ms = append(ms, models.Milestone{
			MilestoneID: m.MilestoneID,
			ProjectID:   m.ProjectID,
			StageID:     m.StageID,
			Order:       m.Order,
			DateStart:   m.DateStart,
			Status:      m.Status.String(),
			DateEnd:     m.DateStop,
			Description: m.Description,
			Title:       m.Title,
			AssignID:    m.AssignID,
			ProcessID:   m.ProcessID,
		})
	}

	c.JSON(http.StatusOK, ms)
}

func (p milestoneHandler) EditMilestone(c *gin.Context) {
	id, err := strconv.ParseInt(c.Param("id"), 10, 64)
	if err != nil {
		p.log.Warnln("Param err: ", err.Error())
		c.JSON(http.StatusBadGateway, gin.H{"error": err.Error()})
		return
	}
	var milestoneReq models.Milestone
	if err := c.ShouldBindJSON(&milestoneReq); err != nil {
		p.log.Warnln("bind error")
		c.JSON(http.StatusBadGateway, gin.H{"error": "bind error"})
		return
	}

	mileDB := milestone.New(p.db.GetDB())
	milestoneEntity := milestone.MilestoneEntity{
		MilestoneID: id,
		StageID:     milestoneReq.StageID,
		Title:       milestoneReq.Title,
		Description: milestoneReq.Description,
		Order:       milestoneReq.Order,
		Status:      milestone.Status(milestoneReq.Status),
		DateStart:   milestoneReq.DateStart,
		DateStop:    milestoneReq.DateEnd,
		AssignID:    milestoneReq.AssignID,
	}

	mile, err := mileDB.Update(milestoneEntity)
	if err != nil {
		p.log.Warnln("update error")
		c.JSON(http.StatusBadGateway, gin.H{"error": "update error"})
		return
	}
	milestoneResp := models.Milestone{
		MilestoneID: mile.MilestoneID,
		AssignID:    milestoneReq.AssignID,
		StageID:     milestoneReq.StageID,
		Title:       mile.Title,
		Description: mile.Description,
		Status:      mile.Status.String(),
		DateStart:   mile.DateStart,
		DateEnd:     mile.DateStop,
		Order:       mile.Order,
	}

	c.JSON(http.StatusOK, milestoneResp)
}

func (p milestoneHandler) DeleteMilestone(c *gin.Context) {
	id, err := strconv.ParseInt(c.Param("id"), 10, 64)
	if err != nil {
		p.log.Warnln("Param err: ", err.Error())
		c.JSON(http.StatusBadGateway, gin.H{"error": err.Error()})
		return
	}

	mileDB := milestone.New(p.db.GetDB())
	if err := mileDB.DeleteByID(id); err != nil {
		p.log.Warnln("delete error ", err)
		c.JSON(http.StatusBadGateway, gin.H{"error": "delete error"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"success": true})
}

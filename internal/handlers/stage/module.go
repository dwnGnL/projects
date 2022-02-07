package stage

import (
	"net/http"
	"projects/internal/database/stage"
	"projects/internal/models"
	"projects/pkg/config"
	"projects/pkg/db"
	"strconv"

	"github.com/gin-gonic/gin"
	"github.com/sirupsen/logrus"
	"go.uber.org/fx"
)

var Module = fx.Provide(NewStageHandler)

type StageHandler interface {
	DeleteStage(c *gin.Context)
	GetStageByAcPLan(c *gin.Context)
	GetStageByProjectID(c *gin.Context)
	UpdateStage(c *gin.Context)
	CreateStage(c *gin.Context)
}

type Params struct {
	fx.In
	db.DbInter
	*config.Tuner
	*logrus.Logger
}

type stageHandler struct {
	db   db.DbInter
	log  *logrus.Logger
	conf *config.Tuner
}

func NewStageHandler(params Params) StageHandler {
	return &stageHandler{db: params.DbInter, log: params.Logger, conf: params.Tuner}
}
func (p stageHandler) CreateStage(c *gin.Context) {
	var stagesReq models.ProjectTemplate
	if err := c.ShouldBindJSON(&stagesReq); err != nil {
		p.log.Warnln(err.Error())
		c.JSON(http.StatusBadGateway, gin.H{"error": "bind error"})
		return
	}
	var stages []stage.StageEntity
	for _, st := range stagesReq.Stage {
		stages = append(stages, stage.StageEntity{
			Title:        st.Title,
			Order:        st.Order,
			ProjectID:    st.ProjectID,
			WorkspaceID:  st.WorkspaceID,
			ActionPlanID: st.ActionPlanID,
		})
	}

	s := stage.New(p.db.GetDB())
	sts, err := s.CreateMany(stages)
	if err != nil {
		p.log.Warnln("Param err: ", err.Error())
		c.JSON(http.StatusBadGateway, gin.H{"error": err.Error()})
		return
	}
	var stageResponse []models.Stage
	for _, entity := range sts {
		stageResponse = append(stageResponse, models.Stage{
			StageID:      entity.StageID,
			Order:        entity.Order,
			Title:        entity.Title,
			ProjectID:    entity.ProjectID,
			ActionPlanID: entity.ActionPlanID,
			WorkspaceID:  entity.WorkspaceID,
			Hidden:       entity.Hidden,
		})
	}

	c.JSON(http.StatusOK, stageResponse)
}

func (p stageHandler) UpdateStage(c *gin.Context) {
	id, err := strconv.ParseInt(c.Param("id"), 10, 64)
	if err != nil {
		p.log.Warnln("Param err: ", err.Error())
		c.JSON(http.StatusBadGateway, gin.H{"error": err.Error()})
		return
	}

	var stageReq models.Stage
	if err := c.ShouldBindJSON(&stageReq); err != nil {
		p.log.Warnln("bind error")
		c.JSON(http.StatusBadGateway, gin.H{"error": "bind error"})
		return
	}

	s := stage.New(p.db.GetDB())
	stageUpdated, err := s.Update(stage.StageEntity{
		StageID:      id,
		Title:        stageReq.Title,
		Order:        stageReq.Order,
		DateStart:    stageReq.DateStart,
		Description:  stageReq.Description,
		DateStop:     stageReq.DateEnd,
		ActionPlanID: stageReq.ActionPlanID,
		Hidden:       stageReq.Hidden,
	})
	if err != nil {
		p.log.Warnln("Param err: ", err.Error())
		c.JSON(http.StatusBadGateway, gin.H{"error": err.Error()})
		return
	}

	stageResponse := models.Stage{
		StageID:      stageUpdated.StageID,
		Order:        stageUpdated.Order,
		Title:        stageUpdated.Title,
		ProjectID:    stageUpdated.ProjectID,
		DateStart:    stageUpdated.DateStart,
		Description:  stageUpdated.Description,
		DateEnd:      stageUpdated.DateStop,
		ActionPlanID: stageUpdated.ActionPlanID,
		WorkspaceID:  stageUpdated.WorkspaceID,
		Hidden:       stageUpdated.Hidden,
	}

	c.JSON(http.StatusOK, stageResponse)
}

func (p stageHandler) GetStageByProjectID(c *gin.Context) {
	id, err := strconv.ParseInt(c.Param("id"), 10, 64)
	if err != nil {
		p.log.Warnln("Param err: ", err.Error())
		c.JSON(http.StatusBadGateway, gin.H{"error": err.Error()})
		return
	}
	s := stage.New(p.db.GetDB())
	stages := s.GetByProjectID(id)
	var stagesResponse []models.Stage
	for _, st := range stages {
		stagesResponse = append(stagesResponse, models.Stage{
			StageID:      st.StageID,
			Order:        st.Order,
			Title:        st.Title,
			DateStart:    st.DateStart,
			Description:  st.Description,
			DateEnd:      st.DateStop,
			ActionPlanID: st.ActionPlanID,
			Hidden:       st.Hidden,
			WorkspaceID:  st.WorkspaceID,
			ProjectID:    st.ProjectID,
		})
	}

	c.JSON(http.StatusOK, stagesResponse)
}

func (p stageHandler) GetStageByAcPLan(c *gin.Context) {
	id, err := strconv.ParseInt(c.Param("id"), 10, 64)
	if err != nil {
		p.log.Warnln("Param err: ", err.Error())
		c.JSON(http.StatusBadGateway, gin.H{"error": err.Error()})
		return
	}
	s := stage.New(p.db.GetDB())
	stages := s.GetByActionPlan(id)
	var stagesResponse []models.Stage
	for _, st := range stages {
		stagesResponse = append(stagesResponse, models.Stage{
			StageID:      st.StageID,
			Order:        st.Order,
			Title:        st.Title,
			DateStart:    st.DateStart,
			Description:  st.Description,
			DateEnd:      st.DateStop,
			ActionPlanID: st.ActionPlanID,
			Hidden:       st.Hidden,
			WorkspaceID:  st.WorkspaceID,
			ProjectID:    st.ProjectID,
		})
	}

	c.JSON(http.StatusOK, stagesResponse)
}

func (p stageHandler) DeleteStage(c *gin.Context) {
	id, err := strconv.ParseInt(c.Param("id"), 10, 64)
	if err != nil {
		p.log.Warnln("Param err: ", err.Error())
		c.JSON(http.StatusBadGateway, gin.H{"error": err.Error()})
		return
	}

	s := stage.New(p.db.GetDB())
	if err := s.DeleteStage(id); err != nil {
		p.log.Warnln("Can't delete stage with err: ", err.Error())
		c.JSON(http.StatusBadGateway, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"result": "success"})
}

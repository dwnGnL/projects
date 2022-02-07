package epic

import (
	"net/http"
	"projects/internal/database/epics"
	"projects/internal/models"
	"projects/pkg/config"
	"projects/pkg/db"
	"strconv"

	"github.com/gin-gonic/gin"
	"github.com/sirupsen/logrus"
	"go.uber.org/fx"
)

var Module = fx.Provide(NewEpicHandler)

type EpicHandler interface {
	CreateEpic(c *gin.Context)
	GetEpics(c *gin.Context)
	UpdateEpic(c *gin.Context)
	DeleteEpic(c *gin.Context)
}

type Params struct {
	fx.In
	db.DbInter
	*config.Tuner
	*logrus.Logger
}

type epicHandler struct {
	db  db.DbInter
	log *logrus.Logger
}

func NewEpicHandler(params Params) EpicHandler {
	return &epicHandler{db: params.DbInter, log: params.Logger}
}

func (p epicHandler) CreateEpic(c *gin.Context) {
	var epic models.EpicRequest
	if err := c.ShouldBindJSON(&epic); err != nil {
		p.log.Warnln("bind error")
		c.JSON(http.StatusBadGateway, err.Error())
		return
	}

	newEpic, err := epics.NewEpic(p.db.GetDB()).CreateEpic(epics.EpicEntity{
		WorkspaceID: epic.WorkspaceID,
		ProjectID:   epic.ProjectID,
		StageID:     epic.StageID,
		MilestoneID: epic.MilestoneID,
		Title:       epic.Title,
		Description: epic.Description,
	})
	if err != nil {
		p.log.Warnln(err)
		c.JSON(http.StatusBadGateway, "can't create epic")
		return
	}

	c.JSON(http.StatusOK, models.EpicResponse{
		ID:          newEpic.ID,
		WorkspaceID: newEpic.WorkspaceID,
		ProjectID:   newEpic.ProjectID,
		StageID:     newEpic.StageID,
		MilestoneID: newEpic.MilestoneID,
		Title:       newEpic.Title,
		Description: newEpic.Description,
	})
}

func (p epicHandler) GetEpics(c *gin.Context) {
	var epic models.EpicRequest
	if err := c.ShouldBindJSON(&epic); err != nil {
		p.log.Warnln("bind error")
		c.JSON(http.StatusBadGateway, err.Error())
		return
	}

	var epicEntity epics.EpicEntity
	if epic.WorkspaceID > 0 {
		epicEntity.WorkspaceID = epic.WorkspaceID
	}
	if epic.ProjectID > 0 {
		epicEntity.ProjectID = epic.ProjectID
	}
	if epic.MilestoneID > 0 {
		epicEntity.MilestoneID = epic.MilestoneID
	}
	if epic.StageID > 0 {
		epicEntity.StageID = epic.StageID
	}
	eps, err := epics.NewEpic(p.db.GetDB()).GetEpic(epicEntity)
	if err != nil {
		p.log.Warnln(err)
		c.JSON(http.StatusBadGateway, "can't get epics")
		return
	}

	var response []models.EpicResponse
	for _, ep := range eps {
		response = append(response, models.EpicResponse{
			ID:          ep.ID,
			WorkspaceID: ep.WorkspaceID,
			MilestoneID: ep.MilestoneID,
			ProjectID:   ep.ProjectID,
			StageID:     ep.StageID,
			Title:       ep.Title,
			Description: ep.Description,
		})
	}

	c.JSON(http.StatusOK, response)
}

func (p epicHandler) UpdateEpic(c *gin.Context) {
	var epic models.EpicRequest
	if err := c.ShouldBindJSON(&epic); err != nil {
		p.log.Warnln("bind error")
		c.JSON(http.StatusBadGateway, err.Error())
		return
	}
	id, err := strconv.ParseInt(c.Param("id"), 10, 64)
	if err != nil || id <= 0 {
		p.log.Warnln("Param err: ", err.Error())
		c.JSON(http.StatusBadGateway, gin.H{"error": err.Error()})
		return
	}
	epic.ID = id

	if err := epics.NewEpic(p.db.GetDB()).UpdateEpic(epics.EpicEntity{
		ID:          epic.ID,
		StageID:     epic.StageID,
		MilestoneID: epic.MilestoneID,
		Title:       epic.Title,
		Description: epic.Description,
	}); err != nil {
		p.log.Warnln(err)
		c.JSON(http.StatusBadGateway, "can't update epic")
		return
	}

	c.JSON(http.StatusOK, gin.H{"success": "ok"})
}

func (p epicHandler) DeleteEpic(c *gin.Context) {
	id, err := strconv.ParseInt(c.Param("id"), 10, 64)
	if err != nil {
		p.log.Warnln("Param err: ", err.Error())
		c.JSON(http.StatusBadGateway, gin.H{"error": err.Error()})
		return
	}

	if err := epics.NewEpic(p.db.GetDB()).DeleteEpic(id); err != nil {
		p.log.Warnln(err)
		c.JSON(http.StatusBadGateway, "can't delete epic")
		return
	}

	c.JSON(http.StatusOK, gin.H{"success": "ok"})
}

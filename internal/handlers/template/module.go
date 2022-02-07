package template

import (
	"net/http"
	"projects/internal/database/milestone"
	"projects/internal/database/stage"
	"projects/internal/models"
	"projects/pkg/config"
	"projects/pkg/db"
	"strconv"

	"github.com/gin-gonic/gin"
	"github.com/sirupsen/logrus"
	"go.uber.org/fx"
)

var Module = fx.Provide(NewtemplateHandler)

type TemplateHandler interface {
	UpdateTemplate(c *gin.Context)
	GetTemplates(c *gin.Context)
}

type Params struct {
	fx.In
	db.DbInter
	*config.Tuner
	*logrus.Logger
}

type templateHandler struct {
	db   db.DbInter
	log  *logrus.Logger
	conf *config.Tuner
}

func NewtemplateHandler(params Params) TemplateHandler {
	return &templateHandler{db: params.DbInter, log: params.Logger, conf: params.Tuner}
}
func (p templateHandler) UpdateTemplate(c *gin.Context) {
	tr := p.db.GetDB().Begin()
	sch := stage.New(tr)
	st := milestone.New(tr)
	id, err := strconv.ParseInt(c.Param("id"), 10, 64)
	if err != nil {
		p.log.Warnln("Param err: ", err.Error())
		c.JSON(http.StatusBadGateway, gin.H{"error": err.Error()})
		return
	}

	var stageReq models.ProjectTemplate
	if err := c.ShouldBindJSON(&stageReq); err != nil {
		p.log.Warnln("bind error")
		c.JSON(http.StatusBadGateway, gin.H{"error": "bind error"})
		return
	}
	schedulesDB := sch.GetByProjectID(id)
	var schedulesForUpdate []stage.StageEntity
	for _, v := range stageReq.Stage {
		schedulesForUpdate = append(schedulesForUpdate, stage.StageEntity{StageID: v.StageID, ProjectID: id, Order: v.Order, Title: v.Title, DateStart: v.DateStart,
			Description: v.Description,
			DateStop:    v.DateEnd})
	}
	for _, v := range schedulesDB {
		deleted := &v
		for _, v2 := range stageReq.Stage {
			if v2.StageID == v.StageID {
				deleted = nil
			}
		}
		if deleted != nil {
			deleted.Hidden = true
			schedulesForUpdate = append(schedulesForUpdate, *deleted)
		}
	}
	for _, scheduleForUpdate := range schedulesForUpdate {
		_, err := sch.Update(scheduleForUpdate)
		if err != nil {
			tr.Rollback()
			p.log.Warnln("schedule update error")
			c.JSON(http.StatusBadGateway, gin.H{"error": "schedule update error"})
			return
		}
	}

	milestonesDB := st.GetByProjectID(id)
	var milestonesForUpdate []milestone.MilestoneEntity
	for _, v := range stageReq.Stage {
		for _, v2 := range v.Milestone {
			milestonesForUpdate = append(milestonesForUpdate, milestone.MilestoneEntity{MilestoneID: v2.MilestoneID, Status: milestone.Status(v2.Status), StageID: v.StageID, ProjectID: id, Order: v2.Order, Title: v2.Title, DateStart: v.DateStart,
				Description: v.Description,
				DateStop:    v.DateEnd})
		}
	}
	for _, v := range milestonesDB {
		deleted := &v
		for _, v2 := range milestonesForUpdate {
			if v2.MilestoneID == v.MilestoneID {
				deleted = nil
			}
		}
		if deleted != nil {
			deleted.Hidden = true
			milestonesForUpdate = append(milestonesForUpdate, *deleted)
		}
	}
	for _, milestoneForUpdate := range milestonesForUpdate {
		_, err := st.Update(milestoneForUpdate)
		if err != nil {
			tr.Rollback()
			p.log.Warnln("milestones update error")
			c.JSON(http.StatusBadGateway, gin.H{"error": "milestones update error"})
			return
		}
	}
	tr.Commit()
	p.log.Println("milestones and schedule update copmliete")
	var projTemplate models.ProjectTemplate

	for _, v := range schedulesForUpdate {
		if !v.Hidden {
			schedules := models.Stage{StageID: v.StageID, Order: v.Order, Title: v.Title, DateStart: v.DateStart,
				Description: v.Description,
				DateEnd:     v.DateStop}
			for _, v2 := range milestonesForUpdate {
				if !v2.Hidden {
					if v2.StageID == v.StageID {
						schedules.Milestone = append(schedules.Milestone, models.Milestone{MilestoneID: v2.MilestoneID, Order: v2.Order, Title: v2.Title, DateStart: v2.DateStart,
							Description: v2.Description,
							DateEnd:     v2.DateStop,
							Status:      v2.Status.String()})
					}
				}
			}
			projTemplate.Stage = append(projTemplate.Stage, schedules)
		}
	}
	c.JSON(http.StatusOK, projTemplate)

}

func (p templateHandler) GetTemplates(c *gin.Context) {
	id, err := strconv.ParseInt(c.Param("id"), 10, 64)
	if err != nil {
		p.log.Warnln("Param err: ", err.Error())
		c.JSON(http.StatusBadGateway, gin.H{"error": err.Error()})
		return
	}
	db := p.db.GetDB()
	sch := stage.New(db)
	st := milestone.New(db)
	schedulesDB := sch.GetByProjectID(id)
	milestonesDB := st.GetByProjectID(id)

	var projTemplate models.ProjectTemplate

	for _, v := range schedulesDB {
		if !v.Hidden {
			schedules := models.Stage{
				StageID: v.StageID, Order: v.Order,
				Title: v.Title, DateStart: v.DateStart,
				Description: v.Description, DateEnd: v.DateStop,
			}
			for _, v2 := range milestonesDB {
				if !v2.Hidden {
					if v2.StageID == v.StageID {
						schedules.Milestone = append(schedules.Milestone, models.Milestone{
							MilestoneID: v2.MilestoneID,
							Order:       v2.Order,
							StageID:     v2.StageID,
							Title:       v2.Title,
							Status:      v2.Status.String(),
							DateStart:   v2.DateStart,
							Description: v2.Description,
							DateEnd:     v2.DateStop,
						})
					}
				}
			}
			projTemplate.Stage = append(projTemplate.Stage, schedules)
		}
	}
	c.JSON(http.StatusOK, projTemplate)
}

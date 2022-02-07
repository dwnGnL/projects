package project

import (
	"net/http"
	"projects/internal/database/actionPlan"
	"projects/internal/database/milestone"
	"projects/internal/database/projects"
	"projects/internal/database/stage"
	"projects/internal/database/workspace"
	"projects/internal/models"
	"projects/pkg/db"
	"strconv"

	"github.com/gin-gonic/gin"
	"github.com/sirupsen/logrus"
	"go.uber.org/fx"
)

var Module = fx.Provide(NewprojectHandler)

type ProjectHandler interface {
	GetProjects(c *gin.Context)
	Project(c *gin.Context)
	UpdateProject(c *gin.Context)
	CreateProject(c *gin.Context)
	DeleteProject(ctx *gin.Context)
}

type Params struct {
	fx.In
	db.DbInter
	*logrus.Logger
}

type projectHandler struct {
	db  db.DbInter
	log *logrus.Logger
}

func NewprojectHandler(params Params) ProjectHandler {
	return &projectHandler{db: params.DbInter, log: params.Logger}
}

func (p projectHandler) GetProjects(c *gin.Context) {
	pr := projects.New(p.db.GetDB())
	var filterReq models.ProjectFilter

	if err := c.BindQuery(&filterReq); err != nil {
		p.log.Warnln("bind err: ", err.Error())
		c.JSON(http.StatusBadGateway, gin.H{"error": err.Error()})
		return
	}
	filter := projects.ProjectFilter{
		Cluster: filterReq.Cluster,
		Type:    filterReq.Type,
		Stage:   filterReq.Stage,
	}

	proj, err := pr.GetAll(filter)
	if err != nil {
		p.log.Warnln("Get projects err: ", err.Error())
		c.JSON(http.StatusBadGateway, gin.H{"error": err.Error()})
		return
	}
	projectsResp := []models.ProjectResp{}
	for _, v := range proj {
		var phases []models.Phase
		for _, phase := range v.Phases {
			phases = append(phases, models.Phase{
				ID:   phase.PhaseID,
				Name: phase.Name,
			})
		}

		projectResp := models.ProjectResp{
			ProjectID:     v.ProjectID,
			Title:         v.Title,
			Description:   v.Description,
			MediaID:       v.MediaID,
			Type:          v.Type.String(),
			BusinessOwner: v.BusinessOwner,
			LegacyEntity:  v.LegacyEntity,
			Cluster:       v.Cluster,
			Phase:         phases,
			Stage:         v.Stage.String(),
			OwnerID:       v.OwnerID,
			Category:      v.Category.String(),
			Created:       v.Created,
			Priority:      v.Priority,
		}
		projectsResp = append(projectsResp, projectResp)
	}
	c.JSON(http.StatusOK, projectsResp)
}

func (p projectHandler) Project(c *gin.Context) {
	pr := projects.New(p.db.GetDB())

	id, err := strconv.ParseInt(c.Param("id"), 10, 64)
	if err != nil {
		p.log.Warnln("Param err: ", err.Error())
		c.JSON(http.StatusBadGateway, gin.H{"error": err.Error()})
		return
	}
	proj, err := pr.Get(id)
	if err != nil {
		p.log.Warnln("Get project err: ", err.Error())
		c.JSON(http.StatusBadGateway, gin.H{"error": err.Error()})
		return
	}

	var phases []models.Phase
	for _, phase := range proj.Phases {
		phases = append(phases, models.Phase{
			ID:   phase.PhaseID,
			Name: phase.Name,
		})
	}
	projectsResp := models.ProjectResp{
		ProjectID:       proj.ProjectID,
		Title:           proj.Title,
		Description:     proj.Description,
		MediaID:         proj.MediaID,
		Type:            proj.Type.String(),
		BusinessOwner:   proj.BusinessOwner,
		LegacyEntity:    proj.LegacyEntity,
		Cluster:         proj.Cluster,
		Phase:           phases,
		Stage:           proj.Stage.String(),
		Status:          proj.Status,
		Region:          proj.Region,
		ProjectManager:  proj.ProjectManager,
		PipelineManager: proj.PipelineManager,
		OwnerID:         proj.OwnerID,
		Category:        proj.Category.String(),
		Created:         proj.Created,
		Priority:        proj.Priority,
	}

	c.JSON(http.StatusOK, projectsResp)
}

func (p projectHandler) UpdateProject(c *gin.Context) {
	pr := projects.New(p.db.GetDB())
	id, err := strconv.ParseInt(c.Param("id"), 10, 64)
	if err != nil {
		p.log.Warnln("Param err: ", err.Error())
		c.JSON(http.StatusBadGateway, gin.H{"error": err.Error()})
		return
	}

	var projectReq models.ProjectReq
	if err := c.ShouldBindJSON(&projectReq); err != nil {
		p.log.Warnln("bind error")
		c.JSON(http.StatusBadGateway, gin.H{"error": "bind error"})
		return
	}
	updateColums := make(map[string]interface{})
	if projectReq.BusinessOwner != nil {
		updateColums["business_owner"] = *projectReq.BusinessOwner
	}
	if projectReq.Title != nil {
		updateColums["title"] = *projectReq.Title
	}
	if projectReq.Description != nil {
		updateColums["description"] = *projectReq.Description
	}
	if projectReq.MediaID != nil {
		updateColums["media_id"] = *projectReq.MediaID
	}
	if projectReq.Type != nil {
		updateColums["type"] = *projectReq.Type
	}
	if projectReq.Stage != nil {
		updateColums["stage"] = *projectReq.Stage
	}
	if projectReq.LegacyEntity != nil {
		updateColums["legacy_entity"] = *projectReq.LegacyEntity
	}
	if projectReq.Cluster != nil {
		updateColums["cluster"] = *projectReq.Cluster
	}
	if projectReq.OwnerID != nil {
		updateColums["owner_id"] = *projectReq.OwnerID
	}
	if projectReq.OwnerPhoto != nil {
		updateColums["owner_photo"] = *projectReq.OwnerPhoto
	}
	if projectReq.Region != nil {
		updateColums["region"] = *projectReq.Region
	}
	if projectReq.Status != nil {
		updateColums["status"] = *projectReq.Status
	}
	if projectReq.Priority != nil {
		updateColums["priority"] = *projectReq.Priority
	}
	if projectReq.ProjectManager != nil {
		updateColums["project_manager"] = *projectReq.ProjectManager
	}
	if projectReq.PipelineManager != nil {
		updateColums["pipeline_manager"] = *projectReq.PipelineManager
	}
	proj, err := pr.Update(id, updateColums)
	if err != nil {
		p.log.Warnln("Update projects err: ", err.Error())
		c.JSON(http.StatusBadGateway, gin.H{"error": err.Error()})
		return
	}

	projectsResp := models.ProjectResp{
		ProjectID:       proj.ProjectID,
		Title:           proj.Title,
		Description:     proj.Description,
		MediaID:         proj.MediaID,
		Type:            proj.Type.String(),
		BusinessOwner:   proj.BusinessOwner,
		LegacyEntity:    proj.LegacyEntity,
		Cluster:         proj.Cluster,
		Stage:           proj.Stage.String(),
		OwnerID:         proj.OwnerID,
		Created:         proj.Created,
		Region:          proj.Region,
		Status:          proj.Status,
		Priority:        proj.Priority,
		ProjectManager:  proj.ProjectManager,
		PipelineManager: proj.PipelineManager,
	}

	c.JSON(http.StatusOK, projectsResp)
}

func (p projectHandler) CreateProject(c *gin.Context) {
	tr := p.db.GetDB().Begin()
	pr := projects.New(tr)

	var projectReq models.ProjectReq
	if err := c.ShouldBindJSON(&projectReq); err != nil {
		p.log.Warnln("bind error")
		c.JSON(http.StatusBadGateway, gin.H{"error": "bind error"})
		return
	}
	projectEntity := projects.ProjectEntity{}
	if projectReq.BusinessOwner != nil {
		projectEntity.BusinessOwner = *projectReq.BusinessOwner
	}
	if projectReq.Title != nil {
		projectEntity.Title = *projectReq.Title
	} else {
		p.log.Warnln("title is required")
		c.JSON(http.StatusBadGateway, gin.H{"error": "title is required"})
		return
	}
	if projectReq.Description != nil {
		projectEntity.Description = *projectReq.Description
	}
	if projectReq.MediaID != nil {
		projectEntity.MediaID = *projectReq.MediaID
	}
	if projectReq.Type != nil {
		if err := projectEntity.Type.Scan(*projectReq.Type); err != nil {
			c.JSON(http.StatusBadGateway, gin.H{"error": err.Error()})
			return
		}
	}
	if projectReq.Stage != nil {
		if err := projectEntity.Stage.Scan(*projectReq.Stage); err != nil {
			c.JSON(http.StatusBadGateway, gin.H{"error": err.Error()})
			return
		}
	}
	if projectReq.LegacyEntity != nil {
		projectEntity.LegacyEntity = *projectReq.LegacyEntity
	}
	if projectReq.Cluster != nil {
		projectEntity.Cluster = *projectReq.Cluster
	}
	if projectReq.OwnerID != nil {
		projectEntity.OwnerID = *projectReq.OwnerID
	}

	if projectReq.Region != nil {
		projectEntity.Region = *projectReq.Region
	}
	if projectReq.Status != nil {
		projectEntity.Status = *projectReq.Status
	}
	if projectReq.Priority != nil {
		projectEntity.Priority = *projectReq.Priority
	}
	proj, err := pr.Create(projectEntity)
	if err != nil {
		tr.Rollback()
		p.log.Warnln("Create projects err: ", err.Error())
		c.JSON(http.StatusBadGateway, gin.H{"error": err.Error()})
		return
	}

	w := workspace.New(tr)
	workSpaceEntity := workspace.WorkspaceEntity{
		ProjectID: proj.ProjectID,
		Type:      workspace.Legal,
		Title:     "main",
	}
	err = w.Create(&workSpaceEntity)
	if err != nil {
		tr.Rollback()
		p.log.Warnln("Create workspace err: ", err.Error())
		c.JSON(http.StatusBadGateway, gin.H{"error": err.Error()})
		return
	}
	aP := actionPlan.New(tr)
	actionPlanEntity := actionPlan.ActionPlanEntity{
		ProjectID:   proj.ProjectID,
		WorkspaceID: workSpaceEntity.WorkspaceID,
		Title:       "main",
	}
	err = aP.Create(&actionPlanEntity)
	if err != nil {
		tr.Rollback()
		p.log.Warnln("Create actionPlan err: ", err.Error())
		c.JSON(http.StatusBadGateway, gin.H{"error": err.Error()})
		return
	}
	sch := stage.New(tr)
	stages := []stage.StageEntity{
		{ProjectID: proj.ProjectID, Title: "Ideation", Order: 1, WorkspaceID: workSpaceEntity.WorkspaceID, ActionPlanID: actionPlanEntity.ActionPlanID},
		{ProjectID: proj.ProjectID, Title: "Concept", Order: 2, WorkspaceID: workSpaceEntity.WorkspaceID, ActionPlanID: actionPlanEntity.ActionPlanID},
		{ProjectID: proj.ProjectID, Title: "Business case", Order: 3, WorkspaceID: workSpaceEntity.WorkspaceID, ActionPlanID: actionPlanEntity.ActionPlanID},
	}
	stages, err = sch.CreateMany(stages)
	if err != nil {
		tr.Rollback()
		p.log.Warnln("Create schedules err: ", err.Error())
		c.JSON(http.StatusBadGateway, gin.H{"error": err.Error()})
		return
	}

	st := milestone.New(tr)
	milestones := []milestone.MilestoneEntity{
		{StageID: stages[0].StageID, ProjectID: proj.ProjectID, Title: "Research", Order: 1, WorkspaceID: workSpaceEntity.WorkspaceID, ActionPlanID: actionPlanEntity.ActionPlanID, Status: milestone.NewStatus},
		{StageID: stages[0].StageID, ProjectID: proj.ProjectID, Title: "Idea description", Order: 2, WorkspaceID: workSpaceEntity.WorkspaceID, ActionPlanID: actionPlanEntity.ActionPlanID, Status: milestone.NewStatus},
		{StageID: stages[0].StageID, ProjectID: proj.ProjectID, Title: "Team forming", Order: 3, WorkspaceID: workSpaceEntity.WorkspaceID, ActionPlanID: actionPlanEntity.ActionPlanID, Status: milestone.NewStatus},
		{StageID: stages[0].StageID, ProjectID: proj.ProjectID, Title: "Highlevel planning", Order: 4, WorkspaceID: workSpaceEntity.WorkspaceID, ActionPlanID: actionPlanEntity.ActionPlanID, Status: milestone.NewStatus},
		{StageID: stages[0].StageID, ProjectID: proj.ProjectID, Title: "Budget forming", Order: 5, WorkspaceID: workSpaceEntity.WorkspaceID, ActionPlanID: actionPlanEntity.ActionPlanID, Status: milestone.NewStatus},
	}
	milestones, err = st.CreateMany(milestones)
	if err != nil {
		tr.Rollback()
		p.log.Warnln("Create stages err: ", err.Error())
		c.JSON(http.StatusBadGateway, gin.H{"error": err.Error()})
		return
	}
	tr.Commit()

	var stageResp []models.Stage
	for _, sc := range stages {
		sheduleResp := models.Stage{StageID: sc.StageID, Order: sc.Order, Title: sc.Title}
		for _, stg := range milestones {
			if stg.StageID == sheduleResp.StageID {
				sheduleResp.Milestone = append(sheduleResp.Milestone, models.Milestone{
					MilestoneID: stg.MilestoneID,
					Order:       stg.Order,
					StageID:     stg.StageID,
					Status:      stg.Status.String(),
					Title:       stg.Title,
					DateStart:   stg.DateStart,
					DateEnd:     stg.DateStop,
				})
			}
		}
		stageResp = append(stageResp, sheduleResp)
	}

	projectsResp := models.ProjectResp{
		ProjectID:     proj.ProjectID,
		Title:         proj.Title,
		Description:   proj.Description,
		MediaID:       proj.MediaID,
		Type:          proj.Type.String(),
		BusinessOwner: proj.BusinessOwner,
		LegacyEntity:  proj.LegacyEntity,
		Cluster:       proj.Cluster,
		Stage:         proj.Stage.String(),
		OwnerID:       proj.OwnerID,
		Created:       proj.Created,
		Region:        proj.Region,
		Status:        proj.Status,
		Priority:      proj.Priority,
		Template:      &models.ProjectTemplate{Stage: stageResp},
	}

	c.JSON(http.StatusOK, projectsResp)
}

func (p projectHandler) DeleteProject(ctx *gin.Context) {
	project := projects.New(p.db.GetDB())
	id, err := strconv.ParseInt(ctx.Param("id"), 10, 64)
	if err != nil {
		p.log.Warnln("Param err: ", err.Error())
		ctx.JSON(http.StatusBadGateway, gin.H{"error": err.Error()})
		return
	}
	if id <= 0 {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": "wrong id"})
		return
	}

	if err := project.Delete(id); err != nil {
		p.log.Warnln("Can't delete project with err: ", err.Error())
		ctx.JSON(http.StatusBadGateway, gin.H{"error": err.Error()})
		return
	}

	ctx.JSON(http.StatusOK, gin.H{"result": "success"})
}

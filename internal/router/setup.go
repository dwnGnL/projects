package router

import (
	"context"
	"net/http"
	"projects/internal/handlers/actionPlan"
	"projects/internal/handlers/epic"
	"projects/internal/handlers/milestone"
	"projects/internal/handlers/processes"
	"projects/internal/handlers/project"
	"projects/internal/handlers/stage"
	"projects/internal/handlers/task"
	"projects/internal/handlers/template"
	"projects/pkg/config"

	"github.com/gin-gonic/gin"
	"github.com/sirupsen/logrus"
	"go.uber.org/fx"
)

var Module = fx.Options(
	fx.Invoke(
		SetupRouter,
	),
)

type Params struct {
	fx.In
	Lifecycle  fx.Lifecycle
	ActionPlan actionPlan.ActionPlanHandler
	Epic       epic.EpicHandler
	Milestone  milestone.MilestoneHandler
	Stage      stage.StageHandler
	Processes  processes.ProcessesHandler
	Project    project.ProjectHandler
	Task       task.TaskHandler
	Template   template.TemplateHandler
	*logrus.Logger
	*config.Tuner
}

func SetupRouter(params Params) {
	r := gin.Default()

	baseRoute := r.Group("/api/projects")
	baseRoute.GET("", params.Project.GetProjects)
	baseRoute.GET("/:id", params.Project.Project)
	baseRoute.POST("/:id", params.Project.UpdateProject)
	baseRoute.PUT("/new", params.Project.CreateProject)
	baseRoute.DELETE("/delete/:id", params.Project.DeleteProject)

	baseRoute.GET("/:id/template", params.Template.GetTemplates)
	baseRoute.POST("/:id/template", params.Template.UpdateTemplate)

	baseRoute.GET("/milestone/:id", params.Milestone.GetMilestoneByID)
	baseRoute.GET("/milestone", params.Milestone.GetMilestones)
	baseRoute.PUT("/milestone", params.Milestone.CreateMilestone)
	baseRoute.POST("/milestone/:id", params.Milestone.EditMilestone)
	baseRoute.DELETE("/milestone/delete/:id", params.Milestone.DeleteMilestone)

	baseRoute.POST("/acplan/create", params.ActionPlan.CreateActionPlan)
	baseRoute.GET("/acplan/download/:id", params.ActionPlan.DownloadActionPlan)
	baseRoute.GET("/:id/acplan", params.ActionPlan.GetAcPlans)
	baseRoute.DELETE("/acplan/delete/:id", params.ActionPlan.DeleteActionPlan)
	baseRoute.POST("/acplan/update/:id", params.ActionPlan.UpdateActionPlan)

	baseRoute.GET("/stage/:id", params.Stage.GetStageByProjectID)
	baseRoute.GET("/stage/acplan/:id", params.Stage.GetStageByAcPLan)
	baseRoute.PUT("/stage/update/:id", params.Stage.UpdateStage)
	baseRoute.POST("/stage/create", params.Stage.CreateStage)
	baseRoute.DELETE("/stage/:id", params.Stage.DeleteStage)

	baseRoute.PUT("/epic", params.Epic.CreateEpic)
	baseRoute.POST("/epic/read", params.Epic.GetEpics)
	baseRoute.POST("/epic/:id", params.Epic.UpdateEpic)
	baseRoute.DELETE("/epic/:id", params.Epic.DeleteEpic)

	taskRoute := baseRoute.Group("/tasks")
	taskRoute.POST("", params.Task.CreateTask)
	// taskRoute.GET("/:task_id", params.Task.GetTask)
	taskRoute.GET("", params.Task.GetTasksBatch)
	taskRoute.GET("/epic/:epic_id", params.Task.GetTaskByEpic)
	taskRoute.GET("/milestone/:milestone_id", params.Task.GetTaskByMilestone)
	taskRoute.GET("/:id", params.Task.GetTaskByID)
	taskRoute.POST("/:id", params.Task.UpdateTask)
	taskRoute.DELETE("/:id", params.Task.DeleteTask)

	processesRoute := baseRoute.Group("/processes")
	processesRoute.PUT("", params.Processes.CreateProcess)
	processesRoute.GET("", params.Processes.ReadProcesses)
	processesRoute.POST("/:id", params.Processes.UpdateProcess)
	processesRoute.DELETE("/:id", params.Processes.DeleteProcess)

	srv := http.Server{
		Addr:    ":" + params.Config.Main.Port,
		Handler: r,
	}
	params.Lifecycle.Append(
		fx.Hook{
			OnStart: func(_ context.Context) error {
				params.Logger.Info("Application started")
				go srv.ListenAndServe()
				return nil
			},
			OnStop: func(ctx context.Context) error {
				params.Logger.Info("Application stopped")
				return srv.Shutdown(ctx)
			},
		},
	)

}

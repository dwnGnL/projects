package handlers

import (
	"projects/internal/handlers/actionPlan"
	"projects/internal/handlers/epic"
	"projects/internal/handlers/milestone"
	"projects/internal/handlers/processes"
	"projects/internal/handlers/project"
	"projects/internal/handlers/stage"
	"projects/internal/handlers/task"
	"projects/internal/handlers/template"

	"go.uber.org/fx"
)

var Modules = fx.Options(
	actionPlan.Module,
	epic.Module,
	milestone.Module,
	stage.Module,
	processes.Module,
	project.Module,
	task.Module,
	template.Module,
)

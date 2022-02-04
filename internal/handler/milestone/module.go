package milestone

import "go.uber.org/fx"

type MilestoneHandler interface {
}

type Params struct {
	fx.In
}

type milestoneHandler struct {
}

func NewMilestoneHandler(params Params) MilestoneHandler {
	return &milestoneHandler{}
}

package actionplan

import "go.uber.org/fx"

type ActionPlanHandler interface {
}

type Params struct {
	fx.In
}

type actionPlanHandler struct {
}

func NewActionPlanHandler(params Params) ActionPlanHandler {
	return &actionPlanHandler{}
}

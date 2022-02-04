package task

import "go.uber.org/fx"

type TaskHandler interface {
}

type Params struct {
	fx.In
}

type taskHandler struct {
}

func NewTaskHandler(params Params) TaskHandler {
	return &taskHandler{}
}

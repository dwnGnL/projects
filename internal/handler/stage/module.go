package stage

import "go.uber.org/fx"

type StageHandler interface {
}

type Params struct {
	fx.In
}

type stageHandler struct {
}

func NewStageHandler(params Params) StageHandler {
	return &stageHandler{}
}

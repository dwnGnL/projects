package epic

import "go.uber.org/fx"

type EpicHandler interface {
}

type Params struct {
	fx.In
}

type epicHandler struct {
}

func NewEpicHandler(params Params) EpicHandler {
	return &epicHandler{}
}

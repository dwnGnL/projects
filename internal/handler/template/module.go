package template

import "go.uber.org/fx"

type TemplateHandler interface {
}

type Params struct {
	fx.In
}

type templateHandler struct {
}

func NewtemplateHandler(params Params) TemplateHandler {
	return &templateHandler{}
}

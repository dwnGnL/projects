package router

import "go.uber.org/fx"

var Module = fx.Options(
	fx.Invoke(
		SetupRouter,
	),
)

type Params struct {
	fx.In
	Lifecycle fx.Lifecycle
}

func SetupRouter()

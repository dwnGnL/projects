package main

import (
	"projects/internal/handlers"
	"projects/internal/router"

	"projects/pkg"

	"go.uber.org/fx"
)

func main() {
	app := fx.New(
		router.Module,
		pkg.Modules,
		handlers.Modules,
	)
	app.Run()
}

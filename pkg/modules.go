package pkg

import (
	"projects/pkg/config"
	"projects/pkg/db"
	"projects/pkg/logger"

	"go.uber.org/fx"
)

var Modules = fx.Options(
	config.Module,
	db.Module,
	logger.Module,
)

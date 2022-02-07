package logger

import (
	"github.com/sirupsen/logrus"
	"go.uber.org/fx"
)

var Module = fx.Provide(SetupLog)

func SetupLog() *logrus.Logger {
	return logrus.New()
}

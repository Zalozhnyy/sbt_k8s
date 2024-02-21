package logger

import (
	"github.com/mdevilliers/go/env"
	helper "github.com/mdevilliers/go/logger"
	"github.com/rs/zerolog"
)

func NewFromEnvironment(fields map[string]interface{}) zerolog.Logger {
	logLevel := env.FromEnvWithDefaultStr("SBT_K8S_LEVEL", "info")
	useConsole := env.FromEnvWithDefaultBool("SBT_K8S_LOG_USE_CONSOLE", false)
	return helper.New(logLevel, useConsole, fields)
}

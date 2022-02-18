/*
Copyright © 2022 NAME HERE <EMAIL ADDRESS>

*/
package main

import (
	"os"

	"github.com/andyfeller/gh-dependency-report/cmd"
	"go.uber.org/zap"
)

func main() {

	// Initlaize global logger
	loggerConfig := zap.Config{
		Level:            zap.NewAtomicLevelAt(zap.InfoLevel),
		Encoding:         "console",
		EncoderConfig:    zap.NewDevelopmentEncoderConfig(),
		OutputPaths:      []string{"stderr"},
		ErrorOutputPaths: []string{"stderr"},
	}

	logger, _ := loggerConfig.Build()
	defer logger.Sync() // nolint:errcheck // not sure how to errcheck a deferred call like this
	zap.ReplaceGlobals(logger)

	// Instantiate and execute root command
	cmd := cmd.NewCmd()

	if err := cmd.Execute(); err != nil {
		os.Exit(1)
	}
}

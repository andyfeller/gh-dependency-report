/*
Copyright © 2022 NAME HERE <EMAIL ADDRESS>

*/
package main

import (
	"os"

	"github.com/andyfeller/gh-dependency-report/cmd"
	"github.com/andyfeller/gh-dependency-report/internal/log"
	"go.uber.org/zap"
)

func main() {

	// Initlaize global logger
	logger, _ := log.NewLogger(false)
	defer logger.Sync() // nolint:errcheck // not sure how to errcheck a deferred call like this
	zap.ReplaceGlobals(logger)

	// Instantiate and execute root command
	cmd := cmd.NewCmd()

	if err := cmd.Execute(); err != nil {
		os.Exit(1)
	}
}

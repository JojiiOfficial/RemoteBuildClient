package main

import (
	"github.com/RemoteBuild/RemoteBuildClient/commands"
)

// Generates a commands.Commanddata object based on the cli parameter
func buildCData(parsed string, appTrimName int) *commands.CommandData {
	// Command data
	commandData := commands.CommandData{
		Command:   parsed,
		Config:    config,
		Yes:       *appYes,
		Force:     *appForce,
		HideTitel: *appNoTitle,
	}

	// Init cdata
	if !commandData.Init() {
		return nil
	}

	// Initialize encryption sources
	return &commandData
}

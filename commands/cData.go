package commands

import (
	"fmt"

	librb "github.com/JojiiOfficial/LibRemotebuild"
	dmConfig "github.com/JojiiOfficial/LibRemotebuild/config"
	"github.com/JojiiOfficial/gaw"
)

// CommandData data for commands
type CommandData struct {
	LibDM   *librb.LibRB
	Command string
	Config  *dmConfig.Config

	NoRedaction, OutputJSON bool
	Yes, Force, Quiet       bool
}

// Init init CommandData
func (cData *CommandData) Init() bool {
	// Get requestconfig
	// Allow setup, register and login command to continue without
	// handling the error

	var config *librb.RequestConfig
	if cData.Config != nil {
		var err error
		config, err = cData.Config.ToRequestConfig()
		if err != nil && !gaw.IsInStringArray(cData.Command, []string{"setup", "register", "login"}) {
			fmt.Println(err)
			return false
		}
	}

	// Create new dmanager lib object
	cData.LibDM = librb.NewLibRB(config)

	// return success
	return true
}

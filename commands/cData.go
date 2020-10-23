package commands

import (
	"fmt"
	"os"

	librb "github.com/RemoteBuild/LibRemotebuild"
	dmConfig "github.com/RemoteBuild/LibRemotebuild/config"
	"github.com/JojiiOfficial/gaw"
)

// CommandData data for commands
type CommandData struct {
	Librb   *librb.LibRB
	Command string
	Config  *dmConfig.Config

	NoRedaction, OutputJSON bool
	Yes, Force, Quiet       bool
	HideTitel               bool
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
	cData.Librb = librb.NewLibRB(config)

	// return success
	return true
}

func (cData *CommandData) upload_given() bool {
	return gaw.IsInStringArray("--uploadTo", os.Args)
}

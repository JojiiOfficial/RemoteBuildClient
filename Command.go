package main

import (
	"time"

	"github.com/JojiiOfficial/RemoteBuildClient/commands"
)

func runCommand(parsed string, commandData *commands.CommandData) {
	// Execute the desired command
	switch parsed {
	case setupCmd.FullCommand():
		commandData.SetupClient(*setupCmdHost, *appCfgFile, *setupCmdIgnoreCert, *setupCmdServerOnly, *setupCmdRegister, *setupCmdLogin, *setupCmdToken, *setupCmdUser)

	case loginCmd.FullCommand():
		commandData.LoginCommand("")

	case registerCmd.FullCommand():
		commandData.RegisterCommand()

	case jobs.FullCommand():
		commandData.ListJobs()

	case aurBuild.FullCommand():
		commandData.CreateAURJob(*aurbuildPackage, *jobUploadTo)

	case jobCancelCmd.FullCommand():
		commandData.CancelJob(*jobCancelID)

	case jobLogsCmd.FullCommand():
		// Get logs starting 20 sec ago
		start := time.Unix(time.Now().Unix()-20, 0)
		commandData.Logs(*jobLogsID, start)
	}
}

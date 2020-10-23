package main

import (
	"time"

	libremotebuild "github.com/RemoteBuild/LibRemotebuild"
	"github.com/RemoteBuild/RemoteBuildClient/commands"
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
		commandData.ListJobs(*jobsn)

	case aurBuild.FullCommand():
		commandData.CreateAURJob(*aurbuildPackage, *jobUploadTo, *jobDisableCcache)

	case jobCancelCmd.FullCommand():
		commandData.CancelJob(*jobCancelID)

	case jobPauseCmd.FullCommand():
		commandData.SetJobState(*jobPauseID, libremotebuild.JobPaused)

	case jobResumeCmd.FullCommand():
		commandData.SetJobState(*jobResumeID, libremotebuild.JobRunning)

	case jobLogsCmd.FullCommand():
		// Get logs starting 20 sec ago
		start := time.Unix(time.Now().Unix()-20, 0)
		commandData.Logs(*jobLogsID, start, true)

	case logsCmd.FullCommand():
		// Get logs starting 20 sec ago
		start := time.Unix(time.Now().Unix()-20, 0)
		commandData.Logs(*logsID, start, true)

	case ccacheClearCmd.FullCommand():
		commandData.ClearCcache()

	case ccacheInfoCmd.FullCommand():
		commandData.QueryCcache()

	case jobInfo.FullCommand():
		commandData.JobInfo(*jobInfoID)

	}

}

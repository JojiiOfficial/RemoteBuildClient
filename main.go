package main

import (
	"fmt"
	"log"
	"os"

	librb "github.com/JojiiOfficial/LibRemotebuild"
	dmConfig "github.com/JojiiOfficial/LibRemotebuild/config"
	"github.com/JojiiOfficial/gaw"

	kingpin "gopkg.in/alecthomas/kingpin.v2"
)

const (
	appName = "rbuild"
	version = "1.0.0"
)

// ...
const (
	// EnVarPrefix prefix for env vars
	EnVarPrefix = "RBUILD"

	// EnVarPrefix prefix of all used env vars
	EnVarLogLevel   = "LOG_LEVEL"
	EnVarNoColor    = "NO_COLOR"
	EnVarNoEmojis   = "NO_EMOJIS"
	EnVarConfigFile = "CONFIG"
)

// Return the variable using the server prefix
func getEnVar(name string) string {
	return fmt.Sprintf("%s_%s", EnVarPrefix, name)
}

// App commands
var (
	app = kingpin.New(appName, "A Remote build client")

	// Global flags
	appYes     = app.Flag("yes", "Skip confirmations").Short('y').Bool()
	appCfgFile = app.Flag("config", "the configuration file for the app").Envar(getEnVar(EnVarConfigFile)).Short('c').String()

	// File related flags
	appForce = app.Flag("force", "Forces an action").Short('f').Bool()

	setupCmd           = app.Command("setup", "Setup the connection to the server")
	setupCmdHost       = setupCmd.Arg("host", "The host").Required().String()
	setupCmdIgnoreCert = setupCmd.Flag("ignore-cert", "Ignore invalid SSL/TLS certs").Bool()
	setupCmdServerOnly = setupCmd.Flag("server-only", "Setup the server only").Bool()
	setupCmdRegister   = setupCmd.Flag("register", "Create an account").Bool()
	setupCmdLogin      = setupCmd.Flag("login", "Login into an existing account").Bool()
	setupCmdToken      = setupCmd.Flag("token", "Use an existing sessiontoken to setup a connection").String()
	setupCmdUser       = setupCmd.Flag("username", "Required if --token is used").String()

	// User commands
	loginCmd    = app.Command("login", "Login into an existing account")
	registerCmd = app.Command("register", "Create a new account")

	jobs  = app.Command("jobs", "List active jobs").Alias("js")
	jobsn = jobs.Flag("limit", "Limit the jobs to display").Short('n').Int()

	// Job commands
	job = app.Command("job", "Job actions").Alias("j")

	// New jobs
	newJobCmd = job.Command("create", "Create a new job").Alias("c")

	jobUploadTo      = app.Flag("uploadTo", "Upload compiled file").Short('u').HintOptions([]string{librb.DataManagerUploadType.String()}...).String()
	jobDisableCcache = app.Flag("disable-ccache", "Don't use ccache to build the specified package").Bool()

	// -- New AUR job
	aurBuild        = newJobCmd.Command("aurbuild", "Build an AUR package")
	aurbuildPackage = aurBuild.Arg("Package", "The AUR package to build").Required().String()

	// Cancel job
	jobCancelCmd = job.Command("cancel", "Cancel a job").Alias("stop").Alias("rm")
	jobCancelID  = jobCancelCmd.Arg("JobID", "ID of job to cancel").Required().Uint()

	// Job logs
	jobLogsCmd = job.Command("logs", "View logs of job").Alias("l").Alias("log")
	jobLogsID  = jobLogsCmd.Arg("JobID", "ID of job to retrieve the logs from").Required().Uint()

	// logs
	logsCmd = app.Command("logs", "View logs of job").Alias("l").Alias("log")
	logsID  = logsCmd.Arg("JobID", "ID of job to retrieve the logs from").Required().Uint()

	// Pause
	jobPauseCmd = job.Command("pause", "Pause a job")
	jobPauseID  = jobPauseCmd.Arg("JobID", "ID of job to pause").Required().Uint()

	// Resume
	jobResumeCmd = job.Command("resume", "Resume a job")
	jobResumeID  = jobResumeCmd.Arg("JobID", "ID of job to resume").Required().Uint()

	// Ccache
	ccacheCmd = app.Command("ccache", "Ccache commands").Alias("cc").Alias("cache").Alias("c")

	// clear
	ccacheClearCmd = ccacheCmd.Command("clear", "Clear ccache on server").Alias("c")

	// getinfo
	ccacheInfoCmd = ccacheCmd.Command("stats", "Clear ccache on server").Alias("i").Alias("query").Alias("q").Alias("s")
)

var (
	config       *dmConfig.Config
	appTrimName  int
	unmodifiedNS string
)

func main() {
	app.HelpFlag.Short('h')
	app.Version(version)

	// Init random seed from gaw
	gaw.Init()

	// Prase cli flags
	parsed := kingpin.MustParse(app.Parse(os.Args[1:]))

	// Init config
	if !initConfig(parsed) {
		return
	}

	// Bulid commandData
	commandData := buildCData(parsed, appTrimName)
	if commandData == nil {
		return
	}

	// Run desired command
	runCommand(parsed, commandData)
}

// Load and init config. Return false on error
func initConfig(parsed string) bool {
	// Init config
	var err error
	config, err = dmConfig.InitConfig(dmConfig.GetDefaultConfigFile(), *appCfgFile)
	if err != nil {
		log.Fatalln(err)
	}

	if config == nil {
		fmt.Println("New config created")
		if parsed != setupCmd.FullCommand() {
			return false
		}
	}

	return true
}

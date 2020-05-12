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

	jobs = app.Command("jobs", "List active jobs")

	// Job commands
	job = app.Command("job", "Job actions")

	// New jobs
	newJobCmd = job.Command("create", "Create a new job")

	jobUploadTo = app.Flag("uploadTo", "Upload compiled file").HintOptions([]string{librb.DataManagerUploadType.String()}...).String()

	// -- New AUR job
	aurBuild        = newJobCmd.Command("aurbuild", "Build an AUR package")
	aurbuildPackage = aurBuild.Arg("Package", "The AUR package to build").Required().String()
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

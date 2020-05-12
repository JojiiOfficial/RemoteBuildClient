package main

import (
	"fmt"
	"log"
	"os"

	dmConfig "github.com/JojiiOfficial/LibRemotebuild/config"
	"github.com/JojiiOfficial/gaw"

	kingpin "gopkg.in/alecthomas/kingpin.v2"
)

const (
	appName = "manager"
	version = "1.0.0"
)

// ...
const (
	// EnVarPrefix prefix for env vars
	EnVarPrefix = "MANAGER"

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
	app = kingpin.New(appName, "A DataManager")

	// Global flags
	appYes     = app.Flag("yes", "Skip confirmations").Short('y').Bool()
	appCfgFile = app.Flag("config", "the configuration file for the app").Envar(getEnVar(EnVarConfigFile)).Short('c').String()

	// File related flags
	appForce = app.Flag("force", "Forces an action").Short('f').Bool()

	setupCmd = app.Command("setup", "Setup the connection to the server")
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
	defer commandData.CloseKeystore()

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

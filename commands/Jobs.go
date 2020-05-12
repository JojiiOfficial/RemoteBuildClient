package commands

import (
	"encoding/base64"
	"fmt"

	dmConfig "github.com/DataManager-Go/libdatamanager/config"
	librb "github.com/JojiiOfficial/LibRemotebuild"
	"github.com/fatih/color"
	clitable "gopkg.in/benweidig/cli-table.v2"
)

// ListJobs list active jobs
func (cData *CommandData) ListJobs() {
	jobs, err := cData.Librb.ListJobs()
	if err != nil {
		printResponseError(err, "retrieving job list")
		return
	}

	if len(jobs.Jobs) == 0 {
		fmt.Println("No jobs found")
		return
	}

	// Create table
	table := clitable.New()
	table.ColSeparator = " "
	table.Padding = 4

	// Build header
	headingColor := color.New(color.FgHiGreen, color.Underline, color.Bold)
	header := []interface{}{headingColor.Sprint("ID"), headingColor.Sprint("Info"), headingColor.Sprint("Pos"), headingColor.Sprint("Job Type"), headingColor.Sprint("Upload Type"), headingColor.Sprint("Status")}
	table.AddRow(header...)

	// Fill table with data
	for _, job := range jobs.Jobs {
		rowitems := []interface{}{
			job.ID,
			job.Info,
			job.Position,
			job.BuildType,
			job.UploadType,
			job.Status,
		}

		table.AddRow(rowitems...)
	}

	// Print table
	fmt.Println(table)
}

// CreateAURJob create an aur build job
func (cData *CommandData) CreateAURJob(pkg, sUploadType string) {
	uploadtype := librb.ParseUploadType(sUploadType)

	if uploadtype == librb.NoUploadType {
		fmt.Println("Warning: not uploading")
	}

	aurBuild := cData.Librb.NewAURBuild(pkg)
	if uploadtype == librb.DataManagerUploadType {
		conf, err := initDMConfig()
		if err != nil {
			fmt.Println("Can't read to DM config:", err)
			return
		}

		// Get token from Dmanager Config/Keyring
		token, err := conf.GetToken()
		if err != nil {
			fmt.Println("Can't retrieve DataManager token. Are you logged in?")
			return
		}

		// Use Dmanager data
		aurBuild.WithDmanager(conf.User.Username, base64.RawStdEncoding.EncodeToString([]byte(token)), conf.Server.URL)
	}

	// Create job
	resp, err := aurBuild.CreateJob()
	if err != nil {
		printResponseError(err, "creating AUR build job")
		return
	}

	printSuccess("created job with ID: %d at Pos %d", resp.ID, resp.Position)
}

// Load and init config. Return false on error
func initDMConfig() (*dmConfig.Config, error) {
	return dmConfig.InitConfig(dmConfig.GetDefaultConfigFile(), "")
}

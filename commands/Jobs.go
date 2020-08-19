package commands

import (
	"encoding/base64"
	"fmt"
	"io"
	"os"
	"strconv"
	"strings"
	"time"

	dmConfig "github.com/DataManager-Go/libdatamanager/config"
	librb "github.com/JojiiOfficial/LibRemotebuild"
	"github.com/fatih/color"
	clitable "gopkg.in/benweidig/cli-table.v2"
)

// ListJobs list active jobs
func (cData *CommandData) ListJobs(limit int) {
	jobs, err := cData.Librb.ListJobs(limit)
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
	var jobWithPos bool

	if !cData.HideTitel {
		header := []interface{}{headingColor.Sprint("ID"), headingColor.Sprint("Info")}
		if hasJobWithPos(jobs.Jobs) {
			header = append(header, headingColor.Sprint("Pos"))
			jobWithPos = true
		}
		header = append(header, []interface{}{headingColor.Sprint("Job Type"), headingColor.Sprint("Upload Type"), headingColor.Sprint("Status"), headingColor.Sprint("Duration")}...)

		table.AddRow(header...)
	}

	// Fill table with data
	for _, job := range jobs.Jobs {
		rowitems := []interface{}{
			job.ID,
			job.Info,
		}

		if jobWithPos {
			if job.Position > 0 {
				rowitems = append(rowitems, job.Position)
			} else {
				rowitems = append(rowitems, "-")
			}
		}

		rowitems = append(rowitems, []interface{}{
			job.BuildType,
			job.UploadType,
		}...)

		if job.Status == librb.JobRunning {
			rowitems = append(rowitems, "Running")
			rowitems = append(rowitems, time.Since(job.RunningSince))
		} else {
			rowitems = append(rowitems, job.Status)
			if job.Duration.Seconds() > 0 {
				rowitems = append(rowitems, job.Duration.String())
			} else {
				rowitems = append(rowitems, "-")
			}
		}

		table.AddRow(rowitems...)
	}

	// Print table
	fmt.Println(table)
}

// CancelJob cancel a job
func (cData *CommandData) CancelJob(jobID uint) {
	if err := cData.Librb.CancelJob(jobID); err != nil {
		printResponseError(err, "canceling job")
		return
	}

	printSuccess("%s %d", "cancelled job", jobID)
}

// CreateAURJob create an aur build job
func (cData *CommandData) CreateAURJob(pkg, sUploadType string, disableCcache bool) {
	if len(sUploadType) == 0 && len(cData.Config.DefaultUploadTo) > 0 {
		sUploadType = cData.Config.DefaultUploadTo
	}
	uploadtype := librb.ParseUploadType(sUploadType)

	if uploadtype == librb.NoUploadType {
		fmt.Println("Warning: not uploading")
	}

	// create aurbuild
	aurBuild := cData.Librb.NewAURBuild(pkg)

	// Disable ccache if desired
	if disableCcache {
		aurBuild.WithoutCcache()
	}

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
		aurBuild.WithDmanager(conf.User.Username, base64.RawStdEncoding.EncodeToString([]byte(token)), conf.Server.URL, cData.Config.GetNamspace(librb.JobAUR))
	}

	// Create job
	resp, err := aurBuild.CreateJob()
	if err != nil {
		printResponseError(err, "creating AUR build job")
		return
	}

	printSuccess("created job with ID: %d at Pos %d", resp.ID, resp.Position)
}

// SetJobState to paused or running
func (cData *CommandData) SetJobState(jobID uint, state librb.JobState) {
	err := cData.Librb.SetJobState(jobID, state)
	if err != nil {
		printResponseError(err, "Setting state")
	}
}

// Logs of job
func (cData *CommandData) Logs(jobID uint, since time.Time, first bool) {
	resp, err := cData.Librb.Logs(jobID, since)
	if err != nil {
		// If job exits while retrieving logs
		if !first && strings.Contains(strings.ToLower(err.Error()), "job not found") {
			return
		}

		printResponseError(err, "retrieving logs")
		return
	}

	// Parse response time
	reqTime, err := strconv.ParseInt(resp.Message, 10, 64)
	if err != nil {
		printError(err, "parsing response")
		return
	}

	// Display logs
	_, err = io.Copy(os.Stdout, resp.Response.Body)
	if err != nil && err != io.EOF {
		fmt.Println("ERR:", err)
		return
	}

	// Server said no more logs
	if reqTime == -1 {
		return
	}

	// wait
	time.Sleep(300 * time.Millisecond)

	// Display new logs
	cData.Logs(jobID, time.Unix(reqTime, 0), false)
}

// Load and init config. Return false on error
func initDMConfig() (*dmConfig.Config, error) {
	return dmConfig.InitConfig(dmConfig.GetDefaultConfigFile(), "")
}

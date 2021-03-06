package commands

import (
	"bufio"
	"encoding/base64"
	"encoding/json"
	"fmt"
	"log"
	"net/url"
	"os"

	librb "github.com/RemoteBuild/LibRemotebuild"
	dmConfig "github.com/RemoteBuild/LibRemotebuild/config"
	"github.com/JojiiOfficial/configService"
	"github.com/JojiiOfficial/gaw"
	"github.com/fatih/color"
)

// ConfigView view config
func ConfigView(cData *CommandData, sessionBase64 bool) {
	token, err := cData.Config.GetToken()
	if err != nil {
		token = cData.Config.User.SessionToken
	}

	if cData.NoRedaction && sessionBase64 {
		token = base64.RawStdEncoding.EncodeToString([]byte(token))
	} else if sessionBase64 {
		token = "<redacted>"
	}

	cData.Config.User.SessionToken = token

	if !cData.OutputJSON {
		// Print human output
		fmt.Println(cData.Config.View(!cData.NoRedaction))
	} else {
		// Redact secrets
		if !cData.NoRedaction {
			cData.Config.User.SessionToken = "<redacted>"
		}

		// Make json
		b, err := json.Marshal(cData.Config)
		if err != nil {
			fmt.Println(err)
			return
		}

		fmt.Println(string(b))
	}
}

// SetupClient sets up client config
func (cData *CommandData) SetupClient(host, configFile string, ignoreCert, serverOnly, register, noLogin bool, token, username string) {
	if len(token)*len(username) == 0 && len(token)+len(username) > 0 {
		fmt.Println("Either --user or --token is missing")
		return
	}

	// Confirm creating a config anyway
	if cData.Config != nil && !cData.Config.IsDefault() && !cData.Yes {
		y, _ := gaw.ConfirmInput("There is already a config. Do you want to overwrite it? [y/n]> ", bufio.NewReader(os.Stdin))
		if !y {
			return
		}
	}

	// Load config
	if cData.Config == nil {
		var err error
		cData.Config, err = dmConfig.InitConfig(dmConfig.GetDefaultConfigFile(), configFile)
		if err != nil {
			printError("loading config", err.Error())
			return
		}
	}

	u := bulidURL(host)

	// Check host and verify response
	if err := checkHost(u.String(), ignoreCert); err != nil {
		printError("checking host", err.Error())
		return
	}

	fmt.Printf("%s connected to server\n", color.HiGreenString("Succesfully"))

	// Set new config values
	cData.Config.Server.URL = u.String()
	cData.Config.Server.IgnoreCert = ignoreCert

	err := configService.Save(cData.Config, cData.Config.File)
	if err != nil {
		printError("saving config", err.Error())
		return
	}

	// If severonly mode is requested, stop here
	if serverOnly {
		return
	}

	// Initialize server connection library instance
	// ignore token error since user might not
	// be logged in after setup process
	config, _ := cData.Config.ToRequestConfig()
	cData.Librb = librb.NewLibRB(config)

	// Insert user directly if token and user is set
	if len(token) > 0 && len(username) > 0 {
		// Decode token
		dec, err := base64.RawStdEncoding.DecodeString(token)
		if err != nil {
			fmt.Println(err)
			return
		}

		token = string(dec)
		cData.Config.InsertUser(username, token)
		cData.Config.Save()
		return
	}

	// In register mode, don't login
	if register {
		noLogin = true
	}

	// if not noLogin, login
	if !noLogin {
		fmt.Println("Login")
		cData.LoginCommand("")
		return
	}

	if register {
		fmt.Println("Create an account")
		cData.RegisterCommand()
	}
}

func bulidURL(host string) *url.URL {
	u, err := url.Parse(host)
	if err != nil {
		log.Fatal(err)
	}

	if u.Scheme == "" {
		u.Scheme = "https"
	}

	// Validate scheme
	if !gaw.IsInStringArray(u.Scheme, []string{"http", "https"}) {
		log.Fatalf("Invalid scheme '%s'. Use http or https\n", u.Scheme)
	}

	return u
}

func checkHost(host string, ignoreCert bool) error {
	rb := librb.NewLibRB(&librb.RequestConfig{
		IgnoreCert: ignoreCert,
		URL:        host,
	})

	_, err := rb.Ping()
	if err != nil {
		return err
	}

	return nil
}

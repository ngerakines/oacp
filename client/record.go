package client

import (
	"fmt"
	"github.com/urfave/cli"
	"net/http"
	"net/url"
	"strings"
	"time"
)

var serverFlag = cli.StringFlag{
	Name:     "server",
	Usage:    "The server to interact with.",
	EnvVar:   "SERVER",
	Value:    "",
	Required: true,
}

var apiUserFlag = cli.StringFlag{
	Name:   "api-user",
	Usage:  "The API user to authenticate as.",
	EnvVar: "API_USER",
	Value:  "aocp",
}

var apiPasswordFlag = cli.StringFlag{
	Name:   "api-password",
	Usage:  "The password for the API user",
	EnvVar: "API_PASSWORD",
	Value:  "aocp",
}

var stateFlag = cli.StringFlag{
	Name:     "state",
	Usage:    "The state to record.",
	Value:    "",
	Required: true,
}

var locationFlag = cli.StringFlag{
	Name:     "location",
	Usage:    "The location to record.",
	Value:    "",
	Required: true,
}

var RecordCommand = cli.Command{
	Name:  "record",
	Usage: "Record a location with a state.",
	Flags: []cli.Flag{
		serverFlag,
		apiUserFlag,
		apiPasswordFlag,
		stateFlag,
		locationFlag,
	},
	Action: recordCommandAction,
}

func recordCommandAction(cliCtx *cli.Context) error {
	server := cliCtx.String("server")

	serverUrl, err := url.Parse(server)
	if err != nil {
		return err
	}
	serverUrl.Path = "/api/locations"

	server = serverUrl.String()

	user := cliCtx.String("api-user")
	password := cliCtx.String("api-password")
	state := cliCtx.String("state")
	location := cliCtx.String("location")

	client := &http.Client{
		Timeout: 10 * time.Second,
	}

	form := url.Values{}
	form.Add("state", state)
	form.Add("location", location)

	req, err := http.NewRequest("POST", server, strings.NewReader(form.Encode()))
	if err != nil {
		return err
	}

	req.Header.Set("User-Agent", userAgent(cliCtx))
	req.Header.Set("Content-Type", "application/x-www-form-urlencoded")
	if len(user) > 0 && len(password) > 0 {
		req.SetBasicAuth(user, password)
	}

	resp, err := client.Do(req)
	if err != nil {
		return err
	}
	if resp.StatusCode != http.StatusCreated {
		return fmt.Errorf("unexpected status code: %d", resp.StatusCode)
	}

	return nil
}

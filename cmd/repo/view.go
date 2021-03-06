// Copyright (c) 2020 Target Brands, Inc. All rights reserved.
//
// Use of this source code is governed by the LICENSE file in this repository.

package repo

import (
	"encoding/json"
	"fmt"

	"github.com/go-vela/sdk-go/vela"

	"github.com/urfave/cli"
	yaml "gopkg.in/yaml.v2"
)

// ViewCmd defines the command for viewing a repository.
var ViewCmd = cli.Command{
	Name:        "repo",
	Description: "Use this command to view a repository.",
	Usage:       "View details of the provided repository",
	Action:      view,
	Before:      validate,
	Flags: []cli.Flag{

		// required flags to be supplied to a command
		cli.StringFlag{
			Name:   "org",
			Usage:  "Provide the organization for the repository",
			EnvVar: "REPO_ORG",
		},
		cli.StringFlag{
			Name:   "repo",
			Usage:  "Provide the repository contained with the organization",
			EnvVar: "REPO_NAME",
		},

		// optional flags that can be supplied to a command
		cli.StringFlag{
			Name:  "output,o",
			Usage: "Print the output in json format",
		},
	},
	CustomHelpTemplate: fmt.Sprintf(`%s
EXAMPLES:
 1. View repository details.
    $ {{.HelpName}} --org github --repo octocat
 2. View repository details with json output.
    $ {{.HelpName}} --org github --repo octocat --output json
 3. View repository details when org and repo config or environment variables are set.
    $ {{.HelpName}}
`, cli.CommandHelpTemplate),
}

// helper function to execute logs cli command
func view(c *cli.Context) error {
	// get org and repo information from cmd flags
	org, repo := c.String("org"), c.String("repo")

	// create a carval client
	client, err := vela.NewClient(c.GlobalString("addr"), nil)
	if err != nil {
		return err
	}

	// set token from global config
	client.Authentication.SetTokenAuth(c.GlobalString("token"))

	repository, _, err := client.Repo.Get(org, repo)
	if err != nil {
		return err
	}

	switch c.String("output") {
	case "json":
		output, err := json.MarshalIndent(repository, "", "    ")
		if err != nil {
			return err
		}

		fmt.Println(string(output))
	default:
		// default output should contain all resources fields
		output, err := yaml.Marshal(repository)
		if err != nil {
			return err
		}

		fmt.Println(string(output))
	}

	return nil
}

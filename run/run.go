package run

import (
	"fmt"
	"os"
	"strings"
	"syscall"

	"github.com/HPI-BP2015H/go-travis/client"
	_ "github.com/HPI-BP2015H/go-travis/commands" // import commands
	"github.com/HPI-BP2015H/go-travis/config"
	"github.com/HPI-BP2015H/go-utils/cli"
	"github.com/HPI-BP2015H/go-utils/pathname"
)

// Run executes the command specified via os.Args.
// Takes a API client constructor as argument to enable passing a fake API.
func Run(clientConstructor func(string, string, bool) *client.Client) int {

	app := cli.AppInstance()
	app.Name = "go-travis"
	app.DefaultCommandName = "help"

	app.RegisterFlag(
		cli.Flag{
			Short: "-r",
			Long:  "--repo",
			Ftype: "REPOSITORY_SLUG",
			Help:  "the repository on GitHub",
		},
	)
	app.RegisterFlag(
		cli.Flag{
			Short: "-t",
			Long:  "--token",
			Ftype: "ACCESS_TOKEN",
			Help:  "access token to use",
		},
	)
	app.RegisterFlag(
		cli.Flag{
			Short: "-e",
			Long:  "--api-endpoint",
			Ftype: "URL",
			Help:  "Travis API server to talk to",
		},
	)
	app.RegisterFlag(
		cli.Flag{
			Long:  "--org",
			Ftype: false,
			Help:  "short-cut for --api-endpoint 'https://api.travis-ci.com/'",
		},
	)
	app.RegisterFlag(
		cli.Flag{
			Long:  "--pro",
			Ftype: false,
			Help:  "short-cut for --api-endpoint 'https://api.travis-ci.com/'",
		},
	)
	app.RegisterFlag(
		cli.Flag{
			Long:  "--staging",
			Ftype: false,
			Help:  "short-cut for --api-endpoint 'https://api-staging.travis-ci.org/'",
		},
	)
	app.RegisterFlag(
		cli.Flag{
			Long:  "--debug",
			Ftype: false,
			Help:  "show API requests",
		},
	)
	app.RegisterFlag(
		cli.Flag{
			Short: "-h",
			Long:  "--help",
			Ftype: false,
			Help:  "show help for the command",
		},
	)

	app.Before = func(cmd *cli.Cmd, cmdName string) string {
		configuration := config.DefaultConfiguration()

		endpoint := configuration.GetDefaultTravisEndpoint()
		if cmd.Parameters.IsProvided("--org") {
			endpoint = config.TravisOrgEndpoint
		}
		if cmd.Parameters.IsProvided("--pro") {
			endpoint = config.TravisProEndpoint
		}
		if cmd.Parameters.IsProvided("--staging") {
			endpoint = config.TravisStagingEndpoint
		}
		endpoint = cmd.Parameters.String("--api-endpoint", endpoint)

		debug := cmd.Parameters.IsProvided("--debug")
		token := cmd.Parameters.String("--token", configuration.GetTravisTokenForEndpoint(endpoint))

		commandConfig := config.TravisCommandConfig{
			Config:   configuration,
			Repo:     cmd.Parameters.String("--repo", config.RepoSlugFromGit()),
			Endpoint: endpoint,
			Token:    token,
			Client:   clientConstructor(endpoint, token, debug),
			Debug:    debug,
		}
		cmd.Env = commandConfig

		if cmd.Parameters.IsProvided("--help") {
			args := []string{cmd.Args.ProgramName(), "help", cmdName}
			args = append(args, cmd.Args.Slice(1)...)
			cmd.Args = cli.NewArgs(args)
			return "help"
		}
		return ""
	}

	app.Fallback = func(cmd *cli.Cmd, cmdName string) int {
		env := cmd.Env.(config.TravisCommandConfig)

		os.Setenv("TRAVIS_REPO", env.Repo)
		os.Setenv("TRAVIS_TOKEN", env.Token)
		os.Setenv("TRAVIS_ENDPOINT", env.Endpoint)
		if env.Debug {
			os.Setenv("TRAVIS_DEBUG", "true")
		}

		exeName := cmd.Args.ProgramName() + "-" + cmdName
		results := pathname.FindInPath(exeName, strings.Split(os.Getenv("PATH"), ":"))

		if len(results) > 0 {
			exeCmd := results[0]

			argv := []string{exeName}
			argv = append(argv, cmd.Args.Slice(1)...)

			err := syscall.Exec(exeCmd.String(), argv, os.Environ())
			if err != nil {
				fmt.Fprintf(os.Stderr, "%s: %s\n", exeName, err)
				//os.Exit(1)
				return 1
			}
		} else {
			fmt.Fprintf(os.Stderr, "%s: command not found\n", exeName)
			//os.Exit(1)
			return 1
		}
		return 0
	}

	return app.Run(os.Args)

}
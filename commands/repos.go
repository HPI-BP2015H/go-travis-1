package commands

import (
	"fmt"

	"github.com/HPI-BP2015H/go-travis/client"
	"github.com/HPI-BP2015H/go-travis/config"
	"github.com/HPI-BP2015H/go-utils/cli"
)

func init() {
	cli.AppInstance().RegisterCommand(
		cli.Command{
			Name:     "repos",
			Info:     "lists repositories the user has certain permissions on",
			Function: reposCmd,
		},
	)
}

func reposCmd(cmd *cli.Cmd) cli.ExitValue {
	if NotLoggedIn(cmd) {
		return cli.Failure
	}
	env := cmd.Env.(config.TravisCommandConfig)
	repositories, err := GetAllRepositories(nil, env.Client)
	if err != nil {
		cmd.Stderr.Println("Error: Could not get Repositories.")
		return cli.Failure
	}
	for _, repo := range repositories.Repositories {
		printRepo(repo, cmd)
	}
	return cli.Success
}

//GetAllRepositories returns all the repositories (also those not active in travis)
//of the currently logged in user. also takes params
func GetAllRepositories(params map[string]string, client client.Client) (Repositories, error) {
	if params == nil {
		params = map[string]string{}
	}
	repositories := Repositories{}
	res, err := client.PerformAction("repositories", "for_current_user", params)
	defer res.Body.Close()
	if err != nil {
		return repositories, err
	}
	if res.StatusCode > 299 {
		return repositories, fmt.Errorf("Error: Unexpected HTTP status: %d\n", res.StatusCode)
	}
	res.Unmarshal(&repositories)
	return repositories, nil
}

func printRepo(repo Repository, cmd *cli.Cmd) {
	if repo.Active {
		cmd.Stdout.Cprint("boldgreen", repo.Slug)
		cmd.Stdout.Cprintf("%C(green) (active: %v, private: %v)%C(reset)\n", repo.Active, repo.Private)
	} else {
		cmd.Stdout.Cprint("boldyellow", repo.Slug)
		cmd.Stdout.Cprintf("%C(yellow) (active: %v, private: %v)%C(reset)\n", repo.Active, repo.Private)
	}
	if repo.HasDescription() {
		cmd.Stdout.Cprint("bold", "Description: ")
		cmd.Stdout.Println(repo.Description)
	}
	cmd.Stdout.Println()
}

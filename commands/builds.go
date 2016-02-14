package commands

import (
	"fmt"

	"github.com/mislav/go-travis/client"
	"github.com/mislav/go-travis/config"
	"github.com/mislav/go-utils/cli"
)

func init() {
	cli.Register("builds", buildsCmd)
}

type Builds struct {
	Builds []Build `json:"builds"`
}

type Build struct {
	Number string  `json:"number"`
	State  string  `json:"state"`
	Branch *Branch `json:"branch"`
}

type Branch struct {
	Name string `json:"name"`
}

func buildsCmd(args *cli.Args) {
	params := map[string]string{
		"repository.slug":  "github/hub",
		"build.event_type": "push",
		"limit":            "10",
	}

	res, err := client.Travis.PerformAction("builds", "find", params)
	if err != nil {
		panic(err)
	}

	builds := Builds{}
	res.Unmarshal(&builds)

	for _, build := range builds.Builds {
		fmt.Printf("#%s: %s (%s)\n", build.Number, build.State, build.Branch.Name)
	}
}

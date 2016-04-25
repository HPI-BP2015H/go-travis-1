package commands

import (
	"fmt"
	"sort"
	"strconv"

	"github.com/HPI-BP2015H/go-utils/cli"
)

func init() {
	cli.AppInstance().RegisterCommand(
		cli.Command{
			Name:     "help",
			Info:     "helps you out when in dire need of information",
			Function: helpCmd,
		},
	)
}

type commandByLength []cli.Command

func (s commandByLength) Len() int {
	return len(s)
}
func (s commandByLength) Swap(i, j int) {
	s[i], s[j] = s[j], s[i]
}
func (s commandByLength) Less(i, j int) bool {
	return len(s[i].Name) > len(s[j].Name)
}

type commandByName []cli.Command

func (s commandByName) Len() int {
	return len(s)
}
func (s commandByName) Swap(i, j int) {
	s[i], s[j] = s[j], s[i]
}
func (s commandByName) Less(i, j int) bool {
	return s[i].Name < s[j].Name
}

type flagByLong []cli.Flag

func (s flagByLong) Len() int {
	return len(s)
}
func (s flagByLong) Swap(i, j int) {
	s[i], s[j] = s[j], s[i]
}
func (s flagByLong) Less(i, j int) bool {
	return len(s[i].Long) > len(s[j].Long)
}

type flagByLength []cli.Flag

func (s flagByLength) Len() int {
	return len(s)
}
func (s flagByLength) Swap(i, j int) {
	s[i], s[j] = s[j], s[i]
}
func (s flagByLength) Less(i, j int) bool {
	return flagLen(s[i]) > flagLen(s[j])
}

func helpCmd(cmd *cli.Cmd) cli.ExitValue {
	args := cmd.Args.SubcommandArgs("help")
	cmdName := args.Peek(0)
	if cmdName == "" {
		printGlobalHelp(cmd)
	} else {
		commands := cli.AppInstance().Commands()
		if command, ok := commands[cmdName]; ok {
			printCommandHelp(command, cmd)
		} else {
			cmd.Stderr.Println("Command " + cmdName + " not found!")
			printGlobalHelp(cmd)
			return cli.Failure
		}
	}
	return cli.Success
}

func printGlobalHelp(cmd *cli.Cmd) {
	cmd.Stdout.Printf("Usage: %s COMMAND [OPTIONS]\n\n", cmd.Args.ProgramName())
	cmd.Stdout.Println("Available commands:")
	printCommands(commands(), cmd.Stdout)
	cmd.Stdout.Println("Available options:")
	printFlagsHelp(globalOptions(), cmd.Stdout)
	cmd.Stdout.Println("Run travis help COMMAND for more infos.")
}

func printCommandHelp(command cli.Command, cmd *cli.Cmd) {
	cmd.Stdout.Println(command.Info)
	if command.Help != "" {
		cmd.Stdout.Println(command.Help)
	}
	cmd.Stdout.Printf("Usage: %s %s [OPTIONS]\n\n", cmd.Args.ProgramName(), command.Name)
	cmd.Stdout.Println("Available options:")
	printFlagsHelp(commandOptions(&command), cmd.Stdout)
	cmd.Stdout.Println("Global options:")
	printFlagsHelp(globalOptions(), cmd.Stdout)
}

func printCommands(commands []cli.Command, out *cli.ColoredWriter) {
	sort.Sort(commandByLength(commands))
	maxLength := len(commands[0].Name)
	sort.Sort(commandByName(commands))

	out.Println()
	for _, command := range commands {
		format := "\t%-" + strconv.Itoa(maxLength+3) + "s"
		out.Printf(format, command.Name)
		out.Cprintln("yellow", command.Info)
	}
	out.Println()
}

func printFlagsHelp(flags []cli.Flag, out *cli.ColoredWriter) {
	if len(flags) == 0 {
		return
	}
	sort.Sort(flagByLength(flags))
	maxLength := flagLen(flags[0])
	sort.Sort(flagByLong(flags))

	out.Println()
	for _, flag := range flags {
		out.Print("\t")
		if flag.Short != "" {
			out.Print(flag.Short + ", ")
		} else {
			out.Print("    ")
		}
		if flag.Ftype != false {
			output := fmt.Sprintf("%v [%v]", flag.Long, flag.Ftype)
			format := "%-" + strconv.Itoa(maxLength+3) + "s"
			out.Printf(format, output)
		} else {
			format := "%-" + strconv.Itoa(maxLength+3) + "s"
			out.Printf(format, flag.Long)
		}
		out.Cprintln("yellow", flag.Help)
	}
	out.Println()
}

func commands() []cli.Command {
	app := cli.AppInstance()
	commands := app.Commands()
	result := make([]cli.Command, 0, len(commands))
	for _, command := range commands {
		result = append(result, command)
	}
	return result
}

func globalOptions() []cli.Flag {
	return flagMapToArray(cli.AppInstance().Flags())
}

func commandOptions(command *cli.Command) []cli.Flag {
	return flagMapToArray(command.Flags())
}

func flagMapToArray(flags map[string]cli.Flag) []cli.Flag {
	result := make([]cli.Flag, 0, len(flags))
	for _, flag := range flags {
		result = append(result, flag)
	}
	return result
}

func flagLen(flag cli.Flag) int {
	result := len(flag.Long)
	if flag.Ftype != false {
		result += len(fmt.Sprintf(" [%v]", flag.Ftype))
	}
	return result
}

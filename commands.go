package main

import (
	"github.com/devopracy/devopracy-cli/command"
	"github.com/mitchellh/cli"
)

// Commands is the mapping of all the available Packer commands.
var Commands map[string]cli.CommandFactory

// CommandMeta is the Meta to use for the commands. This must be written
// before the CLI is started.
var CommandMeta *command.Meta

const ErrorPrefix = "e:"
const OutputPrefix = "o:"

func init() {
	Commands = map[string]cli.CommandFactory{

		"plugin": func() (cli.Command, error) {
			return &command.PluginCommand{
				Meta: *CommandMeta,
			}, nil
		},
	}
}

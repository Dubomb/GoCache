package parser

import (
	"fmt"
	"strings"
)

func ParseCommand(command string) (Command, error) {
	parts := strings.Fields(command)

	if len(parts) == 0 {
		return Command{Type: UnknownCommand}, fmt.Errorf("no input provided")
	}

	parsedCommand := Command{
		Type: ToCommandType(parts[0]),
		Args: parts[1:],
	}

	switch parsedCommand.Type {
	case GetCommand, ExistsCommand, DelCommand:
		if len(parsedCommand.Args) > 0 {
			parsedCommand.Key = parsedCommand.Args[0]
		} else {
			return parsedCommand, fmt.Errorf("%s command missing key argument", parsedCommand.Type)
		}
	case SetCommand:
		if len(parsedCommand.Args) > 1 {
			parsedCommand.Key = parsedCommand.Args[0]
			parsedCommand.Value = parsedCommand.Args[1]
		} else {
			return parsedCommand, fmt.Errorf("%s command missing key and value argument", parsedCommand.Type)
		}
	case UnknownCommand:
		return parsedCommand, fmt.Errorf("unknown command %s", parts[0])
	}

	return parsedCommand, nil
}

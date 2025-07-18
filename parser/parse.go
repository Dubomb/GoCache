package parser

import (
	"fmt"
	"strconv"
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
		if len(parsedCommand.Args) == 4 {
			if parsedCommand.Args[2] != "EX" && parsedCommand.Args[2] != "PX" {
				return parsedCommand, fmt.Errorf("%s command has incorrect TTL argument (must be EX or PX)", parsedCommand.Type)
			}

			if val, err := strconv.ParseUint(parsedCommand.Args[3], 10, 32); err != nil {
				return parsedCommand, fmt.Errorf("%s command TTL is out of range or an invalid integer", parsedCommand.Type)
			} else {
				parsedCommand.TTL = uint(val)
			}

			if parsedCommand.Args[2] == "EX" {
				parsedCommand.TTL *= 1000
			}

			parsedCommand.Key = parsedCommand.Args[0]
			parsedCommand.Value = parsedCommand.Args[1]

		} else if len(parsedCommand.Args) == 2 {
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

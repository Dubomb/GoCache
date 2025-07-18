package parser

import "strings"

type CommandType string

const (
	SetCommand     CommandType = "SET"
	GetCommand     CommandType = "GET"
	DelCommand     CommandType = "DEL"
	ExistsCommand  CommandType = "EXISTS"
	UnknownCommand CommandType = "UNKNOWN"
)

type Command struct {
	Type  CommandType
	Key   string
	Value string
	TTL   uint // Milliseconds
	Args  []string
}

func ToCommandType(s string) CommandType {
	upper := CommandType(strings.ToUpper(s))

	switch upper {
	case SetCommand:
		return SetCommand
	case GetCommand:
		return GetCommand
	case DelCommand:
		return DelCommand
	case ExistsCommand:
		return ExistsCommand
	default:
		return UnknownCommand
	}
}

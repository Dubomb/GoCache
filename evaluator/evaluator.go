package evaluator

import (
	"fmt"
	"gocache/cache"
	"gocache/parser"
)

func Evaluate(command parser.Command, cache *cache.GoCache) string {
	switch command.Type {
	case "SET":
		cache.Set(command.Key, command.Value)

		return "OK"

	case "GET":
		val, exists := cache.Get(command.Key)

		if exists {
			return fmt.Sprintf("VALUE: %s", val)
		}
		return "NOT FOUND"

	case "DEL":
		cache.Del(command.Key)
		return "OK"

	case "EXISTS":
		if cache.Exists(command.Key) {
			return "1"
		} else {
			return "0"
		}

	default:
		return "Error: unknown command"
	}
}

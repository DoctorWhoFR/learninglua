package main

import (
	"fmt"

	lua "github.com/yuin/gopher-lua"
)

func executeMethod(l *lua.LState, state *PlayerState, moduleKey string, f func(*lua.LState, *PlayerState) int) int {
	if !state.AcceptedModules[moduleKey].Owned {
		fmt.Printf("[red]Error, module not owner used: '%s'", moduleKey)
		return 0
	}

	return f(l, state)
}

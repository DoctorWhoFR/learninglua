package main

import (
	"fmt"

	lua "github.com/yuin/gopher-lua"
)

func LoaderTesting(L *lua.LState, state *PlayerState) int {
	// register functions to the table
	mod := L.SetFuncs(L.NewTable(), map[string]lua.LGFunction{
		"coucou": func(l *lua.LState) int {
			return executeMethod(l, state, "test", coucou)
		},
		"toto": func(l *lua.LState) int {
			return executeMethod(l, state, "test", toto)
		},
	})

	// register other stuff
	L.SetField(mod, "name", lua.LString("value"))

	// returns the module
	L.Push(mod)
	return 1
}

func coucou(L *lua.LState, state *PlayerState) int {
	state_Module := state.AcceptedModules["strings"]

	test := lua.LBool(state_Module.Owned)
	L.Push(test)
	return 1
}

func toto(L *lua.LState, state *PlayerState) int {
	fmt.Println("test")
	return 0
}

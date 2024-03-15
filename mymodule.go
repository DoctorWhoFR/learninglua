package main

import (
	lua "github.com/yuin/gopher-lua"
)

func Loader(L *lua.LState, state *PlayerState) int {
	// register functions to the table
	mod := L.SetFuncs(L.NewTable(), map[string]lua.LGFunction{
		"isAuthorized": func(l *lua.LState) int {
			return isAuthorized(l, state)
		},
		"buyModule": func(l *lua.LState) int {
			return buyModule(l, state)
		},
	})

	// register other stuff
	L.SetField(mod, "name", lua.LString("value"))

	// returns the module
	L.Push(mod)
	return 1
}

func isAuthorized(L *lua.LState, state *PlayerState) int {
	state_Module := state.AcceptedModules["strings"]

	test := lua.LBool(state_Module.Owned)
	L.Push(test)
	return 1
}

func buyModule(L *lua.LState, state *PlayerState) int {
	state_Module := state.AcceptedModules["strings"]
	state_Module.Owned = true
	state.AcceptedModules["strings"] = state_Module

	test := lua.LBool(state_Module.Owned)
	L.Push(test)
	return 1
}

package main

import (
	"io"
	"os"

	lua "github.com/yuin/gopher-lua"
)

type AcceptedModule struct {
	Owned    bool
	Price    int
	ShopRune rune
}

type PlayerState struct {
	PlayerMoney     int
	AcceptedModules map[string]AcceptedModule
	OpennedFile     string
}

var ModulesLists = make(map[string]int)

func runCode(state *PlayerState, code string) {
	rescueStdout := os.Stdout
	r, w, _ := os.Pipe()
	os.Stdout = w

	L := lua.NewState()
	defer L.Close()
	L.PreloadModule("strings", func(l *lua.LState) int {
		Loader(l, state)
		return 1
	})
	L.PreloadModule("test", func(l *lua.LState) int {
		LoaderTesting(l, state)
		return 1
	})
	if err := L.DoString(code); err != nil {
		InfoBox.SetText(string(err.Error()))
		w.Close()
		os.Stdout = rescueStdout
		return
	}

	w.Close()
	out, _ := io.ReadAll(r)
	os.Stdout = rescueStdout
	InfoBox.SetText(string(out))
}

func main() {
	states := PlayerState{
		AcceptedModules: map[string]AcceptedModule{
			"strings": {
				Owned:    false,
				Price:    30,
				ShopRune: 'a',
			},
			"test": {
				Owned:    false,
				Price:    10,
				ShopRune: 'b',
			},
		},
		PlayerMoney: 100,
	}
	mainMenu(&states)
}

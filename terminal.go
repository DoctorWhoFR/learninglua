package main

import (
	"fmt"
	"os"
	"path/filepath"

	"github.com/gdamore/tcell/v2"
	"github.com/rivo/tview"
)

var pages = tview.NewPages()
var app = tview.NewApplication()

func mainMenuSelect(menu string) {
	fmt.Println("selected", menu)
	switch menu {
	case "codes":
		pages.SwitchToPage("codes")

		break
	case "servers":
		break
	case "profiles":
		pages.SwitchToPage("profiles")
	default:
		break
	}
}

func selectionsMenu() *tview.List {
	list := tview.NewList()
	list.AddItem("QEditor", "Manage your code and buyed modules", 'q', func() {
		mainMenuSelect("codes")
	})
	list.AddItem("Server Management", "Manage your server and client server.", 's', func() {
		mainMenuSelect("servers")
	})
	list.AddItem("Profiles", "Show your account information", 'p', func() {
		mainMenuSelect("profiles")
	})

	return list
}

type topBar struct {
	state *PlayerState
	View  *tview.TextView
}

var TopBar *topBar

func (bar *topBar) draw() *tview.TextView {
	textViewTopInfo := tview.NewTextView().
		SetDynamicColors(true)
	fmt.Fprintf(textViewTopInfo, "Money: %d", bar.state.PlayerMoney)
	bar.View = textViewTopInfo
	return textViewTopInfo
}

func (bar *topBar) redraw() *tview.TextView {
	textViewTopInfo := bar.View
	textViewTopInfo.SetText("")
	fmt.Fprintf(textViewTopInfo, "Money: %d", bar.state.PlayerMoney)
	bar.View = textViewTopInfo
	return textViewTopInfo
}

var TextArea *tview.TextArea
var Inside *tview.Flex

func centerZone(state *PlayerState) *tview.Flex {
	_topBar := topBar{
		state: state,
	}

	TopBar = &_topBar

	center_text := tview.NewFlex().
		AddItem(_topBar.draw(), 0, 1, false).
		AddItem(tview.NewBox().SetBorder(false), 0, 1, false).
		AddItem(tview.NewBox().SetBorder(false), 0, 1, false)
	center_text.SetBorder(true)

	TextArea = tview.NewTextArea().
		SetPlaceholder("Enter text here...")
	TextArea.SetTitle("Text Area").SetBorder(true)

	executeButton := tview.NewButton("Execute").SetSelectedFunc(func() {
		runCode(state, TextArea.GetText())
	})

	style := tcell.Style.Background(tcell.StyleDefault, tcell.ColorGreen)
	style_dark := tcell.Style.Background(tcell.StyleDefault, tcell.ColorYellow)
	executeButton.SetStyle(style)
	executeButton.SetActivatedStyle(style_dark)

	saveAsButton := tview.NewButton("Save As").SetSelectedFunc(func() {
		topArea := Inside.GetItem(0)
		textArea := Inside.GetItem(1)
		textAreaBtns := Inside.GetItem(2)

		Inside.RemoveItem(textArea)
		Inside.RemoveItem(textAreaBtns)

		path := ""
		form := tview.NewForm().
			AddTextArea("Path:", "", 40, 1, 0, func(text string) {
				path = text
			}).
			AddButton("Save", func() {
				if path == "" {
					InfoBox.SetText("[red] File name can't be empty")
					return
				}
				final_path := fmt.Sprintf("./qcodes/%s", path)
				InfoBox.SetText(final_path)

				os.WriteFile(final_path, []byte(TextArea.GetText()), 0777)
				Inside.Clear()
				Inside.AddItem(topArea, 0, 1, false)
				Inside.AddItem(textArea, 0, 3, false)
				Inside.AddItem(textAreaBtns, 3, 1, false)

				root.ClearChildren()
				addToRoot(root, "./qcodes")

				InfoBox.SetText(fmt.Sprintf("file %s saved", path))
				state.OpennedFile = final_path
			}).
			AddButton("Quit", func() {
				Inside.Clear()
				Inside.AddItem(topArea, 0, 1, false)
				Inside.AddItem(textArea, 0, 3, false)
				Inside.AddItem(textAreaBtns, 3, 1, false)
			})
		Inside.AddItem(form, 0, 4, true)
	})

	saveButton := tview.NewButton("Save").SetSelectedFunc(func() {
		if state.OpennedFile == "" {
			InfoBox.SetText("[red]No file selected, selecte one in file explorer[white]")
			return
		}
		err := os.WriteFile(state.OpennedFile, []byte(TextArea.GetText()), 0777)
		if err != nil {
			InfoBox.SetText(err.Error())
		}
		InfoBox.SetText(fmt.Sprintf("[green]file %s saved", state.OpennedFile))
	})

	reloadBtn := tview.NewButton("RELOAD").SetSelectedFunc(func() {
		if state.OpennedFile == "" {
			InfoBox.SetText("[red]No file selected, can't reload")
			return
		}
		file, err := os.ReadFile(state.OpennedFile)
		if err != nil {
			InfoBox.SetText(err.Error())
		}
		TextArea.SetText(string(file), true)
		InfoBox.SetText(fmt.Sprintf("[green]file %s reloaded", state.OpennedFile))
	})

	bottom := tview.NewFlex().
		AddItem(executeButton, 0, 1, false).
		AddItem(saveButton, 0, 1, false).
		AddItem(saveAsButton, 0, 1, false).
		AddItem(reloadBtn, 0, 1, false)

	inside := tview.NewFlex().SetDirection(tview.FlexRow).
		AddItem(center_text, 3, 0, false).
		AddItem(TextArea, 0, 1, false).
		AddItem(bottom, 3, 1, false)
	Inside = inside

	return inside
}

var ShopsLists *tview.List
var root *tview.TreeNode

func addToRoot(target *tview.TreeNode, path string) {
	files, err := os.ReadDir(path)
	if err != nil {
		panic(err)
	}
	for _, file := range files {
		node := tview.NewTreeNode(file.Name()).
			SetReference(filepath.Join(path, file.Name())).
			SetSelectable(true)
		if file.IsDir() {
			node.SetColor(tcell.ColorGreen)
		}
		target.AddChild(node)
	}
}

func leftZone(state *PlayerState) *tview.Flex {
	list := tview.NewList()
	ShopsLists = list

	for key, module := range state.AcceptedModules {
		if module.Owned {
			list.AddItem(key, "Owned", module.ShopRune, nil)
		} else {
			list.AddItem(key, fmt.Sprintf("Price: %d", module.Price), module.ShopRune, func() {
				_module := state.AcceptedModules[key]
				if _module.Owned {
					return
				}
				if state.PlayerMoney >= module.Price {
					state.PlayerMoney -= module.Price
					module := state.AcceptedModules[key]
					module.Owned = true
					state.AcceptedModules[key] = module
					TopBar.redraw()
				} else {
					InfoBox.SetText("You don't have enought money to buy this item.")
					return
				}

				test := ShopsLists.FindItems(key, fmt.Sprintf("Price: %d", module.Price), true, false)
				ShopsLists.SetItemText(test[0], key, "Owned")
			})
		}
	}

	list.SetBorder(true)
	list.SetTitle("Shops")

	rootDir := "./qcodes"
	root = tview.NewTreeNode(rootDir).
		SetColor(tcell.ColorRed)
	tree := tview.NewTreeView().
		SetRoot(root).
		SetCurrentNode(root)
	addToRoot(root, rootDir)
	tree.SetBorder(true)
	tree.SetTitle("Files")
	tree.SetSelectedFunc(func(node *tview.TreeNode) {
		reference := node.GetReference()
		if reference == nil {
			InfoBox.SetText("test")

			return // Selecting the root node does nothing.
		}
		children := node.GetChildren()
		if len(children) == 0 {
			// Load and show files in this directory.
			path := fmt.Sprintf("%s/%s", rootDir, node.GetText())
			file, err := os.ReadFile(path)
			if err != nil {
				InfoBox.SetText(err.Error())
			}
			state.OpennedFile = path
			TextArea.SetText(string(file), true)
		} else {
			// Collapse if visible, expand if collapsed.
			node.SetExpanded(!node.IsExpanded())

			InfoBox.SetText("test")
		}
	})

	inside := tview.NewFlex().SetDirection(tview.FlexRow).
		AddItem(list, 0, 1, false).
		AddItem(tree, 0, 1, false)

	return inside
}

var InfoBox *tview.TextView

func InfoBoxDraw() *tview.TextView {
	InfoBox = tview.NewTextView().
		SetDynamicColors(true)
	InfoBox.SetBorder(true)
	InfoBox.SetTitle("InfoBox")
	return InfoBox
}

func mainMenu(state *PlayerState) {
	//box := tview.NewBox().SetBorder(true).SetTitle("Hello, world!")
	pages.AddPage("main", selectionsMenu(), true, true)

	flex := tview.NewFlex().
		AddItem(leftZone(state), 0, 1, false).
		AddItem(centerZone(state), 150, 0, false).
		AddItem(InfoBoxDraw(), 0, 1, false)

	pages.AddPage("codes", flex, true, false)

	form := tview.NewForm().
		AddTextView("Pseudo", "Test", 40, 2, true, false).
		AddTextView("Money", fmt.Sprintf("%d$", state.PlayerMoney), 40, 2, true, false).
		AddTextView("Modules:", "", 40, 2, true, false)
	// AddButton("Quit", func() {
	// 	pages.SwitchToPage("main")
	// })

	for key, module := range state.AcceptedModules {
		text := fmt.Sprintf("Not Owned, price is: %d", module.Price)
		if module.Owned {
			text = "Owned"
		}
		form.AddTextView(key, text, 40, 2, true, false)
	}
	form.AddButton("Quit", func() {
		pages.SwitchToPage("main")
	})

	form.SetBorder(true)
	form.SetTitle("Profiles")

	pages.AddPage("profiles", form, true, false)

	app.SetInputCapture(func(event *tcell.EventKey) *tcell.EventKey {
		if event.Key() == tcell.KeyCtrlA {
			pages.SwitchToPage("main")
		}

		return event
	})

	if err := app.SetRoot(pages, true).EnableMouse(true).Run(); err != nil {
		panic(err)
	}

}

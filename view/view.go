package view

import (
	"github.com/gdamore/tcell/v2"
	"github.com/rivo/tview"
)

var mainFlex *tview.Flex
var inputText *tview.InputField
var messageView *tview.TextView
var infoView *tview.TextView
var commandList *tview.DropDown
var app *tview.Application

var commandValues, commandNames []string

var enterFunc func(text string)

func init() {

	commandValues = []string{"", "/sub ", "/aka "}
	commandNames = []string{"Back to Pub to current topic", "/sub: Sub a topic", "/aka: Set AKA"}

	app = tview.NewApplication()
	mainFlex = tview.NewFlex()
	inputText = tview.NewInputField()
	messageView = tview.NewTextView()
	infoView = tview.NewTextView()
	commandList = tview.NewDropDown()

	messageView.
		SetDynamicColors(true).
		SetBorder(false).
		SetBackgroundColor(tcell.ColorBlack)

	infoView.
		SetDynamicColors(true).
		SetBorder(false).
		SetBackgroundColor(tcell.ColorGreenYellow)

	inputText.
		SetFieldBackgroundColor(tcell.ColorBlack).
		SetBackgroundColor(tcell.ColorAntiqueWhite)

	commandList.
		SetOptions(commandNames, nil).
		SetBackgroundColor(tcell.ColorBlanchedAlmond)

	mainFlex.SetDirection(tview.FlexRow).
		AddItem(infoView, 1, 0, false).
		AddItem(messageView, 0, 1, false).
		AddItem(inputText, 1, 0, true)

	setEvents()

}

func setEvents() {
	inputText.SetChangedFunc(func(text string) {
		if text == "/" {
			mainFlex.AddItem(commandList, 1, 0, false)
			app.SetFocus(commandList)
		}
	})

	inputText.SetDoneFunc(func(key tcell.Key) {
		switch key {
		case tcell.KeyEnter:
			if enterFunc != nil {
				enterFunc(inputText.GetText())
			}

			inputText.SetText("")
		case tcell.KeyEscape:
			inputText.SetText("")
		}

	})

	messageView.SetDoneFunc(func(key tcell.Key) {
		app.SetFocus(inputText)
	})

	infoView.SetDoneFunc(func(key tcell.Key) {
		app.SetFocus(inputText)
	})

	commandList.SetSelectedFunc(func(text string, index int) {
		app.SetFocus(inputText)
		mainFlex.RemoveItem(commandList)
		inputText.SetText(commandValues[index])
	})
}

func Run(ef func(text string)) {
	enterFunc = ef
	if err := app.SetRoot(mainFlex, true).EnableMouse(true).SetFocus(inputText).Run(); err != nil {
		panic(err)
	}
}

func AddMessage(text []byte) {

	go func(txt []byte) {
		app.QueueUpdateDraw(func() {
			messageView.Write(append(txt, '\n'))
			messageView.ScrollToEnd()
		})
	}(text)

}

func SetInfoView(info string) {
	go func(inf string) {
		app.QueueUpdateDraw(func() {
			infoView.SetText(inf)
		})
	}(info)
}

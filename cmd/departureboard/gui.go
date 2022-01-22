package main

import (
	"github.com/lxn/walk"
	. "github.com/lxn/walk/declarative"
	"github.com/pkg/browser"
	"github.com/sqweek/dialog"
)

var mainWindow *walk.MainWindow
var statusLink *walk.LinkLabel

func runWindowsUI() {

	var timetableLabel *walk.Label
	var userEdit, passwordEdit, addressEdit *walk.LineEdit
	var connectButton, chooseButton *walk.PushButton

	MainWindow{
		AssignTo: &mainWindow,
		Title:    "SimSig Departure Board",
		Size:     Size{200, 270},
		MinSize:  Size{200, 270},
		MaxSize:  Size{600, 270},
		Layout:   VBox{},
		Children: []Widget{
			Label{
				Text: "SimSig Departure Board",
				Font: Font{PointSize: 18, Bold: true},
			},
			Composite{
				Layout: Grid{Columns: 2},
				Children: []Widget{
					PushButton{
						AssignTo: &chooseButton,
						Text:     "Choose Timetable",
						OnClicked: func() {
							var filename string
							var err error
							filename, err = dialog.File().Title("Select Active SimSig WTT").SetStartDir("C:\\Users\\Public\\Documents\\SimSig\\Timetables").Filter("SimSig WTT(*.WTT)", "wtt").Load()
							if err != nil {
								walk.MsgBox(mainWindow, "Timetable Load Error", err.Error(), walk.MsgBoxOK|walk.MsgBoxIconError)
								return
							}
							result := loadTimetable(filename)
							if result != "" {
								walk.MsgBox(mainWindow, "Timetable Load Error", result, walk.MsgBoxOK|walk.MsgBoxIconError)
							} else {
								timetableLabel.SetText(filename)
								connectButton.SetEnabled(true)
								chooseButton.SetEnabled(false)
							}
						},
					},
					Label{
						AssignTo:     &timetableLabel,
						Text:         "No timetable selected",
						EllipsisMode: EllipsisPath,
					},

					Label{
						Text:          "User:",
						TextAlignment: AlignFar,
					},
					LineEdit{AssignTo: &userEdit},

					Label{
						Text:          "Password:",
						TextAlignment: AlignFar,
					},
					LineEdit{
						AssignTo:     &passwordEdit,
						PasswordMode: true,
					},
					Label{
						Text:          "Interface Gateway:",
						TextAlignment: AlignFar,
					},
					LineEdit{
						AssignTo: &addressEdit,
						Text:     "localhost:51515",
					},
				},
			},
			PushButton{
				AssignTo: &connectButton,
				Text:     "Connect",
				Enabled:  false,
				OnClicked: func() {
					statusLink.SetText("Connecting to " + addressEdit.Text())
					connectButton.SetEnabled(false)
					go gatewayConnection(userEdit.Text(), passwordEdit.Text(), addressEdit.Text())
				},
			},
			LinkLabel{
				AssignTo: &statusLink,
				//MaxSize:  Size{500, 0},
				Text: `Not Connected to SimSig`,
				OnLinkActivated: func(link *walk.LinkLabelLink) {
					browser.OpenURL("http://localhost:8090/")
				},
			},
		},
	}.Run()
}

func updateStatus(status string) {
	mainWindow.Synchronize(func() {
		statusLink.SetText(status)
	})
}

func webInterfaceReady() {
	mainWindow.Synchronize(func() {
		statusLink.SetText("<a href=\"http://localhost:8090\">Open Departure Board</a>")
	})
}

func showError(message string) {
	mainWindow.Synchronize(func() {
		walk.MsgBox(mainWindow, "Error", message, walk.MsgBoxOK|walk.MsgBoxIconError)
	})
}

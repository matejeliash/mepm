package gui

import (
	"fmt"
	"image/color"
	"os"
	"path/filepath"
	"runtime"
	"strings"

	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/app"
	"fyne.io/fyne/v2/canvas"
	"fyne.io/fyne/v2/container"
	"fyne.io/fyne/v2/dialog"
	"fyne.io/fyne/v2/theme"
	"fyne.io/fyne/v2/widget"
	"github.com/matejeliash/mepm/internal/models"
	"github.com/matejeliash/mepm/internal/passmanager"
)

func CreateDbPath(app fyne.App) string {
	if runtime.GOOS == "android" {

		rootUri := app.Storage().RootURI()
		dbPath := filepath.Join(rootUri.Path(), "test.db")
		return dbPath
	} else if runtime.GOOS == "linux" {
		//home := os.Getenv("HOME")
		//return filepath.Join(home, ".test.db")
		return "./data.db"
	} else {
		return "test.db"
	}
}

func RemoveDb(app fyne.App) error {
	if runtime.GOOS == "android" {

		rootUri := app.Storage().RootURI()
		dbPath := filepath.Join(rootUri.Path(), "test.db")
		err := os.Remove(dbPath)
		return err
	} else if runtime.GOOS == "linux" {
		//home := os.Getenv("HOME")
		//return filepath.Join(home, ".test.db")
		err := os.Remove("./data.db")
		return err
	} else {
		err := os.Remove("./data.db")
		return err
	}
}

type GuiManager struct {
	App         fyne.App
	Window      fyne.Window
	ScreenStack []fyne.CanvasObject
	passmanager.PassManager
	FilterPasswordStr string
}

func NewGuiManager() (*GuiManager, error) {

	gm := &GuiManager{}
	gm.App = app.New()
	gm.Window = gm.App.NewWindow("App")
	gm.Window.Resize(fyne.NewSize(500, 1000))

	if runtime.GOOS == "android" {
		gm.App.Settings().SetTheme(&mobileTheme{})
	}

	dbPath := CreateDbPath(gm.App)
	pm, err := passmanager.NewPassManager(dbPath)
	gm.PassManager = *pm
	if err != nil {
		return nil, err
	}
	return gm, nil

}

func (gm *GuiManager) GoBack() {
	if len(gm.ScreenStack) > 1 {
		gm.ScreenStack = gm.ScreenStack[:len(gm.ScreenStack)-1]
		gm.Window.SetContent(gm.ScreenStack[len(gm.ScreenStack)-1])

		prev := gm.ScreenStack[len(gm.ScreenStack)-1]
		gm.Window.SetContent(prev)

	} else {
		gm.Window.Close()
	}

}

func (gm *GuiManager) SetWindowContent(name string, content fyne.CanvasObject) {
	gm.Window.SetContent(content)
	gm.Window.SetTitle(name)
}

func (gm *GuiManager) ShowFirstScreen() {

	// gm.Window.SetCloseIntercept(func() {
	// 	gm.GoBack()
	// })

	if gm.DB.HasMasterTable() {
		gm.ShowLoginScreen()
	} else {
		gm.ShowInitScreen()
	}

	gm.Window.ShowAndRun()
}

func (gm *GuiManager) ShowInitScreen() {

	infoText := widget.NewLabel("Please, create your master password.")
	passwordEntry := widget.NewPasswordEntry()
	passwordEntry.SetPlaceHolder("Password")

	passwordRepeatEntry := widget.NewPasswordEntry()
	passwordRepeatEntry.SetPlaceHolder("Password")

	errorLabel := widget.NewLabel("")

	// Login button handler
	confirmButton := widget.NewButton("confirm", func() {
		password1 := passwordEntry.Text
		password2 := passwordRepeatEntry.Text

		// Simple hardcoded check (replace with real check if needed)
		if password1 == password2 && (password1 != "") {
			// Close login window and open main app window

			err := gm.InitPassManagerDB(password1)
			if err != nil {
				errorLabel.SetText(err.Error())
			} else {
				err = gm.FetchMasterTable()
				if err != nil {
					errorLabel.SetText(err.Error())
				}

				gm.PassManager.Password = password1

				gm.ShowMainScreen("passwords")

			}
		} else {
			dialog.ShowError(fmt.Errorf("Invalid credentials"), gm.Window)
		}
	})

	// Layout for login screen
	loginForm := container.NewVBox(
		infoText,
		passwordRepeatEntry,
		passwordEntry,
		confirmButton,
	)

	content := container.NewVBox(loginForm, errorLabel)

	gm.SetWindowContent("init", content)

}

func (gm *GuiManager) ShowLoginScreen() {
	largeText := canvas.NewText("MEPM", color.White)
	largeText.TextSize = 36 // Set custom font size
	largeText.Alignment = fyne.TextAlignCenter

	passwordEntry := widget.NewPasswordEntry()
	passwordEntry.SetPlaceHolder("Password")
	errorLabel := widget.NewLabel("")

	err := gm.FetchMasterTable()
	if err != nil {
		panic(err)
	}

	// Login button handler
	loginButton := widget.NewButton("Login", func() {
		password := passwordEntry.Text

		// Simple hardcoded check (replace with real check if needed)
		if gm.IsPasswordCorrect(password) {

			// Close login window and open main app window
			gm.PassManager.Password = password
			gm.ShowMainScreen("passwords")
		} else {
			//dialog.ShowError(fmt.Errorf("Invalid credential"), gm.Window)
			errorLabel.SetText("incorrect password")
		}
	})

	// Layout for login screen
	container := container.NewVBox(
		largeText,
		widget.NewLabel("Please login"),
		passwordEntry,
		container.NewHBox(loginButton),
		errorLabel,
	)

	gm.SetWindowContent("login", container)

}

func (gm *GuiManager) ShowMainScreen(activeTab string) {

	notesContent := gm.CreateNoteScreenContent()
	otherContent := gm.CreateOtherScreenContent()
	passwordContent := gm.CreateScreenContent()

	passwordsTab := container.NewTabItem("Passwords", passwordContent)
	notesTab := container.NewTabItem("Notes", notesContent)
	othersTab := container.NewTabItem("Other", otherContent)
	// Tabs
	tabs := container.NewAppTabs(
		passwordsTab,
		notesTab,
		othersTab,
	)

	if activeTab == "passwords" {
		tabs.SelectTab(passwordsTab)
	} else if activeTab == "notes" {
		tabs.SelectTab(notesTab)

	} else {
		tabs.SelectTab(othersTab)
	}

	if runtime.GOOS == "android" {
		tabs.SetTabLocation(container.TabLocationBottom)

	} else {

		tabs.SetTabLocation(container.TabLocationTop)
	}

	gm.SetWindowContent("main", tabs)

}

func (gm *GuiManager) SimpleRecordFilter(records []models.Record) []models.Record {
	filterStr := gm.FilterPasswordStr
	if filterStr == "" {
		return records
	}
	var filteredRecords []models.Record
	for _, r := range records {
		if strings.Contains(r.Info, filterStr) || strings.Contains(r.Username, filterStr) {
			filteredRecords = append(filteredRecords, r)
		}
	}
	return filteredRecords
}

func (gm *GuiManager) CreateScreenContent() *fyne.Container {

	largeText := canvas.NewText("mepm", color.White)
	largeText.TextSize = 36 // Set custom font size
	largeText.Alignment = fyne.TextAlignCenter

	searchEntry := widget.NewEntry()

	searchEntry.SetPlaceHolder("filter")
	filterButton := widget.NewButtonWithIcon("", theme.SearchIcon(), func() {
		gm.FilterPasswordStr = searchEntry.Text
		gm.ShowMainScreen("passwords")

	})

	insertButton := widget.NewButtonWithIcon("insert", theme.ContentAddIcon(), func() {
		gm.ShowInsertScreen()

	})

	errorLabel := widget.NewLabel("")
	//items := []string{"Apple", "Banana", "Cherry"}
	items, err := gm.GetRecords()
	if err != nil {
		errorLabel.SetText(err.Error())
	}

	filteredItems := gm.SimpleRecordFilter(items)

	itemContainer := container.NewVBox()

	for _, item := range filteredItems {
		itemLabel := widget.NewLabel(item.Info)

		//usernameLabel := widget.NewLabel(item.Username)
		usernameLabel := canvas.NewText(item.Username, color.White)
		usernameLabel.TextSize = 20 // Set custom font size
		usernameLabel.Alignment = fyne.TextAlignCenter

		// adding of button for copying password to clipboard
		getPasswordButton := widget.NewButtonWithIcon("", theme.ContentCopyIcon(), func(item models.Record) func() {
			return func() {
				passwordStr, err := gm.DecryptPassword(item)

				if err != nil {
					dialog.ShowError(err, gm.Window)

				} else {
					fyne.Clipboard.SetContent(gm.App.Clipboard(), passwordStr)

				}
			}
		}(item))
		getUsernameButton := widget.NewButtonWithIcon("", theme.AccountIcon(), func(models.Record) func() {
			return func() {
				fyne.Clipboard.SetContent(gm.App.Clipboard(), item.Username)

			}
		}(item))

		editButton := widget.NewButtonWithIcon("", theme.DocumentCreateIcon(), func(models.Record) func() {
			return func() {
				gm.ShowEditScreen(item)
			}
		}(item))

		deleteButton := widget.NewButtonWithIcon("", theme.DeleteIcon(), func() {
			dialog.NewConfirm("Deletion", "Do you want to delete?", func(answer bool) {
				if answer {
					err := gm.DB.RemoveRecord(&item)
					if err != nil {
						dialog.ShowError(err, gm.Window)
					}
				}
				gm.ShowMainScreen("passwords")
			}, gm.Window).Show() // <-- .Show() is critical!
		})
		itemRow := container.NewVBox(
			container.NewHBox(itemLabel),
			container.NewHBox(usernameLabel),
			container.NewHBox(getPasswordButton, getUsernameButton, editButton, deleteButton),
		)

		itemContainer.Add(itemRow)
	}

	top := container.NewVBox(
		largeText,
		searchEntry,
		container.NewHBox(filterButton),
		container.NewHBox(insertButton),
	)

	// Scrollable content
	scroll := container.NewVScroll(itemContainer)

	// Use Border layout: top widgets at the top, scroll fills remaining space
	content := container.NewBorder(top, nil, nil, nil, scroll)

	return content

}

func (gm *GuiManager) CreateNoteScreenContent() *fyne.Container {

	largeText := canvas.NewText("mepm", color.White)
	largeText.TextSize = 36 // Set custom font size
	largeText.Alignment = fyne.TextAlignCenter

	searchEntry := widget.NewEntry()

	searchEntry.SetPlaceHolder("filter")
	filterButton := widget.NewButtonWithIcon("", theme.SearchIcon(), func() {
		gm.ShowMainScreen("notes")

	})

	insertButton := widget.NewButtonWithIcon("insert", theme.ContentAddIcon(), func() {
		gm.ShowNoteInsertScreen()

	})

	errorLabel := widget.NewLabel("")
	items, err := gm.GetNotes()
	if err != nil {
		errorLabel.SetText(err.Error())
	}

	filteredItems := items

	itemContainer := container.NewVBox()

	for _, item := range filteredItems {
		itemLabel := widget.NewLabel(item.Title)

		editButton := widget.NewButtonWithIcon("", theme.DocumentCreateIcon(), func(models.Note) func() {
			return func() {
				//gm.ShowEditScreen(item)
			}
		}(item))

		deleteButton := widget.NewButtonWithIcon("", theme.DeleteIcon(), func() {
			dialog.NewConfirm("Deletion", "Do you want to delete?", func(answer bool) {
				if answer {
					err := gm.DB.RemoveNote(&item)
					if err != nil {
						dialog.ShowError(err, gm.Window)
					}
				}
				gm.ShowMainScreen("notes")
			}, gm.Window).Show() // <-- .Show() is critical!
		})
		itemRow := container.NewVBox(
			container.NewHBox(itemLabel),
			container.NewHBox(editButton, deleteButton),
		)

		itemContainer.Add(itemRow)
	}

	top := container.NewVBox(
		largeText,
		searchEntry,
		container.NewHBox(filterButton),
		container.NewHBox(insertButton),
	)

	// Scrollable content
	scroll := container.NewVScroll(itemContainer)

	// Use Border layout: top widgets at the top, scroll fills remaining space
	content := container.NewBorder(top, nil, nil, nil, scroll)

	return content

}

func (gm *GuiManager) ShowInsertScreen() {

	infoLabel := widget.NewLabel("New password record")

	backButton := widget.NewButtonWithIcon("back", theme.NavigateBackIcon(), func() {
		gm.ShowMainScreen("passwords")

	})

	passwordEntry := widget.NewPasswordEntry()
	passwordEntry.SetPlaceHolder("Password")

	passwordRepeatEntry := widget.NewPasswordEntry()
	passwordRepeatEntry.SetPlaceHolder("Password")

	userEntry := widget.NewEntry()
	userEntry.SetPlaceHolder("username")

	infoEntry := widget.NewEntry()
	infoEntry.SetPlaceHolder("info")
	// Login button handler
	confirmButton := widget.NewButton("confirm", func() {

		password1 := passwordEntry.Text
		password2 := passwordRepeatEntry.Text
		username := userEntry.Text
		info := infoEntry.Text

		// Simple hardcoded check (replace with real check if needed)
		if password1 == password2 && (password1 != "") && info != "" && username != "" {
			// Close login window and open main app window
			gm.PassManager.InsertRecord(password1, username, info)

			gm.ShowMainScreen("passwords")
		} else {
			dialog.ShowError(fmt.Errorf("Invalid credentials"), gm.Window)
		}
	})

	// Layout for login screen
	form := container.NewVBox(
		container.NewCenter(infoLabel),
		container.NewHBox(backButton),
		widget.NewLabel(""),
		infoEntry,
		userEntry,
		passwordRepeatEntry,
		passwordEntry,
		container.NewHBox(confirmButton),
	)

	gm.SetWindowContent("Add new password record:", form)

}

func (gm *GuiManager) ShowEditScreen(record models.Record) {

	infoLabel := widget.NewLabel("Edit password record")

	backButton := widget.NewButtonWithIcon("back", theme.NavigateBackIcon(), func() {
		gm.ShowMainScreen("passwords")

	})

	passwordStr, err := gm.DecryptPassword(record)
	if err != nil {
		dialog.ShowError(err, gm.Window)
	}

	passwordEntry := widget.NewPasswordEntry()
	passwordEntry.SetPlaceHolder("Password")
	passwordEntry.SetText(passwordStr)

	passwordRepeatEntry := widget.NewPasswordEntry()
	passwordRepeatEntry.SetPlaceHolder("Password")
	passwordRepeatEntry.SetText(passwordStr)

	userEntry := widget.NewEntry()
	userEntry.SetPlaceHolder("username")
	userEntry.SetText(record.Username)

	infoEntry := widget.NewEntry()
	infoEntry.SetPlaceHolder("info")
	infoEntry.SetText(record.Info)
	// Login button handler
	confirmButton := widget.NewButton("confirm", func() {

		password1 := passwordEntry.Text
		password2 := passwordRepeatEntry.Text
		username := userEntry.Text
		info := infoEntry.Text

		if password1 == password2 && (password1 != "") && info != "" && username != "" {
			// Close login window and open main app window
			gm.PassManager.UpdateRecord(int(record.ID), info, username, password1)

			gm.ShowMainScreen("passwords")
		} else {
			dialog.ShowError(fmt.Errorf("Invalid values"), gm.Window)
		}
	})

	// Layout for login screen
	form := container.NewVBox(
		container.NewCenter(infoLabel),
		container.NewHBox(backButton),
		widget.NewLabel(""),
		infoEntry,
		userEntry,
		passwordRepeatEntry,
		passwordEntry,
		container.NewHBox(confirmButton),
	)

	gm.SetWindowContent("Edit password record:", form)

}

func (gm *GuiManager) ShowNoteInsertScreen() {

	infoLabel := widget.NewLabel("New note")

	backButton := widget.NewButtonWithIcon("back", theme.NavigateBackIcon(), func() {
		gm.ShowMainScreen("notes")

	})

	titleEntry := widget.NewEntry()
	titleEntry.SetPlaceHolder("Title")
	textEntry := widget.NewMultiLineEntry()
	textEntry.SetPlaceHolder("Write your large text here...")

	// Make it scrollable
	scroll := container.NewScroll(textEntry)
	// Login button handler
	confirmButton := widget.NewButton("confirm", func() {

		text := textEntry.Text
		title := titleEntry.Text

		// Simple hardcoded check (replace with real check if needed)
		if text != "" && title != "" {
			// Close login window and open main app window
			gm.PassManager.InsertNote(title, text)

			gm.ShowMainScreen("notes")
		} else {
			dialog.ShowError(fmt.Errorf("Invalid credentials"), gm.Window)
		}
	})

	// Layout for login screen
	form := container.NewVBox(
		container.NewCenter(infoLabel),
		container.NewHBox(backButton),
		widget.NewLabel(""),
		titleEntry,
		scroll,
		container.NewHBox(confirmButton),
	)

	gm.SetWindowContent("Add new password record:", form)

}

func (gm *GuiManager) CreateOtherScreenContent() *fyne.Container {

	largeText := canvas.NewText("mepm", color.White)
	largeText.TextSize = 36 // Set custom font size
	largeText.Alignment = fyne.TextAlignCenter
	infoLabel := widget.NewLabel(`This app is just demo app. I created this app so i could learn fyne GUI framework.`)

	// removeDbButton := widget.NewButtonWithIcon("", theme.DeleteIcon(), func() {
	// 	dialog.NewConfirm("Deletion", "Do you want to delete?", func(answer bool) {
	// 		if answer {
	// 			err := RemoveDb(gm.App)
	// 			if err != nil {
	// 				dialog.ShowError(fmt.Errorf("Could not delete db."), gm.Window)
	// 				return

	// 			}
	// 		}

	// 		dbPath := CreateDbPath(gm.App)
	// 		pm, err := passmanager.NewPassManager(dbPath)
	// 		gm.PassManager = *pm
	// 		if err != nil {
	// 			dialog.ShowError(fmt.Errorf("Could not create db file"), gm.Window)
	// 		}
	// 		gm.ShowFirstScreen()
	// 	}, gm.Window).Show() // <-- .Show() is critical!
	// })

	top := container.NewVBox(
		largeText,
		infoLabel,
	)

	// Scrollable content

	// Use Border layout: top widgets at the top, scroll fills remaining space

	return top

}

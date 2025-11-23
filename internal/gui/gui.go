package gui

import (
	"fmt"
	"image/color"
	"io"
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

func (gm *GuiManager) RemoveDb() error {
	err := os.Remove(gm.PassManager.DbPath)
	return err
}

type GuiManager struct {
	App         fyne.App
	Window      fyne.Window
	ScreenStack []fyne.CanvasObject
	passmanager.PassManager
	FilterPasswordStr     string
	NotesScreenContent    *fyne.Container
	OtherScreenContent    *fyne.Container
	PasswordScreenContent *fyne.Container
}

type RefreshType int

const (
	REFRESH_NONE RefreshType = iota
	REFRESH_NOTES
	REFRESH_PASSWORDS
	REFRESH_ALL
)

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
	if err != nil {
		return nil, err
	}
	pm.DbPath = dbPath
	gm.PassManager = *pm
	return gm, nil

}

// func (gm *GuiManager) GoBack() {
// 	if len(gm.ScreenStack) > 1 {
// 		gm.ScreenStack = gm.ScreenStack[:len(gm.ScreenStack)-1]
// 		gm.Window.SetContent(gm.ScreenStack[len(gm.ScreenStack)-1])

// 		prev := gm.ScreenStack[len(gm.ScreenStack)-1]
// 		gm.Window.SetContent(prev)

// 	} else {
// 		gm.Window.Close()
// 	}

// }

func (gm *GuiManager) SetWindowContent(name string, content fyne.CanvasObject) {
	gm.Window.SetContent(content)
	gm.Window.SetTitle(name)
}

func (gm *GuiManager) ShowFirstScreen() {

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

				gm.ShowMainScreen("passwords", REFRESH_ALL)

			}
		} else {
			dialog.ShowError(fmt.Errorf("Invalid credentials"), gm.Window)
		}
	})

	importBtn := widget.NewButton("import", func() {
		dialog.NewFileOpen(func(reader fyne.URIReadCloser, err error) {
			if err != nil || reader == nil {
				return
			}
			defer reader.Close()

			//srcName := filepath.Base(reader.URI().Path())
			dstPath := gm.PassManager.DbPath
			// Open your internal file

			dst, err := os.Create(dstPath)
			if err != nil {
				dialog.ShowError(err, gm.Window)
				return
			}
			defer dst.Close()

			// Copy data to chosen destination
			_, err = io.Copy(dst, reader)
			if err != nil {
				dialog.ShowError(err, gm.Window)
			} else {
				dialog.ShowInformation("Success", "File exported!", gm.Window)
			}

			dbPath := CreateDbPath(gm.App)
			pm, _ := passmanager.NewPassManager(dbPath)
			pm.DbPath = dbPath
			gm.PassManager = *pm

			gm.ShowFirstScreen()

		}, gm.Window).Show()

	})

	// Layout for login screen
	loginForm := container.NewVBox(
		infoText,
		passwordRepeatEntry,
		passwordEntry,
		confirmButton,
		importBtn,
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
			gm.ShowMainScreen("passwords", REFRESH_ALL)
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

func (gm *GuiManager) RefreshContent(refresh RefreshType) {
	if refresh == REFRESH_ALL {
		gm.NotesScreenContent = gm.CreateNoteScreenContent()
		gm.OtherScreenContent = gm.CreateOtherScreenContent()
		gm.PasswordScreenContent = gm.CreatePasswordScreenContent()

	} else if refresh == REFRESH_PASSWORDS {
		gm.PasswordScreenContent = gm.CreatePasswordScreenContent()

	} else if refresh == REFRESH_NOTES {
		gm.NotesScreenContent = gm.CreateNoteScreenContent()
	}

}

func (gm *GuiManager) ShowMainScreen(activeTab string, refresh RefreshType) {

	gm.RefreshContent(refresh)

	passwordsTab := container.NewTabItem("Passwords", gm.PasswordScreenContent)
	notesTab := container.NewTabItem("Notes", gm.NotesScreenContent)
	othersTab := container.NewTabItem("Other", gm.OtherScreenContent)
	// Tabs
	tabs := container.NewAppTabs(
		passwordsTab,
		notesTab,
		othersTab,
	)

	// put bar position based on platform
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
	betterFilterStr := strings.ToLower(strings.Replace(gm.FilterPasswordStr, " ", "", -1))
	if betterFilterStr == "" {
		return records
	}
	var filteredRecords []models.Record

	for _, r := range records {
		betterInfo := strings.ToLower(strings.Replace(r.Info, " ", "", -1))
		betterUsername := strings.ToLower(strings.Replace(r.Info, " ", "", -1))
		if strings.Contains(betterUsername, betterFilterStr) || strings.Contains(betterInfo, betterFilterStr) {
			filteredRecords = append(filteredRecords, r)
		}
	}
	return filteredRecords
}

func (gm *GuiManager) CreatePasswordScreenContent() *fyne.Container {

	largeText := canvas.NewText("mepm", color.White)
	largeText.TextSize = 36 // Set custom font size
	largeText.Alignment = fyne.TextAlignCenter

	searchEntry := widget.NewEntry()

	searchEntry.SetPlaceHolder("filter")
	filterButton := widget.NewButtonWithIcon("", theme.SearchIcon(), func() {
		gm.FilterPasswordStr = searchEntry.Text
		gm.ShowMainScreen("passwords", REFRESH_PASSWORDS)

	})

	insertButton := widget.NewButtonWithIcon("insert", theme.ContentAddIcon(), func() {
		gm.ShowPasswordInsertScreen()

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

		// adding username to screen
		usernameLabel := canvas.NewText(item.Username, color.White)
		usernameLabel.TextSize = 20 // Set custom font size
		usernameLabel.Alignment = fyne.TextAlignCenter

		// adding button to copy username
		getUsernameButton := widget.NewButtonWithIcon("", theme.AccountIcon(), func(models.Record) func() {
			return func() {
				fyne.Clipboard.SetContent(gm.App.Clipboard(), item.Username)

			}
		}(item))

		// adding button to copy password
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

		// adding button to edit password entry
		editButton := widget.NewButtonWithIcon("", theme.DocumentCreateIcon(), func(models.Record) func() {
			return func() {
				gm.ShowPasswordEditScreen(item)
			}
		}(item))

		//adding button to delete password entry
		deleteButton := widget.NewButtonWithIcon("", theme.DeleteIcon(), func() {
			dialog.NewConfirm("Deletion", "Do you want to delete?", func(answer bool) {
				if answer {
					err := gm.DB.RemoveRecord(&item)
					if err != nil {
						dialog.ShowError(err, gm.Window)
					}
				}
				gm.ShowMainScreen("passwords", REFRESH_PASSWORDS)
			}, gm.Window).Show() // <-- .Show() is critical!
		})

		// group of all texts and buttons related to password entry
		entryGroup := container.NewVBox(
			container.NewHBox(itemLabel),
			container.NewHBox(usernameLabel),
			container.NewHBox(getPasswordButton, getUsernameButton, editButton, deleteButton),
		)

		itemContainer.Add(entryGroup)
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
		gm.ShowMainScreen("notes", REFRESH_NOTES)

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
				gm.ShowMainScreen("notes", REFRESH_NOTES)
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

func (gm *GuiManager) ShowPasswordInsertScreen() {

	infoLabel := widget.NewLabel("New password record")

	backButton := widget.NewButtonWithIcon("back", theme.NavigateBackIcon(), func() {
		gm.ShowMainScreen("passwords", REFRESH_NONE)

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

			gm.ShowMainScreen("passwords", REFRESH_PASSWORDS)
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

func (gm *GuiManager) ShowPasswordEditScreen(record models.Record) {

	infoLabel := widget.NewLabel("Edit password record")

	backButton := widget.NewButtonWithIcon("back", theme.NavigateBackIcon(), func() {
		gm.ShowMainScreen("passwords", REFRESH_NONE)

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

			gm.ShowMainScreen("passwords", REFRESH_PASSWORDS)
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
		gm.ShowMainScreen("notes", REFRESH_NONE)

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

			gm.ShowMainScreen("notes", REFRESH_NOTES)
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
	infoLabel := widget.NewLabel(`Work in progress !!!`)

	resetAppBtn := widget.NewButton("reset app", func() {
		gm.RemoveDb()
		dbPath := CreateDbPath(gm.App)
		pm, _ := passmanager.NewPassManager(dbPath)
		pm.DbPath = dbPath
		gm.PassManager = *pm

		gm.ShowInitScreen()

	})

	// var startServerBtn *widget.Button
	// var stopServerBtn *widget.Button
	// killServer := make(chan bool, 1)
	// stopServerBtn = widget.NewButton("stop sharing", func() {
	// 	killServer <- true
	// 	fyne.Do(func() {
	// 		stopServerBtn.Hide()
	// 		startServerBtn.Show()
	// 	})

	// })

	// startServerBtn = widget.NewButton("share db", func() {

	// 	fyne.Do(func() {
	// 		startServerBtn.Hide()
	// 		stopServerBtn.Show()
	// 	})

	// 	go func(killServer chan bool) {

	// 		appDir := filepath.Dir(gm.PassManager.DbPath)
	// 		fileServer := http.FileServer(http.Dir(appDir))

	// 		srv := &http.Server{
	// 			Addr:    ":8080",
	// 			Handler: fileServer, // use this handler instead of global mux
	// 		}

	// 		go func() {
	// 			if err := srv.ListenAndServe(); err != nil {
	// 				fmt.Println(err)
	// 			}

	// 		}()

	// 		<-killServer
	// 		srv.Close()

	// 	}(killServer)

	// })
	// stopServerBtn.Hide()

	//download file from url
	//

	exportBtn := widget.NewButton("export", func() {
		dialog.NewFileSave(func(writer fyne.URIWriteCloser, err error) {
			if err != nil || writer == nil {
				return
			}
			defer writer.Close()

			// Open your internal file
			internalFilePath := gm.PassManager.DbPath
			src, err := os.Open(internalFilePath)
			if err != nil {
				dialog.ShowError(err, gm.Window)
				return
			}
			defer src.Close()

			// Copy data to chosen destination
			_, err = io.Copy(writer, src)
			if err != nil {
				dialog.ShowError(err, gm.Window)
			} else {
				dialog.ShowInformation("Success", "File exported!", gm.Window)
			}
		}, gm.Window).Show()

	})

	importBtn := widget.NewButton("import", func() {
		dialog.NewFileOpen(func(reader fyne.URIReadCloser, err error) {
			if err != nil || reader == nil {
				return
			}
			defer reader.Close()

			//srcName := filepath.Base(reader.URI().Path())
			dstPath := gm.PassManager.DbPath
			// Open your internal file

			dst, err := os.Create(dstPath)
			if err != nil {
				dialog.ShowError(err, gm.Window)
				return
			}
			defer dst.Close()

			// Copy data to chosen destination
			_, err = io.Copy(dst, reader)
			if err != nil {
				dialog.ShowError(err, gm.Window)
			} else {
				dialog.ShowInformation("Success", "File exported!", gm.Window)
			}

			dbPath := CreateDbPath(gm.App)
			pm, _ := passmanager.NewPassManager(dbPath)
			pm.DbPath = dbPath
			gm.PassManager = *pm

			gm.ShowFirstScreen()

		}, gm.Window).Show()

	})

	top := container.NewVBox(
		largeText,
		infoLabel,
		resetAppBtn,
		// startServerBtn,
		// stopServerBtn,
		exportBtn,
		importBtn,
	)

	return top

}

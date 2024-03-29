package main

import (
	"database/sql"
	"log"
	"os/user"

	"fyne.io/fyne/v2/app"
	"github.com/ktappdev/filesync/database"
	"github.com/ktappdev/filesync/getFiles"
	"github.com/ktappdev/filesync/parser"

	"github.com/ktappdev/filesync/logging"
	"github.com/ktappdev/filesync/theme"
	"github.com/ktappdev/filesync/ui"
	_ "github.com/mattn/go-sqlite3"
)

func checkError(err error, message string) {
	if err != nil {
		log.Fatal(message, err)
	}
}

func main() {
	alsSourcePath := "/Users/kentaylor/developer/go-projects/livesync/ap/ap/ableton12.als"

	alsData, err := parser.ExtractALS(alsSourcePath)
	if err != nil {
		log.Fatal(err)
	}

	log.Printf("ALS DATA -> %+v", *alsData)
	///////

	usr, err := user.Current()
	if err != nil {
		log.Fatal("Error getting executable path:", err)
	}

	homeDir := usr.HomeDir

	directory := homeDir + "/Desktop/"
	logging.Setup(homeDir + "/livesync")

	// Open a connection to the SQLite database file
	db, err := sql.Open("sqlite3", homeDir+"/livesync/file_manager.db")
	checkError(err, "Failed to open database:")
	defer db.Close()

	// Initialize the database and create the table if it doesn't exist
	err = database.InitDB(db)
	checkError(err, "Failed to init database:")

	database.UpdateFirstRunSetting(db, true)
	checkError(err, "Failed to update first run setting:")

	allFiles, err := getFiles.GetFiles(directory)
	checkError(err, "Failed to get files:")

	err = database.InsertFilesIntoDB(db, allFiles)
	checkError(err, "Failed to insert files into database:")

	allDbFiles, err := database.GetAllFilesFromDB(db)
	checkError(err, "Failed to get files from database:")

	a := app.New()
	a.Settings().SetTheme(theme.NewMyTheme())
	w := ui.NewFileManagerUI(a, allDbFiles)

	w.Run()
}

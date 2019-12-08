package db

import (
	"database/sql"
	"git.bn4t.me/bn4t/dynamic-qr/app/utils"
	_ "github.com/mattn/go-sqlite3"
	"log"
	"os"
)

var Db *sql.DB

func Connect() {
	var err error

	execDir, err := utils.GetExecutionDir()
	if err != nil {
		log.Fatal(err)
	}

	// check if data directory exists and create it if not
	if !utils.DirExists(execDir + "/data") {
		err := os.Mkdir(execDir+"/data", 644)
		if err != nil {
			log.Fatal(err)
		}
	}

	Db, err = sql.Open("sqlite3", execDir+"/data/dynqr.db")
	if err != nil {
		log.Fatal(err.Error())
	}

	_, err = Db.Exec("CREATE TABLE IF NOT EXISTS qrcodes (id INTEGER PRIMARY KEY, password VARCHAR(255), target VARCHAR(255))")
	if err != nil {
		log.Fatal(err.Error())
	}
}

package db

import (
	"database/sql"
	"fmt"
	"os"

	_ "modernc.org/sqlite"
)

func SetupDB(dir string) {
	// Проверяет, существует ли база данных
	_, err := os.Stat(dir)
	if err == nil {
		return
	}

	// Если базы нет, её необходимо создать
	db, err := sql.Open("sqlite", dir)
	if err != nil {
		fmt.Println(err)
		return
	}

	defer db.Close()

	_, err = db.Exec("CREATE TABLE scheduler (id INTEGER PRIMARY KEY AUTOINCREMENT, date CHAR(8), title VARCHAR(255), comment VARCHAR(255), repeat VARCHAR(128));")
	if err != nil {
		fmt.Println(err)
		return
	}

	_, err = db.Exec("CREATE INDEX scheduler_dates ON scheduler (date);")
	if err != nil {
		fmt.Println(err)
		return
	}
}

package database

import (
	"database/sql"
	"fmt"
	"log"
	"os"

	_ "modernc.org/sqlite"
)

// Подключается к существующей базе SQLite или создает новую.
func SetupDB(dir string) *sql.DB {
	// Проверяет, существует ли база данных по указанному пути
	_, errNoDB := os.Stat(dir)

	// Подключение
	db, err := sql.Open("sqlite", dir)
	if err != nil {
		fmt.Println(err)
		return nil
	}

	// Если БД уже существовала, возвращает на неё ссылку
	if errNoDB == nil {
		log.Println("SQLite database detected, connecting...")
		return db
	}

	// В ином случае БД необходимо создать
	log.Println("Creating new SQLite database...")

	_, err = db.Exec("CREATE TABLE scheduler (id INTEGER PRIMARY KEY AUTOINCREMENT, date CHAR(8), title VARCHAR(255), comment VARCHAR(255), repeat VARCHAR(128));")
	if err != nil {
		fmt.Println(err)
		return nil
	}

	_, err = db.Exec("CREATE INDEX scheduler_dates ON scheduler (date);")
	if err != nil {
		fmt.Println(err)
		return nil
	}

	return db
}

package database

import (
	"database/sql"
	"errors"
	"time"

	"github.com/chnmk/todo-list-final-task/internal/services"
)

// Добавляет задачу в БД, возвращает LastInsertId.
func AddTask(DatabaseFile *sql.DB, task services.Task) (int64, error) {
	res, err := DatabaseFile.Exec("INSERT INTO scheduler (date, title, comment, repeat) VALUES (:date, :title, :comment, :repeat)",
		sql.Named("date", task.Date),
		sql.Named("title", task.Title),
		sql.Named("comment", task.Comment),
		sql.Named("repeat", task.Repeat),
	)
	if err != nil {
		return 0, err
	}

	id, err := res.LastInsertId()
	if err != nil {
		return 0, err
	}

	return id, nil
}

// Изменяет задачу в соответствии с входными данными
func UpdateTask(DatabaseFile *sql.DB, task services.Task) error {
	upd, err := DatabaseFile.Exec("UPDATE scheduler SET (date, title, comment, repeat) = (:date, :title, :comment, :repeat) WHERE id = :id",
		sql.Named("id", task.Id),
		sql.Named("date", task.Date),
		sql.Named("title", task.Title),
		sql.Named("comment", task.Comment),
		sql.Named("repeat", task.Repeat),
	)
	if err != nil {
		return err
	}

	rows, err := upd.RowsAffected()
	if err != nil {
		return err
	} else if rows == 0 {
		return errors.New("задача не найдена")
	}

	return nil
}

// Возвращает задачу с указанным id.
func GetTaskById(DatabaseFile *sql.DB, id string) (services.Task, error) {
	var task services.Task

	// Выполнение запроса
	rows, err := DatabaseFile.Query("SELECT * FROM scheduler WHERE id = :id LIMIT 1", sql.Named("id", id))
	if err != nil {
		return task, err
	}

	defer rows.Close()

	// Чтение полученных данных
	for rows.Next() {
		err := rows.Scan(&task.Id, &task.Date, &task.Title, &task.Comment, &task.Repeat)
		if err != nil {
			return task, err
		}
	}

	// Возврат ошибки если задача не найдена
	if task.Id == "" {
		return task, errors.New("задача не найдена")
	}

	return task, nil
}

// Удаляет или переносит на следующую дату задачу с указанным id.
func CompleteTaskById(DatabaseFile *sql.DB, id string) error {
	var task services.Task

	// Поиск нужной записи
	row := DatabaseFile.QueryRow("SELECT date, repeat FROM scheduler WHERE id = :id LIMIT 1", sql.Named("id", id))

	err := row.Scan(&task.Date, &task.Repeat)
	if err != nil {
		return err
	}
	if task.Date == "" {
		return errors.New("задача не найдена")
	}

	// Удаление записи без правил повторения
	if task.Repeat == "" {
		del, err := DatabaseFile.Exec("DELETE FROM scheduler WHERE id = :id",
			sql.Named("id", id),
		)
		if err != nil {
			return err
		}

		rows, err := del.RowsAffected()
		if err != nil {
			return err
		} else if rows == 0 {
			return errors.New("задача не найдена")
		}
	} else {
		// Редактирование записи с правилами повторения
		task.Date, err = services.NextDate(time.Now(), task.Date, task.Repeat)
		if err != nil {
			return err
		}

		upd, err := DatabaseFile.Exec("UPDATE scheduler SET date = :date WHERE id = :id",
			sql.Named("id", id),
			sql.Named("date", task.Date),
		)
		if err != nil {
			return err
		}

		rows, err := upd.RowsAffected()
		if err != nil {
			return err
		} else if rows == 0 {
			return errors.New("задача не найдена")
		}
	}

	return nil
}

// Удаляет задачу с указанным id.
func DeleteTaskById(DatabaseFile *sql.DB, id string) error {
	del, err := DatabaseFile.Exec("DELETE FROM scheduler WHERE id = :id",
		sql.Named("id", id),
	)
	if err != nil {
		return err
	}

	rows, err := del.RowsAffected()
	if err != nil {
		return err
	} else if rows == 0 {
		return errors.New("задача не найдена")
	}

	return nil
}

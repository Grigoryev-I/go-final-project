package db

import (
	"database/sql"
	"errors"
	"fmt"
	"time"
)

type Task struct {
	ID      string `json:"id"`
	Date    string `json:"date"`
	Title   string `json:"title"`
	Comment string `json:"comment"`
	Repeat  string `json:"repeat"`
}

const dateFormat = "20060102"

func Tasks(limit int, search string) ([]*Task, error) {
	query := `SELECT id, date, title, comment, repeat FROM scheduler`

	var (
		rows *sql.Rows
		err  error
	)

	// все задачи
	if search == "" {
		rows, err = db.Query(query+` ORDER BY date LIMIT ?`, limit)
	} else {
		// поиск по дате
		if d, e := time.Parse("02.01.2006", search); e == nil {
			rows, err = db.Query(query+` WHERE date = ? ORDER BY date LIMIT ?`, d.Format(dateFormat), limit)
		} else {
			// поиск в заголовке и комментарии. Подумать, как сделать регистронезависимый поиск
			like := "%" + search + "%"
			rows, err = db.Query(query+` WHERE title LIKE ? OR comment LIKE ? ORDER BY date LIMIT ?`, like, like, limit)
		}
	}
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	tasks := make([]*Task, 0, limit)

	for rows.Next() {
		var t Task
		if err := rows.Scan(&t.ID, &t.Date, &t.Title, &t.Comment, &t.Repeat); err != nil {
			return nil, err
		}
		tasks = append(tasks, &t)
	}

	return tasks, rows.Err()
}

func AddTask(task Task) (int64, error) {
	var id int64
	query := `INSERT INTO scheduler (date, title, comment, repeat) values (?, ?, ?, ?)`

	res, err := db.Exec(query, task.Date, task.Title, task.Comment, task.Repeat)

	if err == nil {
		id, err = res.LastInsertId()
	}

	return id, err
}

func GetTask(id string) (*Task, error) {
	var task Task

	query := `SELECT id, date, title, comment, repeat FROM scheduler WHERE id = ?`

	err := db.QueryRow(query, id).Scan(&task.ID, &task.Date, &task.Title, &task.Comment, &task.Repeat)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, fmt.Errorf("Задача не найдена")
		}
		return nil, err
	}

	return &task, nil
}

func UpdateTask(task *Task) error {
	query := `UPDATE scheduler SET date = ?, title = ?, comment = ?, repeat = ? where id = ?`

	res, err := db.Exec(query, task.Date, task.Title, task.Comment, task.Repeat, task.ID)
	if err != nil {
		return err
	}

	count, err := res.RowsAffected()
	if err != nil {
		return err
	}
	if count == 0 {
		return fmt.Errorf(`incorrect id for updating task`)
	}

	return nil
}

func DeleteTask(id string) error {
	res, err := db.Exec(`DELETE FROM scheduler WHERE id = ?`, id)
	if err != nil {
		return err
	}

	count, err := res.RowsAffected()
	if err != nil {
		return err
	}
	if count == 0 {
		return fmt.Errorf("incorrect id for deleting task")
	}

	return nil
}

func UpdateDate(next string, id string) error {
	res, err := db.Exec(`UPDATE scheduler SET date = ? WHERE id = ?`, next, id)
	if err != nil {
		return err
	}

	count, err := res.RowsAffected()
	if err != nil {
		return err
	}
	if count == 0 {
		return fmt.Errorf("incorrect id for updating date")
	}

	return nil
}

package main

import (
	"database/sql"
	_ "github.com/mattn/go-sqlite3"
)

type DB struct {
	*sql.DB
}

type Task struct {
	Task_id      int       `json:"task_id"`
	Task_name    string    `json:"task_name"`
	Due_date     string    `json:"due_date"`
	Priority     int       `json:"priority"`
	Description  string    `json:"description"`
	Completed    int       `json:"completed"`
	List_id      int       `json:"list_id"`
}

type Todos interface {
	Add(db *DB, t Task) error
	Get(db *DB, id string) Task
	GetAll(db *DB) []Task
	Update(db *DB) int64
	Delete(db *DB) int64
}

type List struct {
	List_id    	 int       `json:"list_id`
	List_name    string    `json:"list_name`
}

func createTables(db *DB) {
	tables := `
	CREATE TABLE IF NOT EXISTS task (
		task_id 	INTEGER PRIMARY KEY,
		task_name 	TEXT,
		due_date 	TEXT,
		priority 	INTEGER,
		description TEXT,
		completed 	INTEGER,
		list_id 	INTEGER,
		FOREIGN KEY (list_id) REFERENCES list(list_id)
	);
	CREATE TABLE IF NOT EXISTS list (
		list_id 	INTEGER PRIMARY KEY,
		list_name   TEXT
	);`
	db.Exec(tables)
}

func (db *DB) GetLists() []List {
	str := `SELECT * FROM list`
	rows, err := db.Query(str)
	checkErr(err)
	defer rows.Close()

	lists := make([]List, 0)

	for rows.Next() {
		list := List{}
		err = rows.Scan(
			&list.List_id,
			&list.List_name,
		)
		checkErr(err)

		lists = append(lists, list)
	}

	return lists
}

func (db *DB) CreateList(name string) error {
	str := `INSERT INTO list (list_id, list_name) VALUES (?, ?)`
	stmt, err := db.Prepare(str)
	stmt.Exec(nil, name)
	return err
}

func (db *DB) CreateTask(t Task) error {
	str := `INSERT INTO task (task_id, task_name, due_date, priority, description, completed, list_id) VALUES (?, ?, ?, ?, ?, ?, ?)`
	stmt, err := db.Prepare(str)
	stmt.Exec(
		nil,
		t.Task_name,
		t.Due_date,
		t.Priority,
		t.Description,
		t.Completed,
		t.List_id,
	)
	return err
}

func (db *DB) GetAll() []Task {
	str := `SELECT * FROM task`
	rows, err := db.Query(str)
	checkErr(err)
	defer rows.Close()

	todos := make([]Task, 0)

	for rows.Next() {
		t := Task{}
		err = rows.Scan(
			&t.Task_id,
			&t.Task_name,
			&t.Due_date,
			&t.Priority,
			&t.Description,
			&t.Completed,
			&t.List_id,
		)
		checkErr(err)

		todos = append(todos, t)
	}

	return todos
}

func (db *DB) GetList(listID string) []Task {
	str := `SELECT * FROM task WHERE list_id = ?`
	rows, err := db.Query(str, listID)
	checkErr(err)
	defer rows.Close()

	todos := make([]Task, 0)

	for rows.Next() {
		t := Task{}
		err = rows.Scan(
			&t.Task_id,
			&t.Task_name,
			&t.Due_date,
			&t.Priority,
			&t.Description,
			&t.Completed,
			&t.List_id,
		)
		checkErr(err)

		todos = append(todos, t)
	}

	return todos
}

func (db *DB) DeleteTask(taskID string) int64 {
	task, err := db.Prepare(`DELETE FROM task WHERE task_id = ?`)
	checkErr(err)
	defer task.Close()

	res, err := task.Exec(taskID)
	checkErr(err)

	aff, err := res.RowsAffected()
	checkErr(err)

	return aff
}

func (db *DB) DeleteList(listID string) int64 {
	task, err := db.Prepare(`DELETE FROM task WHERE list_id = ?`)
	checkErr(err)
	defer task.Close()

	list, err := db.Prepare(`DELETE FROM list WHERE list_id = ?`)
	checkErr(err)
	defer list.Close()

	taskRes, err := task.Exec(listID)
	checkErr(err)

	list.Exec(listID)

	aff, err := taskRes.RowsAffected()
	checkErr(err)

	return aff
}

func (db *DB) UpdateTask(t Task) int64 {
	q := `UPDATE task set task_name = ?, due_date = ?, priority = ?, description = ?, completed = ?, list_id = ? WHERE task_id = ?`
	stmt, err := db.Prepare(q)

	checkErr(err)
	defer stmt.Close()

	res, err := stmt.Exec(
		t.Task_name,
		t.Due_date,
		t.Priority,
		t.Description,
		t.Completed,
		t.List_id,
		t.Task_id,
	)
	checkErr(err)

	aff, err := res.RowsAffected()
	checkErr(err)

	return aff

}

package main

import (
	"fmt"
	"database/sql"
	"log"

	_ "github.com/mattn/go-sqlite3"
)

func main() {
	sqldb, err := sql.Open("sqlite3", cfg.path)
	checkErr(err)
	defer sqldb.Close()

	db := &DB{sqldb}
	db.Exec("PRAGMA foreign_keys = ON")
	createTables(db)

	//tests(db)
	fmt.Println(db.GetAll())
}

func checkErr(err error) {
	if err != nil {
		log.Fatal(err)
	}
}

var cfg = dbConfig {
	path: "./data/todos.db",
}

type dbConfig struct {
	path string
}

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
	List_id int
	List_name string
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

func tests(db *DB) {
	db.CreateList("List 1")
	db.CreateList("List 2")

	for i := 0; i < 3; i++ {
		t := Task {
			Task_name: fmt.Sprintf("Task %d", i),
			Due_date: "12-34-56",
			Priority: 1,
			Description: "Desc",
			Completed: 0,
			List_id: 1,
		}
		db.CreateTask(t)
	}
	for i := 0; i < 2; i++ {
		t := Task {
			Task_name: fmt.Sprintf("Task %d", i),
			Due_date: "12-34-56",
			Priority: 1,
			Description: "Desc",
			Completed: 0,
			List_id: 2,
		}
		db.CreateTask(t)
	}
	fmt.Println(db.GetLists())
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

func (db *DB) GetList(listID int) []Task {
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

func (db *DB) DeleteList(listID int) int64 {
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

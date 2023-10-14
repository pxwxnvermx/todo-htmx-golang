package main

import (
	"database/sql"

	"github.com/jmoiron/sqlx"
	_ "github.com/mattn/go-sqlite3"
)

type Todo struct {
	Id   string
	Name string
	Done bool
}

func (todo *Todo) save(db *sqlx.DB) sql.Result {
	insert_todo_query := `INSERT INTO todos (name, done) VALUES (?, ?);`
	return db.MustExec(insert_todo_query, todo.Name, todo.Done)
}

func get_todos(db *sqlx.DB) ([]Todo, error) {
	var todos []Todo
	rows, err := db.Queryx("SELECT * FROM todos")
	if err != nil {
		return nil, err
	}

	for rows.Next() {
		var todo Todo
		err := rows.StructScan(&todo)
		if err != nil {
			return nil, err
		}
		todos = append(todos, todo)
	}
	return todos, nil
}

func NewDB() *sqlx.DB {
	db, err := sqlx.Open("sqlite3", "sqlite.db")
	if err != nil {
		panic("cant open db")
	}
	todos_schema := `CREATE TABLE IF NOT EXISTS todos (
	id integer primary key,
	name varchar,
	done boolean);`

	if _, err := db.Exec(todos_schema); err != nil {
		panic("cant create table")
	}

	return db
}

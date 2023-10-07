package main

import (
	"database/sql"
	"fmt"
	"html/template"
	"io"

	"github.com/jmoiron/sqlx"
	"github.com/labstack/echo/v4"
	_ "github.com/mattn/go-sqlite3"
)

type Template struct {
	templates *template.Template
}

func (t *Template) Render(w io.Writer, name string, data interface{}, c echo.Context) error {
	tmpl := template.Must(t.templates.Clone())
	return tmpl.ExecuteTemplate(w, name, data)
}

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

func main() {
	db, err := sqlx.Open("sqlite3", "sqlite.db")
	if err != nil {
		panic("cant open db")
	}
	defer db.Close()
	todos_schema := `CREATE TABLE IF NOT EXISTS todos (
	id integer primary key,
	name varchar,
	done boolean);`

	if _, err := db.Exec(todos_schema); err != nil {
		panic("cant create table")
	}

	t := &Template{
		templates: template.Must(template.ParseGlob("./templates/*.html")),
	}
	e := echo.New()

	e.GET("/", func(c echo.Context) error {
		return c.Render(200, "index.html", 1)
	})

	e.GET("/todos", func(c echo.Context) error {
		todos, err := get_todos(db)
		if err != nil {
			return c.NoContent(401)
		}
		return c.Render(200, "todo_list.html", todos)
	})

	e.POST("/todos", func(c echo.Context) error {
		if err := c.Request().ParseForm(); err != nil {
			return c.NoContent(401)
		}
		form := c.Request().PostForm
		todo := Todo{Id: "", Name: form.Get("name"), Done: false}
		todo.save(db)
		todos, err := get_todos(db)
		if err != nil {
			return c.NoContent(401)
		}
		return c.Render(201, "todo_list.html", todos)
	})

	e.DELETE("/todos/:id", func(c echo.Context) error {
		id := c.Param("id")
		fmt.Printf("deleting %v", id)
		db.MustExec("DELETE FROM todos WHERE id = ?", id)
		return c.NoContent(200)
	})

	e.Renderer = t
	e.Logger.Info(e.Start(":8000"))
}

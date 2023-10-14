package main

import (
	"github.com/jmoiron/sqlx"
	"github.com/labstack/echo/v4"
)

func index_handler(db *sqlx.DB) echo.HandlerFunc {
	return func(c echo.Context) error {
		return c.Render(200, "index.html", 1)
	}
}

func get_all_todos_handler(db *sqlx.DB) echo.HandlerFunc {
	return func(c echo.Context) error {
		todos, err := get_todos(db)
		if err != nil {
			return c.NoContent(401)
		}
		return c.Render(200, "todo_list.html", todos)
	}
}

func save_todos_handler(db *sqlx.DB) echo.HandlerFunc {
	return func(c echo.Context) error {
		if err := c.Request().ParseForm(); err != nil {
			return c.NoContent(401)
		}
		form := c.Request().PostForm
		done := false
		if form.Get("done") == "true" {
			done = true
		}
		todo := Todo{Id: "", Name: form.Get("name"), Done: done}
		todo.save(db)
		todos, err := get_todos(db)
		if err != nil {
			return c.NoContent(401)
		}
		return c.Render(201, "todo_list.html", todos)
	}
}

func delete_todos_handler(db *sqlx.DB) echo.HandlerFunc {
	return func(c echo.Context) error {
		id := c.Param("id")
		db.MustExec("DELETE FROM todos WHERE id = ?", id)
		return c.NoContent(200)
	}
}

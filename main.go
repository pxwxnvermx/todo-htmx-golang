package main

import "github.com/labstack/echo/v4"

func main() {
	db := NewDB()
	defer db.Close()
	e := echo.New()

	e.GET("/", index_handler(db))
	e.GET("/todos", get_all_todos_handler(db))
	e.POST("/todos", save_todos_handler(db))
	e.DELETE("/todos/:id", delete_todos_handler(db))

	e.Renderer = NewTemplates()
	e.Logger.Info(e.Start(":8000"))
}

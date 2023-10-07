package main

import (
	"html/template"
	"io"

	"github.com/labstack/echo/v4"
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

func main() {
	t := &Template{
		templates: template.Must(template.ParseGlob("./templates/*.html")),
	}
	e := echo.New()

	var todos []Todo
	todos = append(todos, Todo{
		Id:   "1",
		Name: "some 1",
		Done: false,
	})
	todos = append(todos, Todo{
		Id:   "2",
		Name: "some 2",
		Done: false,
	})
	todos = append(todos, Todo{
		Id:   "3",
		Name: "some 3",
		Done: false,
	})

	e.GET("/", func(c echo.Context) error {
		return c.Render(200, "index.html", 1)
	})

	e.GET("/todos", func(c echo.Context) error {
		return c.Render(200, "todo_list.html", todos)
	})

	e.POST("/todos", func(c echo.Context) error {
		if err := c.Request().ParseForm(); err != nil {
			return c.NoContent(401)
		}
		form := c.Request().PostForm
		todos = append(todos, Todo{
			Id:   "3",
			Name: form.Get("name"),
			Done: false,
		})
		return c.Render(201, "todo_list.html", todos)
	})

	e.DELETE("/todos/:id", func(c echo.Context) error {
		return c.NoContent(201)
	})

	e.Renderer = t
	e.Logger.Info(e.Start(":8000"))
}

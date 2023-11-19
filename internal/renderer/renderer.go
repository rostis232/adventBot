package renderer

import (
	"fmt"
	"html/template"
	"io"

	"github.com/labstack/echo/v4"
)

const pathToTemplates = "./templates"

type Template struct {
}

func (t *Template) Render(w io.Writer, name string, data interface{}, c echo.Context) error {
	partials := []string{
		fmt.Sprintf("%s/%s", pathToTemplates, name),
		fmt.Sprintf("%s/base.layout.html", pathToTemplates),
		fmt.Sprintf("%s/header.partial.html", pathToTemplates),
		fmt.Sprintf("%s/navbar.partial.html", pathToTemplates),
		fmt.Sprintf("%s/footer.partial.html", pathToTemplates),
		fmt.Sprintf("%s/alerts.partial.html", pathToTemplates),
	}
	tmpl, err := template.ParseFiles(partials...)
	if err != nil {
		return err
	}
	
	return tmpl.Execute(w, data)
}

var Tmps = &Template{
}
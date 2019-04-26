package view

import (
	"html/template"
	"io"
	"path/filepath"

	"alexandria.app/models"
)

// A View represents a template which can be rendered to HTML.
type View struct {
	layout   string
	template string
	config   *models.Config
}

type viewDataWrapper struct {
	Config *models.Config
	Data   interface{}
	User   *models.User
}

// Render exectutes the template and writes the resulting HTML to the io.Writer.
func (v *View) Render(w io.Writer, user *models.User, data interface{}) error {
	tmpl, err := template.New(v.layout).ParseFiles(
		filepath.Join(v.config.TemplateDirectory, v.layout+".html"),
		filepath.Join(v.config.TemplateDirectory, v.template+".html"),
	)

	if err != nil {
		return err
	}

	return tmpl.Execute(w, &viewDataWrapper{
		Config: v.config,
		Data:   data,
		User:   user,
	})
}

// New creates a new View struct from the layout, template and config.
func New(layout, template string, config *models.Config) *View {
	return &View{
		layout:   layout,
		template: template,
		config:   config,
	}
}

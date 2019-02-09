package view

import (
	"html/template"
	"io"
	"path/filepath"

	"alexandria.app/models"
)

type View struct {
	layout   string
	template string
	config   *models.Config
	data     interface{}
}

type viewDataWrapper struct {
	Config *models.Config
	Data   interface{}
}

func (v *View) Render(w io.Writer) error {
	tmpl, err := template.New(v.layout).ParseFiles(
		filepath.Join(v.config.TemplateDirectory, v.layout+".html"),
		filepath.Join(v.config.TemplateDirectory, v.template+".html"),
	)

	if err != nil {
		return err
	}

	return tmpl.Execute(w, &viewDataWrapper{
		Config: v.config,
		Data:   v.data,
	})
}

func New(layout, template string, config *models.Config, data interface{}) *View {
	return &View{
		layout:   layout,
		template: template,
		config:   config,
		data:     data,
	}
}

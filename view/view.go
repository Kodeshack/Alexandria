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
	Styles   []string
	Scripts  []string
}

type viewDataWrapper struct {
	Config  *models.Config
	Data    interface{}
	User    *models.User
	Styles  []string
	Scripts []string
}

func (v *View) Render(w io.Writer, user *models.User, data interface{}) error {
	tmpl, err := template.New(v.layout).ParseFiles(
		filepath.Join(v.config.TemplateDirectory, v.layout+".html"),
		filepath.Join(v.config.TemplateDirectory, v.template+".html"),
	)

	if err != nil {
		return err
	}

	return tmpl.Execute(w, &viewDataWrapper{
		Config:  v.config,
		Data:    data,
		User:    user,
		Styles:  v.Styles,
		Scripts: v.Scripts,
	})
}

func New(layout, template string, config *models.Config) *View {
	return &View{
		layout:   layout,
		template: template,
		config:   config,
	}
}

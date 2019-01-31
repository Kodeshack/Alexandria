package view

import (
	"html/template"
	"io"
	"path/filepath"
)

type View struct {
	layout   string
	template string
	dir      string
	data     interface{}
}

func (v *View) Render(w io.Writer) error {
	tmpl, err := template.New(v.layout).ParseFiles(filepath.Join(v.dir, v.layout+".html"), filepath.Join(v.dir, v.template+".html"))
	if err != nil {
		return err
	}

	return tmpl.Execute(w, v.data)
}

func New(layout, template, dir string, data interface{}) *View {
	return &View{
		layout:   layout,
		template: template,
		dir:      dir,
		data:     data,
	}
}

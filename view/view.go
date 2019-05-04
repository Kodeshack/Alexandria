package view

import (
	"fmt"
	"html/template"
	"io"
	"path/filepath"
	"strings"

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

func renderArticleContent(article *models.Article) (template.HTML, error) {
	body, err := article.ContentHTML()
	if err != nil {
		return "", err
	}

	return template.HTML(body), nil
}

func createArticleEditURL(config *models.Config) func(article models.Article) string {
	return func(article models.Article) string {
		path := filepath.Join(article.Category, article.File)
		path = strings.Replace(path, ".md", "", 1)
		return fmt.Sprintf("%sarticles/edit/%s", config.BaseURL, path)
	}
}

func stripContentPathPrefix(config *models.Config) func(string) string {
	return func(path string) string {
		return strings.Replace(path, config.ContentPath, "", 1)
	}
}

func bytesToString(b []byte) string {
	return string(b)
}

func templateFunctions(config *models.Config) template.FuncMap {
	return map[string]interface{}{
		"stripContentPathPrefix": stripContentPathPrefix(config),
		"articleEditURL":         createArticleEditURL(config),
		"articleContent":         renderArticleContent,
		"bytesToString":          bytesToString,
	}
}

// Render exectutes the template and writes the resulting HTML to the io.Writer.
func (v *View) Render(w io.Writer, user *models.User, data interface{}) error {
	tmpl, err := template.New(v.layout).Funcs(templateFunctions(v.config)).ParseFiles(
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

package models

import (
	"bytes"
	"io/ioutil"
	"os"
	"path"
	"path/filepath"
	"time"

	blackfriday "github.com/russross/blackfriday/v2"

	"github.com/BurntSushi/toml"
)

const (
	delimiter = "\n+++\n"
)

// Metadata for an article.
// Extracted from the TOML data at the beginning of an article file.
type Metadata struct {
	Title        string
	LastEditedAt int64
}

// Article contains all the relevant data for displaying an article.
type Article struct {
	Dir      string
	File     string
	Category string
	Meta     Metadata
	parsed   bool
	Content  []byte
}

// Read the article data from disk and parse the TOML at the beginning of the file.
func (a *Article) Read() error {
	data, err := ioutil.ReadFile(a.Path())
	if err != nil {
		return err
	}

	index := bytes.Index(data, []byte(delimiter))

	a.Content = data[index+len(delimiter):]

	metadata := data[:index]

	err = toml.Unmarshal(metadata, &a.Meta)

	if err != nil {
		return err
	}

	a.parsed = true
	return nil
}

// ContentHTML converts the article's content from markdown to HTML.
func (a *Article) ContentHTML() ([]byte, error) {
	options := blackfriday.WithExtensions(blackfriday.CommonExtensions)

	output := blackfriday.Run(a.Content, options)

	return output, nil
}

// Write the article's content back to disk. Also creates all relevant directories.
func (a *Article) Write() error {
	err := os.MkdirAll(filepath.Join(a.Dir, a.Category), os.ModePerm)
	if err != nil {
		return err
	}

	file, err := os.OpenFile(a.Path(), os.O_RDWR|os.O_CREATE|os.O_TRUNC, os.ModePerm)
	if err != nil {
		return err
	}

	enc := toml.NewEncoder(file)

	err = enc.Encode(&a.Meta)
	if err != nil {
		return err
	}

	_, err = file.Write([]byte("\n" + delimiter))
	if err != nil {
		return err
	}

	_, err = file.Write(a.Content)
	if err != nil {
		return err
	}

	return file.Close()
}

// Path returns the full path on disk.
func (a *Article) Path() string {
	return filepath.Join(a.Dir, a.Category, a.File)
}

// LoadArticle loads the article (contents) at the specified path from disk.
func LoadArticle(dir, category, title string) (*Article, error) {
	path := path.Join(dir, category, title+".md")
	_, err := os.Stat(path)
	if err != nil {
		return nil, err
	}

	a := Article{
		Dir:      dir,
		File:     title + ".md",
		Category: category,
		Meta:     Metadata{},
	}

	err = a.Read()
	if err != nil {
		return nil, err
	}

	return &a, nil
}

// NewArticle is a convenience function to create a new Article struct.
// The article's LastEditedAt field will be set to the current time.
func NewArticle(category, title, content, dir string) *Article {
	return &Article{
		Dir:      dir,
		File:     title + ".md",
		Category: category,
		Content:  []byte(content),
		parsed:   true,
		Meta: Metadata{
			Title:        title,
			LastEditedAt: time.Now().Unix(),
		},
	}
}

// RemoveArticle will remove the article file at the given path.
// Note: The full path (including the content path prefix) is expected here.
func RemoveArticle(path string) error {
	_, err := os.Stat(path)

	if err != nil {
		return err
	}

	return os.Remove(path)
}

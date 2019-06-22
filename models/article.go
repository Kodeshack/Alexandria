package models

import (
	"bytes"
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
	Path     string
	Meta     Metadata
	parsed   bool
	Content  []byte
	Category *Category
}

// Parse the article data.
func (a *Article) Parse(content []byte) error {
	index := bytes.Index(content, []byte(delimiter))

	a.Content = content[index+len(delimiter):]

	metadata := content[:index]

	err := toml.Unmarshal(metadata, &a.Meta)

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

// Serialize the article's content.
func (a *Article) Serialize() ([]byte, error) {
	buffer := bytes.NewBuffer([]byte{})
	enc := toml.NewEncoder(buffer)

	err := enc.Encode(&a.Meta)
	if err != nil {
		return []byte{}, err
	}

	_, err = buffer.Write([]byte("\n" + delimiter))
	if err != nil {
		return []byte{}, err
	}

	_, err = buffer.Write(a.Content)
	if err != nil {
		return []byte{}, err
	}

	return buffer.Bytes(), nil
}

// NewArticle is a convenience function to create a new Article struct.
// The article's LastEditedAt field will be set to the current time.
func NewArticle(title, content, dir string) *Article {
	return &Article{
		Path:    filepath.Join(dir, title+".md"),
		Content: []byte(content),
		parsed:  true,
		Meta: Metadata{
			Title:        title,
			LastEditedAt: time.Now().Unix(),
		},
	}
}

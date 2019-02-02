package models

import (
	"bytes"
	"io/ioutil"
	"os"
	"path/filepath"
	"time"

	blackfriday "github.com/russross/blackfriday/v2"

	"github.com/BurntSushi/toml"
)

const (
	delimiter = "\n+++\n"
)

type Metadata struct {
	Title        string
	LastEditedAt int64
}

type Article struct {
	Path    string
	Meta    Metadata
	parsed  bool
	Content []byte
}

func (a *Article) Read() error {
	data, err := ioutil.ReadFile(a.Path)
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

func (a *Article) ContentHTML() ([]byte, error) {
	options := blackfriday.WithExtensions(blackfriday.CommonExtensions | blackfriday.HardLineBreak)

	output := blackfriday.Run(a.Content, options)

	return output, nil
}

func (a *Article) Write() error {
	err := os.MkdirAll(filepath.Dir(a.Path), os.ModePerm)
	if err != nil {
		return err
	}

	file, err := os.OpenFile(a.Path, os.O_RDWR|os.O_CREATE, os.ModePerm)
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

func LoadArticle(path string) (*Article, error) {
	_, err := os.Stat(path)
	if err != nil {
		return nil, err
	}

	a := Article{
		Path: path,
		Meta: Metadata{},
	}

	err = a.Read()
	if err != nil {
		return nil, err
	}

	return &a, nil
}

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

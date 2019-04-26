package models

import (
	"io/ioutil"
	"path/filepath"
	"strings"
)

// A Category is essentially a folder which contains either more folders
// or article files.
type Category struct {
	Name    string
	Parent  string
	Path    string
	Entries []string
}

// NewCategory is a convenience function to create a new Category.
func NewCategory(name, path string) *Category {
	d, _ := filepath.Split(name)
	parent := d

	if name == d {
		parent = ""
	}

	return &Category{
		Name:   name,
		Parent: parent,
		Path:   path,
	}
}

// ScanEntries reads the entries of the folder at which the category points.
func (c *Category) ScanEntries() error {
	files, err := ioutil.ReadDir(c.Path)
	if err != nil {
		return err
	}

	entries := make([]string, len(files))

	for i, file := range files {
		// We also remove the extensions here because those are not relevant.
		name := file.Name()
		entries[i] = strings.Replace(name, filepath.Ext(name), "", -1)
	}

	c.Entries = entries

	return nil
}

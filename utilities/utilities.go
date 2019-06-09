package utilities

import (
	"os"
)

// CategoryExists will stat the directory at path to determine whether the category exists.
func CategoryExists(path string) bool {
	stat, _ := os.Stat(path)
	return stat != nil && stat.IsDir()
}

// ArticleExists will stat the file at path to check if it exists.
func ArticleExists(path string) bool {
	_, err := os.Stat(path)
	return !os.IsNotExist(err)
}

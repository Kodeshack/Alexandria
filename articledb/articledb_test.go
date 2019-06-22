package articledb

import (
	"bytes"
	"io/ioutil"
	"os"
	"path/filepath"
	"testing"

	"alexandria.app/models"
)

var testContent = []byte(`Title = "Test"
LastEditedAt = 1557593441


+++

# Test Heading
This is just some test text
`)

func writeTestFile(path string) {
	file, err := os.OpenFile(path, os.O_RDWR|os.O_CREATE|os.O_TRUNC, os.ModePerm)
	if err != nil {
		panic(err)
	}

	_, err = file.Write(testContent)
	if err != nil {
		panic(err)
	}
}

func removeFile(path string) {
	_, err := os.Stat(path)

	if os.IsNotExist(err) {
		return
	}

	os.Remove(path)
}

func TestRead(t *testing.T) {
	tmpDir := os.TempDir()
	path := filepath.Join(tmpDir, "_TestRead.md")
	defer removeFile(path)
	writeTestFile(path)

	adb := New(tmpDir)

	a, c, err := adb.Load("_TestRead")
	if err != nil {
		t.Fatal(err)
	}

	if a.Meta.Title != "Test" {
		t.Errorf(`Title is "%v", should be %v`, a.Meta.Title, "Test")
	}

	if c.Name != "." {
		t.Errorf(`Category name is "%v", should be %v`, c.Name, ".")
	}
}

func TestWrite(t *testing.T) {
	tmpDir := os.TempDir()
	path := filepath.Join(tmpDir, "_TestWrite.md")
	defer removeFile(path)

	adb := New(tmpDir)

	article := &models.Article{
		Path: "_TestWrite",
		Meta: models.Metadata{
			Title:        "Test",
			LastEditedAt: 1557593441,
		},
	}

	article.Content = []byte(`
# Test Heading
This is just some test text
`)

	err := adb.Write(article)
	if err != nil {
		t.Error(err)
	}

	file, err := os.Open(path)
	if err != nil {
		t.Fatal(err)
	}

	fileContent, err := ioutil.ReadAll(file)
	if err != nil {
		t.Fatal(err)
	}

	if !bytes.Equal(fileContent, testContent) {
		t.Errorf("Content of %v is wrong!", path)
	}
}

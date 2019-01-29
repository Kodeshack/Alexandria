package models

import (
	"os"
	"path/filepath"
	"testing"
	"time"
)

func removeFile(path string) {
	_, err := os.Stat(path)

	if os.IsNotExist(err) {
		return
	}

	os.Remove(path)
}

func TestWriteFile(t *testing.T) {
	path := filepath.Join(os.TempDir(), "_TestWriteFile")
	defer removeFile(path)

	title := "Foo"
	now := time.Now().Unix()

	article := &Article{
		Path: path,
		Meta: Metadata{
			Title:        title,
			LastEditedAt: now,
		},
	}

	article.Content = []byte(`
# Test Heading
This is just some test text
`)

	err := article.Write()
	if err != nil {
		t.Error(err)
	}

	article, err = NewArticle(path)
	if err != nil {
		t.Error(err)
	}

	if article.Meta.Title != "Foo" {
		t.Errorf(`Title is "%v", should be %v`, article.Meta.Title, title)
	}

	if article.Meta.LastEditedAt != now {
		t.Errorf(`LastEditedAt is %v, should be %v`, article.Meta.LastEditedAt, now)
	}

	shouldHTML := `<h1>Test Heading</h1>

<p>This is just some test text</p>
`

	html, err := article.ContentHTML()
	if err != nil {
		t.Error(err)
	}

	if string(html) != shouldHTML {
		t.Errorf(`HTML is "%v", should be "%v"`, string(html), shouldHTML)
	}
}

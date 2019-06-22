package articledb

import (
	"fmt"
	"io/ioutil"
	"os"
	"path/filepath"
	"sync"

	"alexandria.app/models"
)

func main() {
	fmt.Println("vim-go")
}

// ArticleDB is a thin layer on top of the file system which
// guarantees full consitency.
// Only one ArticleDB per rootDir (and child directory of rootDir) must exist.
type ArticleDB interface {
	Load(path string) (*models.Article, *models.Category, error)
	Write(article *models.Article) error
}

type articledb struct {
	rootDir string
	mutex   *sync.RWMutex
}

// New creates a new ArticleDB.
func New(rootDir string) ArticleDB {
	return &articledb{
		rootDir: rootDir,
		mutex:   &sync.RWMutex{},
	}
}

func (adb *articledb) startReadTransaction() {
	adb.mutex.RLock()
}

func (adb *articledb) endReadTransaction() {
	adb.mutex.RUnlock()
}

func (adb *articledb) startWriteTransaction() {
	adb.mutex.Lock()
}

func (adb *articledb) endWriteTransaction() {
	adb.mutex.Unlock()
}

// Load an article at the given location from disk.
func (adb *articledb) Load(path string) (*models.Article, *models.Category, error) {
	adb.startReadTransaction()
	defer adb.endReadTransaction()
	return adb.load(path)
}

func (adb *articledb) constructCategory(path string) (*models.Category, error) {
	absPath := filepath.Join(adb.rootDir, path)
	stat, err := os.Stat(absPath)
	if err != nil && !os.IsNotExist(err) {
		return nil, err
	}

	var c models.Category

	if stat != nil && stat.IsDir() {
		c = models.Category{
			Name:   filepath.Base(path),
			Parent: "",
			Path:   path,
		}
	} else {
		dir := filepath.Dir(path)

		c = models.Category{
			Name:   filepath.Base(dir),
			Parent: "",
			Path:   dir,
		}
	}

	files, err := ioutil.ReadDir(c.Path)
	if err != nil {
		return nil, err
	}

	c.SetEntries(files)

	return &c, nil
}

func (adb *articledb) load(path string) (*models.Article, *models.Category, error) {
	absPath := filepath.Join(adb.rootDir, path)

	c, err := adb.constructCategory(path)
	if err != nil {
		return nil, nil, err
	}

	if c.Path == path {
		return nil, c, nil
	}

	a := models.Article{
		Path:     path,
		Meta:     models.Metadata{},
		Category: c,
	}

	_, err = os.Stat(absPath + ".md")
	if err != nil {
		return nil, nil, err
	}

	articleData, err := ioutil.ReadFile(absPath + ".md")
	if err != nil {
		return nil, nil, err
	}

	err = a.Parse(articleData)
	if err != nil {
		return nil, nil, err
	}

	return &a, c, nil
}

// Write the article to disk. Also creates all relevant directories.
func (adb *articledb) Write(article *models.Article) error {
	adb.startWriteTransaction()
	defer adb.endWriteTransaction()
	return adb.write(article)
}

func (adb *articledb) write(article *models.Article) error {
	path := filepath.Join(adb.rootDir, article.Path+".md")
	err := os.MkdirAll(filepath.Dir(path), os.ModePerm)
	if err != nil {
		return err
	}

	file, err := os.OpenFile(path, os.O_RDWR|os.O_CREATE|os.O_TRUNC, os.ModePerm)
	defer file.Close()
	if err != nil {
		return err
	}

	serialized, err := article.Serialize()
	if err != nil {
		return err
	}

	_, err = file.Write(serialized)
	if err != nil {
		return err
	}

	return nil
}

package models

import (
	"os"
	"path/filepath"
	"testing"
)

func removeTestDB(path string) {
	_, err := os.Stat(path)

	if os.IsNotExist(err) {
		return
	}

	os.Remove(path)
}

func TestEncodeDecodeUser(t *testing.T) {
	testStoragePath := filepath.Join(os.TempDir(), "_TestEncodeDecodeUser.db")
	defer removeTestDB(testStoragePath)

	user, err := NewUser("test@example.com", "Not Bob", "123456789", false)
	if err != nil {
		t.Error(err)
	}

	ustr, err := LoadUserStorage(testStoragePath)
	if err != nil {
		t.Error(err)
	}

	err = ustr.AddUser(user)
	if err != nil {
		t.Error(err)
	}

	err = ustr.Save()
	if err != nil {
		t.Error(err)
	}

	ustr, err = LoadUserStorage(testStoragePath)
	if err != nil {
		t.Error(err)
	}

	err = ustr.AddUser(user)
	if err.Error() != "User Already Exists" {
		t.Error(err)
	}
}

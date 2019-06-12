package models

import (
	"fmt"
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

func TestConcurrency(t *testing.T) {
	const numUsers = 5
	dones := make(chan *User)
	users := []*User{}

	testStoragePath := filepath.Join(os.TempDir(), "_TestConcurrency.db")
	defer removeTestDB(testStoragePath)

	ustr, err := LoadUserStorage(testStoragePath)
	if err != nil {
		t.Error(err)
	}

	for i := 0; i < numUsers; i++ {
		user, err := NewUser(fmt.Sprintf("test_%v@example.com", i), "Not Bob", "123456789", false)
		if err != nil {
			t.Error(err)
		}

		users = append(users, user)
	}

	for i := 0; i < numUsers; i++ {
		user := users[i]
		go func() {
			err := ustr.AddUser(user)
			if err != nil {
				t.Error(err)
			}

			err = ustr.Save()
			if err != nil {
				t.Error(err)
			}

			dones <- user
		}()
	}

	d := 0
	for range dones {
		d++
		if d == numUsers {
			break
		}
	}

	ustr, err = LoadUserStorage(testStoragePath)
	if err != nil {
		t.Error(err)
	}

	if len(ustr.GetUsers()) != numUsers {
		t.Errorf("Number of users in store should be %v, but is %v", numUsers, len(ustr.GetUsers()))
	}
}

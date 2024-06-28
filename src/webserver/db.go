package webserver

import (
	"encoding/json"
	"fmt"
	"math/rand"
	"os"
	"strings"
	"time"

	"github.com/google/uuid"
)

var users = make([]*User, 0, 10)
var logs = make([]*Log, 0, 10)

func init() {
	if err := loadUsers(); err != nil {
		panic(err)
	}
	createUser("admin", "admin")
	createUser("test", "test")
}

func createProfile() string {
	v := rand.Int31n(70) + 1
	return "https://i.pravatar.cc/150?img=" + fmt.Sprint(v)
}

func createUser(username, password string) (*User, error) {
	for _, u := range users {
		if u.Username == username {
			return nil, fmt.Errorf("user already exists")
		}
	}

	secret := createSecret()

	new_user := User{
		ID:        uuid.New(),
		Username:  strings.ToLower(username),
		Password:  strings.ToLower(password),
		CreatedAt: time.Now(),
		Secret:    secret,
		Avatar:    createProfile(),
	}

	users = append(users, &new_user)

	go func() {
		for _, conn := range userConnections {
			conn.C <- new_user
		}
	}()

	go func() {
		presistUsers()
	}()

	return &new_user, nil
}

func deleteUser(id uuid.UUID) error {

	filteredUsers := make([]*User, 0, len(users))
	found := false

	for _, u := range users {
		if u.ID == id {
			found = true
			continue
		}
		filteredUsers = append(filteredUsers, u)
	}

	if !found {
		return fmt.Errorf("user not found")
	}

	users = filteredUsers

	go func() {
		presistUsers()
	}()

	return nil
}

const _USER_FILE_PATH = "./data/users.json"

func presistUsers() error {
	d, err := json.Marshal(users)
	if err != nil {
		return err
	}

	os.WriteFile(_USER_FILE_PATH, d, 0600)

	return nil
}

func loadUsers() error {
	raw, err := os.ReadFile(_USER_FILE_PATH)
	if err != nil {
		if err := os.WriteFile(_USER_FILE_PATH, []byte("[]"), 0600); err != nil {
			return err
		}
		return nil
	}

	err = json.Unmarshal(raw, &users)
	if err != nil {
		if err := os.WriteFile(_USER_FILE_PATH, []byte("[]"), 0600); err != nil {
			return err
		}
		return nil
	}

	return nil
}

package mongoclient

import (
	"encoding/json"
	"fmt"
	"io"
	"io/ioutil"
	"net/http/httptest"
	"os"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestLogin(t *testing.T) {
	w := httptest.NewRecorder()
	cred := Credential{
		Id:       "soku",
		Password: "mypwd",
	}
	login(w, cred)
	assert.Equal(t, w.Code, 200)
}

func TestBadLogin(t *testing.T) {
	w := httptest.NewRecorder()
	cred := Credential{
		Id:       "soku",
		Password: "badPwd",
	}
	login(w, cred)
	assert.NotEqual(t, w.Code, 200)
}

func TestGetList(t *testing.T) {
	TestLogin(t)
	w := httptest.NewRecorder()
	getUsersList(w)
	assert.Equal(t, w.Code, 200)

	b, err := io.ReadAll(w.Body)
	assert.NoError(t, err)

	fmt.Println(string(b))
}

func TestGetUser(t *testing.T) {
	TestLogin(t)
	w := httptest.NewRecorder()
	getUser(w, "KFPy6bbMWtPZ7CSqrQ7Qy3ilM2EYV8VzKPotC7SURNtSDhm1N2Q2POO94MMWuEyLjMPstBOAyQX0JBVsDQFgFrchfz42ObW3NpLE")
	assert.Equal(t, w.Code, 200)

	b, err := io.ReadAll(w.Body)
	assert.NoError(t, err)

	fmt.Println(string(b))
}

func TestGeUpdate(t *testing.T) {
	TestLogin(t)
	w := httptest.NewRecorder()
	mUser := map[string]interface{}{
		"data": "shortData",
	}
	_ = os.MkdirAll(DIRECTORY, os.ModePerm)
	updateUser(w, "KFPy6bbMWtPZ7CSqrQ7Qy3ilM2EYV8VzKPotC7SURNtSDhm1N2Q2POO94MMWuEyLjMPstBOAyQX0JBVsDQFgFrchfz42ObW3NpLE", mUser)
	assert.Equal(t, w.Code, 200)

	getUser(w, "KFPy6bbMWtPZ7CSqrQ7Qy3ilM2EYV8VzKPotC7SURNtSDhm1N2Q2POO94MMWuEyLjMPstBOAyQX0JBVsDQFgFrchfz42ObW3NpLE")
	b, err := io.ReadAll(w.Body)
	assert.NoError(t, err)

	fmt.Println(string(b))
}

func TestDelete(t *testing.T) {
	TestLogin(t)
	w := httptest.NewRecorder()
	deleteUser(w, "1qS9OI4YX8daKvHpwvhrUt6PVnG6MLQMemeFirBdqzEjwibcE1y1EZJELvXWi6w7hU9GwHMQ0RgVc3uWEOEJBbwolVD7rqIUgcwN")
	assert.Equal(t, w.Code, 200)

	getUser(w, "1qS9OI4YX8daKvHpwvhrUt6PVnG6MLQMemeFirBdqzEjwibcE1y1EZJELvXWi6w7hU9GwHMQ0RgVc3uWEOEJBbwolVD7rqIUgcwN")
	assert.NotEqual(t, w.Code, 200)
}

func TestAdd(t *testing.T) {
	TestLogin(t)
	w := httptest.NewRecorder()
	jsonFile, err := os.Open("Dataset")
	assert.NoError(t, err)

	byteValue, err := ioutil.ReadAll(jsonFile)
	mUsers := make([]map[string]interface{}, 0)
	err = json.Unmarshal(byteValue, &mUsers)
	assert.NoError(t, err)

	addUsers(w, mUsers)
	assert.Equal(t, w.Code, 200)
}

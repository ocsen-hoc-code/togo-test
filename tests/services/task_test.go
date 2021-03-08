package task

import (
	"bytes"
	"database/sql"
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"net/http/httptest"
	"testing"
	"time"

	"github.com/manabie-com/togo/internal/services"
	"github.com/manabie-com/togo/internal/storages"
	sqllite "github.com/manabie-com/togo/internal/storages/sqlite"
	_ "github.com/mattn/go-sqlite3"
	"github.com/stretchr/testify/assert"
)

func TestLogin(t *testing.T) {
	db, err := sql.Open("sqlite3", "./../../data.db")
	if err != nil {
		log.Fatal("error opening db", err)
	}

	server := httptest.NewServer(&services.ToDoService{
		JWTKey: "wqGyEBBfPK9w3Lxw",
		Store: &sqllite.LiteDB{
			DB: db,
		},
	})
	defer server.Close()

	client := &http.Client{}
	req, _ := http.NewRequest("GET", server.URL+"/login", nil)
	query := req.URL.Query()
	query.Add("user_id", "firstUser")
	query.Add("password", "example")
	req.URL.RawQuery = query.Encode()
	res, _ := client.Do(req)
	var body map[string]interface{}

	json.NewDecoder(res.Body).Decode(&body)
	assert.Equal(t, 200, res.StatusCode, "Request failed")
	assert.NotEmpty(t, body["data"], "Token is empty")
}

func TestPostTask(t *testing.T) {
	db, err := sql.Open("sqlite3", "./../../data.db")
	if err != nil {
		log.Fatal("error opening db", err)
	}

	server := httptest.NewServer(&services.ToDoService{
		JWTKey: "wqGyEBBfPK9w3Lxw",
		Store: &sqllite.LiteDB{
			DB: db,
		},
	})
	defer server.Close()

	userID := "firstUser"
	password := "example"
	content := "Binh Minh Dep Trai"
	client := &http.Client{}
	req, _ := http.NewRequest("GET", server.URL+"/login", nil)
	query := req.URL.Query()
	query.Add("user_id", userID)
	query.Add("password", password)
	req.URL.RawQuery = query.Encode()
	res, _ := client.Do(req)
	var body map[string]interface{}

	json.NewDecoder(res.Body).Decode(&body)
	assert.Equal(t, 200, res.StatusCode, "Request failed")
	assert.NotEmpty(t, body["data"], "Token is empty")

	values := map[string]string{"content": content}
	jsonData, _ := json.Marshal(values)
	req, _ = http.NewRequest("POST", server.URL+"/tasks", bytes.NewBuffer(jsonData))
	req.Header.Set("Authorization", fmt.Sprintf("%s", body["data"]))
	res, _ = client.Do(req)
	json.NewDecoder(res.Body).Decode(&body)
	assert.Equal(t, 200, res.StatusCode, "Request failed")
	assert.NotEmpty(t, body["data"], "Data not empty")
	task := storages.Task{}
	dataByte, _ := json.Marshal(body["data"])
	json.Unmarshal(dataByte, &task)
	assert.Equal(t, userID, task.UserID, "UserID is not correctly")
	assert.Equal(t, content, task.Content, "Content is not correctly")
}

func TestGetTask(t *testing.T) {
	db, err := sql.Open("sqlite3", "./../../data.db")
	if err != nil {
		log.Fatal("error opening db", err)
	}

	server := httptest.NewServer(&services.ToDoService{
		JWTKey: "wqGyEBBfPK9w3Lxw",
		Store: &sqllite.LiteDB{
			DB: db,
		},
	})
	defer server.Close()

	userID := "firstUser"
	password := "example"
	createdDate := time.Now().Format("2006-01-02")
	client := &http.Client{}
	req, _ := http.NewRequest("GET", server.URL+"/login", nil)
	query := req.URL.Query()
	query.Add("user_id", userID)
	query.Add("password", password)
	req.URL.RawQuery = query.Encode()
	res, _ := client.Do(req)
	var body map[string]interface{}

	json.NewDecoder(res.Body).Decode(&body)
	assert.Equal(t, res.StatusCode, 200, "Request failed")
	assert.NotEmpty(t, body["data"], "Token is empty")

	req, _ = http.NewRequest("GET", server.URL+"/tasks", nil)
	req.Header.Set("Authorization", fmt.Sprintf("%s", body["data"]))
	query = req.URL.Query()
	query.Add("created_date", createdDate)
	req.URL.RawQuery = query.Encode()
	res, _ = client.Do(req)
	json.NewDecoder(res.Body).Decode(&body)
	assert.Equal(t, 200, res.StatusCode, "Request failed")
	assert.NotEmpty(t, body["data"], "Data not empty")
}

package tests

import (
	"bytes"
	"context"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"
	"time"

	"github.com/gorilla/mux"
	"myapp/internal/handlers"
	"myapp/internal/models"
	"github.com/stretchr/testify/assert"
)

func TestCreateCommand(t *testing.T) {
	router := mux.NewRouter()
	handlers.RegisterHandlers(router)

	cmd := models.Command{Command: "echo Hello"}
	body, _ := json.Marshal(cmd)
	req, _ := http.NewRequest("POST", "/commands", bytes.NewBuffer(body))
	resp := httptest.NewRecorder()
	router.ServeHTTP(resp, req)

	assert.Equal(t, http.StatusCreated, resp.Code)
	var createdCmd models.Command
	json.NewDecoder(resp.Body).Decode(&createdCmd)
	assert.Equal(t, "echo Hello", createdCmd.Command)
	assert.Equal(t, "running", createdCmd.Status)
}

func TestGetCommands(t *testing.T) {
	router := mux.NewRouter()
	handlers.RegisterHandlers(router)

	req, _ := http.NewRequest("GET", "/commands", nil)
	resp := httptest.NewRecorder()
	router.ServeHTTP(resp, req)

	assert.Equal(t, http.StatusOK, resp.Code)
	var commands []models.Command
	json.NewDecoder(resp.Body).Decode(&commands)
	assert.NotEmpty(t, commands)
}

func TestGetCommand(t *testing.T) {
	router := mux.NewRouter()
	handlers.RegisterHandlers(router)

	cmd := models.Command{Command: "sleep 2"}
	body, _ := json.Marshal(cmd)
	req, _ := http.NewRequest("POST", "/commands", bytes.NewBuffer(body))
	resp := httptest.NewRecorder()
	router.ServeHTTP(resp, req)

	assert.Equal(t, http.StatusCreated, resp.Code)
	var createdCmd models.Command
	json.NewDecoder(resp.Body).Decode(&createdCmd)
	time.Sleep(3 * time.Second) // Wait for the command to finish

	req, _ = http.NewRequest("GET", "/commands/"+createdCmd.ID, nil)
	resp = httptest.NewRecorder()
	router.ServeHTTP(resp, req)

	assert.Equal(t, http.StatusOK, resp.Code)
	var fetchedCmd models.Command
	json.NewDecoder(resp.Body).Decode(&fetchedCmd)
	assert.Equal(t, createdCmd.ID, fetchedCmd.ID)
	assert.Equal(t, "success", fetchedCmd.Status)
}

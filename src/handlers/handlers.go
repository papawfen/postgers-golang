package handlers

import (
	"context"
	"encoding/json"
	"net/http"
	"os/exec"
	"sync"
	"time"

	"github.com/gorilla/mux"
	"github.com/google/uuid"
	"github.com/papawfen/postgers-golang/src/db"
	"github.com/papawfen/postgers-golang/src/models"
)

var mu sync.Mutex

func RegisterHandlers(router *mux.Router) {
	router.HandleFunc("/commands", createCommand).Methods("POST")
	router.HandleFunc("/commands", getCommands).Methods("GET")
	router.HandleFunc("/commands/{id}", getCommand).Methods("GET")
}

func createCommand(w http.ResponseWriter, r *http.Request) {
	var cmd models.Command
	err := json.NewDecoder(r.Body).Decode(&cmd)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	cmd.ID = uuid.New().String()
	cmd.Status = "running"
	mu.Lock()
	_, err = db.DB.Exec(r.Context(), "INSERT INTO commands (id, command, status) VALUES ($1, $2, $3)", cmd.ID, cmd.Command, cmd.Status)
	mu.Unlock()
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	go executeCommand(cmd)

	w.WriteHeader(http.StatusCreated)
	json.NewEncoder(w).Encode(cmd)
}

func executeCommand(cmd models.Command) {
	ctx, cancel := context.WithTimeout(context.Background(), 1*time.Minute)
	defer cancel()

	outputChan := make(chan string)
	go func() {
		out, err := exec.CommandContext(ctx, "bash", "-c", cmd.Command).CombinedOutput()
		if err != nil {
			outputChan <- "time limit exceed"
		} else {
			outputChan <- string(out)
		}
	}()

	select {
	case <-ctx.Done():
		cmd.Status = "time limit exceed"
	case out := <-outputChan:
		cmd.Output = out
		if ctx.Err() == context.DeadlineExceeded {
			cmd.Status = "time limit exceed"
		} else {
			cmd.Status = "success"
		}
	}

	mu.Lock()
	defer mu.Unlock()
	_, dbErr := db.DB.Exec(context.Background(), "UPDATE commands SET output=$1, status=$2 WHERE id=$3", cmd.Output, cmd.Status, cmd.ID)
	if dbErr != nil {
		// TODO: database errors
	}
}

func getCommands(w http.ResponseWriter, r *http.Request) {
	rows, err := db.DB.Query(r.Context(), "SELECT id, command, output, status, created_at FROM commands")
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	defer rows.Close()

	var commands []models.Command
	for rows.Next() {
		var cmd models.Command
		if err := rows.Scan(&cmd.ID, &cmd.Command, &cmd.Output, &cmd.Status, &cmd.CreatedAt); err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}
		commands = append(commands, cmd)
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(commands)
}

func getCommand(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	id := vars["id"]

	var cmd models.Command
	err := db.DB.QueryRow(r.Context(), "SELECT id, command, output, status, created_at FROM commands WHERE id=$1", id).Scan(&cmd.ID, &cmd.Command, &cmd.Output, &cmd.Status, &cmd.CreatedAt)
	if err != nil {
		http.Error(w, err.Error(), http.StatusNotFound)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(cmd)
}

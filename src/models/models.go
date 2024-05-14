package models

import "time"

type Command struct {
	ID        string    `json:"id"`
	Command   string    `json:"command"`
	Output    string    `json:"output"`
	Status    string    `json:"status"`
	CreatedAt time.Time `json:"created_at"`
}

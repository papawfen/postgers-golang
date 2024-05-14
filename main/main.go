package main

import (
	"log"
	"net/http"

	"github.com/gorilla/mux"
	"github.com/papawfen/postgers-golang/src/db"
	"github.com/papawfen/postgers-golang/src/handlers"
)

func main() {
	db.InitDB()

	router := mux.NewRouter()
	handlers.RegisterHandlers(router)

	log.Println("Server started at :8080")
	log.Fatal(http.ListenAndServe(":8080", router))
}

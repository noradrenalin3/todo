package main

import (
	"fmt"
	"database/sql"
	"log"
	_ "github.com/mattn/go-sqlite3"

	"strings"
	"net/http"
	"encoding/json"
	"github.com/rs/cors"
)
/*

func jsonResponse(w http.ResponseWriter, statusCode int, data interface{}) {
}
*/
type Response struct {
	Status    string    `json:"status"`
	Lists     []List    `json:"lists"`
	Todos     []Task 	`json:"todos"`
}

func main() {
	sqldb, err := sql.Open("sqlite3", "./data/todos.db")
	checkErr(err)
	defer sqldb.Close()

	db := &DB{sqldb}
	db.Exec("PRAGMA foreign_keys = ON")
	createTables(db)

	//fmt.Println(db.GetAll())

	mux := routes(db)
	log.Print("Listening...")

	c := cors.New(cors.Options {
		AllowedOrigins: []string{"*"},
		AllowedMethods: []string{
			http.MethodGet,
			http.MethodPost,
			http.MethodPut,
			http.MethodDelete,
		},
	})

	handler := c.Handler(mux)
	http.ListenAndServe(":3000", handler)
}

func routes(db *DB) *http.ServeMux {
	mux := http.NewServeMux()

	// ALL
	mux.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json")
		lists := db.GetLists()
		todos := db.GetAll()
		json.NewEncoder(w).Encode(&Response{"success", lists, todos})
	})

	// LIST
	mux.HandleFunc("GET /list/{listID}",
	func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json")
		listID := r.PathValue("listID")

		lists := db.GetLists()
		todos := db.GetList(listID)

		json.NewEncoder(w).Encode(&Response{"success", lists, todos})
	})

	// ADD
	mux.HandleFunc("POST /list",
	func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Access-Control-Allow-Origin", "*")
		ct := r.Header.Get("Content-Type")
		mt := strings.ToLower(strings.TrimSpace(strings.Split(ct, ";")[0]))

		if ct != "" && mt != "application/json" {
			msg := "Content-Type header is not application/json"
			http.Error(w, msg, http.StatusUnsupportedMediaType)
			return
		}

		r.Body = http.MaxBytesReader(w, r.Body, 1048576)

		var newTask Task
		json.NewDecoder(r.Body).Decode(&newTask)
		db.CreateTask(newTask)

		fmt.Println("Added: ")
		print(newTask)
	})

	// UPDATE
	mux.HandleFunc("PUT /list",
	func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Access-Control-Allow-Origin", "*")
		ct := r.Header.Get("Content-Type")
		mt := strings.ToLower(strings.TrimSpace(strings.Split(ct, ";")[0]))

		if ct != "" && mt != "application/json" {
			msg := "Content-Type header is not application/json"
			http.Error(w, msg, http.StatusUnsupportedMediaType)
			return
		}

		r.Body = http.MaxBytesReader(w, r.Body, 1048576)

		var updated Task
		json.NewDecoder(r.Body).Decode(&updated)
		db.UpdateTask(updated)

		fmt.Println("Updated: ")
		print(updated)
	})

	// DELETE task
	mux.HandleFunc("DELETE /task/{taskID}",
	func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Access-Control-Allow-Origin", "*")
		w.Header().Set("Content-Type", "application/json")
		id := r.PathValue("taskID")

		db.DeleteTask(id)
	})

	// DELETE list
	mux.HandleFunc("DELETE /list/{listID}",
	func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Access-Control-Allow-Origin", "*")
		w.Header().Set("Content-Type", "application/json")
		id := r.PathValue("listID")

		db.DeleteList(id)
	})

	return mux
}


func checkErr(err error) {
	if err != nil {
		log.Fatal(err)
	}
}

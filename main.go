package main

import (
	"fmt"
	"net/http"
	"database/sql"
	_ "github.com/lib/pq"
    "errors"
	"log"
	"io/ioutil"
)

func openDB() *sql.DB {
	db, err := sql.Open("postgres", "user=app dbname=encryptbox password=pass")
	if err != nil {
		log.Fatal(err)
	}
	return db
}

func handleInsert(key string, value string) bool {
	db := openDB()
	_, err := db.Query(`INSERT INTO entries VALUES ($1, $2)`, key, value)
	if err != nil {
		log.Fatal(err)
		return false
	}
	return true
}

func handleSelect(key string) (string, error) {
	db := openDB()

	rows, err := db.Query(`SELECT value FROM entries WHERE key = $1`, key)
	if err != nil {
		log.Fatal(err)
	}
	if rows.Next() {
		var value string
		rows.Scan(&value)
		rows.Close()
		return value, nil
	}
	return nil, errors.New("Link not found")
}

func retrieveHandler(w http.ResponseWriter, r *http.Request) {
	t := string(ioutil.ReadFile("decrypt.html"))
	fmt.Fprintf(w, t)
}

func retrievePasteHandler(w http.ResponseWriter, r *http.Request) {
	r.ParseForm()
	key := r.FormValue("key")
	if len(key) != 44 || r.Header.Get("Content-Type") != "application/json" {
		w.WriteHeader(400)
		fmt.Fprintf(w, "Bad Request")
		return
	}
	value, err := handleSelect(key)
	if err != nil {
		w.WriteHeader(404)
		fmt.Fprintf(w, err.Error())
		return
	}
	w.WriteHeader(200)
	w.Header().Set("Content-Type", "application/json")
	response := `{"value":"` + value + `"}`
	fmt.Fprintf(w, response)
}

func postPasteHandler(w http.ResponseWriter, r *http.Request) {
	r.ParseForm()
	key := r.FormValue("key")
	value := r.FormValue("value")
	if len(key) != 44 || value == "" {
		w.WriteHeader(400)
		fmt.Fprintf(w, "Bad Request")
		return
	}
	if handleInsert(key, value) {
		w.WriteHeader(201)
		fmt.Fprintf(w, "Created")
	} else {
		w.WriteHeader(500)
		fmt.Fprintf(w, "Internal Server Error")
	}
}

func createHandler(w http.ResponseWriter, r *http.Request) {
	t := string(ioutil.ReadFile("entry.html"))
	fmt.Fprintf(w, t)
}

func main() {
	http.HandleFunc("/", createHandler)
	http.HandleFunc("/paste/", retrieveHandler)
	http.HandleFunc("/api/post/", postPasteHandler)
	http.HandleFunc("/api/retrieve/", retrievePasteHandler)
	http.ListenAndServe(":8080", nil)
}
package main

import (
	"fmt"
	"net/http"
	"database/sql"
	_ "github.com/lib/pq"
    "encoding/base64" 
    "errors"
)

func openDB() *sql.DB {
	db, err := sql.Open("postgres", "user=app dbname=encryptbox password=pass")
	if err != nil {
		log.Fatal(err)
	}

}

func handleInsert(key string, value string) {
	db := openDB()
	_, err := db.Query(`INSERT INTO entries VALUES ($1, $2)`, key, value)
	if err != nil {
		log.Fatal(err)
	}
}

func handleSelect(key string) string, error {
	db := openDB()

	rows, err := db.Query(`SELECT value FROM entries WHERE key = $1`, key)
	if err != nil {
		log.Fatal(err)
	}
	if len(rows) > 0 {
		return rows[0], nil
	}
	return nil, errors.new("Link not found")
}

func retrieveHandler(w http.ResponseWriter, r *http.Request) {
	link := r.URL.Path[6:]
	value, err := handleSelect(link)
	if err != nil {
		log.Fatal(err)
	}
	

}

func ajaxHandler(w http.ResponseWriter, r *http.Request) {

}

func handler(w http.ResponseWriter, r *http.Request) {

}

func main() {
	http.HandleFunc("/", handler)
	http.HandleFunc("/paste/", retrieveHandler)
	http.HandleFunc("/prime/", primeSite)
	http.ListenAndServe(":8080", nil)
}
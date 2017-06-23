package main

import (
	"fmt"
	"net/http"
	"database/sql"
	_ "github.com/lib/pq"
    "errors"
	"log"
	"io/ioutil"
	"encoding/json"
)

type Note struct {
	Key string `json="key"`
	Value string `json="value"`
}

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
		log.Println(err)
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
	return "", errors.New("Link not found")
}

func retrieveHandler(w http.ResponseWriter, r *http.Request) {
	t, err := ioutil.ReadFile("decrypt.html")
	if err != nil {
		log.Fatal(err)
		return
	}
	fmt.Fprintf(w, string(t))
}

func retrievePasteHandler(w http.ResponseWriter, r *http.Request) {
	var note Note
	//read body
	defer r.Body.Close()
	body, err := ioutil.ReadAll(r.Body)
	
	//get json
	err = json.Unmarshal(body, &note)
	if err != nil {
		fmt.Print(err)
		fmt.Print("error!")
		return
	}
	if len(note.Key) != 44 || r.Header.Get("Content-Type") != "application/json" {
		w.WriteHeader(400)
		log.Println(note.Key)
		fmt.Fprintf(w, "Bad Request")
		return
	}
	value, err := handleSelect(note.Key)
	if err != nil {
		w.WriteHeader(404)
		fmt.Fprintf(w, err.Error())
		return
	}
	w.WriteHeader(200)
	w.Header().Set("Content-Type", "application/json")
	response := `{"value":` + value + `}`
	fmt.Fprintf(w, response)
}

func postPasteHandler(w http.ResponseWriter, r *http.Request) {
	var note Note
	//read body
	defer r.Body.Close()
	body, err := ioutil.ReadAll(r.Body)
	
	//get json
	err = json.Unmarshal(body, &note)
	if err != nil {
		fmt.Print(err)
		fmt.Print("error!")
		return
	}
	if len(note.Key) != 44 || note.Value == "" || err != nil{
		w.WriteHeader(400)
		fmt.Fprintf(w, "Bad Request Key, %d", len(note.Key))
		return
	}
	if handleInsert(note.Key, note.Value) {
		w.WriteHeader(201)
		fmt.Fprintf(w, "Created")
	} else {
		w.WriteHeader(500)
		fmt.Fprintf(w, "Internal Server Error")
	}
}

func createHandler(w http.ResponseWriter, r *http.Request) {
	t, err := ioutil.ReadFile("entry.html")
	if err != nil {
		log.Fatal(err)
		return
	}
	fmt.Fprintf(w, string(t))
}

func jsHandler(w http.ResponseWriter, r *http.Request) {
	filePath := r.URL.Path[1:]
	t, err := ioutil.ReadFile(filePath)
	if err != nil {
		w.WriteHeader(404)
		fmt.Fprintf(w, "Page not found")
		return
	}
	w.WriteHeader(200)
	fmt.Fprintf(w, string(t))
}

func main() {
	http.HandleFunc("/", createHandler)
	http.HandleFunc("/paste/", retrieveHandler)
	http.HandleFunc("/api/post/", postPasteHandler)
	http.HandleFunc("/api/retrieve/", retrievePasteHandler)
	http.Handle("/js/", http.FileServer(http.Dir("./public")))
	http.ListenAndServe(":8080", nil)
}
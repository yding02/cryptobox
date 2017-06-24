package main

import (
	"database/sql"
	"log"
	"encoding/json"
	"io/ioutil"
	"net/http"
)

func openDB() *sql.DB {
	db, err := sql.Open("postgres", "user=app dbname=encryptbox password=pass")
	if err != nil {
		log.Println(err)
	}
	return db
}


func readBodyAndStore(r *http.Request, note *Note) error {
	defer r.Body.Close()
	body, err := ioutil.ReadAll(r.Body)
	
	//get json
	err = json.Unmarshal(body, note)
	if err != nil {
		log.Println(err)
		return err
	}
	return nil
}

/**
  Copyright 2017 Florin Patan

  Licensed under the Apache License, Version 2.0 (the "License");
  you may not use this file except in compliance with the License.
  You may obtain a copy of the License at

      http://www.apache.org/licenses/LICENSE-2.0

  Unless required by applicable law or agreed to in writing, software
  distributed under the License is distributed on an "AS IS" BASIS,
  WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
  See the License for the specific language governing permissions and
  limitations under the License.
*/

package main

import (
	"database/sql"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"time"

	"github.com/gorilla/mux"
	_ "github.com/lib/pq"
)

type (
	toodo struct {
		ID        string    `json:"id"`
		Title     string    `json:"title"`
		Completed bool      `json:"completed"`
		CreatedAt time.Time `json:"-"`

		found bool
	}

	toodos struct {
		db *sql.DB
	}
)

const (
	dbInit = `
CREATE TABLE IF NOT EXISTS toodo (
	ID        TEXT,
	Title     TEXT,
	Completed BOOLEAN,
	CreatedAt TIMESTAMP WITH TIME ZONE
);
`
	QueryInsert    = `INSERT INTO toodo (ID, Title, Completed, CreatedAt) VALUES ($1, $2, $3, $4)`
	QuerySelectAll = `SELECT * FROM toodo ORDER BY CreatedAt DESC`
	QueryExists    = `SELECT count(*) as c FROM toodo WHERE ID = $1`
	QueryUpdate    = `UPDATE toodo set ID=$1, Title=$2, Completed=$3 WHERE ID = $4`
	QueryDelete    = `DELETE FROM toodo WHERE ID = $1`
)

func (*toodos) homeHandler(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("cache-control", "no-cache")
	w.Header().Set("expires", "0")
	w.Header().Set("pragma", "no-cache")

	w.Header().Set("Location", "/ui")
	w.WriteHeader(http.StatusTemporaryRedirect)
}

func (t *toodos) exists(id string) (bool, error) {
	var c int
	row := t.db.QueryRow(QueryExists, id)
	err := row.Scan(&c)
	return c > 0, err
}

func (t *toodos) update(todo toodo) (bool, error) {
	res, err := t.db.Exec(QueryUpdate, todo.ID, todo.Title, todo.Completed, todo.ID)
	if err != nil {
		return false, err
	}
	ra, err := res.RowsAffected()
	return ra == 1, err
}

func (t *toodos) loadExisting() ([]toodo, error) {
	toodos := []toodo{}

	rows, err := t.db.Query(QuerySelectAll)
	if err != nil {
		return nil, err
	}
	// if something happens here, we want to make sure the rows are Closed
	defer rows.Close()

	for rows.Next() {
		var ID string
		var Title string
		var Completed bool
		var CreatedAt time.Time
		err := rows.Scan(&ID, &Title, &Completed, &CreatedAt)
		if err != nil {
			return nil, err
		}
		toodos = append(toodos, toodo{ID, Title, Completed, CreatedAt, true})
	}
	err = rows.Err()
	if err != nil {
		return nil, err
	}

	return toodos, err
}

func (t *toodos) delete(id string) error {
	_, err := t.db.Exec(QueryDelete, id)
	if err != nil {
		log.Printf("got error: %v\n", err)
	}
	return err
}

func (t *toodos) postToodo(w http.ResponseWriter, r *http.Request) {
	body, err := ioutil.ReadAll(r.Body)
	if err != nil {
		log.Printf("got error: %v\n", err)
		w.WriteHeader(http.StatusInternalServerError)
		fmt.Fprintf(w, "%v", err)
		return
	}
	r.Body.Close()

	toodos := []toodo{}
	err = json.Unmarshal(body, &toodos)
	if err != nil {
		log.Printf("got error: %v\n", err)
		w.WriteHeader(http.StatusBadRequest)
		fmt.Fprintf(w, "%v", err)
		return
	}

	if len(toodos) == 0 {
		log.Printf("no new toodo to add")
		t.getToodos(w, r)
		return
	}

	existingTodos, err := t.loadExisting()
	if err != nil {
		log.Printf("got error: %v\n", err)
		w.WriteHeader(http.StatusInternalServerError)
		fmt.Fprintf(w, "%v", err)
		return
	}

	markFound := func(id string) bool {
		for idx := range existingTodos {
			if existingTodos[idx].ID == id {
				existingTodos[idx].found = true
				return true
			}
		}

		return false
	}

	lenNotFound := 0
	for idx := range toodos {
		if toodos[idx].CreatedAt.IsZero() {
			toodos[idx].CreatedAt = time.Now()
		}

		if markFound(toodos[idx].ID) {
			_, err := t.update(toodos[idx])
			if err != nil {
				log.Printf("got error: %v\n", err)
				w.WriteHeader(http.StatusInternalServerError)
				fmt.Fprintf(w, "%v", err)
				return
			}
			continue
		} else {
			lenNotFound++
		}

		_, err = t.db.Exec(QueryInsert, toodos[idx].ID, toodos[idx].Title, toodos[idx].Completed, toodos[idx].CreatedAt)
		if err != nil {
			log.Printf("got error: %v\n", err)
			w.WriteHeader(http.StatusInternalServerError)
			fmt.Fprintf(w, "%v", err)
			return
		}
	}

	if lenNotFound != len(existingTodos) {
		for idx := range existingTodos {
			if !existingTodos[idx].found {
				t.delete(existingTodos[idx].ID)
			}
		}
	}

	t.getToodos(w, r)
}

func (t *toodos) getToodos(w http.ResponseWriter, r *http.Request) {
	toodos, err := t.loadExisting()
	if err != nil {
		log.Printf("got error: %v\n", err)
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	resp, err := json.Marshal(toodos)
	if err != nil {
		log.Printf("got error: %v\n", err)
	}
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	w.Write(resp)
}

func main() {
	log.SetFlags(log.LstdFlags | log.Lshortfile)
	db, err := sql.Open("postgres", "postgres://postgres:postgres@db:5432/toodo?sslmode=disable")
	if err != nil {
		log.Fatalln(err)
	}

	_, err = db.Exec(dbInit)
	if err != nil {
		log.Fatalf("	got error: %v\n", err)
	}

	toodoo := &toodos{
		db: db,
	}

	r := mux.NewRouter()

	r.PathPrefix("/ui/").Handler(http.StripPrefix("/ui", http.FileServer(http.Dir("/ui/"))))

	r.HandleFunc("/toodo", toodoo.getToodos).
		Methods("GET")

	r.HandleFunc("/toodo", toodoo.postToodo).
		Methods("POST")

	r.HandleFunc("/", toodoo.homeHandler)

	http.Handle("/", r)

	log.Println("starting server...")
	log.Println(http.ListenAndServe(":8000", r))
}

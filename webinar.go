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
	"context"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"time"

	"github.com/gorilla/mux"
	"github.com/jmoiron/sqlx"
	_ "github.com/lib/pq"
)

type (
	toodo struct {
		ID        string    `json:"id"`
		Title     string    `json:"title"`
		Completed bool      `json:"completed"`
		CreatedAt time.Time `json:"-"`
	}

	toodos struct {
		db *sqlx.DB
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
)

func (*toodos) homeHandler(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("cache-control", "no-cache")
	w.Header().Set("expires", "0")
	w.Header().Set("pragma", "no-cache")

	w.Header().Set("Location", "/ui")
	w.WriteHeader(http.StatusTemporaryRedirect)
}

func (t *toodos) exists(ctx context.Context, id string) (bool, error) {
	var c int
	err := t.db.SelectContext(ctx, &c, QueryExists, id)
	return c > 0, err
}

func (t *toodos) postToodo(w http.ResponseWriter, r *http.Request) {
	body, err := ioutil.ReadAll(r.Body)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		fmt.Fprintf(w, "%v", err)
		return
	}
	r.Body.Close()

	toodos := []toodo{}
	err = json.Unmarshal(body, &toodos)
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		fmt.Fprintf(w, "%v", err)
		return
	}

	if len(toodos) == 0 {
		t.getToodos(w, r)
		return
	}

	for idx := range toodos {
		if toodos[idx].CreatedAt.IsZero() {
			toodos[idx].CreatedAt = time.Now()
		}

		if t.exists(r.Context(), toodos[idx].ID) {
			continue
		}

		_, err = t.db.ExecContext(r.Context(), QueryInsert, toodos[idx].ID, toodos[idx].Title, toodos[idx].Completed, toodos[idx].CreatedAt)
		if err != nil {
			w.WriteHeader(http.StatusInternalServerError)
			fmt.Fprintf(w, "%v", err)
			return
		}
	}

	t.getToodos(w, r)
}

func (t *toodos) getToodos(w http.ResponseWriter, r *http.Request) {
	toodos := []toodo{}
	t.db.SelectContext(r.Context(), &toodos, QuerySelectAll)

	resp, _ := json.Marshal(toodos)
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	w.Write(resp)
}

func main() {
	db, err := sqlx.Connect("postgres", "postgres://postgres:postgres@db:5432/toodo?sslmode=disable")
	if err != nil {
		log.Fatalln(err)
	}

	db.MustExec(dbInit)

	toodoo := &toodos{
		db: db,
	}

	r := mux.NewRouter()

	r.PathPrefix("/ui/").HandlerFunc(http.StripPrefix("/ui", http.FileServer(http.Dir("/ui/"))).ServeHTTP)

	r.HandleFunc("/toodo", toodoo.getToodos).
		Methods("GET")

	r.HandleFunc("/toodo", toodoo.postToodo).
		Methods("POST")

	r.HandleFunc("/", toodoo.homeHandler)

	http.Handle("/", r)

	log.Println("starting server...")
	log.Println(http.ListenAndServe(":8000", r))
}

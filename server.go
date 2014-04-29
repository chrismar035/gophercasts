package main

import (
	"database/sql"
	"fmt"
	"github.com/go-martini/martini"
	_ "github.com/lib/pq"
	"net/http"
)

func SetupDB() *sql.DB {
	db, err := sql.Open("postgres", "user=go password=gopass host=localhost dbname=lesson4 sslmode=disable")
	PanicIf(err)

	return db
}

func PanicIf(err error) {
	if err != nil {
		panic(err)
	}
}

func main() {
	m := martini.Classic()
	m.Map(SetupDB())

	m.Get("/", func(db *sql.DB, r *http.Request, rw http.ResponseWriter) {
		rows, err := db.Query("SELECT title, author, description FROM books")
		PanicIf(err)
		defer rows.Close()

		var title, author, description string
		for rows.Next() {
			err := rows.Scan(&title, &author, &description)
			PanicIf(err)
			fmt.Fprintf(rw, "Title: %s\nAuthor: %s\nDescription: %s\n\n", title, author, description)
		}
	})

	m.Get("/hello/:name", func(params martini.Params) string {
		return "Hello " + params["name"]
	})
	m.Run()
}

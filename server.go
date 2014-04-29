package main

import (
	"database/sql"
	"github.com/go-martini/martini"
	_ "github.com/lib/pq"
	"net/http"
	"github.com/martini-contrib/render"
)

type Book struct {
	Title		string
	Author		string
	Description string
}

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
	m.Use(render.Renderer(render.Options{
		Layout: "layout",
	}))

	m.Get("/", ShowBooks)
	m.Post("/books", CreateBook)
	m.Get("/create", NewBooks)

	m.Get("/hello/:name", func(params martini.Params) string {
		return "Hello " + params["name"]
	})
	m.Run()
}

func NewBooks(r render.Render) {
	r.HTML(200, "create", nil)
}

func CreateBook(ren render.Render, r *http.Request, db *sql.DB) {
	_, err := db.Query("INSERT INTO books (title, author, description) values ($1, $2, $3)",
			r.FormValue("title"),
			r.FormValue("author"),
			r.FormValue("description"))

	PanicIf(err)
	ren.Redirect("/")
}

func ShowBooks(db *sql.DB, r *http.Request, ren render.Render) {
		search := "%" + r.URL.Query().Get("search") + "%"
		rows, err := db.Query(`SELECT title, author, description FROM books
					WHERE title ILIKE $1
					OR author ILIKE $1
					OR description ILIKE $1`, search)
		PanicIf(err)
		defer rows.Close()

		books := []Book{}
		for rows.Next() {
			book := Book{}
			err := rows.Scan(&book.Title, &book.Author, &book.Description)
			PanicIf(err)
			books = append(books, book)
		}

		ren.HTML(200, "books", books)
	}

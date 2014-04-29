package main

import "github.com/go-martini/martini"

func main() {
	m := martini.Classic()
	m.Get("/", func() (int, string) {
		return 418, "i'm a teapot"
	})
	m.Run()
}

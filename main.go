package main

import (
	"net/http"

	"github.com/go-martini/martini"
)

func main() {
	m := martini.Classic()
	m.Get("/", func() string {
		return "Hello world!"
	})
	http.Handle("/", m)

	m.Run()
}

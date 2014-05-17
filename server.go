// package main of server launches a Martini server with REST endpoints configured to manipulate Words (/words).
// Available endpoints are:
// GET /words/  	- fetch all Words.
// POST /words/		- add a Word or increment its count if present.
// GET /word/:word	- fetch a specific Word.
package main

import (
	"app/words"
	"github.com/go-martini/martini"
	"github.com/martini-contrib/binding"
	"github.com/martini-contrib/render"
	"log"
	"os"
)

var logger *log.Logger
var welcome ApiDescription

type ApiDescription struct {
	Status    string   `json:"status"`
	Endpoints []string `json:"endpoints"`
}

func init() {
	logger = log.New(os.Stderr, "[server] ", log.Lshortfile)
	welcome = ApiDescription{"ok", []string{"/words"}}
	logger.Println(words.dbmap)

}

func main() {
	m := martini.Classic()
	m.Use(render.Renderer())

	m.Get("/", func(r render.Render) {
		logger.Println("GET / called")
		r.JSON(200, welcome)
	})

	m.Get("/words", func(r render.Render) {
		list, err := words.FetchWords()

		if err != nil {
			r.JSON(404, err)
		} else {
			r.JSON(200, list)
		}
	})

	m.Get("/words/:word", func(params martini.Params, r render.Render) {
		w, err := words.FetchWord(params["word"])
		if err != nil {
			r.JSON(404, err)
		} else {
			r.JSON(200, w)
		}
	})

	m.Post("/words", binding.Bind(words.Word{}), func(word words.Word, r render.Render) {
		//TODO : check for no spaces (one word only)

		//update count for this work
		err := words.WordUpsert(word)
		//return okay
		if err != nil {
			r.JSON(404, err)
		} else {
			r.JSON(201, map[string]interface{}{"okay": "true"})
		}
	})

	m.Run()
}

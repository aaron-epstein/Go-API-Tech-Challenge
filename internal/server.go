package internal

import (
	"log"
	"net/http"
	"os"

	"github.com/go-chi/chi/middleware"
	"github.com/go-chi/chi/v5"
)

var err error

func RunServer() {
	r := InitServer()
	runServer(r)
}

func InitServer() *chi.Mux {

	r := chi.NewRouter()
	r.Use(middleware.Logger)

	r.Route("/api", func(r chi.Router) {
		r.Route("/course", func(r chi.Router) {
			r.Get("/", GetCourses)
			r.Get("/{id}", GetCourse)
			r.Post("/", CreateCourse)
			r.Put("/{id}", UpdateCourse)
			r.Delete("/{id}", DeleteCourse)
		})
		r.Route("/person", func(r chi.Router) {
			r.Get("/", GetPersons)
			r.Get("/{name}", GetPerson)
			r.Post("/", CreatePerson)
			r.Put("/{name}", UpdatePerson)
			r.Delete("/{name}", DeletePerson)
		})
	})

	return r
}

func runServer(r *chi.Mux) {
	HTTP_PORT := os.Getenv("HTTP_PORT")

	Outf("Starting server on port %v", HTTP_PORT)
	err = http.ListenAndServe(HTTP_PORT, r)
	if err != nil {
		log.Fatal("Error running server")
	}
}

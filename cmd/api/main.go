package main

import (
	"email-extractor/api/mails/routes"
	"flag"
	"log"
	"net/http"

	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
	"github.com/joho/godotenv"
	"github.com/rs/cors"
)

func main() {

	var port *string = flag.String("port", "3002", "Set a different port")

	flag.Parse()

	er := godotenv.Load("../../.env")

	if er != nil {
		log.Panic("env not found")
	}

	s := CreateServer()

	s.MountHandlers()

	println("Starting server in port: " + *port)

	http.ListenAndServe(":"+*port, s.Router)
}

type Server struct {
	Router *chi.Mux
}

type Data struct {
	Hits []any `json:"hits"`
}

func CreateServer() *Server {
	s := &Server{}
	s.Router = chi.NewRouter()
	return s
}

func (s *Server) MountHandlers() {
	myCors := cors.New(cors.Options{
		AllowedOrigins:   []string{"https://*", "http://*"},
		AllowCredentials: false,
		AllowedMethods:   []string{"GET", "POST", "PUT", "DELETE", "OPTIONS"},
		AllowedHeaders:   []string{"Accept", "Authorization", "Content-Type", "X-CSRF-Token"},
	})

	s.Router.Use(middleware.Logger)

	s.Router.Use(myCors.Handler)

	s.Router.Mount("/", routes.UserRoutes())

	s.Router.NotFound(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(404)
		w.Write([]byte("Route does not exist"))
	})

	s.Router.MethodNotAllowed(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(405)
		w.Write([]byte("Method is not valid"))
	})
}

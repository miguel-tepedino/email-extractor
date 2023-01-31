package main

import (
	"encoding/json"
	"errors"
	"flag"
	"fmt"
	"io"
	"log"
	"net/http"
	"os"
	"strconv"
	"strings"

	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
	"github.com/joho/godotenv"
	"github.com/rs/cors"
)

type Query struct {
	SearchType string `json:"search_type"`
	From       int    `json:"from"`
	MaxResults int    `json:"max_results"`
	Source     []any  `json:"_source"`
	Search     any    `json:"query"`
}

type Search struct {
	Term   string `json:"term"`
	Offset uint   `json:"offset"`
}

func main() {

	var port *string = flag.String("port", "3002", "Set a different port")

	flag.Parse()

	er := godotenv.Load("../.env")

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

	s.Router.Get("/getmails", getMails)

	s.Router.Get("/mails/{offset}", getOffsetMails)

	s.Router.Post("/search", wordSearch)

	s.Router.NotFound(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(404)
		w.Write([]byte("route does not exist"))
	})
	s.Router.MethodNotAllowed(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(405)
		w.Write([]byte("method is not valid"))
	})

}

func getOffsetMails(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")

	offsetParam := chi.URLParam(r, "offset")

	if offsetParam == "" {
		w.WriteHeader(500)
		w.Write([]byte("No offset value found"))
		return
	}

	integer, integerErr := strconv.Atoi(offsetParam)

	if integerErr != nil {
		w.WriteHeader(500)
		w.Write([]byte("Invalid offset value"))
		return
	}

	query := Query{
		SearchType: "alldocuments",
		From:       integer,
		MaxResults: 10,
		Source:     make([]any, 0),
	}

	queryjson, marshalErr := json.Marshal(query)

	if marshalErr != nil {
		w.WriteHeader(500)
		w.Write([]byte("Internal server error"))
		return
	}

	res, err := HttpRequest("http://localhost:4080/api/enron/_search", "POST", string(queryjson))

	if err != nil {
		w.WriteHeader(500)
		w.Write([]byte("Error doing query"))
		return
	}

	w.Write(res)
}

func wordSearch(w http.ResponseWriter, r *http.Request) {

	word := new(Search)

	errL := json.NewDecoder(r.Body).Decode(word)

	if errL != nil {
		fmt.Println(errL)
		w.WriteHeader(500)
		w.Write([]byte(errL.Error()))
		return
	}

	query := Query{
		SearchType: "match",
		From:       int(word.Offset),
		MaxResults: 10,
		Source:     make([]any, 0),
		Search:     word,
	}

	queryformated, errMarshal := json.Marshal(query)

	if errMarshal != nil {
		w.WriteHeader(500)
		w.Write([]byte("Error generating query"))
		return
	}

	res, err := HttpRequest("http://localhost:4080/api/enron/_search", "POST", string(queryformated))

	if err != nil {
		w.WriteHeader(500)
		w.Write([]byte("Error doing query"))
		return
	}

	fmt.Println(string(res))

	w.Write(res)
}

func getMails(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")

	query := Query{
		SearchType: "alldocuments",
		From:       0,
		MaxResults: 10,
		Source:     make([]any, 0),
	}

	queryFormated, errMarshal := json.Marshal(query)

	if errMarshal != nil {
		w.WriteHeader(500)
		w.Write([]byte("Error generating query"))
		return
	}

	res, err := HttpRequest("http://localhost:4080/api/enron/_search", "POST", string(queryFormated))
	if err != nil {
		w.WriteHeader(500)
		w.Write([]byte("Error doing query"))
		return
	}

	w.Write(res)
}

func HttpRequest(url string, method string, data string) ([]byte, error) {

	req, err := http.NewRequest(method, url, strings.NewReader(data))

	if err != nil {
		return nil, errors.New(err.Error())
	}

	req.SetBasicAuth(os.Getenv("ZINCSEARCH_USERNAME"), os.Getenv("ZINCSEARCH_PASS"))
	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("User-Agent", "Mozilla/5.0 (Macintosh; Intel Mac OS X 10_15_4) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/81.0.4044.138 Safari/537.36")

	res, resErr := http.DefaultClient.Do(req)

	if err != nil {
		return nil, errors.New(resErr.Error())
	}

	defer res.Body.Close()

	body, dataErr := io.ReadAll(res.Body)

	if dataErr != nil {
		return nil, errors.New(dataErr.Error())
	}

	return body, nil
}

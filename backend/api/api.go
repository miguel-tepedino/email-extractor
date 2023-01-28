package main

import (
	"encoding/json"
	"errors"
	"flag"
	"io"
	"net/http"
	"strconv"
	"strings"

	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
)

type Query struct {
	SearchType string `json:"search_type"`
	From       int    `json:"from"`
	MaxResults int    `json:"max_results"`
	Source     []any  `json:"_source"`
}

func main() {

	var port *string = flag.String("port", "3000", "Set a different port")

	flag.Parse()

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
	s.Router.Use(middleware.Logger)

	s.Router.Get("/getmails", getMails)

	s.Router.Get("/mails/{offset}", getOffsetMails)
}

func getOffsetMails(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Access-Control-Allow-Origin", "*")
	w.Header().Set("Content-Type", "application/json")

	offsetParam := chi.URLParam(r, "offset")

	if offsetParam == "" {
		w.WriteHeader(400)
		w.Write([]byte("No offset value found"))
		return
	}

	integer, integerErr := strconv.Atoi(offsetParam)

	if integerErr != nil {
		w.WriteHeader(400)
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

func getMails(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Access-Control-Allow-Origin", "*")
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

	req.SetBasicAuth("lambda", "05111998")
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

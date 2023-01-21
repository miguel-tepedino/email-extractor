package main

import (
	"errors"
	"flag"
	"io"
	"net/http"
	"strings"

	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
)

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

func CreateServer() *Server {
	s := &Server{}
	s.Router = chi.NewRouter()
	return s
}

func (s *Server) MountHandlers() {
	s.Router.Use(middleware.Logger)

	s.Router.Get("/", HelloWorld)

	s.Router.Get("/getmails", GetMails)
}

func HelloWorld(w http.ResponseWriter, r *http.Request) {
	w.Write([]byte("Hello World!"))
}

func GetMails(w http.ResponseWriter, r *http.Request) {

	query := `{
        "search_type": "matchall",
        "from": 0,
        "max_results": 20,
        "_source": []
    }`

	res, err := HttpRequest("http://localhost:4080/api/games3/_search", "POST", query)
	if err != nil {
		w.Write([]byte(err.Error()))
	}
	w.Write(res)
}

func HttpRequest(url string, method string, data string) ([]byte, error) {

	req, err := http.NewRequest(method, url, strings.NewReader(data))

	if err != nil {
		return nil, errors.New(err.Error())
	}

	req.SetBasicAuth("lambda", "051111998")

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

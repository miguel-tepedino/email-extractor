package routes

import (
	"email-extractor/api/mails/types"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"net/http"
	"os"
	"strconv"
	"strings"

	"github.com/go-chi/chi/v5"
)

func UserRoutes() *chi.Mux {
	router := chi.NewRouter()

	router.Get("/getmails", getMails)

	router.Get("/mails/{offset}", getOffsetMails)

	router.Post("/search", wordSearch)

	return router
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

	query := types.Query{
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

	word := new(types.Search)

	errL := json.NewDecoder(r.Body).Decode(word)

	if errL != nil {
		fmt.Println(errL)
		w.WriteHeader(400)
		w.Write([]byte(errL.Error()))
		return
	}

	query := types.Query{
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

	w.Write(res)
}

func getMails(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")

	query := types.Query{
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
		w.Write([]byte(err.Error()))
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

	fmt.Println(res.Status)

	if dataErr != nil {
		return nil, errors.New(dataErr.Error())
	}

	if res.StatusCode == http.StatusBadRequest {
		return nil, errors.New(string(body))
	}

	return body, nil
}

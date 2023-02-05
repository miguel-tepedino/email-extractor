package downloader

import (
	"flag"
	"fmt"
	"io"
	"log"
	"net/http"
	"os"
	"strconv"
	"sync"
	"time"
)

var url string = "http://ipv4.download.thinkbroadband.com/5MB.zip"

var client = &http.Client{
	Timeout: time.Second * 30,
}

var connN = flag.String("connections", "16", "Number of connections")

func Downloader() {

	direc := os.Mkdir("parts", 0777)

	if direc != nil && !(direc != os.ErrExist) {
		log.Fatalln(direc)
	}

	res, errQ := http.Head(url)

	if errQ != nil {
		log.Fatalln(errQ)
	}

	acceptRanges := res.Header.Get("accept-ranges")

	lenght, errGCL := strconv.Atoi(res.Header.Get("content-length"))

	if errGCL != nil {
		log.Fatalln(errGCL)
	}

	if acceptRanges == "none" {
		*connN = "1"
		fmt.Println("Connection only accept 1 connection")
	}

	conns, err := strconv.Atoi(*connN)

	if err != nil {
		log.Panic("Invalid connections input")
	}

	length_limit := lenght / conns

	wg := &sync.WaitGroup{}

	for i := 0; i < conns; i++ {
		wg.Add(1)
		initialBytes := length_limit * i
		finalBytes := length_limit * (i + 1)

		if i-1 == conns {
			finalBytes = lenght
		}

		client := &http.Client{}

		go func(minBytes int, maxBytes int, i int, wg *sync.WaitGroup, client *http.Client) {

			defer wg.Done()

			req, errReq := http.NewRequest("GET", url, nil)

			if errReq != nil {
				log.Fatalln(errReq)
			}

			range_header := "bytes=" + strconv.Itoa(minBytes) + "-" + strconv.Itoa(maxBytes-1)

			req.Header.Set("Range", range_header)

			res, errRes := client.Do(req)

			if errRes != nil {
				log.Fatalln(errRes)
			}

			defer res.Body.Close()

			fileName := "part" + strconv.Itoa(i)

			file, errFile := os.Create("./parts/" + fileName)

			if errFile != nil {
				log.Fatal(errFile)
			}

			defer file.Close()

			written, errC := io.Copy(file, res.Body)

			for {
				fmt.Println("Part", i, " ", written)
				if errC == nil {
					return
				}
			}

		}(initialBytes, finalBytes, i, wg, client)
	}

	wg.Wait()

	fmt.Println(acceptRanges)

	fmt.Println(lenght)

}

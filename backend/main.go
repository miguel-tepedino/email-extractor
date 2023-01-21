package main

import (
	"bufio"
	"encoding/json"
	"flag"
	"fmt"
	"log"
	"net/http"
	"os"
	"runtime/pprof"
	"strings"
	"sync"
)

type Email struct {
	Date    string `json:"Date"`
	From    string `json:"From"`
	To      string `json:"To"`
	Subject string `json:"Subject"`
	Body    string `json:"Body"`
	XFrom   string `json:"X-From"`
}

var cpuprofile = flag.String("cpuprofile", "", "write cpu profile to file")

func main() {

	flag.Parse()

	if *cpuprofile != "" {
		f, err := os.Create(*cpuprofile)
		if err != nil {
			log.Fatal(err)
		}
		pprof.StartCPUProfile(f)
		defer pprof.StopCPUProfile()
	}

	myfilepath := os.Args[1:]

	if myfilepath == nil {
		log.Fatal("Please add file direction")
	}

	data, err := os.Open(myfilepath[0])

	if err != nil {
		log.Fatal(err)
	}

	defer data.Close()

	fileScanner := bufio.NewScanner(data)

	emails := make(chan Email)

	isFinished := make(chan bool)

	var wg sync.WaitGroup

	wg.Add(1)

	go ProcessLine(fileScanner, emails, &wg, isFinished)

	go func() {
		wg.Wait()
		close(emails)
	}()

	var wg2 sync.WaitGroup

	for email := range emails {
		wg2.Add(1)
		fmt.Println("----------------------------------------------------")
		go ZincSearchIngestion(email, &wg2)
		fmt.Println("----------------------------------------------------")
	}

	wg2.Wait()

	fmt.Println("Program finished")
}

func ProcessLine(buffer *bufio.Scanner, emails chan<- Email, wg *sync.WaitGroup, isFinished chan<- bool) {
	var date, from, to, xfrom, subject, body string
	var isfirst bool = true
	for buffer.Scan() {
		wg.Add(1)
		if strings.HasPrefix(buffer.Text(), "Date: ") {
			date = strings.TrimPrefix(buffer.Text(), "Date: ")
			if !isfirst {
				email := Email{Date: date, From: from, To: to, XFrom: xfrom, Body: body, Subject: subject}
				emails <- email
				body = ""
				isfirst = true

			}
		} else if strings.HasPrefix(buffer.Text(), "To: ") {
			to = FormatText(buffer.Text(), "To: ")
		} else if strings.HasPrefix(buffer.Text(), "From: ") {
			from = FormatText(buffer.Text(), "From: ")
		} else if strings.HasPrefix(buffer.Text(), "X-From: ") {
			xfrom = FormatText(buffer.Text(), "X-From: ")
		} else if strings.HasPrefix(buffer.Text(), "Subject: ") {
			subject = FormatText(buffer.Text(), "Subject: ")
		} else if buffer.Text() != "" {
			body += buffer.Text() + "\n"
			isfirst = false
		}
		wg.Done()
	}
	wg.Done()
	isFinished <- true
}

func FormatText(str string, prefix string) string {
	line := strings.TrimPrefix(str, prefix)
	return line
}

func ZincSearchIngestion(email Email, wg2 *sync.WaitGroup) {
	data, err := json.Marshal(email)

	if err != nil {
		fmt.Println("Error formating data")
	}

	req, err := http.NewRequest("POST", "http://localhost:4080/api/games3/_doc", strings.NewReader(string(data)))
	if err != nil {
		fmt.Println(err)
	}
	req.SetBasicAuth("lambda", "051111998")
	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("User-Agent", "Mozilla/5.0 (Macintosh; Intel Mac OS X 10_15_4) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/81.0.4044.138 Safari/537.36")

	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		fmt.Println(err)
	}
	defer resp.Body.Close()
	wg2.Done()
	log.Println(resp.StatusCode)
}

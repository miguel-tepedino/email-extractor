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
}

var cpuprofile = flag.String("cpuprofile", "", "write cpu profile to file")

var file = flag.String("file", "", "File to extract")

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

	myfilepath := *file

	if myfilepath == "" {
		log.Fatal("Please add file direction")
	}

	data, err := os.Open(myfilepath)

	if err != nil {
		log.Fatal(err)
	}

	defer data.Close()

	fileScanner := bufio.NewScanner(data)

	fileScanner.Split(bufio.ScanLines)

	emails := make(chan Email)

	var wg sync.WaitGroup

	newmail := &Email{
		From:    "",
		To:      "",
		Subject: "",
		Date:    "",
		Body:    "",
	}

	var isNewDateLine *bool

	*isNewDateLine = true

	wg.Add(1)

	go func() {
		for fileScanner.Scan() {
			wg.Add(1)
			newmail.ProcessLine(fileScanner.Text(), emails, &wg, isNewDateLine)
			wg.Done()
		}
		wg.Done()
	}()

	go func() {
		wg.Wait()
		close(emails)
	}()

	var wg2 sync.WaitGroup

	for email := range emails {
		go ZincSearchIngestion(email, &wg2)
	}

	wg2.Wait()

	fmt.Println("Program finished")
}

func (mail *Email) ProcessLine(line string, emails chan<- Email, wg *sync.WaitGroup, isNewLine *bool) {
	if strings.HasPrefix(line, "Date: ") {
		mail.Date = strings.TrimPrefix(line, "Date: ")
		if !*isNewLine {
			emails <- *mail
			mail.Body = ""
			*isNewLine = true
		}
	} else if strings.HasPrefix(line, "To: ") {
		mail.To = FormatText(line, "To: ")
	} else if strings.HasPrefix(line, "From: ") {
		mail.From = FormatText(line, "From: ")
	} else if strings.HasPrefix(line, "Subject: ") {
		mail.Subject = FormatText(line, "Subject: ")
	} else if line != "" {
		mail.Body += line
		*isNewLine = false
	}
}

func FormatText(str string, prefix string) string {
	line := strings.TrimPrefix(str, prefix)
	return line
}

func ZincSearchIngestion(email Email, wg2 *sync.WaitGroup) {
	wg2.Add(1)
	data, err := json.Marshal(email)

	if err != nil {
		fmt.Println("Error formating data")
		return
	}

	req, err := http.NewRequest("POST", "http://localhost:4080/api/games3/_doc", strings.NewReader(string(data)))
	if err != nil {
		fmt.Println(err)
		return
	}
	req.SetBasicAuth("lambda", "05111998")
	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("User-Agent", "Mozilla/5.0 (Macintosh; Intel Mac OS X 10_15_4) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/81.0.4044.138 Safari/537.36")

	res, err := http.DefaultClient.Do(req)
	if err != nil {
		fmt.Println(err)
		return
	}

	defer res.Body.Close()

	fmt.Println(res.StatusCode)
	wg2.Done()
}

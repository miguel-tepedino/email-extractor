package main

import (
	"bufio"
	"encoding/json"
	"flag"
	"fmt"
	"io"
	"log"
	"net/http"
	"net/mail"
	"os"
	"runtime/pprof"
	"strings"

	"github.com/joho/godotenv"
)

type Email struct {
	Date    string `json:"Date"`
	From    string `json:"From"`
	To      string `json:"To"`
	Subject string `json:"Subject"`
	Body    string `json:"Body"`
}

var cpuprofileArg = flag.String("cpuprofile", "", "write cpu profile to file")

var fileArg = flag.String("file", "", "File to extract")

func main() {

	isNewDateLine := new(bool)

	*isNewDateLine = true

	emails := make(chan string)

	emailsParsed := make(chan Email)

	mailLines := new(string)

	*mailLines = ""

	flag.Parse()

	if *cpuprofileArg != "" {
		f, err := os.Create(*cpuprofileArg)
		if err != nil {
			log.Fatal(err)
		}
		pprof.StartCPUProfile(f)
		defer pprof.StopCPUProfile()
	}

	if *fileArg == "" {
		log.Fatal("Please add file direction")
	}

	godotenv.Load("../.env")

	fmt.Println(os.Getenv("ZINCSEARCH_PASS"))

	data, err := os.Open(*fileArg)

	if err != nil {
		log.Fatal(err)
	}

	defer data.Close()

	fileScanner := bufio.NewScanner(data)

	go func() {
		for fileScanner.Scan() {
			ProcessLine(fileScanner.Text(), emails, isNewDateLine, mailLines)
		}
		close(emails)
	}()

	go func() {
		for email := range emails {
			go ParseEmail(email, emailsParsed)
		}
		close(emailsParsed)
	}()

	for emailparsed := range emailsParsed {
		go ZincSearchIngestion(emailparsed)
	}

	fmt.Println("Program finished")
}

func ProcessLine(line string, emails chan<- string, isNewLine *bool, mailLines *string) {
	if strings.Contains(line, "Message-ID:") {
		if !*isNewLine {
			emails <- *mailLines
			*mailLines = ""
			*isNewLine = true
		}
		*mailLines += line + "\n"
	} else {
		*mailLines += line + "\n"
		*isNewLine = false
	}
}

func FormatText(str string, prefix string) string {
	line := strings.TrimPrefix(str, prefix)
	return line
}

func ZincSearchIngestion(email Email) {
	data, err := json.Marshal(email)

	if err != nil {
		fmt.Println("Error formating data")
		return
	}

	req, err := http.NewRequest("POST", "http://localhost:4080/api/enron/_doc", strings.NewReader(string(data)))
	if err != nil {
		fmt.Println(err)
		return
	}
	req.SetBasicAuth(os.Getenv("ZINCSEARCH_USERNAME"), os.Getenv("ZINCSEARCH_PASS"))
	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("User-Agent", "Mozilla/5.0 (Macintosh; Intel Mac OS X 10_15_4) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/81.0.4044.138 Safari/537.36")

	res, errDo := http.DefaultClient.Do(req)
	if errDo != nil {
		fmt.Println(errDo)
		return
	}
	defer res.Body.Close()

	fmt.Println(res.Status)
}

func ParseEmail(email string, emailParsed chan<- Email) {
	msg, err := mail.ReadMessage(strings.NewReader(email))
	if err != nil {
		fmt.Println(err)
		return
	}

	body, readingErr := io.ReadAll(msg.Body)
	if readingErr != nil {
		fmt.Println(readingErr)
		return
	}

	emailParsed <- Email{
		Date:    msg.Header.Get("Date"),
		From:    msg.Header.Get("From"),
		To:      msg.Header.Get("To"),
		Subject: msg.Header.Get("Subject"),
		Body:    string(body),
	}
}

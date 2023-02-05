package main

import (
	"archive/tar"
	"compress/gzip"
	"encoding/json"
	"flag"
	"fmt"
	"io"
	"io/fs"
	"log"
	"net/http"
	"net/mail"
	"os"
	"path/filepath"
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

type Myfs func(path string, d fs.DirEntry, err error, x chan string) error

var cpuprofileArg = flag.String("cpuprofile", "", "write cpu profile to file")

var download = flag.Bool("download", false, "Download de enrondataset")

var counter int64 = 0

func main() {

	isNewDateLine := new(bool)

	*isNewDateLine = true

	mailLines := new(string)

	*mailLines = ""

	var fileArg = os.Args[1]

	flag.Parse()

	// if *download {
	// 	fmt.Println("Starting Download")
	// 	downloader.Downloader()
	// }

	if *cpuprofileArg != "" {
		f, err := os.Create(*cpuprofileArg)
		if err != nil {
			log.Fatal(err)
		}
		pprof.StartCPUProfile(f)
		defer pprof.StopCPUProfile()
	}

	if fileArg == "" {
		log.Fatal("Please add file direction")
	}

	godotenv.Load("../.env")

	file, err := os.Open(fileArg)

	if err != nil {
		log.Fatal(err)
	}

	defer file.Close()

	filepath.WalkDir(strings.TrimSuffix(file.Name(), ".tgz")+"/maildir", doSomething)

	fmt.Println(counter)
}

func createEmails(pathfile string) {

	file, errFile := os.Open(pathfile)

	if errFile != nil {
		fmt.Println("Error ", errFile)
		return
	}

	msg, errMail := mail.ReadMessage(file)

	if errMail != nil {
		fmt.Println(errMail)
		return
	}
	counter++
	parseEmailAndSend(msg)

	defer file.Close()
}

func parseEmailAndSend(msg *mail.Message) {

	body, errRB := io.ReadAll(msg.Body)

	if errRB != nil {
		fmt.Println(errRB)
		return
	}

	mail := Email{
		Date:    msg.Header.Get("Date"),
		From:    msg.Header.Get("From"),
		To:      msg.Header.Get("To"),
		Subject: msg.Header.Get("Subject"),
		Body:    string(body),
	}

	ZincSearchIngestion(mail)
}

func doSomething(path string, d fs.DirEntry, err error) error {

	if err != nil {
		fmt.Println(err)
	}

	if !d.IsDir() {
		createEmails(path)
	}

	return err
}

func extractTgz(file *os.File) {
	untgzStream, errTgzStream := gzip.NewReader(file)

	if errTgzStream != nil {
		fmt.Println(errTgzStream)
		os.Exit(1)
	}

	defer untgzStream.Close()

	tarReader := tar.NewReader(untgzStream)

	for {
		header, errHeader := tarReader.Next()
		if errHeader == io.EOF {
			break
		}

		if errHeader != nil {
			fmt.Println(errHeader)
			os.Exit(1)
		}

		switch header.Typeflag {
		case tar.TypeDir:
			errDir := os.Mkdir(header.Name, 0755)
			if os.IsExist(errDir) {
				continue
			}
			if errDir != nil && errDir != os.ErrExist {
				fmt.Println("falle aqui")
				log.Fatal(errDir.Error())
			}
		case tar.TypeReg:
			newfile, errCreate := os.Create(header.Name)
			if errCreate == os.ErrExist {
				continue
			}
			if errCreate != nil {
				log.Fatal("Extracttgz failed", errCreate)
			}
			if _, errCopy := io.Copy(newfile, tarReader); errCopy != nil {
				log.Fatal("Extracttgz failed", errCopy)
			}
			newfile.Close()
		}
	}
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

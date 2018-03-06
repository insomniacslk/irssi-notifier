package main

import (
	"bufio"
	"flag"
	"fmt"
	"io"
	"log"
	"net/http"
	"os"
	"os/user"
	"time"
)

// if you want to use self-signed key/cert, generate them with openssl, e.g.:
//   openssl genrsa -out server.key 2048
//   openssl ecparam -genkey -name secp384r1 -out server.key
//   openssl req -new -x509 -sha256 -key server.key -out server.crt -days 3650

var (
	certFile    = flag.String("cert", "", "The certificate file")
	keyFile     = flag.String("key", "", "The key file")
	fnotifyFile = flag.String("fnotify-path", "", "The fnotify file (default: ~/.irssi/fnotify)")
)

func follow(reader bufio.Reader, writer io.Writer) error {
	for {
		line, _, err := reader.ReadLine()
		if err != nil && err != io.EOF {
			return fmt.Errorf("Error when reading: %v", err)
		}
		if len(line) == 0 {
			time.Sleep(time.Second)
			continue
		}
		if _, err = writer.Write(append(line, '\n')); err != nil {
			return fmt.Errorf("Error when writing: %v", err)
		}
		if f, ok := writer.(http.Flusher); ok {
			f.Flush()
		}
	}
	return nil
}

func Tailer(w http.ResponseWriter, req *http.Request) {
	log.Printf("Request from %v", req.RemoteAddr)
	fd, err := os.Open(*fnotifyFile)
	if err != nil {
		log.Printf("Error opening file: %v", err)
		return
	}
	defer fd.Close()
	if _, err = fd.Seek(0, os.SEEK_END); err != nil {
		log.Printf("Error seeking to the end of the file: %v", err)
	}
	reader := bufio.NewReader(fd)
	if err = follow(*reader, w); err != nil {
		log.Print(err)
	}
}

func parseFlags() {
	flag.Parse()
	if *certFile == "" {
		log.Fatal("Error: empty certificate file")
	}
	if *keyFile == "" {
		log.Fatal("Error: empty key file")
	}
	if *fnotifyFile == "" {
		usr, err := user.Current()
		if err != nil {
			log.Fatal(err)
		}
		*fnotifyFile = fmt.Sprintf("%s/.irssi/fnotify", usr.HomeDir)
	}
}
func main() {
	parseFlags()

	log.Print("Starting irssi tailer")
	http.HandleFunc("/tail", Tailer)
	if err := http.ListenAndServeTLS(":8080", *certFile, *keyFile, nil); err != nil {
		log.Fatal(err)
	}
}

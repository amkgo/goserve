package main

import (
	"flag"
	"fmt"
	"log"
	"net/http"
	"os"
	"os/signal"
	"strconv"
)

var hostName = flag.String("host", "localhost", "The hostname: angang.ca or 127.0.0.1")
var portNo = flag.Int("port", 8088, "The HTTP Port number to listen")
var rootDir = flag.String("dir", ".", "The directory to serve")

func main() {
	flag.Usage = func() {
		fmt.Fprintf(os.Stderr, "Usage of %s:\n", os.Args[0])
		flag.PrintDefaults()
	}
	flag.Parse()

	// current folder
	var oldDir = ""
	var err error
	if oldDir, err = os.Getwd(); err != nil {
		log.Fatalf("Can't get current directory: %v", err)
	}

	// serve the folder "dir"
	if err = os.Chdir(*rootDir); err != nil {
		if os.IsNotExist(err) {
			log.Fatalf("folder does not exist: %v", err)
		} else {
			log.Fatalf("change directory failed: %v", err)
		}
	}

	ch := make(chan os.Signal, 1)
	signal.Notify(ch, os.Interrupt)
	signal.Notify(ch, os.Kill)

	port := strconv.Itoa(*portNo)
	fmt.Printf("Serving folder %q at http://%s:%s\n", *rootDir, *hostName, port)
	go serveDir(*hostName, port, *rootDir)

	// wait for Ctrl + C
	fmt.Println("Press Ctrl + C to exit ...")
	<-ch

	// return to the previous folder
	if err = os.Chdir(oldDir); err != nil {
		log.Fatalf("Can't return to %v\n", oldDir)
	}

	fmt.Println("Done.")
}

func serveDir(host string, port string, dir string) {
	err := http.ListenAndServe(host+":"+port, http.FileServer(http.Dir(dir)))
	if err != nil {
		log.Fatalf("failed to serve %s at %s:%s\n", dir, host, port)
	}
}

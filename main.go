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
	go serveDir(port, *rootDir)

	// wait for Ctrl + C
	<-ch

	// return to the previous folder
	if err = os.Chdir(oldDir); err != nil {
		log.Fatalf("Can't return to %v", oldDir)
	}

	fmt.Println("Done.")
}

func serveDir(port string, dir string) {
	fmt.Printf("Serve folder %q at http://localhost:%s ...\n", dir, port)
	err := http.ListenAndServe(":"+port, http.FileServer(http.Dir(dir)))
	if err != nil {
		log.Fatal("failed to serve")
	}
}

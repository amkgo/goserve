package main

import (
	"flag"
	"fmt"
	"net/http"
	"os"
	"strconv"
)

var portNo = flag.Int("port", 8088, "The HTTP Port number to listen")
var dir = flag.String("dir", ".", "The directory to serve")

func main() {
	flag.Usage = func() {
		fmt.Fprintf(os.Stderr, "Usage of %s:\n", os.Args[0])
		flag.PrintDefaults()
	}
	flag.Parse()

	port := strconv.Itoa(*portNo)

	fmt.Println("Serve folder " + *dir + " at http://localhost:" + port)
	http.ListenAndServe(":"+port, http.FileServer(http.Dir(*dir)))
}

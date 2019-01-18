package main

import (
	"flag"
	"fmt"
	"log"
	"net/http"
)

func main() {
	portPtr := flag.Int("p", 8080, "The port to run the gopull-api on")
	flag.Parse()

	router := NewRouter()
	log.Fatal(http.ListenAndServe(fmt.Sprintf(":%v", *portPtr), router))
}

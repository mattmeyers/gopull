package main

import (
	"flag"
	"fmt"
	"log"
	"net/http"

	"github.com/mattmeyers/gopull"
	"github.com/spf13/viper"
)

func main() {
	gopull.NewConfig()
	viper.WatchConfig()

	portPtr := flag.Int("p", 8080, "The port to run the gopull-api on")
	flag.Parse()

	router := NewRouter()
	log.Fatal(http.ListenAndServe(fmt.Sprintf(":%v", *portPtr), router))
}

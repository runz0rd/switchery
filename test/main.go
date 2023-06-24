package main

import (
	"flag"
	"log"
	"math/rand"
	"net/http"
	"os"
	"time"
)

func main() {
	var flagAddress string
	flag.StringVar(&flagAddress, "address", ":80", "address to listen on")
	flag.Parse()

	handler := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		greeting := "there!"
		version := os.Getenv("VERSION")
		if version != "" {
			greeting = version
		}
		rs := rand.NewSource(42)
		rnd := rand.New(rs)
		time.Sleep(time.Duration(rnd.Intn(2000) * int(time.Millisecond)))
		_, _ = w.Write([]byte(greeting))
	})
	log.Printf("listening on %v", flagAddress)
	log.Fatal(http.ListenAndServe(flagAddress, handler))
}

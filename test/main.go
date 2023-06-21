package main

import (
	"flag"
	"log"
	"math/rand"
	"net/http"
	"time"
)

func main() {
	var flagAddress string
	var flagGreeting string
	flag.StringVar(&flagAddress, "address", ":80", "address to listen to ")
	flag.StringVar(&flagGreeting, "greeting", "HELLO", "sets the greeting message")
	flag.Parse()

	rs := rand.NewSource(42)
	rnd := rand.New(rs)
	handler := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		alt := r.URL.Query().Get("greeting")
		if alt != "" {
			flagGreeting = alt
		}
		time.Sleep(time.Duration(rnd.Intn(1500) * int(time.Millisecond)))
		_, _ = w.Write([]byte(flagGreeting))
	})
	log.Printf("listening on %v", flagAddress)
	log.Fatal(http.ListenAndServe(flagAddress, handler))
}

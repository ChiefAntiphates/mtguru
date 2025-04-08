package main

import (
	"encoding/json"
	"fmt"
	"net/http"
	"time"

	"github.com/rs/cors"
)

var listenerPort string = "8888"

func initHandler() http.Handler {
	mux := http.NewServeMux()

	// Add route handlers
	mux.HandleFunc("/time", getTime)

	mux.HandleFunc("POST /clicks", postNumberOfClicks)

	// cors.Default() setup the middleware with default options being
	// all origins accepted with simple methods (GET, POST). See
	// documentation below for more options.
	return cors.Default().Handler(mux)
}

func getTime(w http.ResponseWriter, r *http.Request) {
	fmt.Fprintf(w, time.Now().String())
}

func postNumberOfClicks(w http.ResponseWriter, r *http.Request) {
	type Payload struct {
		Clicks int `json:"clicks"`
	}
	var payload Payload
	err := json.NewDecoder(r.Body).Decode(&payload)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}
	fmt.Println("There have been", payload.Clicks, "clicks on the client!")
	w.WriteHeader(http.StatusNoContent)
}

func main() {

	handler := initHandler()

	fmt.Println("Server listening at port", listenerPort)
	http.ListenAndServe(":"+listenerPort, handler)
}

package main

import (
	"github.com/gorilla/mux"
	"log"
	"net/http"
)

func main() {
	r := mux.NewRouter()
	r.HandleFunc("/screenshot", ScreenshotHandler).Methods("GET")

	log.Println("Server started at :80.")
	log.Fatal(http.ListenAndServe(":80", r))
}

func ScreenshotHandler(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Server", "Go Server")
	w.WriteHeader(http.StatusOK)
	str := []byte("Server successfully working!")
	_, _ = w.Write(str)
}

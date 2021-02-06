package main

import (
	"github.com/gorilla/mux"
	"log"
	"net/http"
)

func main() {
	r := mux.NewRouter()
	r.HandleFunc("/screenshot", ScreenshotHandler).Methods("GET")
	r.Handle("/", &Server{r})

	log.Println("Server started at :80.")
	log.Fatal(http.ListenAndServe(":80", r))
}

type Server struct {
	r *mux.Router
}

func (s *Server) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	if origin := r.Header.Get("origin"); origin != "" {
		r.Header.Set("Access-Control-Allow-Origin", origin)
		r.Header.Set("Access-Control-Allow-Methods", "GET")
		r.Header.Set("Server", "Go Server")
	}

	s.r.ServeHTTP(w, r)
}

func ScreenshotHandler(w http.ResponseWriter, r *http.Request) {
	v := r.URL.Query()

	var (
		url       string
		scrollY   string
		scrollSel string
		height    string
		weight    string
	)

	if url = v.Get("url"); url == "" {
		w.WriteHeader(http.StatusBadRequest)
		_, _ = w.Write([]byte("Bad Request: URL not found"))
		return
	}
	scrollY = v.Get("y")
	scrollSel = v.Get("sel")
	if weight = v.Get("weight"); weight == "" {
		weight = "1080"
	}
	if height = v.Get("height"); height == "" {
		height = "720"
	}

	w.WriteHeader(http.StatusOK)
	str := []byte("Params: " + url + ", " + scrollY + ", " + scrollSel + ", " + height + ", " + weight)
	_, _ = w.Write(str)
}

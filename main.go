package main

import (
	"github.com/gorilla/mux"
	"image/color"
	"log"
	"net/http"
	"os"
	"strconv"
	"strings"
)

func main() {
	r := mux.NewRouter()
	r.HandleFunc("/screenshot", ScreenshotHandler).Methods("GET")

	fs := http.FileServer(http.Dir("./static/"))
	r.Handle("/", &Server{r}).Handler(fs)
	r.PathPrefix("/css").Handler(fs)
	r.PathPrefix("/js").Handler(fs)
	r.PathPrefix("/images").Handler(fs)

	r.NotFoundHandler = http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		http.Redirect(w, r, "/", http.StatusPermanentRedirect)
	})

	port := os.Getenv("PORT")
	if port == "" {
		port = "80"
	}
	log.Fatal(http.ListenAndServe(":"+port, r))
}

type Server struct {
	r *mux.Router
}

type ScreenshotOption struct {
	url        string
	width      int64
	height     int64
	scrollY    int64
	selector   string
	quality    int64
	hide       []string
	background BackgroundOption
}

type BackgroundOption struct {
	color  color.RGBA
	width  int
	height int
}

type ImageBuffer = []byte

const (
	defaultWidth     int64 = 1920
	defaultHeight    int64 = 1080
	defaultQuality   int64 = 100
	defaultMinWidth  int64 = 1820
	defaultMinHeight int64 = 980
)

func (s *Server) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	if origin := r.Header.Get("origin"); origin != "" {
		r.Header.Set("Access-Control-Allow-Origin", origin)
		r.Header.Set("Access-Control-Allow-Methods", "GET")
	}

	s.r.ServeHTTP(w, r)
}

func ScreenshotHandler(w http.ResponseWriter, r *http.Request) {
	v := r.URL.Query()

	opts := ScreenshotOption{}
	if opts.url = v.Get("url"); opts.url == "" {
		w.WriteHeader(http.StatusBadRequest)
		_, _ = w.Write([]byte("Bad Request: URL not found"))
		return
	}

	quality := v.Get("quality")
	if opts.quality = defaultQuality; quality != "" {
		if i, err := strconv.ParseInt(quality, 10, 64); err == nil {
			opts.quality = i
		}
	}

	if hide := v.Get("hide"); hide != "" {
		opts.hide = strings.Split(hide, "\r\n")
	}

	opts.selector = v.Get("selector")
	if scrollY := v.Get("scrollY"); scrollY != "" {
		if i, err := strconv.ParseInt(scrollY, 10, 64); err == nil {
			opts.scrollY = i
		}
	}

	width := v.Get("width")
	if opts.width = defaultWidth; width != "" {
		if i, err := strconv.ParseInt(width, 10, 64); err == nil {
			opts.width = i
		}
	}

	height := v.Get("height")
	if opts.height = defaultHeight; height != "" {
		if i, err := strconv.ParseInt(height, 10, 64); err == nil {
			opts.height = i
		}
	}

	if bgColor := v.Get("bgColor"); bgColor != "" {
		c, err := ParseHexColor(bgColor)
		if err == nil && !IsNoBackgroundColor(c) {
			opts.width = defaultMinWidth
			opts.height = defaultMinHeight
			opts.background.color = c
			opts.background.width = int(defaultWidth)
			opts.background.height = int(defaultHeight)
		}
	}

	buffer, err := TakeScreenshot(opts)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		_, _ = w.Write([]byte("Internal Server Error: " + err.Error()))
		return
	}

	w.Header().Set("Content-Length", strconv.Itoa(len(buffer)))
	w.WriteHeader(http.StatusOK)
	_, _ = w.Write(buffer)
}

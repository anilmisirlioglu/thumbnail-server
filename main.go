package main

import (
	"context"
	"github.com/chromedp/cdproto/emulation"
	"github.com/chromedp/cdproto/page"
	"github.com/chromedp/chromedp"
	"github.com/gorilla/mux"
	"log"
	"net/http"
	"strconv"
	"time"
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

type ScreenshotOption struct {
	url            string
	scrollY        int64
	scrollSelector string
	height         int64
	width          int64
}

type ImageBuffer = []byte

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

	opts.scrollSelector = v.Get("sel")
	if scrollY := v.Get("y"); scrollY != "" {
		if i, err := strconv.ParseInt(scrollY, 10, 64); err == nil {
			opts.scrollY = i
		}
	}

	width := v.Get("width")
	if opts.width = 1080; width != "" {
		if i, err := strconv.ParseInt(width, 10, 64); err == nil {
			opts.width = i
		}
	}

	height := v.Get("height")
	if opts.height = 720; height != "" {
		if i, err := strconv.ParseInt(height, 10, 64); err == nil {
			opts.height = i
		}
	}

	buffer, err := TakeScreenshot(opts)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		_, _ = w.Write([]byte("Internal Server Error: " + err.Error()))
		return
	}

	w.Header().Set("Content-Type", "image/png")
	w.WriteHeader(http.StatusOK)
	_, _ = w.Write(buffer)
}

func TakeScreenshot(opts ScreenshotOption) (ImageBuffer, error) {
	ctx, cancel := chromedp.NewContext(context.Background())
	defer cancel()

	ctx, cancel = context.WithTimeout(ctx, 10*time.Second)
	defer cancel()

	var buffer ImageBuffer
	if err := chromedp.Run(ctx, ScreenshotTasks(opts, &buffer)); err != nil {
		return nil, err
	}

	return buffer, nil
}

func ScreenshotTasks(opts ScreenshotOption, buffer *ImageBuffer) chromedp.Tasks {
	actions := chromedp.Tasks{
		chromedp.Navigate(opts.url),
		chromedp.WaitVisible("html", chromedp.ByQuery),
		chromedp.EmulateViewport(opts.width, opts.height, func(params *emulation.SetDeviceMetricsOverrideParams, _ *emulation.SetTouchEmulationEnabledParams) {
			params.WithPositionY(opts.scrollY)
		}),
	}

	if opts.scrollSelector != "" {
		actions = append(actions, chromedp.ScrollIntoView(opts.scrollSelector))
	}

	actions = append(actions, chromedp.ActionFunc(func(ctx context.Context) (err error) {
		*buffer, err = page.
			CaptureScreenshot().
			WithQuality(100).
			Do(ctx)

		return err
	}))

	return actions
}

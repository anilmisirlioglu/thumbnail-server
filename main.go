package main

import (
	"context"
	"fmt"
	"github.com/chromedp/cdproto/emulation"
	"github.com/chromedp/cdproto/page"
	"github.com/chromedp/cdproto/runtime"
	"github.com/chromedp/chromedp"
	"github.com/gorilla/mux"
	"log"
	"net/http"
	"strconv"
	"strings"
	"time"
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
	hide           []string
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

	if hide := v.Get("hide"); hide != "" {
		opts.hide = strings.Split(hide, "\r\n")
	}

	opts.scrollSelector = v.Get("selector")
	if scrollY := v.Get("scrollY"); scrollY != "" {
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
	}

	var scripts []string
	if len(opts.hide) != 0 {
		scripts = append(scripts, BuildInvisibleScript(opts.hide))
	}

	if opts.scrollY != 0 {
		scripts = append(scripts, BuildWindowScrollScript(opts.scrollY))
	}

	for _, script := range scripts {
		actions = append(actions, chromedp.ActionFunc(func(ctx context.Context) error {
			_, exp, err := runtime.Evaluate(script).Do(ctx)
			if err != nil {
				return err
			}

			if exp != nil {
				return exp
			}

			return nil
		}))
	}

	if opts.scrollSelector != "" {
		actions = append(actions, chromedp.ScrollIntoView(opts.scrollSelector))
	}

	actions = append(actions, chromedp.ActionFunc(func(ctx context.Context) (err error) {
		_, _, size, err := page.GetLayoutMetrics().Do(ctx)
		if err != nil {
			return err
		}

		err = emulation.
			SetDeviceMetricsOverride(opts.width, opts.height, 1, false).
			WithScreenOrientation(&emulation.ScreenOrientation{
				Type:  emulation.OrientationTypePortraitPrimary,
				Angle: 0,
			}).
			Do(ctx)
		if err != nil {
			return err
		}

		*buffer, err = page.
			CaptureScreenshot().
			WithQuality(100).
			WithClip(&page.Viewport{
				X:      size.X,
				Y:      size.Y + float64(opts.scrollY),
				Width:  float64(opts.width),
				Height: float64(opts.height),
				Scale:  1,
			}).
			Do(ctx)

		return err
	}))

	return actions
}

func BuildInvisibleScript(hide []string) string {
	return fmt.Sprintf(`[%s].forEach(item => {
		document.querySelector(item).remove()
	})`, fmt.Sprintf(
		"'%s'",
		strings.Join(hide, `', '`)),
	)
}

func BuildWindowScrollScript(scrollY int64) string {
	return fmt.Sprintf(`window.scrollTo(0, %d)`, scrollY)
}

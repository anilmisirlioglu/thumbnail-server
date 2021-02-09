package main

import (
	"bytes"
	"context"
	"fmt"
	"github.com/chromedp/cdproto/emulation"
	"github.com/chromedp/cdproto/page"
	"github.com/chromedp/cdproto/runtime"
	"github.com/chromedp/chromedp"
	"image/png"
	"strings"
	"time"
)

func TakeScreenshot(opts ScreenshotOption) (ImageBuffer, error) {
	ctx, cancel := chromedp.NewContext(context.Background())
	defer cancel()

	ctx, cancel = context.WithTimeout(ctx, 10*time.Second)
	defer cancel()

	var byteArray ImageBuffer
	if err := chromedp.Run(ctx, ScreenshotTasks(opts, &byteArray)); err != nil {
		return nil, err
	}

	if !IsNoBackgroundColor(opts.background.color) {
		img, err := ServeFrames(byteArray)
		if err != nil {
			return nil, err
		}

		draw := DrawBackground(opts, img)

		buffer := new(bytes.Buffer)
		if err := png.Encode(buffer, draw); err != nil {
			return nil, err
		}

		return buffer.Bytes(), nil
	}

	return byteArray, nil
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

	if opts.selector != "" {
		actions = append(actions, chromedp.ScrollIntoView(opts.selector))
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
			WithQuality(opts.quality).
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

/* SCRIPTS */
func BuildInvisibleScript(hide []string) string {
	return fmt.Sprintf(`[%s].forEach(item => {
		const element = document.querySelector(item)
		if(element) element.remove()
	})`, fmt.Sprintf(
		"'%s'",
		strings.Join(hide, `', '`),
	),
	)
}

func BuildWindowScrollScript(scrollY int64) string {
	return fmt.Sprintf(`window.scrollTo(0, %d)`, scrollY)
}

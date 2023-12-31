package naviga

import (
	"os"
	"time"

	"github.com/go-rod/rod"
	"github.com/go-rod/rod/lib/launcher"
)

func launchInLambda() *launcher.Launcher {
	if os.Getenv("APP_ENV") == "dev" {
		return launcher.New().Headless(false).Devtools(true)
	}

	return launcher.New().
		// Lambda extracts the layer contents into the /opt directory in your
		// function’s execution environment
		// Ref: https://docs.aws.amazon.com/lambda/latest/dg/chapter-layers.html
		Bin("/opt/chromium").

		// recommended flags to run in serverless environments
		// see https://github.com/alixaxel/chrome-aws-lambda/blob/master/source/index.ts
		Set("allow-running-insecure-content").
		Set("autoplay-policy", "user-gesture-required").
		Set("disable-component-update").
		Set("disable-domain-reliability").
		Set("disable-features", "AudioServiceOutOfProcess", "IsolateOrigins", "site-per-process").
		Set("disable-print-preview").
		Set("disable-setuid-sandbox").
		Set("disable-site-isolation-trials").
		Set("disable-speech-api").
		Set("disable-web-security").
		Set("disk-cache-size", "33554432").
		Set("enable-features", "SharedArrayBuffer").
		Set("hide-scrollbars").
		Set("ignore-gpu-blocklist").
		Set("in-process-gpu").
		Set("mute-audio").
		Set("no-default-browser-check").
		Set("no-pings").
		Set("no-sandbox").
		Set("no-zygote").
		Set("single-process").
		Set("use-gl", "swiftshader").
		Set("window-size", "1920", "1080")
}

// TakeScreenshot navigates to a specified URL using a web browser and captures a screenshot.
// The screenshot is then saved to a file at the provided path.
//
// Parameters:
//
//	url: The URL to navigate to. It should be a valid URL starting with 'http://' or 'https://'.
//	path: The absolute or relative file path where the screenshot should be saved.
//	      This should include the desired file extension (e.g., '.png', '.jpg').
//
// Returns:
//
//	error: Returns an error if any step of the operation fails.
//
// Possible Errors:
//   - Invalid URL
//   - Issues initializing or controlling the web browser
//   - File path issues
//
// Example:
//
//	err := TakeScreenshot("https://www.example.com", "/path/to/screenshot.png")
//	if err != nil {
//	  log.Fatal(err)
//	}
//
// Note: Ensure that the required browser driver is installed and accessible in your PATH.
func TakeScreenshot(website, screenshotPath string) error {
	// If Rod fails, it needs to correctly timeout before the timeout we set as the lambda fn's timeout
	// this ensures that the browser instance is properly killed and cleaned up
	//
	// these timeouts should collectively be less than the timeout we set for the lambda
	const (
		navigateTimeout    = 10 * time.Second
		navigationTimeout  = 10 * time.Second
		requestIdleTimeout = 10 * time.Second
		htmlTimeout        = 15 * time.Second
	)

	err := rod.Try(func() {
		// instantiate the chromium launcher
		launcher := launchInLambda()

		// lambda warm starts reuse environments:
		//
		// we must delete data generated by the browser,
		// otherwise repeated calls to this fn will eat up storage
		// and the lambda will fail
		defer launcher.Cleanup()

		// likewise, browser.close() will leave a zombie process
		// so we must kill the chromium processes completely
		// otherwise memory consumption will be much higher
		defer launcher.Kill()

		u := launcher.MustLaunch()

		// create a browser instance
		browser := rod.New().ControlURL(u).MustConnect()

		// open a page
		page := browser.MustPage()

		// go to the url
		page.Timeout(navigateTimeout).MustNavigate(website)

		waitNavigation := page.Timeout(navigationTimeout).MustWaitNavigation()
		waitNavigation()

		waitRequestIdle := page.Timeout(requestIdleTimeout).MustWaitRequestIdle()
		waitRequestIdle()

		page.Timeout(htmlTimeout).MustScreenshot(screenshotPath)
	})

	return err
}

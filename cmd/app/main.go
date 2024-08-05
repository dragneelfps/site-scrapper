package main

import (
	"flag"
	"log/slog"
	"os"

	"github.com/dragneelfps/site-scrapper/pkg/scrapper/factory"
	"github.com/dragneelfps/site-scrapper/pkg/scrapper/site"
	"github.com/dragneelfps/site-scrapper/pkg/selenium"
	"github.com/dragneelfps/site-scrapper/pkg/writer"
)

const (
	defaultDriverPath        = "./tools/chromedriver"
	defaultBrowserBinaryPath = "./tools/chrome-linux64/chrome"
	defaultHeadlessBrowser   = true
	defaultOutputPath        = "./build/out.csv"
)

var (
	driverPath   = flag.String("driver-path", defaultDriverPath, "chrome driver path")
	browserPath  = flag.String("chrome-binary-path", defaultBrowserBinaryPath, "chrome binary path")
	headlessMode = flag.Bool("headless-mode", defaultHeadlessBrowser, "run browser in headless mode")
	outputPath   = flag.String("output-path", defaultOutputPath, "output path for csv")
	entityId     = flag.String("entity-id", "", "entity id")
	siteStr      = flag.String("site", "", "site such as fapello")
)

func main() {
	flag.Parse()

	if len(*siteStr) == 0 {
		slog.Error("site not provided", "availableSites", site.Sites())
		os.Exit(1)
	}

	siteFlag, ok := site.FromString(*siteStr)
	if !ok {
		slog.Error("invalid site provided", "availableSites", site.Sites(), "provided", *siteStr)
		os.Exit(1)
	}

	if len(*entityId) == 0 {
		slog.Error("entity id not provided")
		os.Exit(1)
	}

	webDriver, err := selenium.New(selenium.Config{
		DriverPath:        *driverPath,
		BrowserBinaryPath: *browserPath,
		HeadlessBrowser:   *headlessMode,
	})
	defer webDriver.Stop() //nolint:errcheck
	if err != nil {
		slog.Error("unable to start driver", "error", err)
		os.Exit(1)
	}

	s, ok := factory.NewScrapperFactory().Get(siteFlag)
	if !ok {
		slog.Error("scrapper not found", "name", siteFlag)
		os.Exit(1)
	}

	entityData, err := s.GetEntity(webDriver, *entityId)
	if err != nil {
		slog.Error("get entity error", "entityID", "monika-balsai", "error", err)
		os.Exit(1)
	}

	slog.Info("got entity data", "data", entityData)

	csvWriter := writer.CSVWriter{}
	err = csvWriter.Write(entityData, *outputPath)
	if err != nil {
		slog.Error("failed to save in csv", "entityID", "monika-balsai", "error", err)
		os.Exit(1)
	}
}

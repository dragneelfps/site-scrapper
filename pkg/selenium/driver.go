package selenium

import (
	"errors"
	"fmt"
	"strings"

	"github.com/PuerkitoBio/goquery"
	"github.com/tebeka/selenium"
	"github.com/tebeka/selenium/chrome"
)

const (
	_defaultDriverPort = 4444
)

type Config struct {
	DriverPath        string
	DriverPort        int
	BrowserBinaryPath string
	BrowserArgs       []string
	HeadlessBrowser   bool
	SeleniumURL       string
}

type Driver struct {
	service *selenium.Service
	driver  selenium.WebDriver
}

func New(c Config) (*Driver, error) {
	// initialize a Chrome browser at given port or default
	driverPath := strings.TrimSpace(c.DriverPath)
	if len(driverPath) == 0 {
		return nil, errors.New("driver path cannot be empty")
	}
	driverPort := c.DriverPort
	if driverPort == 0 {
		driverPort = _defaultDriverPort
	}
	service, err := selenium.NewChromeDriverService(driverPath, driverPort)
	if err != nil {
		return nil, fmt.Errorf("create driver service: %w", err)
	}

	// configure the browser options
	browserBinaryPath := strings.TrimSpace(c.BrowserBinaryPath)
	if len(browserBinaryPath) == 0 {
		return nil, errors.New("browser binary path cannot be empty")
	}
	browserArgs := c.BrowserArgs
	if c.HeadlessBrowser {
		browserArgs = append(browserArgs, "--headless=new")
	}
	caps := selenium.Capabilities{}
	caps.AddChrome(chrome.Capabilities{
		Path: browserBinaryPath,
		Args: browserArgs,
	})

	// create a new remote client with the specified options
	driver, err := selenium.NewRemote(caps, "")
	if err != nil {
		return nil, fmt.Errorf("new remote: %w", err)
	}
	return &Driver{
		service: service,
		driver:  driver,
	}, err
}

func (d Driver) Stop() error {
	return d.service.Stop()
}

func (d Driver) GetPageSource(url string) (*goquery.Document, error) {
	err := d.driver.Get(url)
	if err != nil {
		return nil, fmt.Errorf("driver get: %w", err)
	}
	rawHtml, err := d.driver.PageSource()
	if err != nil {
		return nil, fmt.Errorf("driver get page source: %w", err)
	}
	doc, err := goquery.NewDocumentFromReader(strings.NewReader(rawHtml))
	if err != nil {
		return nil, fmt.Errorf("goquery get document from source: %w", err)
	}
	return doc, nil
}

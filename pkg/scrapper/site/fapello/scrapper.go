package fapello

import (
	"errors"
	"fmt"
	"log/slog"
	"strings"

	"github.com/PuerkitoBio/goquery"
	"github.com/dragneelfps/site-scrapper/pkg/scrapper"
	"github.com/dragneelfps/site-scrapper/pkg/selenium"
)

const (
	_entityPageFormat = "https://fapello.com/%s"
)

type FapelloScrapper struct{}

var _ scrapper.Scrapper = FapelloScrapper{}

func (s FapelloScrapper) GetEntity(webDriver *selenium.Driver, entityID string) (data scrapper.EntityData, err error) {
	if webDriver == nil {
		return data, errors.New("web driver cannot be nil")
	}
	if strings.TrimSpace(entityID) == "" {
		return data, errors.New("entity id cannot be empty")
	}

	doc, err := webDriver.GetPageSource(fmt.Sprintf(_entityPageFormat, entityID))
	if err != nil {
		return data, fmt.Errorf("web driver get page source: %w", err)
	}

	mediaSelections := doc.Find("div#content > div > a")

	errorList := make([]error, 0, len(mediaSelections.Nodes))

	mediaSelections.Each(func(_ int, mediaSel *goquery.Selection) {
		mediaHref, mediaHrefFound := mediaSel.Attr("href")
		if !mediaHrefFound || strings.TrimSpace(mediaHref) == "" {
			slog.Warn("no href attr exists", "node", mediaSel)
			return
		}

		mediaDoc, err := webDriver.GetPageSource(mediaHref)
		if err != nil {
			errorList = append(errorList, err)
			return
		}

		media, err := s.getMedia(mediaDoc)
		if err != nil {
			errorList = append(errorList, err)
			return
		}

		data.TotalMediaCount += 1
		data.MediaList = append(data.MediaList, media)
	})

	return data, errors.Join(errorList...)
}

func (s FapelloScrapper) getMedia(doc *goquery.Document) (scrapper.EntityMediaData, error) {
	imageHref, imageFound := getImage(doc)
	if imageFound {
		return scrapper.EntityMediaData{
			ID:        imageHref,
			URL:       imageHref,
			MediaType: scrapper.EntityMediaTypeImage,
		}, nil
	}
	videoHref, videoFound := getVideo(doc)
	if videoFound {
		return scrapper.EntityMediaData{
			ID:        videoHref,
			URL:       videoHref,
			MediaType: scrapper.EntityMediaTypeVideo,
		}, nil
	}
	return scrapper.EntityMediaData{}, scrapper.NewEntityNotFoundError(doc.Url.String())
}

func getImage(doc *goquery.Document) (string, bool) {
	return doc.Find("div.items-center.justify-between > a > img").Attr("src")
}

func getVideo(doc *goquery.Document) (string, bool) {
	return doc.Find("div > video > source").Attr("src")
}

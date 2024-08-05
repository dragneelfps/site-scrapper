package factory

import (
	"github.com/dragneelfps/site-scrapper/pkg/scrapper"
	"github.com/dragneelfps/site-scrapper/pkg/scrapper/site"
	"github.com/dragneelfps/site-scrapper/pkg/scrapper/site/fapello"
)

type ScrapperFactory struct {
	scrappers map[site.Site]scrapper.Scrapper
}

func NewScrapperFactory() ScrapperFactory {
	f := ScrapperFactory{
		scrappers: map[site.Site]scrapper.Scrapper{
			site.SiteFapello: fapello.FapelloScrapper{},
		},
	}

	return f
}

func (f ScrapperFactory) Get(site site.Site) (scrapper.Scrapper, bool) {
	s, ok := f.scrappers[site]
	return s, ok
}

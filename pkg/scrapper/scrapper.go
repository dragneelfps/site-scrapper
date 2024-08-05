package scrapper

import (
	"github.com/dragneelfps/site-scrapper/pkg/selenium"
)

type Scrapper interface {
	GetEntity(webDriver *selenium.Driver, entityID string) (EntityData, error)
}

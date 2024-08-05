package scrapper

import "fmt"

type EntityData struct {
	TotalMediaCount int
	URL             string
	MediaList       []EntityMediaData
}

type EntityMediaData struct {
	ID        string          `csv:"id"`
	URL       string          `csv:"url"`
	MediaType EntityMediaType `csv:"media_type"`
}

type EntityMediaType int

const (
	EntityMediaTypeImage EntityMediaType = iota + 1
	EntityMediaTypeVideo
	EntityMediaTypeText
)

func (t EntityMediaType) String() string {
	return [...]string{"image", "video"}[t-1]
}

type EntityNotFoundError struct {
	URL string
}

func NewEntityNotFoundError(url string) EntityNotFoundError {
	return EntityNotFoundError{URL: url}
}

func (e EntityNotFoundError) Error() string {
	return fmt.Sprintf("not found: %s", e.URL)
}

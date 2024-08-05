package writer

import (
	"fmt"
	"os"
	"path/filepath"

	"github.com/dragneelfps/site-scrapper/pkg/csv"
	"github.com/dragneelfps/site-scrapper/pkg/scrapper"
)

type CSVWriter struct{}

func (w CSVWriter) Write(data scrapper.EntityData, outputPath string) error {
	dirPath := filepath.Dir(outputPath)
	err := os.MkdirAll(dirPath, os.ModePerm)
	if err != nil {
		return fmt.Errorf("create dir %s: %w", dirPath, err)
	}
	return csv.WriteCSV(data.MediaList, outputPath)
}

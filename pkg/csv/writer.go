package csv

import (
	"encoding/csv"
	"errors"
	"fmt"
	"os"
	"reflect"
)

type Header []HeaderKey

type HeaderKey struct {
  Name string
  Key string
  Type reflect.Type
}

func GetHeader[T any](data T) (Header, error) {
  t := reflect.TypeOf(data)
  if t.Kind() == reflect.Slice || t.Kind() == reflect.Array {
    t = t.Elem()
  }
  if t.Kind() != reflect.Struct {
    return nil, errors.New("data should be of type Struct or []Struct")
  }
  header := make(Header, 0)
  for i := range t.NumField() {
    field := t.Field(i)
    name := field.Name
    tag := field.Tag
    nameOverride := tag.Get("csv")

    key := name
    if nameOverride != "" {
      key = nameOverride
    }

    header = append(header, HeaderKey{
      Name: name,
      Key: key,
      Type: field.Type,
    })
  }

  if len(header) == 0 {
    return nil, errors.New("header could not be extracted")
  }

  return header, nil
}

func (h Header) Keys() []string {
  if len(h) == 0 {
    return nil
  }
  keys := make([]string, 0, len(h))
  for _, header := range h {
    keys = append(keys, header.Key)
  }
  return keys
}

func GetAllValuesByHeader[T any](data []T, header Header) [][]string {
  if len(data) == 0 {
    return nil
  }
  rows := make([][]string, 0, len(data))
  for _, datum := range data {
    row := GetValuesByHeader(datum, header) 
    rows = append(rows, row)
  }
  return rows
}

func GetValuesByHeader[T any](datum T, header Header) []string {
  if len(header) == 0 {
    return nil
  }
  datumValue := reflect.ValueOf(datum)
  values := make([]string, 0, len(header))
  for _, key := range header {
    fieldValue := datumValue.FieldByName(key.Name)
    fieldValueStr := fieldValue.String()
    if fieldValueStringer, ok := fieldValue.Interface().(fmt.Stringer); ok {
      fieldValueStr = fieldValueStringer.String()
    }

    values = append(values, fieldValueStr)
  }

  return values
}


func WriteCSV[T any](data []T, outputPath string) (error) {
  header, err := GetHeader(data)
  if err != nil {
    return fmt.Errorf("get header: %w", err)
  }
  
  file, err := os.Create(outputPath)
  if err != nil {
    return fmt.Errorf("create file: %w", err)
  }


  csvFile := csv.NewWriter(file)
  
  csvFile.Write(header.Keys())
  csvFile.WriteAll(GetAllValuesByHeader(data, header))
  csvFile.Flush()
  if err := csvFile.Error(); err != nil {
    return fmt.Errorf("write csv file: %w", err)
  }
  
  return nil
}

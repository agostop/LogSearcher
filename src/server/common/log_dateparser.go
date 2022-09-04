package common

import (
    "time"

    "github.com/blevesearch/bleve/v2/analysis"
    "github.com/blevesearch/bleve/v2/registry"
)

const (
    Name = "logTimeParser"
    TimeFormat = "2006-01-02 15:04:05"
)

type MyTimeParser struct {
    layouts []string
}

func New(layouts []string) *MyTimeParser {
    return &MyTimeParser{
        layouts: layouts,
    }
}

func (p *MyTimeParser) ParseDateTime(input string) (time.Time, error) {
    for _, layout := range p.layouts {
        rv, err := time.ParseInLocation(layout, input, time.FixedZone("CST", 8*3600))
        if err == nil {
            return rv, nil
        }
    }
    return time.Time{}, analysis.ErrInvalidDateTime
}

func DateTimeParserConstructor(config map[string]interface{}, cache *registry.Cache) (analysis.DateTimeParser, error) {
    var layoutStrs []string
    layoutStrs = append(layoutStrs, TimeFormat)
    return New(layoutStrs), nil
}

func init() {
    registry.RegisterDateTimeParser(Name, DateTimeParserConstructor)
}
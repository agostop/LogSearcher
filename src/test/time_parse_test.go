package test

import (
	"fmt"
	"testing"
	"time"

	"github.com/blevesearch/bleve/v2"
)

func TestIndexTimeParse(t *testing.T) {
	idxMapping := bleve.NewIndexMapping()
	if err := idxMapping.AddCustomDateTimeParser("custom_datetime", map[string]interface{}{
		"type": "flexiblego",
		"layouts": []interface{}{
			"2006-Jan-02",
		},
	}); err != nil {
	      fmt.Println(err)
              return
	}

	timeStampMapping := bleve.NewDateTimeFieldMapping()
        timeStampMapping.DateFormat = "custom_datetime"

	docMapping := bleve.NewDocumentMapping()
	docMapping.AddFieldMappingsAt("timestamp", timeStampMapping)
	idxMapping.AddDocumentMapping("log", docMapping)

	idx, err := bleve.New("./bleve", idxMapping)
	if err != nil {
                fmt.Println(err)
                return
	}

	defer func() {
		err := idx.Close()
		if err != nil {
			fmt.Println(err)
		}
	}()

	if err = idx.Index("doc", map[string]string{
		"_type":     "log",
		"timestamp": "2022-Aug-03",
	}); err != nil {
		fmt.Println(err)
                 return  
	}

	start, _ := time.Parse(time.RFC3339, "2022-08-01T17:00:00Z")
	end, _ := time.Parse(time.RFC3339, "2022-08-04T18:00:00Z")
	dataRangeQuery := bleve.NewDateRangeQuery(start, end)
	dataRangeQuery.SetField("timestamp")

	sr := bleve.NewSearchRequest(dataRangeQuery)
	results, err := idx.Search(sr)
	if err != nil {
		fmt.Println(err)
                 return
	}

	fmt.Printf("result:%v",results.Hits)

	if results == nil || len(results.Hits) != 1 || results.Hits[0].ID != "doc" {
		fmt.Println("Expected the 1 hit")
	}
}

func TestTimestamp(t *testing.T) {
	t1 ,_ := time.Parse("2006-01-02 15:04:05", "2022-09-04 10:13:00")
	fmt.Println("time is: ", t1.Unix())
	fmt.Println("format: ", t1.String())
}
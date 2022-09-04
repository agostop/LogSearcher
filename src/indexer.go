package main

import (
	"encoding/json"
	"errors"
	"log"
	"os"
	"searchEngine/src/db"
	"searchEngine/src/server/common"
	"searchEngine/src/util"
	"strconv"
	"time"

	"github.com/blevesearch/bleve/v2"
	"github.com/blevesearch/bleve/v2/search"
)

type Indexer struct {
    indexer bleve.Index
    ldb     *db.LevelDbHandle
}

var (
    snowflakeGen, _  = util.NewNode(1)
    containerNameKey = []byte("container")
)

func NewIndexerInstance(basePath string) (*Indexer, error) {
    ldb := db.InitLevelDb(basePath + "/levelDb")

    Indexer := &Indexer{}
    Indexer.ldb = ldb
    Indexer.InitBleveIndexer(basePath + "/indexer")

    return Indexer, nil
}

func (s *Indexer) Accept(data []byte) error {
    msg := &common.LogMessage{}
    if err := json.Unmarshal(data, msg); err != nil {
        log.Printf("unmarshal receive data error: %s", err.Error())
        log.Printf("origin data: %s", string(data))
        return err
    }

    id := snowflakeGen.Generate()
    s.indexer.Index(id.String(), msg)

    s.saveToLevelDb(&id, msg)

    log.Printf("save success. msg: %s", id.String())
    return nil
}

func (s *Indexer) InitBleveIndexer(arg ...interface{}) error {

    var path = "/mnt/data/test.bleve"
    if len(arg) > 0 {
        path = arg[0].(string)
    }

    i, _ := openPathIfExists(path)
    if i != nil {
        s.indexer = i
        return nil
    }

    i, err := makeNewBleve(path)
    if err != nil {
        panic(err)
    }

    s.indexer = i

    return nil
}

func (s *Indexer) Search(containerName, searchText string, start, end int64) ([]string, error) {
    isAllEmpty := 0

    conjunctionQuery := bleve.NewConjunctionQuery()
    if len(searchText) > 0 {
        query := bleve.NewMatchQuery(searchText)
        query.SetField("message")
        conjunctionQuery.AddQuery(query)
    } else {
        isAllEmpty++
    }


    if len(containerName) > 0 {
        containerQuery := bleve.NewMatchQuery(containerName)
        containerQuery.SetField("container")
        conjunctionQuery.AddQuery(containerQuery)
    } else {
        isAllEmpty++
    }

    if start > 0 && end > 0 {
        startInc := true
        endInc := true
        dataRangeQuery := bleve.NewDateRangeInclusiveQuery(time.Unix(start, 0), time.Unix(end, 0), &startInc, &endInc)
        dataRangeQuery.SetField("timestamp")
        conjunctionQuery.AddQuery(dataRangeQuery)
    } else {
        isAllEmpty++
    }

    if isAllEmpty == 3 {
        log.Print("all parameter is emtpy, return")
        return nil, errors.New("parameter error")
    }

    searchReq := bleve.NewSearchRequest(conjunctionQuery)
    searchReq.Size = 1000
    searchReq.Sort = []search.SearchSort{ &search.SortDocID{ Desc: false } }
    searchResults, err := s.indexer.Search(searchReq)
    if err != nil {
        log.Print("search error.")
        return nil, err
    }
    log.Print("search result: ", searchResults)
    dmc := searchResults.Hits
    var resultContents []string
    for _, v := range dmc {
        id := v.ID
        i, err2 := strconv.ParseInt(id, 10, 64)
        if err2 != nil {
            log.Printf("got error. %s", err2.Error())
            return nil, err2
        }
        b, e := s.ldb.Get(util.ConvertIntToByte(uint64(i)))
        if e != nil {
            log.Printf("ldb get error. %v", e.Error())
        }

        resultContents = append(resultContents, string(b))
    }
    return resultContents, nil
}

func (s *Indexer) saveToLevelDb(id *util.ID, msg *common.LogMessage) error {
    if id == nil {
        return errors.New("must have ID parameter")
    }

    if msg == nil {
        return errors.New("must have msg parameter")
    }

    data, _ := json.Marshal(msg)
    s.ldb.Put(util.ConvertIntToByte(uint64((*id).Int64())), data)
    b, err := s.ldb.Get(containerNameKey)
    if err != nil {
        return err
    }

    containers := make(map[string]bool, 0)
    if b != nil {
        if err := json.Unmarshal(b, &containers); err != nil {
            return err
        }
    }
    containers[msg.Container] = true
    containersBytes, err := json.Marshal(containers)
    if err != nil {
        return err
    }
    s.ldb.Put(containerNameKey, containersBytes)
    return nil
}

func (s *Indexer) getAllContainerName() ([]string, error) {
    b, err := s.ldb.Get(containerNameKey)
    if err != nil {
        return nil, err
    }
    if len(b) == 0 {
        return []string{}, nil
    }
    nameMap := make(map[string]bool)
    err = json.Unmarshal(b, &nameMap)
    if err != nil {
        return nil, err
    }

    names := make([]string, 0, len(nameMap))
    for k := range nameMap {
        names = append(names, k)
    }

    return names, nil
}

func openPathIfExists(path string) (bleve.Index, error) {
    _, err := os.Stat(path)
    if os.IsNotExist(err) {
        log.Print("path not exists.")
        return nil, err
    }

    i, err := bleve.Open(path)
    if err != nil {
        log.Printf("bleve open exists dir failed. %s", err.Error())
        return nil, err
    }

    return i, nil
}

func makeNewBleve(path string) (bleve.Index, error) {
    container := bleve.NewTextFieldMapping()
    container.Store = false
    time := bleve.NewDateTimeFieldMapping()
    time.Store = false
    // time.DateFormat = "date_parser"
    message := bleve.NewTextFieldMapping()
    message.Store = false

    docMapping := bleve.NewDocumentMapping()
    docMapping.AddFieldMappingsAt("container", container)
    docMapping.AddFieldMappingsAt("timestamp", time)
    docMapping.AddFieldMappingsAt("message", message)

    indexMapping := bleve.NewIndexMapping()
    indexMapping.DefaultDateTimeParser = "logTimeParser"
    // indexMapping.AddCustomDateTimeParser("date_parser", map[string]interface{}{
    //     "type": "logTimeParser",
    //     "layouts": []interface{}{
    //         common.TimeFormat,
    //     },
    // })
    indexMapping.AddDocumentMapping("log", docMapping)
    index, err := bleve.New(path, indexMapping)
    if err != nil {
        log.Printf("bleve config error, %s", err.Error())
        return nil, err
    }

    return index, nil
}

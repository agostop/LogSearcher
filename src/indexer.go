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
	"github.com/blevesearch/bleve/v2/search/query"
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
        return err
    }

    s.indexer = i

    return nil
}

func (s *Indexer) Search(containerName, searchText *string, start, end *time.Time) ([]string, error) {

    if searchText == nil {
        return nil, errors.New("except searchText, but not got it")
    }

    var conjunctionQuery *query.ConjunctionQuery

    query := bleve.NewMatchQuery(*searchText)
    query.SetField("message")
    if containerName != nil {
        containerQuery := bleve.NewMatchQuery(*containerName)
        containerQuery.SetField("container")
        conjunctionQuery = bleve.NewConjunctionQuery(query, containerQuery)
    } else {
        conjunctionQuery = bleve.NewConjunctionQuery(query)
    }

    search := bleve.NewSearchRequest(conjunctionQuery)
    searchResults, err := s.indexer.Search(search)
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

    var containers []string
    if b != nil {
        if err := json.Unmarshal(b, &containers); err != nil {
            return err
        }
    }
    containers = append(containers, msg.Container)
    containersBytes, err := json.Marshal(containers)
    if err != nil {
        return err
    }
    s.ldb.Put(containerNameKey, containersBytes)
    return nil
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
    message := bleve.NewTextFieldMapping()
    message.Store = false

    docMapping := bleve.NewDocumentMapping()
    docMapping.AddFieldMappingsAt("container", container)
    docMapping.AddFieldMappingsAt("timestamp", time)
    docMapping.AddFieldMappingsAt("message", message)

    indexMapping := bleve.NewIndexMapping()
    indexMapping.AddDocumentMapping("log", docMapping)
    index, err := bleve.New(path, indexMapping)
    if err != nil {
        log.Printf("bleve config error, %s", err.Error())
        return nil, err
    }

    return index, nil
}



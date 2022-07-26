package main

import (
	"encoding/json"
	"log"
	"os"
	"searchEngine/src/db"
	"searchEngine/src/util"
	"time"

	"github.com/blevesearch/bleve/v2"
)

type logMessage struct {
    Id string `json:"id"`
    Container string `json:"container"`
    Timestamp msgDate `json:"timestamp"`
    Message string `json:"message"`
}

type msgDate time.Time

func (m *msgDate) MarshalJSON() ([]byte, error)  {
    return json.Marshal(time.Time(*m))
}

func (m *msgDate) UnmarshalJSON(b []byte) error {
    now, err := time.ParseInLocation(`"`+timeFormat+`"`, string(b), time.Local)
    if err != nil {
        return err
    }
    *m = msgDate(now)
    return nil
}

type Indexer struct {
    indexer bleve.Index
    ldb     *db.LevelDbHandle
}

var (
    timeFormat = "2006-01-02 15:04:05"
    snowflakeGen, _ = util.NewNode(1)
)

func NewIndexerInstance(basePath string) (*Indexer, error) {
    ldb := db.InitLevelDb(basePath + "/levelDb")

    Indexer := &Indexer{}
    Indexer.ldb = ldb
    Indexer.InitBleveIndexer(basePath + "/indexer")

    return Indexer, nil
}

func (s *Indexer) Accept(data []byte) error {
    msg := &logMessage{}
    
    if err := json.Unmarshal(data, msg); err != nil {
        log.Printf("unmarshal receive data error: %s", err.Error())
        log.Printf("origin data: %s", string(data))
        return err
    }

    msg.Id = snowflakeGen.Generate().String()
    s.Save(msg)

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

func (s *Indexer) Search(text string) ([]string, error) {
    query := bleve.NewMatchQuery(text)
    search := bleve.NewSearchRequest(query)
    searchResults, err := s.indexer.Search(search)
    if err != nil {
        log.Print("search error.")
        return nil, err
    }
    log.Print("search result: ", searchResults)
    return nil, nil
}

func (s *Indexer) Save(msg *logMessage) {
    s.indexer.Index(msg.Id, msg)
    log.Printf("save success. msg: %s", msg.Id)
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
    id := bleve.NewTextFieldMapping()
    id.Store = false
    container := bleve.NewTextFieldMapping()
    container.Store = false
    time := bleve.NewDateTimeFieldMapping()
    time.Store = false
    message := bleve.NewTextFieldMapping()
    message.Store = false

    docMapping := bleve.NewDocumentMapping()
    docMapping.AddFieldMappingsAt("id", id)
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

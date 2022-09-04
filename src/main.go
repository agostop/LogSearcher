package main

import (
	"encoding/json"
	"log"
	"searchEngine/src/server/common"
	"searchEngine/src/server/search"
	_ "searchEngine/src/server/search/http"
	"searchEngine/src/server/transport/protocol"
	_ "searchEngine/src/server/transport/tcp"
)

type searchReqBody struct {
    ContainerName string         `json:"containerName"`
    SearchText    string         `json:"searchText"`
    StartTime     int64 `json:"startTime"`
    EndTime       int64 `json:"endTime"`
}

func main() {
    s, err := NewIndexerInstance("./bin")
    if err != nil {
        panic(err)
    }

    go startSearchServer(s)
    go startTransportServer(s)

    select {}
}

func startSearchServer(s *Indexer) {
    http := common.SearchServerFactoryIns.GetServer("http")
    if http == nil {
        log.Fatal("search server not found")
        return
    }

    http.Config("0.0.0.0", 9008)
    http.AddHandle(searchMessage(s))
    http.AddHandle(getAllContainer(s))
    http.Start()

}


func startTransportServer(s *Indexer) {
    serv := common.TransServerFactoryIns.GetServer("tcp")
    if serv == nil {
        log.Print("server not found.")
        return
    }

    if err := serv.Config("0.0.0.0", 10010); err != nil {
        log.Printf("config error, %s", err.Error())
        return
    }

    var protoParser protocol.DelimiterProto = "\n"
    serv.SetProto(protoParser)
    serv.AddCallback(s)
    serv.Start()
}

func searchMessage(s *Indexer) search.UrlHandler {
    body := searchReqBody{}
    return search.UrlHandler{
        Url:     "/v1/search",
        Method:  "POST",
        ReqBody: &body,
        PostHandleFunc: func(bodyBytes []byte) ([]interface{}, error) {
            err := json.Unmarshal(bodyBytes, &body)
            if err != nil {
                return nil, err
            }
            var resultArr []interface{}
            result, err := s.Search(body.ContainerName, body.SearchText, body.StartTime, body.EndTime)
            if err != nil {
                return nil, err
            }
            msg := &common.LogMessage{}
            for _, v := range result {
                if err := json.Unmarshal([]byte(v), msg); err != nil {
                    return nil, err
                }
                resultArr = append(resultArr, *msg)
            }
            return resultArr, nil
        },
    }

}

func getAllContainer(s *Indexer) search.UrlHandler {
    return search.UrlHandler{
        Url:     "/v1/containers",
        Method:  "GET",
        ReqParam: nil,
        GetHandleFunc: func(param ...interface{}) ([]interface{}, error) {
            names, err := s.getAllContainerName()
            if err != nil {
                log.Printf("got container name failed: %v", err.Error())
                return nil, err
            }
            var resultArr []interface{} 
            for _, v := range names {
                resultArr = append(resultArr, v)
            }
            return resultArr, nil
        },
    }


}

// func waitSigs() {

//     sigs := make(chan os.Signal, 1)
//     done := make(chan bool, 1)

//     signal.Notify(sigs, syscall.SIGINT, syscall.SIGTERM)
//     go func()  {
//         sig := <-sigs
//         log.Println()
//         log.Println(sig)
//         done <- true
//     }()

//     log.Print("server up.")
//     <-done
//     log.Print("exiting")
// }

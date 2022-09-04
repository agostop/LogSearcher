package gin

import (
    "encoding/json"
    "errors"
    "fmt"
    "io/ioutil"
    "log"
    "net/http"

    "searchEngine/src/server/common"
    "searchEngine/src/server/search"
)

func init() {
    http := &httpServer{}
    common.SearchServerFactoryIns.Register("http", http)
}

type httpServer struct {
    host        string
    port        int
    handleMap map[string]search.UrlHandler
}

func (s *httpServer) Config(args ...interface{}) error {

    if len(args) > 1 {
        s.host = args[0].(string)
        s.port = args[1].(int)
    }

    return nil
}

func (s *httpServer) AddHandle(f search.UrlHandler) {
    if s.handleMap == nil {
        s.handleMap = make(map[string]search.UrlHandler)
    }
    s.handleMap[f.Url] = f
}

func (s *httpServer) Start(args ...interface{}) error {
    if len(s.host) == 0 {
        s.host = "localhost"
    }

    if s.port == 0 {
        s.port = 8899
    }

    if len(s.handleMap) == 0 {
        return errors.New("not have any handle")
    }

    http.Handle("/", s)

    log.Printf("search server listen on %v:%v", s.host, s.port)
    err := http.ListenAndServe(fmt.Sprintf("%v:%v", s.host, s.port), nil)
    if err != nil {
        log.Panicf("http server start got error: %v", err.Error())
        return err
    }

    return nil

}

func (s *httpServer) ServeHTTP(w http.ResponseWriter, req *http.Request) {
    path := req.URL.Path
    handle := s.handleMap[path]

    var searchResult []interface{}
    var err error

    if handle.Method == "GET" {
        var params []string
        for _, p := range handle.ReqParam {
            parameter := req.URL.Query().Get(p)
            params = append(params, parameter)
        }
        searchResult, err = handle.GetHandleFunc(params)
        if err != nil {
            responseError(w, err)
            return
        }
    } else if handle.Method == "POST" {
        b, err := ioutil.ReadAll(req.Body)
        if err != nil {
            log.Fatal(err)
            responseError(w, err)
            return
        }
        searchResult, err = handle.PostHandleFunc(b)
        if err != nil {
            responseError(w, err)
            return
        }
    } else {
        responseError(w, errors.New("not support method"))
        return
    }
    responseSucc(w, searchResult)
}


func responseSucc(resp http.ResponseWriter, searchResult []interface{}) {
    respData, _ := json.Marshal(&common.ResponseListData{
        Msg:  searchResult,
        Code: 0,
    })
    resp.Header().Add("content-type", "application/json")
    resp.Write(respData)
}

func responseError(resp http.ResponseWriter, err error) {
    b, _ := json.Marshal(&common.ResponseData{
        Msg:  err.Error(),
        Code: -1,
    })
    resp.Header().Add("content-type", "application/json")
    resp.Write(b)
}

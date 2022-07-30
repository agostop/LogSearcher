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

func init()  {
   http := &httpServer{}
   common.SearchServerFactoryIns.Register("http", http)
}

type httpServer struct {
    host string
    port int
    handleArray []search.UrlHandler
}

func (s *httpServer) Config(args ...interface{}) error {
    if len(args) < 1 {
        log.Panicf("param error, except 2 got %d.", len(args))
        return errors.New("parameter error.")
    }

    s.host = args[0].(string)
    s.port = args[1].(int)

    return nil
}

func (s *httpServer) AddHandle(f search.UrlHandler) {
    s.handleArray = append(s.handleArray, f)
}

func (s *httpServer) Start(args ...interface{}) error {
    if len(s.host) == 0 {
        s.host = "localhost"
    }

    if s.port == 0 {
        s.port = 8899;
    }

    if len(s.handleArray) == 0 {
        return errors.New("not have any handle")
    }

    for _, handleInst := range s.handleArray {
        http.HandleFunc(handleInst.Url, func(resp http.ResponseWriter, req *http.Request) {
            var searchResult []interface{}
            var err error
            if req.Method == "GET" {
                parameterValue := req.URL.Query().Get(handleInst.Parameter)
                searchResult, err = handleInst.SearchString(parameterValue)
                if err != nil {
                    responseError(resp, err)
                    return
                }
            } else if req.Method == "POST" {
                b, err := ioutil.ReadAll(req.Body)
                if err != nil {
                    log.Fatal(err)
                    responseError(resp, err)
                    return
                }
                m := make(map[string]string)
                err = json.Unmarshal(b, &m)
                if err != nil {
                    responseError(resp, err)
                    return
                }
                searchResult, err = handleInst.SearchString(m[handleInst.Parameter])
                if err != nil {
                    responseError(resp,err)
                    return
                }
            } else {
                responseError(resp, errors.New("not support method"))
                return
            }
            responseSucc(resp, searchResult)
        })
    }

    log.Printf("search server listen on %v:%v", s.host, s.port)
    err := http.ListenAndServe(fmt.Sprintf("%v:%v", s.host, s.port), nil)
    if err != nil {
        log.Panicf("http server start got error: %v", err.Error())
        return err
    }

    return nil

}

func responseSucc(resp http.ResponseWriter, searchResult []interface{}) {
    respData, _ := json.Marshal(&common.ResponseListData{
        Msg: searchResult,
        Code: 0,
    })
    resp.Header().Add("content-type", "application/json")
    resp.Write(respData)
}


func responseError(resp http.ResponseWriter, err error) {
    b, _ := json.Marshal(&common.ResponseData{
        Msg: err.Error(),
        Code: -1,
    })
    resp.Header().Add("content-type", "application/json")
    resp.Write(b)
}
package common

import (
	"searchEngine/src/server/search"
	"searchEngine/src/server/transport"
)

type LogMessage struct {
    Container string `json:"container"`
    Timestamp string `json:"timestamp"`
    Message string `json:"message"`
}

type ResponseListData struct {
    Msg  []interface{} 
    Code int
}

type ResponseData struct {
    Msg  interface{} 
    Code int
}

var TransServerFactoryIns = &transServerFactory{}
var SearchServerFactoryIns = &searchServerFactory{}

type transServerFactory struct {
    allServerMap map[string]transport.DataTransServer
}

func (sf *transServerFactory) Register(servType string, server transport.DataTransServer)  {
    if sf.allServerMap == nil {
        sf.allServerMap = make(map[string]transport.DataTransServer)
    }
    sf.allServerMap[servType] = server
}

func (sf *transServerFactory) GetServer(servType string) transport.DataTransServer {
    return sf.allServerMap[servType]
}

type searchServerFactory struct {
    allServerMap map[string]search.DataSearchServer
}

func (sf *searchServerFactory) Register(servType string, server search.DataSearchServer)  {
    if sf.allServerMap == nil {
        sf.allServerMap = make(map[string]search.DataSearchServer)
    }
    sf.allServerMap[servType] = server
}

func (sf *searchServerFactory) GetServer(servType string) search.DataSearchServer{
    return sf.allServerMap[servType]
}
package server

import "searchEngine/src/server/protocol"

type DataTransServer interface {
    Config(args ...interface{}) error
    Start(args ...interface{}) error
    AddCallback(cb ...CallbackInf)
    SetProto(proto protocol.ProtoInf)
    IsReady() bool 
}

type CallbackInf interface {
    Accept([] byte) error
}

var ServerFactoryIns = &serverFactory{}


type serverFactory struct {
    allServerMap map[string]DataTransServer
}

func (sf *serverFactory) Register(servType string, server DataTransServer)  {
    if sf.allServerMap == nil {
        sf.allServerMap = make(map[string]DataTransServer)
    }
    sf.allServerMap[servType] = server
}

func (sf *serverFactory) GetServer(servType string) DataTransServer {
    return sf.allServerMap[servType]
}
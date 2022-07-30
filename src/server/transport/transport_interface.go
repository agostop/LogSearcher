package transport

import "searchEngine/src/server/transport/protocol"

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

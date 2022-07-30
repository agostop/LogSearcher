package tcp

import (
	"encoding/json"
	"fmt"
	"log"
	"net"
    "searchEngine/src/server/common"
	"searchEngine/src/server/transport"
	"searchEngine/src/server/transport/protocol"
)

func init() {
    theServer := &tcpServer{}
    common.TransServerFactoryIns.Register("tcp", theServer)
}

var (
    succMsg, _ = json.Marshal(&common.ResponseData{
        Msg:  "ok",
        Code: 0,
    })
    maxPakcageLen = 10240
)

type tcpServer struct {
    address     string
    port        int
    callback    []transport.CallbackInf
    ready       bool
    protoParser protocol.ProtoInf
}

func (s *tcpServer) SetProto(proto protocol.ProtoInf) {
    s.protoParser = proto
}

func (s *tcpServer) Start(arg ...interface{}) error {
    tcpAddr, err := net.ResolveTCPAddr("tcp4", fmt.Sprintf("%s:%d", s.address, s.port))
    if err != nil {
        return err
    }

    t, err := net.ListenTCP("tcp4", tcpAddr)
    if err != nil {
        return err
    }

    s.ready = true

    log.Printf("transport server listen on %v:%v", s.address, s.port)
    for {
        conn, err := t.AcceptTCP()

        if err != nil {
            log.Printf("accept tcp connect failed. %s", err.Error())
            return err
        }

        go func() {

            var recvBuffer []byte
            var dataBuffer []byte
            defer conn.Close()
            for {
                // remoteAddr := conn.RemoteAddr()
                recvBuffer = make([]byte, 1024)
                n, err := conn.Read(recvBuffer)
                if err != nil {
                    log.Printf("error during read: %s", err)
                    break
                }
                markIndex := s.protoParser.Input(recvBuffer[:n])
                dataBuffer = append(dataBuffer, recvBuffer[:n]...)
                if markIndex <= 0 {
                    log.Print("data not complete, continue to read...")
                    if len(dataBuffer) > maxPakcageLen {
                        log.Printf("error data, length is too long. len: %d", markIndex)
                        break
                    }
                    continue
                } else if markIndex > 0 && markIndex < maxPakcageLen {
                    packedData, remainData, err := s.protoParser.Package(dataBuffer)
                    if err != nil {
                        log.Printf("package data error: %v", err.Error())
                        dataBuffer = nil
                        remainData = nil
                        continue
                    }

                    log.Printf("dataBuffer: %v", string(dataBuffer) )
                    log.Printf("remainData: %v", string(remainData) )
                    onMessage(s.callback, packedData)

                    if len(remainData) > 0 {
                        dataBuffer = remainData
                    } else {
                        dataBuffer = nil
                    }

                } else {
                    log.Printf("error data, length is too long. len: %d", markIndex)
                    break
                }

                // log.Printf("<%s> %s\n", remoteAddr, data[:n])
                conn.Write(succMsg)
            }
        }()
    }

}

func (s *tcpServer) Config(arg ...interface{}) error {
    if len(arg) > 0 {
        s.address = arg[0].(string)
    } else {
        s.address = "localhost"
    }

    if len(arg) > 1 {
        s.port = arg[1].(int)
    } else {
        s.port = 8010
    }

    return nil
}

func (s *tcpServer) AddCallback(cb ...transport.CallbackInf) {
    s.callback = cb
}

func (s *tcpServer) IsReady() bool {
    return s.ready
}

func onMessage(cbArr []transport.CallbackInf, packedDatas [][]byte) {
    log.Printf("onMessage receive: %v message", len(packedDatas))
    for _, packedData := range packedDatas {
        for _, cb := range cbArr {
            cb.Accept(packedData)
        }
    }
}

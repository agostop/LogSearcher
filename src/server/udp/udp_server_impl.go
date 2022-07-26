package udp

import (
	"encoding/json"
	"fmt"
	"log"
	"net"
	"searchEngine/src/server"
	"searchEngine/src/server/protocol"
)

type responseData struct {
    Msg  string
    Code int
}

func init() {
    theServer := &udpServer{}

    if err := theServer.Config("0.0.0.0", 10010); err != nil {
        log.Printf("config error, %s", err.Error())
        return
    }

    server.ServerFactoryIns.Register("udp", theServer)
    
}

var (
    succMsg, _ = json.Marshal(&responseData{
        Msg:  "ok",
        Code: 0,
    })
)

type udpServer struct {
    address  string
    port     int
    callback []server.CallbackInf
    conn     *net.UDPConn
    ready    bool
    stop bool
    protoParser protocol.ProtoInf
}

func (s *udpServer) Start(arg ...interface{}) error {
    udpAddr, err := net.ResolveUDPAddr("udp", fmt.Sprintf("%s:%d",s.address,s.port));
    if err != nil {
        return err
    }

    u, err := net.ListenUDP("udp", udpAddr)
    if err != nil {
        return err
    }
    s.conn = u
    s.ready = true
    s.stop = false
    go func() {

        data := make([]byte, 1024)
        for ;!s.stop; {
            n, remoteAddr, err := s.conn.ReadFromUDP(data)
            if err != nil {
                log.Printf("error during read: %s", err)
                continue
            }

            if s.callback == nil || len(s.callback) == 0 {
                log.Printf("callback is null, data not process. %s", data)
                continue
            }
            
            for _, cb := range s.callback {
                cb.Accept(data)
            }

            log.Printf("<%s> %s\n", remoteAddr, data[:n])
            s.conn.WriteToUDP(succMsg, remoteAddr)
        }
    }()

    return nil
}

func (s *udpServer) Config(arg ...interface{}) error {
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

func (s *udpServer) AddCallback(cb ...server.CallbackInf) {
    s.callback = cb
}

func (s *udpServer) IsReady() bool {
    return s.ready
}
func (s *udpServer) SetProto(proto protocol.ProtoInf) {
    s.protoParser = proto
}


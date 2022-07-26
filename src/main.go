package main

import (
	"log"
	"searchEngine/src/server"
	"searchEngine/src/server/protocol"
	_ "searchEngine/src/server/tcp"
)

func main() {
    s, err := NewIndexerInstance("./bin")
    if err != nil {
        panic(err)
    }

    serv := server.ServerFactoryIns.GetServer("tcp")
    if serv == nil {
        log.Print("server not found.")
        return
    }
    var protoParser protocol.DelimiterProto = "\n"
    serv.SetProto(protoParser)
    serv.AddCallback(s)
    serv.Start()

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
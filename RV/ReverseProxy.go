package main

import (
        "io"
        "log"
        "net"
        "fmt"
)

func main() {
        front := ":19530"
        back := "com1.mcmullin.org:19530"
        In, err := net.Listen("tcp4", front)
        if err != nil {
                log.Fatalf("failed to setup listener: %v", err)
        }

        log.Printf("listening on %s, balancing %s", front, back)
        log.Println("//////////////////LISTENING//////////////////")
        i := 1
        for {
                Inconn, err := In.Accept()
                fmt.Println(i)
                go ReProx(Inconn, back)
                if err != nil {
                        log.Printf("failed to accept: %s", err)
                        continue
                }
                i++

        }
}
func copy(src net.Conn, dst net.Conn, stop chan bool) {
        io.Copy(dst, src)
        dst.Close()
        src.Close()
        stop <- true
        return
}


func ReProx(src net.Conn, server string) {

        dst, err := net.Dial("tcp4", server)
        if err != nil {
                src.Close()
                log.Printf("failed to dial %s: %s", server, err)
                return
        }

        stop := make(chan bool)

        go copy(dst, src, stop)
        go copy(src, dst, stop)

        select {
        case <-stop:
                return
        }

}


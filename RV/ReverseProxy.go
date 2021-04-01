package load_balance

import (
	"io"
	"log"
	"net"
)

func copy(src net.Conn, dst net.Conn) {
	io.Copy(dst, src)
	dst.Close()
	src.Close()
	//stop <- true
	return
}

func ReProx(src net.Conn, server string) {

	dst, err := net.Dial("tcp4", server)
	if err != nil {
		src.Close()
		log.Printf("failed to dial %s: %s", server, err)
		return
	}

	//stop := make(chan bool)

	go copy(dst, src)
	go copy(src, dst)

	// select {
	// case <-stop:
	//         return
	// }

}

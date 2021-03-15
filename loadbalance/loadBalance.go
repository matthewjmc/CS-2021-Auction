package main

import (
	"flag"
)

// Data structure storing backend info
type ServerList struct {
	Servers []string
	Latest  int
	Alive   bool
}

// TODO - Manager (all frontend responsibility) TCP Server listening on a port

// TODO - Least Connection routes request to backend that have least number of connection

// TODO - RoundRobin request to backend with round robin fashion

func (server *ServerList) route() string {
	i := server.Latest % len(server.Servers)
	server.Latest++
	//stri := strings.Join(server.Servers,",")
	return server.Servers[i]
}

func init() {
	flag.Parse()

}

func main() {

}

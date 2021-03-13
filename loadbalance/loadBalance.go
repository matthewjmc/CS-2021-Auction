package main

func main() {

}

// TODO - Data structure storing backend info
type Server struct {
	Route string
	Alive bool
}

type ServerList struct {
	Servers []Server
	Latest  int
}

// TODO - Manager (all frontend responsibility) TCP Server listening on a port

// TODO - Least Connection routes request to backend that have least number of connection

// TODO - RoundRobin request to backend with round robin fashion

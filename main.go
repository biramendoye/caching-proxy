package main

import (
	"flag"

	"github.com/biramendoye/caching-proxy/server"
)

func main() {
	port := flag.String("port", "3000", "the port on which the caching proxy server will run")
	origin := flag.String("origin", "", "the URL of the server to which the requests will be forwarded")
	// clearCache := flag.Bool("clear-cache", false, "clear the cache")
	flag.Parse()

	srv := server.NewServer("localhost", *port, *origin)

	if err := srv.Run(); err != nil {
		panic(err)
	}
}

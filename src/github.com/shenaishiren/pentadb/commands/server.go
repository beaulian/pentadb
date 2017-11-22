// Contains the implementation of server-command of levelDB
package main

import (
	"log"
	"net"
	"flag"
	"net/http"
	"net/rpc"
	"github.com/syndtr/goleveldb/leveldb"
	"github.com/shenaishiren/pentadb/opt"
	"github.com/shenaishiren/pentadb/server"
	"fmt"
)

var helpPrompt = `Usage: pentadb [--port <port>] [--path <path>] [options]

A PentaDB rpc server, backed by LevelDB

Options:
	--help           		Display this help message and exit
	--port <port>    		The port to listen on (default: 4567)
	--path <path>    		The path to use for the LevelDB store
`

type Server struct {
	Node *server.Node
}

func (s *Server) listen(port string, path string) error {
	s.Node = server.NewNode("127.0.0.1:" + port)
	db, err := leveldb.OpenFile(path, nil)
	defer db.Close()

	if err != nil {
		return err
	}
	s.Node.DB = db
	rpc.Register(s.Node)
	rpc.HandleHTTP()      // bind prc to http service

	l, e := net.Listen("tcp", ":" + port)
	if e != nil {
		log.Fatal("listen error:", e)
	}
	log.Printf("%c[1;40;32m%s%c[0m",
		0x1B, "listening at http://0.0.0.0:" + port,
		0x1B)
	http.Serve(l, nil)
	return nil
}


func main() {
	var (
		help bool
		port string
		path string
	)
	flag.BoolVar(&help, "h", false, "Display this help message and exit")
	flag.StringVar(&port, "p", "4567", "The port to listen on (default: 4567)")
	flag.StringVar(&path, "a", opt.DeafultPath, "The path to use for the LevelDB store")

	// run
	flag.Parse()

	// help command
	if help {
		fmt.Print(helpPrompt)
	} else {
		server := new(Server)
		server.listen(port, path)
	}
}
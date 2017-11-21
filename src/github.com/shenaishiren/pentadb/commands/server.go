// Contains the implementation of server-command of levelDB
package commands

import (
	"os"
	"log"
	"net"
	"net/http"
	"net/rpc"
	"github.com/urfave/cli"
	"github.com/syndtr/goleveldb/leveldb"
	"github.com/shenaishiren/pentadb/opt"
	"github.com/shenaishiren/pentadb/server"
)

var helpPrompt = `Usage: pentadb [--port <port>] [--path <path>] [options]

A PentaDB rpc server, backed by LevelDB

Options:
	--help           		Display this help message and exit
	--port <port>    		The port to listen on (default: 4567)
	--path <path>    		The path to use for the LevelDB store
`

type Server struct {
	node *server.Node
}

func (s *Server) initRPC(path string) error {
	s.node = server.NewNode("")
	db, err := leveldb.OpenFile(path, nil)
	defer db.Close()

	if err != nil {
		return err
	}
	s.node.DB = db
	rpc.Register(s.node)
	rpc.HandleHTTP()      // bind prc to http service

	return nil
}

func (s *Server) listen(port string, path string) {
	s.initRPC(path)

	l, e := net.Listen("tcp", ":" + port)
	if e != nil {
		log.Fatal("listen error:", e)
	}
	log.Printf("%c[1;40;32m%s%c[0m",
		0x1B, "listening at http://0.0.0.0:" + port,
		0x1B)
	http.Serve(l, nil)
}


func main() {
	app := cli.NewApp()

	app.Flags = []cli.Flag {
		cli.BoolFlag{
			Name: "help, h",
			Hidden: true,
			Usage: "Help prompt",
		},
		cli.StringFlag{
			Name: "port, p",
			Value: "4567",
			Usage: "Port for listening",
		},
		cli.StringFlag{
			Name: "path, a",
			Value: opt.DeafultPath,
			Usage: "Path for levelDB",
		},
	}

	app.Action = func(c *cli.Context) error {
		if c.Bool("help") {
			log.Print(helpPrompt)
		} else {
			server := new(Server)
			server.listen(c.String("port"), c.String("path"))
		}

		return nil
	}
	app.Run(os.Args)
}
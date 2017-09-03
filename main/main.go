package main


import (
	"os"
	log "github.com/sirupsen/logrus"
	cli "github.com/urfave/cli"
	"../server/restful"
	"../server/kcp_server"
	"../world"
	"net/http"
	"time"
)

func main() {
	app:=&cli.App{
		Name:"battle server",
		Usage:"frame sync server for battle",
		Version:"0.0.1",
		Flags:[]cli.Flag{
			&cli.StringFlag{
				Name:"kcp",
				Value:":9090",
				Usage:"listen kcp",
			},
			&cli.StringFlag{
				Name:"rpc",
				Value:":9092",
				Usage:"listen tcp",
			},
		},
	};

	app.Action=func(c *cli.Context) error{
		w:=world.NewWorld();
		kcp_server.NewKcpServer(c.String("tcp")).StartAt(w);
		restful.NewRoomWS(w);
		http.ListenAndServe(c.String("rpc"),nil);
		return nil;
	}
	log.Info("server started at ",time.Now());
	app.Run(os.Args);
	log.Info("server stoped at ",time.Now());
}

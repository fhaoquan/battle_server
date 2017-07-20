package main


import (
	"os"
	log "github.com/sirupsen/logrus"
	cli "github.com/urfave/cli"
	"../udp"
	"../restful_service/room_service"
	"../sessions"
	"../world"
	"../builder"
	"../room"
	"net"
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
				Name:"tcp",
				Value:":9090",
				Usage:"listen tcp",
			},
			&cli.StringFlag{
				Name:"udp",
				Value:":9091",
				Usage:"listen tcp",
			},
			&cli.StringFlag{
				Name:"rpc",
				Value:":9092",
				Usage:"listen tcp",
			},
		},
	};
	//new_world();
	app.Action=func(c *cli.Context) error{
		log.Info(c.String("tcp"));
		w:=world.NewWorld();
		room_service.NewRoomWS(w);
		udp.StartNewKcpServer(c.String("tcp"),func(conn net.Conn){
			builder.BuildKcpSession(conn,w);
		});
		log.Fatal(http.ListenAndServe(":9091",nil));
		return nil;
	}
	log.Info("server started at ",time.Now());
	app.Run(os.Args);
	log.Info("server stoped at ",time.Now());
}

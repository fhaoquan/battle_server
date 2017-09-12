package main


import (
	"os"
	log "github.com/sirupsen/logrus"
	cli "github.com/urfave/cli"
	"../server/restful"
	"../server/kcp_server"
	"../server/udp_server"
	"../world"
	"net/http"
	"time"
	//"../test"
	"runtime"
)

func main() {
	runtime.GOMAXPROCS(runtime.NumCPU());
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
				Name:"udp",
				Value:":9091",
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
		kcp_server.StartGateway(c.String("kcp"),func(uid,rid uint32,session *kcp_server.KcpSession){
			defer func(){
				recover();
			}()
			if r:=w.FindRoom(rid);r!=nil{
				r.OnKcpSession(uid,session);
			}else{
				log.Error("can not find room ",rid," at session",session.RemoteAddr);
				session.Close(false);
			}
		});
		udp_server.StartGateway(c.String("udp"));
		restful.NewRoomWS(w);
		http.ListenAndServe(c.String("rpc"),nil);
		return nil;
	}
	log.Info("server started at ",time.Now());
	app.Run(os.Args);
	log.Info("server stoped at ",time.Now());
}

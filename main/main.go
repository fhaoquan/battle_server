package main

import (
	"os"
	"github.com/sirupsen/logrus"
	cli "github.com/urfave/cli"
	"../rest_api"
	"net/http"
	"time"
	"../gateway"
	//"../service_discovery"
	"runtime"
)

func main() {
	runtime.GOMAXPROCS(runtime.NumCPU());
	logrus.StandardLogger().Formatter.(*logrus.TextFormatter).TimestampFormat=time.StampMilli;
	app:=&cli.App{
		Name:"battle server",
		Usage:"frame sync server for battle",
		Version:"0.0.1",
		Flags:[]cli.Flag{
			&cli.StringFlag{
				Name:"bind",
				Value:"",
				Usage:"service bind address",
			},
			&cli.StringFlag{
				Name:"kcp",
				Value:"9090",
				Usage:"listen kcp",
			},
			&cli.StringFlag{
				Name:"rpc",
				Value:"9092",
				Usage:"listen tcp",
			},
			&cli.StringFlag{
				Name:"consul",
				Value:"10.0.0.101:8500",
				Usage:"consul address",
			},
		},
	};

	app.Action=func(c *cli.Context) error{
		gateway.Start(":"+c.String("kcp"))
		rest_api.NewRoomWS();
		rest_api.NewConsulCheck();
		/*
		service_discovery.RegisteServiceToConsul(
			c.String("consul"),
			c.String("bind")+":"+c.String("kcp")+","+c.String("udp")+","+c.String("rpc"),
			//c.String("bind")+":"+c.String("rpc")+"/consul/check",
			"10.0.0.6"+":"+c.String("rpc")+"/consul/check",
			);
		*/
		http.ListenAndServe(":"+c.String("rpc"),nil);
		return nil;
	}
	logrus.Info("server started at ",time.Now());
	app.Run(os.Args);
	logrus.Info("server stoped at ",time.Now());
}

package restful

import (
	"github.com/emicklei/go-restful"
	"github.com/sirupsen/logrus"
)

func NewConsulCheck(){
	ws:=new(restful.WebService);
	ws.
	Path("/consul").
		Consumes(restful.MIME_JSON).
		Produces(restful.MIME_JSON);
	ws.Route(
		ws.	GET("/check").
			To(
			func(req *restful.Request, res *restful.Response){
				logrus.Error("consul check")
			}));

	restful.Add(ws);
}

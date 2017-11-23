package rest_api

import (
	"../room_1v1"
	"github.com/emicklei/go-restful"
)

func count_room(req *restful.Request, res *restful.Response){
	res.WriteEntity(&struct{
		Count int;
	}{room_1v1.CountRoom()});
}

package room_service

import (
	"../../world"
	"github.com/emicklei/go-restful"
)

func count_room(req *restful.Request, res *restful.Response,wld *world.World){
	res.WriteEntity(&struct{
		Count int;
	}{wld.CountRoom()});
}

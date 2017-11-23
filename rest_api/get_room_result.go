package rest_api

import (
	"../result_cache"
	"github.com/emicklei/go-restful"
	)

func get_room_result(req *restful.Request, res *restful.Response){
	gid:=req.PathParameter("room_guid");
	dt,ok:=result_cache.GetResult(gid);
	if ok {
		res.WriteEntity(dt);
	}else {
		res.WriteEntity(nil);
	}
}

package room_service

import (
	"github.com/emicklei/go-restful"
	"../../world"
)

func get_all_rooms(req *restful.Request, resp *restful.Response){
}
func get_room(req *restful.Request, resp *restful.Response){
}

func del_room(req *restful.Request, resp *restful.Response){
}


func NewRoomWS(w *world.World){
	ws:=new(restful.WebService);
	ws.
	Path("/room").
		Consumes(restful.MIME_JSON).
		Produces(restful.MIME_JSON);
	ws.Route(
		ws.	GET("/").
			To(get_all_rooms));

	ws.Route(
		ws.	GET("/{room_id}").
			To(get_room));

	ws.Route(
		ws.	PUT("/").
			To(
			func(req *restful.Request, res *restful.Response){
				s:=&s_new_room_request_json{};
				if err:=req.ReadEntity(s);err!=nil{
					res.WriteEntity(&struct {
						RoomID int;
						Err error;
					}{-1,err})
				}else {
					id,err:=new_room(w,s)
					res.WriteEntity(&struct {
						RoomID int;
						Err error;
					}{id,err})
				}
			}));

	ws.Route(
		ws.	DELETE("/{room_id}").
			To(del_room));

	restful.Add(ws);
}
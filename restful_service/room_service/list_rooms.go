package room_service

import (
	"../../world"
	"../../room"
	"github.com/emicklei/go-restful"
)

func list_rooms(req *restful.Request, res *restful.Response,wld *world.World){
	s:=&struct {
		ids []int;
	}{}
	wld.ForEachRoom(func(room *room.Room) bool {
		s.ids=append(s.ids, (int)(room.GetID()));
		return true;
	})
	res.WriteEntity(s);
}

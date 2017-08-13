package room_service

import (
	"../../builder"
	"../../world"
	"github.com/pkg/errors"
	"fmt"
	"github.com/emicklei/go-restful"
)

type s_player_info struct {
	id int;
	name string;
}
func (me *s_player_info)GetPlayerID()uint32{
	return uint32(me.id);
}
func (me *s_player_info)GetPlayerName()string{
	return me.name;
}
type new_room_request_json struct{
	room_type int;
	room_players []s_player_info;
}
type new_room_info struct{
	id int;
	room_type int;
}
func build_room(w *world.World,param *new_room_request_json)(int,error){
	switch(param.room_type){
	case 1:
		if(len(param.room_players)!=2){
			return -1,errors.New("player count must 2");
		}
		r,e:=builder.BuildRoom1V1().
			AtWorld(w).
			WaitPlayers(&param.room_players[0],&param.room_players[1]).
			DoBuild();
		if e!=nil{
			return int(r.GetID()),nil;
		}else{
			return -1,e;
		}
	}
	return -1,errors.New(fmt.Sprint("unknown room type=",param.room_type));
}
func new_room(req *restful.Request, res *restful.Response,wld *world.World){
	s:=&new_room_request_json{};
	if err:=req.ReadEntity(s);err!=nil{
		res.WriteEntity(&struct {
			RoomID int;
			Err error;
		}{-1,err})
	}else {
		id,err:=build_room(wld,s)
		res.WriteEntity(&struct {
			RoomID int;
			Err error;
		}{id,err})
	}
}

package restful

import (
	"../../room"
	"../../world"
	"../../battle"
	"github.com/pkg/errors"
	"fmt"
	"github.com/emicklei/go-restful"
	"github.com/sirupsen/logrus"
)
type s_player_info struct {
	Id int;
	Name string;
	Units []battle.Unit;
}
func (me *s_player_info)GetPlayerID()uint32{
	return uint32(me.Id);
}
func (me *s_player_info)GetPlayerName()string{
	return me.Name;
}
func (me *s_player_info)GetUnits()[]battle.Unit{
	return me.Units;
}
type new_room_request_json struct{
	Room_type int;
	Room_players []s_player_info;
}
type new_room_info struct{
	id int;
	room_type int;
}
func build_room(w *world.World,param *new_room_request_json)(int,error){
	switch(param.Room_type){
	case 1:
		if(len(param.Room_players)!=2){
			return -1,errors.New("player count must 2");
		}
		r,e:=room.BuildRoom1v1(&param.Room_players[0],&param.Room_players[1]);
		if e==nil{
			w.AddNewRoom(r);
			return int(r.GetID()),nil;
		}else{
			return -1,e;
		}
	}
	return -1,errors.New(fmt.Sprint("unknown room type=",param.Room_type));
}
func new_room(req *restful.Request, res *restful.Response,wld *world.World){
	s:=&new_room_request_json{};
	if err:=req.ReadEntity(s);err!=nil{
		logrus.Error(err);
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

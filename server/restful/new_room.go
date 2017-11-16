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
	Room_type		int;
	Lifecycle		int;
	SuddenDeath		int;
	WinScore		int;
	Room_players 	[]s_player_info;
}
type new_room_info struct{
	id int;
	room_type int;
}
func build_room(w *world.World,param *new_room_request_json)(*room.Room1v1,error){
	switch(param.Room_type){
	case 1:
		if(len(param.Room_players)!=2){
			return nil,errors.New("player count must 2");
		}
		r,e:=room.BuildRoom1v1(
			&room.RoomBuildContext{param.Lifecycle,param.SuddenDeath,param.WinScore},
			&param.Room_players[0],
			&param.Room_players[1]);
		logrus.Error("build new room id= ",r.GetID()," guid=",r.GetGuid());
		if e!=nil{
			return nil,e;
		}
		w.AddNewRoom(r);
		return r,nil;
	}
	return nil,errors.New(fmt.Sprint("unknown room type=",param.Room_type));
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
		r,e:=build_room(wld,s)
		if e!=nil{
			res.WriteEntity(&struct {
				RoomID int;
				Err error;
			}{-1,err})
		}else{
			res.WriteEntity(&struct {
				RoomID int;
				Guid	string;
				Err error;
			}{int(r.GetID()),r.GetGuid(),err})
		}
	}
}

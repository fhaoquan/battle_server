package room_service

import (
	"../../builder"
	"../../world"
	"github.com/pkg/errors"
	"fmt"
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
type s_new_room_request_json struct{
	room_type int;
	room_players []s_player_info;
}
type s_new_room_info struct{
	id int;
	room_type int;
}
func new_room(w *world.World,param *s_new_room_request_json)(int,error){
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

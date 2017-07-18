package room_service

import (
	"../../room"
	"github.com/pkg/errors"
	"fmt"
)

type s_player_info struct {
	V_id int;
	V_name string;
}
func (me *s_player_info)GetPlayerID()uint32{
	return uint32(me.V_id);
}
func (me *s_player_info)GetPlayerName()string{
	return me.V_name;
}
type s_new_room_request_json struct{
	V_type int;
	V_players []s_player_info;
}
type s_new_room_info struct{
	V_id int;
	V_type int;
}
func new_room(context I_room_service_context,param *s_new_room_request_json)(int,error){
	switch(param.V_type){
	case 1:
		if(len(param.V_players)!=2){
			return -1,errors.New("player count must 2");
		}else{
			return (int)(
				context.AddNewRoom(room.NewRoom1v1().
					Player(&param.V_players[0],&param.V_players[1]).
					Build())),
				nil;
		}
	}
	return -1,errors.New(fmt.Sprint("unknown room type=",param.V_type));
}

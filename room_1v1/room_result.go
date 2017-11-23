package room_1v1
type PlayerResult struct{
	Pid		int;
	Score	int;
}
type RoomResult struct{
	Guid	string;
	Result	[]*PlayerResult;
}
func new_player_result(pid int,Score int)(*PlayerResult){
	return &PlayerResult{
		pid,Score,
	}
}

func new_room_result(guid string)(*RoomResult){
	return &RoomResult{
		guid,
		make([]*PlayerResult,0),
	}
}

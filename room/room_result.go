package room
type PlayerResult struct{
	Pid		int;
	Score	int;
}
type RoomResult struct{
	Guid	string;
	Result	[]*PlayerResult;
}
func newPlayerResult(pid int,Score int)(*PlayerResult){
	return &PlayerResult{
		pid,Score,
	}
}

func newRoomResult(guid string)(*RoomResult){
	return &RoomResult{
		guid,
		make([]*PlayerResult,0),
	}
}

package command

func UnitAttackStart(
	data []byte,
	broadcast func([]byte,int)){
	broadcast(data,len(data));
}

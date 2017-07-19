package command

type IUnitMovementDataGetter interface {
	GetID()uint16;
	GetX()uint16;
	GetY()uint16;
	GetSpeed()uint16;
	GetFace()uint16;
	GetAimingFace()uint16;
}

func BroadcastBattleData(
	s_buf []byte,
	foreach_unit_do func(func(IUnitMovementDataGetter)),
	broadcast func([]byte,int))(err error){
	w:=&packet_encoder{
		s_buf,
		0,
	}
	defer func(){
		if(err!=nil){
			broadcast(s_buf,w.pos+1);
		}
	}();
	p1:=w.get_uint8_placeholder();
	count:=0;
	foreach_unit_do(func(u IUnitMovementDataGetter){
		w.	write_unit_id(u.GetID()).
			write_unit_x(u.GetX()).
			write_unit_y(u.GetY()).
			write_unit_speed(u.GetSpeed()).
			write_unit_face(u.GetFace()).
			write_unit_aiming_face(u.GetAimingFace());
		count++;
	});
	p1(uint8(count));
	return nil;
}

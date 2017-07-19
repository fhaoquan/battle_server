package command

import "github.com/pkg/errors"

type IUnitAttackHelper interface {
	GetAttackPower()uint16;
	GetHP()uint16;
	SetHP(uint16);
	GetID()uint16;
}

func UnitAttackDone(
	data []byte,
	find_unit func(uint16)IUnitAttackHelper,
	s_buf []byte,
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
	r:=&packet_decoder{
		data:data,
		pos:0,
	}
	u1:=find_unit(r.read_unit_id());
	if(u1==nil){
		return errors.New("cant find src unit id ");
	}

	w.write_unit_id(u1.GetID());
	count:=(int)(r.read_unit_count());
	w.write_unit_count((uint8)(count));
	for i:=0;i<count;i++{
		id:=r.read_unit_id();
		w.write_uint16(id);
		u2:=find_unit(id);
		if(u2!=nil){
			if(u2.GetHP()>u1.GetAttackPower()){
				u2.SetHP(u2.GetHP()-u1.GetAttackPower());
			}else{
				u2.SetHP(0);
			}
			w.write_uint16(u2.GetHP());
		}else{
			w.write_uint16(0);
		}
	}
	return nil;
}

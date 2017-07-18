package command

import (
	"encoding/binary"
)

type packet_reader struct{
	data []byte;
	pos int;
}
func (r *packet_reader)read_uint32()(v uint32){
	v=binary.BigEndian.Uint32(r.data[r.pos:r.pos+2])
	r.pos+=4;
	return;
}
func (r *packet_reader)read_uint16()(v uint16){
	v=binary.BigEndian.Uint16(r.data[r.pos:r.pos+2])
	r.pos+=2;
	return;
}
func (r *packet_reader)read_uint8()(v uint8){
	v=r.data[r.pos];
	r.pos+=1;
	return;
}

func (r *packet_reader)read_unit_id()(uint16){
	return r.read_uint16();
}

func (r *packet_reader)read_unit_location_x()(uint16){
	return r.read_uint16();
}

func (r *packet_reader)read_unit_location_y()(uint16){
	return r.read_uint16();
}

func (r *packet_reader)read_unit_speed()(uint16){
	return r.read_uint16();
}

func (r *packet_reader)read_unit_face()(uint16){
	return r.read_uint16();
}

func (r *packet_reader)read_unit_aiming_face()(uint16){
	return r.read_uint16();
}

func (r *packet_reader)read_unit_camps()(uint8){
	return r.read_uint8();
}

func (r *packet_reader)read_unit_attk_power()(uint16){
	return r.read_uint16();
}

func (r *packet_reader)read_unit_hp()(uint16){
	return r.read_uint16();
}

func (r *packet_reader)read_unit_stat()(uint16){
	return r.read_uint16();
}

func (r *packet_reader)read_unit_count()(uint8){
	return r.read_uint8();
}
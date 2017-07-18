package command

import "encoding/binary"

type packet_writer struct{
	data []byte;
	pos int;
}
func (w *packet_writer)write_uint32(v uint32)(*packet_writer){
	binary.BigEndian.PutUint32(w.data[w.pos:],v);
	w.pos+=4;
	return w;
}
func (w *packet_writer)write_uint16(v uint16)(*packet_writer){
	binary.BigEndian.PutUint16(w.data[w.pos:],v);
	w.pos+=2;
	return w;
}
func (w *packet_writer)write_uint8(v uint8)(*packet_writer){
	w.data[w.pos]=v;
	w.pos+=1;
	return w;
}
func (w *packet_writer)write_unit_x(v uint16)(*packet_writer){
	return w.write_uint16(v);
}
func (w *packet_writer)write_unit_y(v uint16)(*packet_writer){
	return w.write_uint16(v);
}
func (w *packet_writer)write_unit_speed(v uint16)(*packet_writer){
	return w.write_uint16(v);
}
func (w *packet_writer)write_unit_face(v uint16)(*packet_writer){
	return w.write_uint16(v);
}
func (w *packet_writer)write_unit_aiming_face(v uint16)(*packet_writer){
	return w.write_uint16(v);
}
func (w *packet_writer)write_unit_id(v uint16)(*packet_writer){
	return w.write_uint16(v);
}
func (w *packet_writer)write_unit_count(v uint8)(*packet_writer){
	return w.write_uint8(v);
}
func (w *packet_writer)get_uint8_placeholder()(writer func(v uint8)){
	p:=w.pos;
	writer=func(v uint8){
		w.data[p]=v;
	}
	w.pos+=1;
	return ;
}
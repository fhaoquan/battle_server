package battle

import "encoding/binary"

type packet_encoder struct{
	data []byte;
	pos int;
}
func (w *packet_encoder)write_uint64(v uint64)(*packet_encoder){
	binary.LittleEndian.PutUint64(w.data[w.pos:],v);
	w.pos+=8;
	return w;
}
func (w *packet_encoder)write_uint32(v uint32)(*packet_encoder){
	binary.LittleEndian.PutUint32(w.data[w.pos:],v);
	w.pos+=4;
	return w;
}
func (w *packet_encoder)write_uint16(v uint16)(*packet_encoder){
	binary.LittleEndian.PutUint16(w.data[w.pos:],v);
	w.pos+=2;
	return w;
}
func (w *packet_encoder)write_uint8(v uint8)(*packet_encoder){
	w.data[w.pos]=v;
	w.pos+=1;
	return w;
}
func (w *packet_encoder)write_unit_x(v uint16)(*packet_encoder){
	return w.write_uint16(v);
}
func (w *packet_encoder)write_unit_y(v uint16)(*packet_encoder){
	return w.write_uint16(v);
}
func (w *packet_encoder)write_unit_speed(v uint16)(*packet_encoder){
	return w.write_uint16(v);
}
func (w *packet_encoder)write_unit_face(v uint16)(*packet_encoder){
	return w.write_uint16(v);
}
func (w *packet_encoder)write_unit_aiming_face(v uint16)(*packet_encoder){
	return w.write_uint16(v);
}
func (w *packet_encoder)write_unit_id(v uint16)(*packet_encoder){
	return w.write_uint16(v);
}
func (w *packet_encoder)write_unit_count(v uint8)(*packet_encoder){
	return w.write_uint8(v);
}
func (w *packet_encoder)write_bytes(d []byte)(*packet_encoder){
	w.pos+=copy(w.data[w.pos:],d);
	return w;
}
func (w *packet_encoder)get_uint08_placeholder()(writer func(v uint8)){
	p:=w.pos;
	writer=func(v uint8){
		w.data[p]=v;
	}
	w.pos+=1;
	return ;
}
func (w *packet_encoder)get_uint16_placeholder()(writer func(v uint16)){
	p:=w.pos;
	writer=func(v uint16){
		binary.LittleEndian.PutUint16(w.data[p:],v);
	}
	w.pos+=2;
	return ;
}
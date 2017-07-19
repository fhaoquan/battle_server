package utils

type Packet struct {
	I_cached_data;
	LEN uint16;
	UID uint32;
	RID uint32;
	BDY []byte;
}
func (me *Packet)Cast(f func(*Packet))(*Packet){
	f(me);
	return me;
}
func NewPacket(pool *PacketPool)(*Packet){
	return pool.GetEmptyPkt().(*Packet);
}
func rrrr(){
	go NewPacket(nil).Cast(func(packet *Packet) {
		packet.LEN=10;
	}).Cast(func(packet *Packet) {

	}).ReturnToPool();
}

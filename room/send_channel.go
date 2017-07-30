package room


type send_channel struct {
	kcp_channel map[uint32](chan []byte);
	udp_channel chan int;
}
func (*send_channel)SendUDP(uid uint32,bdy []byte){

}
func (*send_channel)SendKCP(uid uint32,bdy []byte){

}
func (*send_channel)BroadcastUDP(bdy []byte){

}
func (*send_channel)BroadcastKCP(bdy []byte){

}

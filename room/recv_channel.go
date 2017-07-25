package room

import "../utils"

type room_message struct{
	utils.I_cached_data;
	len uint16;
	uid uint32;
	rid uint32;
	bdy []byte;
}
func (pkt *room_message)Clear(){
	pkt.len=0;
	pkt.uid=0;
	pkt.rid=0;
}

type recv_channel struct{
	recv_msg_chan chan *room_message;
	cache *utils.PacketPool
}

func (r *recv_channel)OnPkt(uid uint32,rid uint32,body []byte){
	msg:=r.cache.GetEmptyPkt().(*room_message);
	msg.len=uint16(len(body));
	msg.uid=uid;
	msg.rid=rid;
	copy(msg.bdy,body);
	r.recv_msg_chan<-msg;
}

func new_recv_channel()(*recv_channel){
	return &recv_channel{
		make(chan *room_message,32),
		utils.NewPacketPool(128,func(i utils.I_cached_data)utils.I_cached_data{
			return &room_message{
				i,0,0,0,make([]byte,utils.MaxPktSize),
			}
		}),
	}
}
package proc

import "../utils"
type i_net_message interface {

}
type i_net_message_pool interface {
	GetEmptyMsg()i_net_message;
}
type i_net_connection interface {
	ReadMsg(msg i_net_message);
}
type i_message_handle interface {
	HandleMsg(msg i_net_message);
}
func StartRecvProc(
	conn i_net_connection,pool *utils.MemoryPool)chan *utils.CachedData{

	return nil;
}
func start_recv_proc(conn i_net_connection,pool i_net_message_pool,dest i_message_handle){
	recv_chan:=make(chan i_net_message,32);
	go func(){
		for{
			msg:=pool.GetEmptyMsg();
			conn.ReadMsg(msg);
			recv_chan<-msg;
		}
	}();
	go func(){
		for{
			select {
			case msg:=<-recv_chan:
				dest.HandleMsg(msg);
			}
		}
	}();
	go func(){

	}();
}
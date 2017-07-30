package room

type player struct{
	flag int;
	id uint32;
	name string;
	send_channel chan []byte;
}
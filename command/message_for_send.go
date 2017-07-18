package command

import "../utils"

type I_message_for_send interface {
	utils.I_cached_data;
	GetMessageBody()[]byte;
	SetMessageLen(v uint16);
	BroadCast();
	Send(player_id uint32);
}

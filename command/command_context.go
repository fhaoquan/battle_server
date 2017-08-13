package command

import (
	"../room"
	"../utils"
)

type CommandContext struct {
	base_room *room.Room;
	kcp_res_pool *utils.MemoryPool;
	udp_res_pool *utils.MemoryPool;
}
func (cmd *CommandContext)SetRoom(r *room.Room){
	cmd.base_room=r;
}
func NewCommandContext()(*CommandContext){
	return &CommandContext{
		nil,
		utils.NewMemoryPool(8, func(impl utils.ICachedData)utils.ICachedData{
			return &kcp_response{
				impl,false,0,make([]byte,utils.MaxPktSize),
			}
		}),
		utils.NewMemoryPool(8, func(impl utils.ICachedData)utils.ICachedData{
			return &kcp_response{}
		}),
	}
}

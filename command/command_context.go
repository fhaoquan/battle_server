package command

import (
	"../battle"
	"../utils"
)
type iRoom interface {
	GetBattle()(*battle.Battle)
}
type CommandContext struct {
	base_room iRoom;
	kcp_res_pool *utils.MemoryPool;
	udp_res_pool *utils.MemoryPool;
}
func (cmd *CommandContext)SetRoom(r iRoom){
	cmd.base_room=r;
}
func NewCommandContext()(*CommandContext){
	return &CommandContext{
		nil,
		utils.NewMemoryPool(8, func(impl utils.ICachedData)utils.ICachedData{
			return &kcp_response{
				impl,false,0,0,make([]byte,utils.MaxPktSize),
			}
		}),
		utils.NewMemoryPool(8, func(impl utils.ICachedData)utils.ICachedData{
			return &kcp_response{}
		}),
	}
}

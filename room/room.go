package room

import (
	"../battle"
	"../sessions/packet"
	"sync"
	"github.com/sirupsen/logrus"
	"fmt"
	"net"
)

const MAX_CMD_ID  = 255;
type players struct{
	m sync.RWMutex;
	all map[uint32]*Player;
}
func (me *players)ForEachPlayer(f func(player *Player)bool){
	defer func(){
		me.m.RUnlock();
	}();
	me.m.RLock();
	for _,p:=range me.all{
		if(!f(p)){
			return ;
		}
	}
}
func (me *players)ForOnePlayer(uid uint32,f func(player *Player)){
	defer func(){
		me.m.RUnlock();
	}();
	me.m.RLock();
	f(me.all[uid]);
}
type Room struct{
	id uint32;
	*players;
	the_battle *battle.Battle;
	cmd_handlers []func([]byte)interface{};
	timer_handlers []func();
}
func (r *Room)GetID()uint32{
	return r.id;
}
func (r *Room)SetID(v uint32)*Room{
	r.id=v;
	return r;
}
func (r *Room)GetBattle()*battle.Battle{
	return r.the_battle;
}
func (r *Room)GetCommand(id byte)(func([]byte)interface{}){
	return r.cmd_handlers[id];
}
func (r *Room)SendKCP(response packet.IKcpResponse){
	r.ForOnePlayer(response.GetUID(), func(player *Player) {
		player.SendKCP(response);
	})
}
func (r *Room)SendUDP(response packet.IUdpResponse){
	r.ForOnePlayer(response.GetUID(), func(player *Player) {
		player.SendUDP(response);
	})
}
func (r *Room)OnUDP(adr net.Addr,len uint16,uid uint32,rid uint32,bdy []byte){
	b:=false;
	r.ForOnePlayer(uid, func(player *Player) {
		player.SetUDPAddr(adr);
		b=true;
	})
	if(b){
		r.OnPacket(bdy);
	}
}
func (r *Room)OnKCP(len uint16,uid uint32,rid uint32,bdy []byte){
	r.OnPacket(bdy);
}
func (r *Room)OnPacket(bdy []byte){
	switch rtn:=r.GetCommand(bdy[0])(bdy[1:]);rtn.(type){
	case nil:
		return ;
	case packet.IKcpResponse:
		r.SendKCP(rtn.(packet.IKcpResponse));
	case []packet.IKcpResponse:
		for _,v:=range rtn.([]packet.IKcpResponse){
			r.SendKCP(v);
		};
	case packet.IUdpResponse:
		r.SendUDP(rtn.(packet.IUdpResponse));
	case []packet.IUdpResponse:
		for _,v:=range rtn.([]packet.IUdpResponse){
			r.SendUDP(v);
		};
	case error:
		logrus.Error(rtn.(error));
	default:
		logrus.Error(fmt.Sprint("unknown command response type! ",bdy[0]));
		return ;
	}
}
func NewRoom()(*RoomBuilder){
	return &RoomBuilder{
		&Room{
			0,
			&players{
				all:make(map[uint32]*Player),
			},
			battle.NewBattle(),
			make([]func([]byte)interface{},MAX_CMD_ID),
			make([]func(),10),
		},
	}
}
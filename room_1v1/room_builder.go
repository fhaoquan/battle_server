package room_1v1

import (
	"sync/atomic"
	"../utils"
	"../battle"
	"time"
	"sync"
)

func BuildRoom(
	context *RoomBuildContext,
	plrs ...interface{
		GetPlayerID()uint32;
		GetPlayerName()string;
		GetUnits()[]battle.Unit;
	},
)(*room,error){
	id:=atomic.AddUint32(&room_id_seed,1);
	r:=&room{
		utils.NewGuid().StringUpper(),
		id,
		time.Now(),
		time.Second*(time.Duration(context.Lifecycle)),
		context.WinScore,
		battle.NewBattle1v1(),
		make(chan interface{},1),
		make(chan interface{},1),
		make(chan utils.I_REQ,16),
		new(sync.Once),
		new(sync.Once),
		new(sync.WaitGroup),
		time.Second*(time.Duration(context.SuddenDeath)),
		1,
		&player{plrs[0].GetPlayerID(),plrs[0].GetPlayerName(),nil},
		&player{plrs[1].GetPlayerID(),plrs[1].GetPlayerName(),nil},
	}
	for _,p:=range plrs{
		for _,u:=range p.GetUnits(){
			r.the_battle.CreateUnitDo(func(unit *battle.Unit){
				unit.SetAll(&u);
				r.the_battle.AddMainBaseID(unit.ID);
			});
		}
	}
	add_room(id,r);
	return r,nil;
}

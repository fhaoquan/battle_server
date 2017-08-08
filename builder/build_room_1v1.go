package builder

import (
	"../room"
	"../command"
	"../world"
)

type BuildRoom1V1Context struct {
	w *world.World;
	p []interface{
		GetPlayerID()uint32;
		GetPlayerName()string;
	};
}
func(context *BuildRoom1V1Context)AtWorld(w *world.World)(*BuildRoom1V1Context){
	context.w=w;
	return context;
}
func(context *BuildRoom1V1Context)WaitPlayers(plrs ...interface{
	GetPlayerID()uint32;
	GetPlayerName()string;
})(*BuildRoom1V1Context){
	context.p=plrs;
	return context;
}
func(context *BuildRoom1V1Context)DoBuild()(*room.Room,error){
	c:=&command.Commamd{};
	f:=func()(*room.Room){
		return room.NewRoom().
			RouteCommand(002,c.CreateUnit).
			RouteCommand(003,c.UnitAttackStart).
			RouteCommand(004,c.UnitAttackDone).
			RouteCommand(005,c.UpdateUnitMovement).
			RouteTimer(001,c.BroadcastBattleData).
			WithPlayers(context.p[0],context.p[1]).
			Build();
	}
	err:=error(nil);
	res:=(*room.Room)(nil);
	context.w.AddNewRoom(
		func(id uint32)*room.Room{
			if b,e:=BuildUdpSession(int(id));e==nil{
				r:=f();
				r.SetID(b(r));
				res=r;
				c.SetRoom(r);
				return r;
			}else{
				err=e;
				return nil;
			}
		});
	return res,nil;
}

func BuildRoom1V1()*BuildRoom1V1Context{
	return &BuildRoom1V1Context{};
}
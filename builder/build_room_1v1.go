package builder

import (
	"../room"
	"../command"
	"../world"
	"../sessions/udp_session"
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
	s,e:=udp_session.NewSession();
	if e!=nil{
		return nil,e;
	}
	c:=&command.NewCommandContext();
	r:=room.NewRoom().
		RouteCommand(002,c.CreateUnit).
		RouteCommand(003,c.UnitAttackStart).
		RouteCommand(004,c.UnitAttackDone).
		RouteCommand(005,c.UpdateUnitMovement).
		RouteTimer(001,c.BroadcastBattleMovementData).
		WithPlayers(context.p[0],context.p[1]).
		Build();
	r.SetID((uint32(s.GetAddr().Port)));
	c.SetRoom(r);
	s.StartAt(r);
	return r,nil;
}

func BuildRoom1V1()*BuildRoom1V1Context{
	return &BuildRoom1V1Context{};
}
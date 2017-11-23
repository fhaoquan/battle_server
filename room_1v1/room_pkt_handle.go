package room_1v1

import (
	"../utils"
	"../battle"
	"github.com/sirupsen/logrus"
	"fmt"
	"runtime/debug"
)

func (me *room) send_udp_response(r utils.I_RES) {
	defer func() {
		r.Return()
		if e := recover(); e != nil {
			logrus.Error(e)
			logrus.Error(fmt.Sprintf("%s", debug.Stack()))
		}
	}()
	switch r.GetUID() {
	case 0:
		me.p1.send_udp(r.GetSendData());
		me.p2.send_udp(r.GetSendData());
	case me.p1.uid:
		me.p1.send_udp(r.GetSendData());
	case me.p2.uid:
		me.p2.send_udp(r.GetSendData());
	}
}
func (me *room) send_kcp_response(r utils.I_RES) {
	defer func() {
		r.Return()
		if e := recover(); e != nil {
			logrus.Error(e)
			logrus.Error(fmt.Sprintf("%s", debug.Stack()))
		}
	}()
	switch r.GetUID() {
	case 0:
		me.p1.send_kcp(r.GetSendData());
		me.p2.send_kcp(r.GetSendData());
	case me.p1.uid:
		me.p1.send_kcp(r.GetSendData());
	case me.p2.uid:
		me.p2.send_kcp(r.GetSendData());
	}
}
func (me *room) send_response(pkt utils.I_RES){
	switch pkt.Protocol() {
	case utils.Protocol_KCP:
		me.send_kcp_response(pkt.(utils.I_RES))
	case utils.Protocol_UDP:
		me.send_udp_response(pkt.(utils.I_RES))
	}
}
func (me *room) handle_battle_result(who uint32, rtn interface{}) {
	switch rtn.(type) {
	case nil:
		return
	case utils.I_RES:
		me.send_response(rtn.(utils.I_RES))
		rtn.(utils.I_RES).Return();
	case []utils.I_RES:
		for _, r := range rtn.([]utils.I_RES) {
			me.send_response(r.(utils.I_RES))
			r.Return();
		}
	case *battle.BattlePanicError:
		me.room_log_err(rtn.(*battle.BattlePanicError).E)
	case error:
		me.room_log_err(rtn.(error))
	default:
		me.room_log_err("unknown command response type! ")
		return
	}
}
func (me *room)on_packet(who uint32, bdy []byte){
	switch bdy[0] {
	case utils.CMD_battle_start:
		me.handle_battle_result(who,me.the_battle.BroadcastBattleStart(who))
		me.handle_battle_result(who,me.the_battle.BroadcastBattleAll(who))
	case utils.CMD_pingpong:
		me.handle_battle_result(who, me.the_battle.Pong(who, bdy[1:]))
	case utils.CMD_unit_movment:
		me.handle_battle_result(who, me.the_battle.UpdateUnitMovement(bdy[1:]))
	case utils.CMD_attack_start:
		me.handle_battle_result(who, me.the_battle.UnitAttackStart(bdy[1:]))
	case utils.CMD_attack_done:
		me.handle_battle_result(who, me.the_battle.UnitAttackDone(who, bdy[1:]))
	case utils.CMD_create_unit:
		me.handle_battle_result(who, me.the_battle.CreateUnit(who, bdy[1:]))
	case utils.CMD_unit_destory:
		me.handle_battle_result(who, me.the_battle.UnitDestory(bdy[1:]))
	}
}

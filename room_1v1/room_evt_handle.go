package room_1v1

import (
	"../gateway"
	)
type event_session_connected struct {
	session gateway.Session;
}
type event_session_closed struct {
	session gateway.Session;
}

func (me *room) on_event_session_connected(event *event_session_connected){
	switch event.session.UserID(){
	case me.p1.uid:
		me.p1.session=event.session;
	case me.p2.uid:
		me.p2.session=event.session;
	default:
		return ;
	}
	event.session.Start(me.packet_chan);
}
func (me *room) on_event_session_closeed(event *event_session_closed){
	switch event.session.UserID(){
	case me.p1.uid:
		me.p1.session=nil;
	case me.p2.uid:
		me.p2.session=nil;
	default:
		return ;
	}
}

func (me *room) on_event(event interface{}) {
	switch event.(type) {
	case *event_session_connected:
		me.on_event_session_connected(event.(*event_session_connected))
	case *event_session_closed:
		me.on_event_session_closeed(event.(*event_session_closed))
	}
}

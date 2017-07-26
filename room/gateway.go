package room

import (
	"../sessions/kcp_session"
	"sync"
)
type yyy struct {
	c chan int;
	next *yyy;
}
type Gateway struct {
	m sync.RWMutex;
	room_id_seed uint32;
	rooms map[uint32]*Room;
}
func (y *yyy)Each(f func(*yyy))int{

}
func (y *yyy)Filter(f func(*yyy)bool)*yyy{
	if f(y){
		if(y.next!=nil){
			return &yyy{c:y.c,next:y.next.Filter(f)};
		}else{
			return &yyy{y.c,nil}
		}
	}else{
		if(y.next!=nil){
			return y.next.Filter(f);
		}else{
			return nil;
		}
	}
}
func (me *Gateway)HandleUnRoomedSession(){
	o:=make([]int,1024);
	i:=0;
	f:=&yyy{nil,nil}.
		Filter(func(t *yyy)bool{
		return (len(t.c)>0);
	}).
		Filter(func(t *yyy)bool{
		return (len(t.c)>0);
	}).
		Each(func(t *yyy){
		o[i]=<-t.c;
		i++;
	})
	cs:=make([]yyy,15);
	select {
	case cs[0].c:
	}

}

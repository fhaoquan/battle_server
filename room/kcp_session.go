package room

import (
	"net"
	"fmt"
	"errors"
	"../utils"
)

type kcp_session struct {
	conn net.Conn;
	uid uint32;
}
func (me *kcp_session)Close(){
	me.conn.Close();
}
func (me *kcp_session)Send(b []byte)(err error){
	defer func(){
		if e:=recover();e!=nil{
			err=errors.New(fmt.Sprint(e));
		}
	}()
	_,err=me.conn.Write(b);
	return ;
}
func (me *kcp_session)Recv(empty utils.IKcpRequest)(err error){
	defer func(){
		if e:=recover();e!=nil{
			err=errors.New(fmt.Sprint(e));
		}
	}()
	if e:=empty.ReadAt(me.conn);e!=nil{
		return e;
	}
	if empty.GetUID()!=me.uid{
		return errors.New("error uid");
	}
	return nil;
}
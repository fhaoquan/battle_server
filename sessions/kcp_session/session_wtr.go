package kcp_session

import (
	"../../utils"
	"github.com/pkg/errors"
	"net"
	"encoding/binary"
	"fmt"
)
func SendUtilError(conn net.Conn,pop func([]byte)int)(e error){
	defer func(){
		if err:=recover();err!=nil{
			e=errors.New(fmt.Sprint(err));
		}
	}();
	buf:=make([]byte,utils.MaxPktSize);
	for{
		l:=pop(buf[6:]);
		if l<=0{
			return errors.New("stoped for pop len")
		}
		binary.BigEndian.PutUint16(buf[4:6],(uint16)(l+2));
		binary.BigEndian.PutUint32(buf[0:4],123454321);
		if _,e:=conn.Write(buf[0:l+6]);e!=nil{
			return e;
		}
	}
}
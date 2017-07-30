package udp_session

import (
	"../../utils"
	"fmt"
	"net"
	"encoding/binary"
	"github.com/pkg/errors"
	"github.com/sirupsen/logrus"
)

func RecvUntilFalse(conn net.PacketConn,receiver func(net.Addr,uint16,uint32,uint32,[]byte)(bool))(e error){
	defer func(){
		if err:=recover();err!=nil{
			e=errors.New(fmt.Sprint(err));
		}
	}();
	buf:=make([]byte,utils.MaxPktSize)
	for{
		_,a,e:=conn.ReadFrom(buf);
		if e!=nil{
			logrus.Error(e);
			continue;
		}
		if(binary.BigEndian.Uint32(buf[0:4])!=123454321){
			logrus.Error(errors.New("packet read error!"));
			continue;
		}
		l:=binary.BigEndian.Uint16(buf[4:6]);
		u:=binary.BigEndian.Uint32(buf[6:10]);
		r:=binary.BigEndian.Uint32(buf[10:14]);
		if e:=receiver(a,l,u,r,buf[14:l+4]);!e{
			return errors.New(fmt.Sprintf("stoped for receiver error=",e))
		}
	}
	return nil;
}

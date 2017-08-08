package proc

import (
)

type ISession interface {

	Read();
}

func drecv_proc(c chan ISession){
	for{
		s:=<-c;
		s.Read();
		c<-s;
	}
}
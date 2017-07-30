package room

import "../utils"

type owner struct {

}

func StartPlayerProc(r *Room){
	msg:=make(chan utils.IDataOwner,16);
	sig:=make(chan int,1);
	snd:=make(chan utils.IDataOwner,16);
	for{
		select {
		case <-msg:
		case <-snd:
		case <-sig:
		}
	}
}

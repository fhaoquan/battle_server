package battle_1v1

import "../utils"
func player_session_proc(IN chan interface{}){
	for{
		select {
		case i,ok:=<-IN:
			if ok{
				switch i.(type) {
				case utils.KcpReq:
					i.(utils.KcpReq).Return();
				case utils.UdpReq:
					i.(utils.UdpReq).Return();
				}
			}else{
				return ;
			}

		}
	}

	return ;
}

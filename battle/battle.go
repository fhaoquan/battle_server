package battle

import (
	"../utils"
	"container/list"
)

const unit_id_offset	=1000;
const max_unit_count	=2000;
type Battle struct {
	kcp_res_pool *utils.MemoryPool;
	udp_res_pool *utils.MemoryPool;
	all_units []*Unit;
	living_units *list.List;
}
func (context *Battle)GetFreeID()uint16{
	for i,u:=range context.all_units{
		if u==nil{
			return uint16(1000+i);
		}
	}
	return 0;
}
func (context *Battle)AllUnit()[]*Unit{
	return context.all_units;
}
func (context *Battle)FindUnit(id uint16)*Unit{
	return context.all_units[id-1000];
}
func (context *Battle)NewUnit(id uint16)*Unit{
	context.all_units[id-1000]=NewUnit(id);
	context.living_units.PushBack(id);
	return context.all_units[id-1000];
}
func (context *Battle)ForEachUnitDo(f func(*Unit)(bool)){
	for _,u:=range context.all_units{
		if u==nil{
			continue;
		}
		if !f(u){
			return ;
		}
	}
}
func (context *Battle)FindUnitDo(id uint16,f func(*Unit)){
	f(context.FindUnit(id));
}
func (context *Battle)CreateUnitDo(f func(*Unit)){
	f(context.NewUnit(context.GetFreeID()));
}
func NewBattle()*Battle{
	return &Battle{
		utils.NewMemoryPool(8, func(impl utils.ICachedData)utils.ICachedData{
			return &utils.KcpRes{
				impl,false,0,make([]byte,utils.MaxPktSize),
			}
		}),
		utils.NewMemoryPool(8, func(impl utils.ICachedData)utils.ICachedData{
			return &utils.UdpRes{
				&utils.KcpRes{
					impl,false,0,make([]byte,utils.MaxPktSize),
				},
			}
		}),
		make([]*Unit,max_unit_count),
		list.New(),
	};
}
package battle

const unit_id_offset	=1000;
const max_unit_count	=2000;
type Battle struct {
	all_units []*Unit;
}
func (context *Battle)AllUnit()[]*Unit{
	return context.all_units;
}
func (context *Battle)FindUnit(id uint16)*Unit{
	return context.all_units[id-1000];
}
func (context *Battle)NewUnit(id uint16)*Unit{
	context.all_units[id-1000]=NewUnit(id);
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
func (context *Battle)CreateUnitDo(id uint16,f func(*Unit)){
	f(context.NewUnit(id));
}
type FindUnit struct{
	B *Battle;
};
func NewBattle()*Battle{
	return &Battle{
		make([]*Unit,max_unit_count),
	};
}